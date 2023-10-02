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

func GetAllProjects() ([]models.Project, error) {
	db := utils.GetDB()
	rows, err := db.Query("SELECT * FROM projects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		if err := rows.Scan(&project.ID, &project.Name /* other fields */); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}
