package models

type Project struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	OrganizationID int    `json:"organization_id"`
	ToolkitID      int    `json:"toolkit_id"`
}
