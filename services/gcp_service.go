package services

import (
	"context"
	"fmt"
	"os"
	"packlify-cloud-backend/models"

	artifactregistry "cloud.google.com/go/artifactregistry/apiv1"
	artifactregistrypb "cloud.google.com/go/artifactregistry/apiv1/artifactregistrypb"
	cloudbuild "cloud.google.com/go/cloudbuild/apiv1/v2"
	"cloud.google.com/go/cloudbuild/apiv1/v2/cloudbuildpb"
	run "cloud.google.com/go/run/apiv2"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

func CreateArtifactRepository(project models.Project) error {
	ctx := context.Background()

	// Construct credentials JSON from environment variables
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

	// Parse credentials JSON
	creds, err := google.CredentialsFromJSON(ctx, []byte(credsJSON), artifactregistry.DefaultAuthScopes()...)
	if err != nil {
		return err
	}

	// Create an instance of Artifact Registry
	arClient, err := artifactregistry.NewClient(
		ctx,
		option.WithCredentials(creds),
	)

	if err != nil {
		return err
	}

	// Create a new repository
	repo, err := arClient.CreateRepository(ctx, &artifactregistrypb.CreateRepositoryRequest{
		Parent: fmt.Sprintf("projects/%s/locations/us-central1", os.Getenv("GCP_PROJECT_ID")),
		Repository: &artifactregistrypb.Repository{
			Format:      artifactregistrypb.Repository_DOCKER,
			Name:        project.Name + "-docker",
			Description: "Test repo created by Packlify Cloud",
		},
		RepositoryId: project.Name + "-docker",
	})

	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(repo)

	defer arClient.Close()

	return nil
}

func CreateCloudRun(project models.Project) error {
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

	// Parse credentials JSON
	creds, err := google.CredentialsFromJSON(ctx, []byte(credsJSON), run.DefaultAuthScopes()...)
	if err != nil {
		return err
	}

	// Create an instance of Cloud Build
	cbClient, err := cloudbuild.NewClient(
		ctx,
		option.WithCredentials(creds),
	)

	if err != nil {
		return err
	}

	// Create an instance of Cloud Run
	runClient, err := run.NewServicesClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return err
	}

	triggerOp, err := cbClient.CreateBuildTrigger(ctx, &cloudbuildpb.CreateBuildTriggerRequest{
		Parent:    fmt.Sprintf("projects/%s/locations/global", os.Getenv("GCP_PROJECT_ID")),
		ProjectId: os.Getenv("GCP_PROJECT_ID"),
		Trigger: &cloudbuildpb.BuildTrigger{
			Name:        "test-trigger" + project.Name,
			Description: "Trigger for test-repo generated by Packlify Cloud",
			Github: &cloudbuildpb.GitHubEventsConfig{
				Owner: os.Getenv("GITHUB_OWNER"),
				Name:  "test-repo-7",
				Event: &cloudbuildpb.GitHubEventsConfig_Push{
					Push: &cloudbuildpb.PushFilter{
						GitRef: &cloudbuildpb.PushFilter_Branch{
							Branch: "main",
						},
					},
				},
			},
			BuildTemplate: &cloudbuildpb.BuildTrigger_Build{
				Build: &cloudbuildpb.Build{
					Id: "test-build" + project.Name,
					Steps: []*cloudbuildpb.BuildStep{
						{
							Id:   "Build",
							Name: "gcr.io/cloud-builders/docker",
							Args: []string{
								"build",
								"-t",
								"us-central1-docker.pkg.dev/shopinpack-com/test-repo-7-docker/test-repo-7:$COMMIT_SHA",
								".",
								"-f",
								"Dockerfile",
							},
						},
						{
							Id:   "Push",
							Name: "gcr.io/cloud-builders/docker",
							WaitFor: []string{
								"Build",
							},
							Args: []string{
								"push",
								"us-central1-docker.pkg.dev/shopinpack-com/test-repo-7-docker/test-repo-7:$COMMIT_SHA",
							},
						},
						{
							Id:   "Deploy",
							Name: "gcr.io/cloud-builders/gcloud",
							WaitFor: []string{
								"Push",
							},
							Args: []string{
								"run",
								"deploy",
								"new-test-repo" + project.Name,
								"--image",
								"us-central1-docker.pkg.dev/shopinpack-com/test-repo-7-docker/test-repo-7:$COMMIT_SHA",
								"--region",
								"us-central1",
								"--platform",
								"managed",
								"--allow-unauthenticated",
							},
						},
					}},
			},
		},
	})

	if err != nil {
		fmt.Println(err)

		return err
	}

	fmt.Println(triggerOp)

	// op, err := runClient.CreateService(ctx, &runpb.CreateServiceRequest{
	// 	Parent: fmt.Sprintf("projects/%s/locations/us-central1", os.Getenv("GCP_PROJECT_ID")),
	// 	Service: &runpb.Service{
	// 		Name: "test-repo-service",
	// 	},
	// 	ServiceId: "test-repo-service",
	// })
	// if err != nil {
	// 	// TODO: Handle error.

	// 	fmt.Printf("Error: %v\n", err)
	// 	return err
	// }

	// resp, err := op.Wait(ctx)
	// if err != nil {
	// 	// TODO: Handle error.
	// 	fmt.Printf("Error: %v\n", err)
	// 	return err
	// }

	// fmt.Printf("Operation result: %v\n", resp)

	defer cbClient.Close()
	defer runClient.Close()

	return nil
}
