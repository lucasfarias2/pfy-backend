package gcp

import (
	"context"
	"fmt"
	"os"
	"packlify-cloud-backend/models"

	cloudbuild "cloud.google.com/go/cloudbuild/apiv1"
	"cloud.google.com/go/cloudbuild/apiv1/v2/cloudbuildpb"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

func CreateBuildTrigger(project models.Project) (*cloudbuildpb.BuildTrigger, error) {
	ctx := context.Background()

	gcpProjectId := os.Getenv("GCP_PROJECT_ID")
	gcpRegion := os.Getenv("GCP_REGION")
	githubOwner := os.Getenv("GITHUB_OWNER")
	githubRepo := "test-repo-7"
	triggerName := project.Name + "-trigger"
	imageName := fmt.Sprintf("%s-docker.pkg.dev/%s/%s/%s:$COMMIT_SHA", gcpRegion, gcpProjectId, project.Name+"-docker", project.Name)

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
		gcpProjectId,
		os.Getenv("GCP_PRIVATE_KEY_ID"),
		os.Getenv("GCP_PRIVATE_KEY"),
		os.Getenv("GCP_CLIENT_EMAIL"),
		os.Getenv("GCP_CLIENT_ID"),
		os.Getenv("GCP_CLIENT_EMAIL"),
	)

	creds, err := google.CredentialsFromJSON(ctx, []byte(credsJSON), cloudbuild.DefaultAuthScopes()...)
	if err != nil {
		return nil, err
	}

	cbClient, err := cloudbuild.NewClient(
		ctx,
		option.WithCredentials(creds),
	)

	if err != nil {
		return nil, err
	}

	triggerOp, err := cbClient.CreateBuildTrigger(ctx, &cloudbuildpb.CreateBuildTriggerRequest{
		Parent:    fmt.Sprintf("projects/%s/locations/global", gcpProjectId),
		ProjectId: gcpProjectId,
		Trigger: &cloudbuildpb.BuildTrigger{
			Name:        triggerName,
			Description: "Trigger for test-repo generated by Packlify Cloud",
			Github: &cloudbuildpb.GitHubEventsConfig{
				Owner: githubOwner,
				Name:  githubRepo,
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
								imageName,
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
								imageName,
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
								project.Name + "-run",
								"--image",
								imageName,
								"--region",
								gcpRegion,
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
	}

	fmt.Println(triggerOp)

	defer cbClient.Close()

	return triggerOp, err
}

func RunBuildTrigger(buildTrigger *cloudbuildpb.BuildTrigger) error {
	ctx := context.Background()

	gcpProjectId := os.Getenv("GCP_PROJECT_ID")

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
		gcpProjectId,
		os.Getenv("GCP_PRIVATE_KEY_ID"),
		os.Getenv("GCP_PRIVATE_KEY"),
		os.Getenv("GCP_CLIENT_EMAIL"),
		os.Getenv("GCP_CLIENT_ID"),
		os.Getenv("GCP_CLIENT_EMAIL"),
	)

	creds, err := google.CredentialsFromJSON(ctx, []byte(credsJSON), cloudbuild.DefaultAuthScopes()...)
	if err != nil {
		return err
	}

	cbClient, err := cloudbuild.NewClient(
		ctx,
		option.WithCredentials(creds),
	)

	if err != nil {
		return err
	}

	fmt.Println(buildTrigger)

	triggerOp, err := cbClient.RunBuildTrigger(ctx, &cloudbuildpb.RunBuildTriggerRequest{
		ProjectId: gcpProjectId,
		TriggerId: buildTrigger.Id,
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(triggerOp)

	defer cbClient.Close()

	return err
}
