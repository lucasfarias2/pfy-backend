package tasks_models

import "cloud.google.com/go/cloudbuild/apiv1/v2/cloudbuildpb"

type BuildTriggerData struct {
	IsSuccess bool
	Trigger   *cloudbuildpb.BuildTrigger
}
