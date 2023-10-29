package project_tasks

import (
	"packlify-cloud-backend/models"
	"packlify-cloud-backend/models/constants"
	"packlify-cloud-backend/services"
)

func GenerateFilesFromToolkitTask(tm *services.TaskManager, project models.Project, projectGenerateFilesFromToolkit chan bool, errs chan error) {
	task, err := tm.CreateTask(project.ID, constants.Running, "", string(constants.PROJECT_GENERATE_FILES))

	err = services.CreateSDKApp(project.Name)
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
		return
	}

	projectGenerateFilesFromToolkit <- true
}
