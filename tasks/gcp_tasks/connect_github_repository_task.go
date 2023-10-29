package gcp_tasks

import (
	"packlify-cloud-backend/models"
	"packlify-cloud-backend/models/constants"
	"packlify-cloud-backend/services"
	"packlify-cloud-backend/services/gcp"
)

func ConnectGithubRepositoryTask(tm *services.TaskManager, newProject chan models.Project, gcpConnectNewRepository chan bool, errs chan error) {
	project := <-newProject

	task, err := tm.CreateTask(project.ID, constants.Running, "", string(constants.GCP_CONNECT_REPOSITORY))
	if err != nil {
		return
	}

	err = gcp.ConnectGithubRepository(project)
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

	gcpConnectNewRepository <- true
}
