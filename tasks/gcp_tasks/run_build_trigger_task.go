package gcp_tasks

import (
	"packlify-cloud-backend/models"
	"packlify-cloud-backend/models/constants"
	"packlify-cloud-backend/models/tasks_models"
	"packlify-cloud-backend/services"
	"packlify-cloud-backend/services/gcp"
)

func RunBuildTriggerTask(tm *services.TaskManager, newProject chan models.Project, gcpCreateBuildTrigger chan tasks_models.BuildTriggerData, gcpRunBuildTrigger chan bool, errs chan error) {
	gcpCreateBuildData := <-gcpCreateBuildTrigger
	project := <-newProject

	task, err := tm.CreateTask(project.ID, constants.Running, "", string(constants.GCP_RUN_BUILD_TRIGGER))
	if err != nil {
		errs <- err
		return
	}

	err = gcp.RunBuildTrigger(project, gcpCreateBuildData.Trigger)

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
}
