package models

type Project struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	GitHubRepo     string `json:"github_repo_name"`
	OrganizationID *int   `json:"organization_id,omitempty"`
	ToolkitID      *int   `json:"toolkit_id,omitempty"`
}
