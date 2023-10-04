package models

type Project struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	OrganizationID *int    `json:"organization_id,omitempty"`
	ToolkitID      *int    `json:"toolkit_id,omitempty"`
}
