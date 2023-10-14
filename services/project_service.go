package services

import (
	"database/sql"
	"log"
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
		if err != nil {
			return models.Project{}, err
		}

		// Create the SDK App using the project name
		err := CreateSDKApp(project.Name)

		// Create a new GitHub repo
		githubToken := os.Getenv("GITHUB_ACCESS_TOKEN")
		cloneURL, err := CreateGitHubRepo(project.Name, githubToken)
		if err != nil {
			log.Fatalf("Error creating Github repo: %s", err)
		}

		// Push the project to the new repo
		err = PushToGitHubRepo(project.Name, cloneURL)
		if err != nil {
			log.Fatalf("Error pushing Github repo: %s", err)
		}

		if err != nil {
			return models.Project{}, err
		}
	}

	return project, nil
}

func GetAllProjects() ([]models.Project, error) {
	db := utils.GetDB()
	rows, err := db.Query("SELECT id, name FROM projects")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Print("Error getting projects")
		}
	}(rows)

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

func GetProjectById(id int) (models.Project, error) {
	db := utils.GetDB()
	var project models.Project

	// Query the database for the project with the specified ID
	err := db.QueryRow("SELECT id, name FROM projects WHERE id=$1", id).Scan(&project.ID, &project.Name)
	if err != nil {
		log.Print("Error querying projects")
		return models.Project{}, err
	}

	return project, nil
}
