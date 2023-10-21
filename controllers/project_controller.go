package controllers

import (
	"github.com/gofiber/fiber/v2"
	"packlify-cloud-backend/models"
	"packlify-cloud-backend/models/constants"
	"packlify-cloud-backend/services"
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

		time.Sleep(20 * time.Second) // Faking the time it takes to complete the task

		// err = services.CreateCloudBuild(newProject)
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
