package controllers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"packlify-cloud-backend/models"
	"packlify-cloud-backend/models/constants"
	"packlify-cloud-backend/services"
	"packlify-cloud-backend/services/gcp"
	"strconv"
	"time"

	"cloud.google.com/go/cloudbuild/apiv1/v2/cloudbuildpb"
	"github.com/gofiber/fiber/v2"
)

type BuildTriggerData struct {
	IsSuccess bool
	Trigger   *cloudbuildpb.BuildTrigger
}

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
	createProjectDone := make(chan bool)
	gcpCreateArtifactRepository := make(chan bool)
	gcpGetGitHubAppInstallationId := make(chan int)
	gcpCreateBuildTrigger := make(chan BuildTriggerData)
	gcpRunBuildTrigger := make(chan bool)
	errs := make(chan error)

	go func() {
		githubToken := os.Getenv("GITHUB_ACCESS_TOKEN")

		appInstallationId, err := services.FetchAppInstallations(githubToken)
		if err != nil {
			return
		}

		gcpGetGitHubAppInstallationId <- appInstallationId
	}()

	go func() {

		appInstallationId, err := services.ConnectGitHubWithCloudBuild()
		if err != nil {
			return
		}

		gcpGetGitHubAppInstallationId <- appInstallationId
	}()

	go func() {
		createProjectDone <- true
	}()

	go func() {
		<-createProjectDone
		task, err := tm.CreateTask(newProject.ID, constants.Running, "", string(constants.GCP_CREATE_ARTIFACT_REPOSITORY))

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

	go func() {
		<-gcpCreateArtifactRepository
		task, err := tm.CreateTask(newProject.ID, constants.Running, "", string(constants.GCP_CREATE_BUILD_TRIGGER))
		if err != nil {
			errs <- err
			return
		}

		trigger, err := gcp.CreateBuildTrigger(newProject)

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

		gcpCreateBuildTrigger <- BuildTriggerData{
			IsSuccess: true,
			Trigger:   trigger,
		}
	}()

	go func() {
		gcpCreateBuildData := <-gcpCreateBuildTrigger

		task, err := tm.CreateTask(newProject.ID, constants.Running, "", string(constants.GCP_RUN_BUILD_TRIGGER))
		if err != nil {
			errs <- err
			return
		}

		err = gcp.RunBuildTrigger(newProject, gcpCreateBuildData.Trigger)

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
	organization_id := c.Query("organization_id")
	if organization_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid organization ID"})
	}
	projects, err := services.GetAllProjects(organization_id)
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
