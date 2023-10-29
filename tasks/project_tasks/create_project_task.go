package project_tasks

import (
	"fmt"
	"os"
	"packlify-cloud-backend/models"
	"packlify-cloud-backend/models/constants"
	"packlify-cloud-backend/services"
	"packlify-cloud-backend/utils"
)

func CreateProjectTask(tm *services.TaskManager, projectRequest models.Project, errs chan error) (models.Project, error) {
	db := utils.GetDB()
	task, err := tm.CreateTask(projectRequest.ID, constants.Running, "", string(constants.PROJECT_CREATE))
	var newProject models.Project

	githubRepo := fmt.Sprintf("https://github.com/%s/%s.git", os.Getenv("GITHUB_OWNER"), projectRequest.Name)

	err = db.QueryRow("INSERT INTO projects(name, organization_id, toolkit, github_repo) VALUES($1, $2, $3, $4) RETURNING id, name, organization_id, toolkit, github_repo", projectRequest.Name, projectRequest.OrganizationID, projectRequest.Toolkit, githubRepo).Scan(&newProject.ID, &newProject.Name, &newProject.OrganizationID, &newProject.Toolkit, &newProject.GitHubRepo)
	if err != nil {
		err := tm.UpdateTaskStatus(task.ID, "Failed", err.Error())
		if err != nil {
			return models.Project{}, nil
		}
		errs <- err
		return models.Project{}, nil
	}

	err = tm.UpdateTaskStatus(task.ID, "Success", "Project created successfully")

	return newProject, err
}
