package gcp

import (
	"context"
	"fmt"
	"os"
	"packlify-cloud-backend/models"

	artifactregistry "cloud.google.com/go/artifactregistry/apiv1"
	"cloud.google.com/go/artifactregistry/apiv1/artifactregistrypb"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

func CreateArtifactRepository(project models.Project) error {
	ctx := context.Background()

	credsJSON := fmt.Sprintf(
		`{
			"type": "service_account",
			"project_id": "%s",
			"private_key_id": "%s",
			"private_key": "%s",
			"client_email": "%s",
			"client_id": "%s",
			"auth_uri": "https://accounts.google.com/o/oauth2/auth",
			"token_uri": "https://oauth2.googleapis.com/token",
			"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
			"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/%s"
		}`,
		os.Getenv("GCP_PROJECT_ID"),
		os.Getenv("GCP_PRIVATE_KEY_ID"),
		os.Getenv("GCP_PRIVATE_KEY"),
		os.Getenv("GCP_CLIENT_EMAIL"),
		os.Getenv("GCP_CLIENT_ID"),
		os.Getenv("GCP_CLIENT_EMAIL"),
	)

	creds, err := google.CredentialsFromJSON(ctx, []byte(credsJSON), artifactregistry.DefaultAuthScopes()...)
	if err != nil {
		return err
	}

	arClient, err := artifactregistry.NewClient(
		ctx,
		option.WithCredentials(creds),
	)

	if err != nil {
		return err
	}

	repositoryName := project.Name + "-docker"

	repo, err := arClient.CreateRepository(ctx, &artifactregistrypb.CreateRepositoryRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", os.Getenv("GCP_PROJECT_ID"), os.Getenv("GCP_REGION")),
		Repository: &artifactregistrypb.Repository{
			Format:      artifactregistrypb.Repository_DOCKER,
			Name:        repositoryName,
			Description: "Test repo created by Packlify Cloud",
		},
		RepositoryId: repositoryName,
	})

	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(repo)

	defer arClient.Close()

	return nil
}
