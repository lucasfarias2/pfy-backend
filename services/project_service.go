package services

import (
	"fmt"
	"os"
	"packlify-cloud-backend/models"
	"packlify-cloud-backend/utils"
)

func CreateProject(project models.Project) (models.Project, error) {
	db := utils.GetDB()

	err := db.QueryRow("INSERT INTO projects(name, organization_id, toolkit_id) VALUES($1, $2, $3) RETURNING id", project.Name, project.OrganizationID, project.ToolkitID).Scan(&project.ID)
	if err != nil {
		return models.Project{}, err
	}

	var toolkitName string
	err = db.QueryRow("SELECT name FROM toolkits WHERE id=$1", project.ToolkitID).Scan(&toolkitName)
	if err != nil {
		return models.Project{}, err
	}

	// Check if the selected toolkit is React
	if toolkitName == "React" {
		awsConfig := map[string]string{
			"ami":          "ami-0f3164307ee5d695a",
			"instanceType": "t2.micro",
			"name":         project.Name,
		}

		reservation, err := CreateAWSInstance(awsConfig)
		if err != nil {
			return models.Project{}, err
		}

		fmt.Println("AWS Instance Created:", reservation)

		githubAccessToken := os.Getenv("GITHUB_ACCESS_TOKEN")
		newRepo, err := CreateGithubRepo(githubAccessToken, project.Name)
		if err != nil {
			return models.Project{}, err
		}

		fmt.Println("GitHub Repository Created:", *newRepo.CloneURL)
	}

	return project, nil
}

func GetAllProjects() ([]models.Project, error) {
	db := utils.GetDB()
	rows, err := db.Query("SELECT id, name FROM projects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		if err := rows.Scan(&project.ID, &project.Name); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}
