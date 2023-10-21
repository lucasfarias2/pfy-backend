package constants

type TaskStatus string

const (
	Pending TaskStatus = "Pending"
	Running TaskStatus = "Running"
	Success TaskStatus = "Success"
	Failed  TaskStatus = "Failed"
)
