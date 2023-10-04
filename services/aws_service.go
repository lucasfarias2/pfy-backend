// awsService.go

package services

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func CreateAWSInstance(config map[string]string) (*ec2.Reservation, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})

	if err != nil {
		return nil, err
	}

	svc := ec2.New(sess)

	input := &ec2.RunInstancesInput{
		ImageId:      aws.String(config["ami"]),
		InstanceType: aws.String(config["instanceType"]),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(config["name"]),
					},
				},
			},
		},
	}

	result, err := svc.RunInstances(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}
