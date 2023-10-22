package controllers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"packlify-cloud-backend/models"
	"packlify-cloud-backend/models/constants"
	"packlify-cloud-backend/services"
	"strconv"
	"time"
)

func CreateProject(c *fiber.Ctx) error {
	project := new(models.Project)
	if err := c.BodyParser(project); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	newProject, err := services.CreateProject(*project)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	tm := services.NewTaskManager()
	createCloudBuildDone := make(chan bool)
	createCloudRunDone := make(chan bool)
	errs := make(chan error)

	go func() {
		task, err := tm.CreateTask(newProject.ID, constants.Running, "Creating Cloud Build", 4)

		if err != nil {
			errs <- err
			return
		}

		// time.Sleep(20 * time.Second) // Faking the time it takes to complete the task

		err = services.CreateCloudRun(newProject)
		if err != nil {
			err := tm.UpdateTaskStatus(task.ID, "Failed", err.Error())
			if err != nil {
				return
			}
			errs <- err
			return
		}

		err = tm.UpdateTaskStatus(task.ID, constants.Success, "Cloud build was created successfully")
		if err != nil {
			return
		}
		createCloudBuildDone <- true
	}()

	go func() {
		<-createCloudBuildDone // wait for deployment to finish

		// Create a new task for integration
		task, err := tm.CreateTask(newProject.ID, constants.Running, "Creating Cloud Run", 5)
		if err != nil {
			errs <- err
			return
		}

		//err = services.CreateCloudRun(newProject)
		if err != nil {
			err := tm.UpdateTaskStatus(task.ID, "Failed", err.Error())
			errs <- err
			return
		}

		// Update task status to Success
		err = tm.UpdateTaskStatus(task.ID, constants.Success, "Cloud Run was created successfully")
		createCloudRunDone <- true
	}()

	return c.JSON(newProject)
}

func GetAllProjects(c *fiber.Ctx) error {
	projects, err := services.GetAllProjects()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(projects)
}

func GetProjectStatus(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, "text/event-stream")
	c.Set(fiber.HeaderCacheControl, "no-cache")
	c.Set(fiber.HeaderConnection, "keep-alive")
	c.Set("Access-Control-Allow-Origin", "*")

	projectId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		for {
			tasks, err := services.GetProjectStatusById(projectId)
			if err != nil {
				continue
			}

			data, err := json.Marshal(tasks)
			if err != nil {
				continue
			}

			msg := fmt.Sprintf("data: %s\n\n", data)
			_, err = w.WriteString(msg)

			err = w.Flush()
			if err != nil {
				break
			}
			time.Sleep(2 * time.Second)
		}
	})

	return nil
}
