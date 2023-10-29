package gcp_tasks

import (
	"packlify-cloud-backend/models"
	"packlify-cloud-backend/models/constants"
	"packlify-cloud-backend/models/tasks_models"
	"packlify-cloud-backend/services"
	"packlify-cloud-backend/services/gcp"
)

func CreateBuildTriggerTask(tm *services.TaskManager, newProject models.Project, gcpCreateBuildTrigger chan tasks_models.BuildTriggerData, gcpCreateArtifactRepository chan bool, errs chan error) {
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

	gcpCreateBuildTrigger <- tasks_models.BuildTriggerData{
		Trigger: trigger,
	}
}
