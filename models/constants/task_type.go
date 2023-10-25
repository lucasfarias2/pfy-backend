package constants

type TaskType string

const (
	PROJECT_CREATE                 TaskType = "Create project in Packlify Cloud"
	PROJECT_CREATE_GITHUB          TaskType = "Create repository in Github"
	PROJECT_CREATE_SDK_APP         TaskType = "Generate code with SDK in Packlify Cloud"
	PROJECT_PUSH_GITHUB            TaskType = "Push project to repository in Github"
	GCP_CREATE_ARTIFACT_REPOSITORY TaskType = "Create Artifact Repository in GCP"
	GCP_CREATE_BUILD_TRIGGER       TaskType = "Create Build Trigger in GCP"
	GCP_RUN_BUILD_TRIGGER          TaskType = "Run Build Trigger in GCP"
)
