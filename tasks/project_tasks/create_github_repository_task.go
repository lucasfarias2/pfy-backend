package project_tasks

import (
	"packlify-cloud-backend/models"
	"packlify-cloud-backend/models/constants"
	"packlify-cloud-backend/services"
)

func CreateGithubRepositoryTask(tm *services.TaskManager, newProject models.Project, projectGenerateFilesFromToolkit chan bool, projectCreateGithubRepository chan bool, errs chan error) {
	<-projectGenerateFilesFromToolkit

	task, _ := tm.CreateTask(newProject.ID, constants.Running, "", string(constants.PROJECT_CREATE_GITHUB))

	updatedProject, err := services.CreateGitHubRepo(newProject)
	if err != nil {
		err := tm.UpdateTaskStatus(task.ID, "Failed", err.Error())
		if err != nil {
			return
		}
		errs <- err
		return
	}

	projectCreateGithubRepository <- true
}
