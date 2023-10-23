package controllers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"packlify-cloud-backend/models"
	"packlify-cloud-backend/models/constants"
	"packlify-cloud-backend/services"
	"packlify-cloud-backend/services/gcp"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
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
	gcpCreateArtifactRepository := make(chan bool)
	gcpCreateBuildTrigger := make(chan bool)
	gcpCreateCloudRun := make(chan bool)
	gcpRunBuildTrigger := make(chan bool)
	errs := make(chan error)

	// Create Artifact Repository Task
	go func() {
		task, err := tm.CreateTask(newProject.ID, constants.Running, "", 4)

		if err != nil {
			errs <- err
			return
		}

		err = gcp.CreateArtifactRepository(newProject)
		if err != nil {
			err := tm.UpdateTaskStatus(task.ID, "Failed", err.Error())
			if err != nil {
				return
			}
			errs <- err
			return
		}

		err = tm.UpdateTaskStatus(task.ID, constants.Success, "")
		if err != nil {
			errs <- err
			return
		}

		gcpCreateArtifactRepository <- true
	}()

	// Create Build Trigger - runs in parallel with Create Artifact Repository
	go func() {
		task, err := tm.CreateTask(newProject.ID, constants.Running, "", 5)
		if err != nil {
			errs <- err
			return
		}

		// err = gcp.CreateBuildTrigger(newProject)
		if err != nil {
			err := tm.UpdateTaskStatus(task.ID, "Failed", err.Error())
			errs <- err
			return
		}

		err = tm.UpdateTaskStatus(task.ID, constants.Success, "")
		if err != nil {
			errs <- err
			return
		}

		gcpCreateBuildTrigger <- true
	}()

	// Create Cloud Run - depends on Build Trigger
	go func() {
		<-gcpCreateBuildTrigger

		task, err := tm.CreateTask(newProject.ID, constants.Running, "", 6)
		if err != nil {
			errs <- err
			return
		}

		// err = gcp.CreateCloudRun(newProject)
		if err != nil {
			err := tm.UpdateTaskStatus(task.ID, "Failed", err.Error())
			errs <- err
			return
		}

		// Update task status to Success
		err = tm.UpdateTaskStatus(task.ID, constants.Success, "")
		if err != nil {
			errs <- err
			return
		}

		gcpCreateCloudRun <- true
	}()

	// Run Build Trigger - depends on Create Cloud Run
	go func() {
		<-gcpCreateCloudRun

		task, err := tm.CreateTask(newProject.ID, constants.Running, "", 7)
		if err != nil {
			errs <- err
			return
		}

		// err = gcp.CreateCloudRun(newProject)
		if err != nil {
			err := tm.UpdateTaskStatus(task.ID, "Failed", err.Error())
			errs <- err
			return
		}

		err = tm.UpdateTaskStatus(task.ID, constants.Success, "")
		if err != nil {
			errs <- err
			return
		}

		gcpRunBuildTrigger <- true
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
