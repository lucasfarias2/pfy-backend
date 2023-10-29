package project_tasks

import (
	"log"
	"packlify-cloud-backend/models"
	"packlify-cloud-backend/models/constants"
	"packlify-cloud-backend/services"
)

func PushToGithubTask(tm *services.TaskManager, newProject models.Project, projectCreateGithubRepository chan bool, projectPushToGithubRepository chan bool, errs chan error) {
	<-projectCreateGithubRepository

	task, _ := tm.CreateTask(newProject.ID, "Running", "", string(constants.PROJECT_PUSH_GITHUB))

	err := services.PushToGitHubRepo(newProject.Name)
	if err != nil {
		err := tm.UpdateTaskStatus(task.ID, "Failed", err.Error())
		errs <- err
		return
	}

	_ = tm.UpdateTaskStatus(task.ID, "Success", "")

	err = services.CleanClonedFolder()
	if err != nil {
		log.Print("Error cleaning cloned folder")
	}

	projectPushToGithubRepository <- true
}
