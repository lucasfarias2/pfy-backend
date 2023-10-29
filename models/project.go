package models

import "packlify-cloud-backend/models/constants"

type Project struct {
	ID             int               `json:"id"`
	Name           string            `json:"name"`
	GitHubRepo     string            `json:"github_repo"`
	OrganizationID *int              `json:"organization_id,omitempty"`
	Toolkit        constants.Toolkit `json:"toolkit"`
}
