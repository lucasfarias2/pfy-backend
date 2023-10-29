package gcp_tasks

import (
	"packlify-cloud-backend/models"
	"packlify-cloud-backend/models/constants"
	"packlify-cloud-backend/models/tasks_models"
	"packlify-cloud-backend/services"
	"packlify-cloud-backend/services/gcp"
)

func CreateBuildTriggerTask(tm *services.TaskManager, newProject chan models.Project, gcpCreateBuildTrigger chan tasks_models.BuildTriggerData, gcpCreateArtifactRepository chan bool, gcpConnectNewRepository chan bool, errs chan error) {
	<-gcpCreateArtifactRepository
	<-gcpConnectNewRepository
	project := <-newProject

	task, err := tm.CreateTask(project.ID, constants.Running, "", string(constants.GCP_CREATE_BUILD_TRIGGER))
	if err != nil {
		errs <- err
		return
	}

	trigger, err := gcp.CreateBuildTrigger(project)

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

	gcpCreateBuildTrigger <- tasks_models.BuildTriggerData{
		IsSuccess: true,
		Trigger:   trigger,
	}
}
