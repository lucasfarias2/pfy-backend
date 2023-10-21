package services

import (
	"database/sql"
	"packlify-cloud-backend/models"
	"packlify-cloud-backend/models/constants"
	"packlify-cloud-backend/utils"
	"time"
)

type TaskManager struct {
	db *sql.DB
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		db: utils.GetDB(),
	}
}

func (tm *TaskManager) CreateTask(projectID int, status constants.TaskStatus, message string, taskTypeID int) (models.Task, error) {
	task := models.Task{ProjectID: projectID, Status: status, Message: message, TaskTypeID: taskTypeID}

	err := tm.db.QueryRow("INSERT INTO tasks(project_id, status, message, task_type_id) VALUES($1, $2, $3, $4) RETURNING id, created_at, updated_at, message, task_type_id",
		task.ProjectID, task.Status, task.Message, task.TaskTypeID).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt, &task.Message, &task.TaskTypeID)
	return task, err
}

//err = db.QueryRow("INSERT INTO projects(name, organization_id, toolkit_id) VALUES($1, $2, $3) RETURNING id", project.Name, project.OrganizationID, project.ToolkitID).Scan(&project.ID)

func (tm *TaskManager) UpdateTaskStatus(taskID int, status constants.TaskStatus, message string) error {
	_, err := tm.db.Exec("UPDATE tasks SET status = $1, message = $2, updated_at = $3 WHERE id = $4",
		status, message, time.Now(), taskID)
	return err
}
