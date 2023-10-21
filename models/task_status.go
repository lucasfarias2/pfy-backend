package models

type TaskStatus int

const (
	Pending TaskStatus = iota
	Running
	Success
	Failed
)
