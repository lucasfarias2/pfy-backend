package controllers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"packlify-cloud-backend/models"
	"packlify-cloud-backend/models/tasks_models"
	"packlify-cloud-backend/services"
	"packlify-cloud-backend/tasks/gcp_tasks"
	"packlify-cloud-backend/tasks/project_tasks"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateProject(c *fiber.Ctx) error {
	projectRequest := new(models.Project)
	if err := c.BodyParser(projectRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	tm := services.NewTaskManager()
	errs := make(chan error)

	projectGenerateFilesFromToolkit := make(chan bool)
	projectCreateGithubRepository := make(chan bool)
	projectPushToGithubRepository := make(chan bool)

	gcpConnectNewRepository := make(chan bool)
	gcpCreateArtifactRepository := make(chan bool)
	gcpCreateBuildTrigger := make(chan tasks_models.BuildTriggerData)
	gcpRunBuildTrigger := make(chan bool)

	// Project for Frontend Toolkit tasks
	newProject, err := project_tasks.CreateProjectTask(tm, *projectRequest, errs)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Error creating project"})
	}

	go project_tasks.GenerateFilesFromToolkitTask(tm, newProject, projectGenerateFilesFromToolkit, errs)
	go project_tasks.CreateGithubRepositoryTask(tm, newProject, projectGenerateFilesFromToolkit, projectCreateGithubRepository, errs)
	go project_tasks.PushToGithubTask(tm, newProject, projectCreateGithubRepository, projectPushToGithubRepository, errs)

	// Google Cloud Platform for Frontend Toolkit tasks
	go gcp_tasks.ConnectGithubRepositoryTask(tm, newProject, gcpConnectNewRepository, errs)
	go gcp_tasks.CreateArtifactRepositoryTask(tm, newProject, gcpCreateArtifactRepository, errs)
	go gcp_tasks.CreateBuildTriggerTask(tm, newProject, gcpCreateBuildTrigger, gcpCreateArtifactRepository, gcpConnectNewRepository, errs)
	go gcp_tasks.RunBuildTriggerTask(tm, newProject, gcpCreateBuildTrigger, gcpRunBuildTrigger, errs)

	return c.JSON(newProject)
}

func GetAllProjects(c *fiber.Ctx) error {
	organizationId := c.Query("organization_id")
	if organizationId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid organization ID"})
	}
	projects, err := services.GetAllProjects(organizationId)
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
			if err != nil {
				break
			}

			err = w.Flush()
			if err != nil {
				break
			}
			time.Sleep(2 * time.Second)
		}
	})

	return nil
}
