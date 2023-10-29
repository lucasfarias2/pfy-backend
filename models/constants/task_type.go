package constants

type TaskType string

const (
	PROJECT_GENERATE_FILES         TaskType = "Generate files from toolkit"
	PROJECT_CREATE_GITHUB          TaskType = "Create repository in Github"
	PROJECT_PUSH_GITHUB            TaskType = "Push project to repository in Github"
	GCP_CONNECT_REPOSITORY         TaskType = "Connect Github Repository to GCP"
	GCP_CREATE_ARTIFACT_REPOSITORY TaskType = "Create Artifact Repository in GCP"
	GCP_CREATE_BUILD_TRIGGER       TaskType = "Create Build Trigger in GCP"
	GCP_RUN_BUILD_TRIGGER          TaskType = "Run Build Trigger in GCP"
)
