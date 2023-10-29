package gcp_tasks

import (
	"packlify-cloud-backend/models"
	"packlify-cloud-backend/models/constants"
	"packlify-cloud-backend/services"
	"packlify-cloud-backend/services/gcp"
)

func CreateArtifactRepositoryTask(tm *services.TaskManager, newProject models.Project, gcpConnectNewRepository chan bool, gcpCreateArtifactRepository chan bool, errs chan error) {
	<-gcpConnectNewRepository

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
}
