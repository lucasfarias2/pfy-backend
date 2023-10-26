package services

import (
	"database/sql"
	"log"
	"os"
	"packlify-cloud-backend/models"
	"packlify-cloud-backend/models/constants"
	"packlify-cloud-backend/utils"
)

func CreateProject(project models.Project) (models.Project, error) {
	db := utils.GetDB()
	tm := NewTaskManager()

	err := db.QueryRow("INSERT INTO projects(name, organization_id, toolkit_id) VALUES($1, $2, $3) RETURNING id", project.Name, project.OrganizationID, project.ToolkitID).Scan(&project.ID)
	if err != nil {
		return models.Project{}, err
	}

	task, err := tm.CreateTask(project.ID, constants.Running, "", string(constants.PROJECT_CREATE))
	if err != nil {
		return models.Project{}, err
	}

	var toolkitName string
	err = db.QueryRow("SELECT name FROM toolkits WHERE id=$1", project.ToolkitID).Scan(&toolkitName)
	if err != nil {
		return models.Project{}, err
	}

	if toolkitName == "React" {
		if err != nil {
			return models.Project{}, err
		}

		err := CreateSDKApp(project.Name)
		if err != nil {
			return models.Project{}, err
		}

		githubToken := os.Getenv("GITHUB_ACCESS_TOKEN")
		cloneURL, err := CreateGitHubRepo(project.Name, githubToken)
		if err != nil {
			log.Fatalf("Error creating Github repo: %s", err)
		}

		err = PushToGitHubRepo(project.Name, cloneURL)
		if err != nil {
			log.Fatalf("Error pushing Github repo: %s", err)
		}

		err = tm.UpdateTaskStatus(task.ID, constants.Success, "Project created successfully")
		if err != nil {
			return models.Project{}, err
		}
	}

	return project, nil
}

func GetAllProjects(organizationId string) ([]models.Project, error) {
	db := utils.GetDB()
	rows, err := db.Query("SELECT id, name FROM projects WHERE organization_id = $1", organizationId)
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

func GetProjectStatusById(id int) ([]models.Task, error) {
	db := utils.GetDB()
	query := `
        SELECT id, project_id, task_name, status, message, created_at, updated_at
        FROM tasks
        WHERE project_id = $1
        ORDER BY created_at DESC
    `

	// Execute the query
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.ProjectID, &task.TaskName, &task.Status, &task.Message, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func UpdateProjectRepoName(id int, newRepoName string) (models.Project, error) {
	db := utils.GetDB()
	var project models.Project

	_, err := db.Exec("UPDATE projects SET github_repo = $1 WHERE id = $2", newRepoName, id)
	if err != nil {
		return models.Project{}, err
	}

	err = db.QueryRow("SELECT id, name, github_repo, organization_id, toolkit_id FROM projects WHERE id = $1", id).Scan(&project.ID, &project.Name, &project.GitHubRepo, &project.OrganizationID, &project.ToolkitID)
	if err != nil {
		return models.Project{}, err
	}

	return project, nil
}
