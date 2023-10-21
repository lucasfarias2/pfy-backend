package services

import (
	"context"

	run "cloud.google.com/go/run/apiv2"

	runpb "cloud.google.com/go/run/apiv2/runpb"
)

func CreateCloudRun() {
	ctx := context.Background()

	c, err := run.NewServicesClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()

	req := &runpb.CreateServiceRequest{
		// TODO: Fill request struct fields.
		// See https://pkg.go.dev/cloud.google.com/go/run/apiv2/runpb#CreateServiceRequest.
	}
	op, err := c.CreateService(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}
