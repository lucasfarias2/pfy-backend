package controllers

import (
	"github.com/gofiber/fiber/v2"
	"packlify-cloud-backend/models"
	"packlify-cloud-backend/services"
)

func CreateProject(c *fiber.Ctx) error {
	project := new(models.Project)
	if err := c.BodyParser(project); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	newProject, err := services.CreateProject(*project)

	tm := services.NewTaskManager()
	createCloudBuildDone := make(chan bool)
	//createCloudRunDone := make(chan bool)
	errs := make(chan error)

	go func() {
		task, err := tm.CreateTask(newProject.ID, "Pending", "Task will run soon...", 4)

		if err != nil {
			errs <- err
			return
		}

		// err = services.CreateCloudBuild(newProject)
		if err != nil {
			err := tm.UpdateTaskStatus(task.ID, "Failed", err.Error())
			if err != nil {
				return
			}
			errs <- err
			return
		}

		err = tm.UpdateTaskStatus(task.ID, "Success", "Cloud build was created successfully")
		if err != nil {
			return
		}
		createCloudBuildDone <- true
	}()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(newProject)
}

func GetAllProjects(c *fiber.Ctx) error {
	projects, err := services.GetAllProjects()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(projects)
}
