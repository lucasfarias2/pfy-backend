package services

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
)

func CreateElasticBeanstalkEnvironment(config map[string]string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})

	if err != nil {
		return err
	}

	optionSettings := []*elasticbeanstalk.ConfigurationOptionSetting{
		{
			Namespace:  aws.String("aws:autoscaling:launchconfiguration"),
			OptionName: aws.String("IamInstanceProfile"),
			Value:      aws.String("aws-elasticbeanstalk-ec2-role"),
		},
	}

	svc := elasticbeanstalk.New(sess)
	input := &elasticbeanstalk.CreateEnvironmentInput{
		ApplicationName:   aws.String(config["applicationName"]),
		EnvironmentName:   aws.String(config["environmentName"]),
		SolutionStackName: aws.String(config["solutionStackName"]),
		CNAMEPrefix:       aws.String(config["cnamePrefix"]),
		OptionSettings:    optionSettings,
	}

	_, err = svc.CreateEnvironment(input)
	if err != nil {
		return err
	}
	return nil
}

func CreateElasticBeanstalkApplication(config map[string]string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})

	if err != nil {
		return err
	}

	svc := elasticbeanstalk.New(sess)
	input := &elasticbeanstalk.CreateApplicationInput{
		ApplicationName: aws.String(config["applicationName"]),
		Description:     aws.String("Packlify generated application"),
	}

	_, err = svc.CreateApplication(input)
	if err != nil {
		return err
	}

	return nil
}
