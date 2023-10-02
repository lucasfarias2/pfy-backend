package services

import (
	"packlify-cloud-backend/models"
	"packlify-cloud-backend/utils"
)

func CreateProject(project models.Project) (models.Project, error) {
	db := utils.GetDB()
	err := db.QueryRow("INSERT INTO projects(name) VALUES($1) RETURNING id", project.Name).Scan(&project.ID)
	if err != nil {
		return models.Project{}, err
	}
	return project, nil
}
