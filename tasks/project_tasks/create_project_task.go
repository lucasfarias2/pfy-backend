package project_tasks

import (
	"fmt"
	"os"
	"packlify-cloud-backend/models"
	"packlify-cloud-backend/utils"
)

func CreateProjectTask(projectRequest models.Project) (models.Project, error) {
	db := utils.GetDB()
	var newProject models.Project

	githubRepo := fmt.Sprintf("https://github.com/%s/%s.git", os.Getenv("GITHUB_OWNER"), projectRequest.Name)

	db.QueryRow("INSERT INTO projects(name, organization_id, toolkit, github_repo) VALUES($1, $2, $3, $4) RETURNING id, name, organization_id, toolkit, github_repo", projectRequest.Name, projectRequest.OrganizationID, projectRequest.Toolkit, githubRepo).Scan(&newProject.ID, &newProject.Name, &newProject.OrganizationID, &newProject.Toolkit, &newProject.GitHubRepo)

	return newProject, nil
}
