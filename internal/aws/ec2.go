package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws" // Add this import for aws.Config
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	_ "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// EC2Dependencies holds the dependencies of an EC2 instance
type EC2Dependencies struct {
	InstanceID     string
	InstanceState  string
	VPCID          string
	SubnetID       string
	SecurityGroups []string
	IAMRole        string
}

// GetEC2Dependencies fetches the dependencies of an EC2 instance
func GetEC2Dependencies(ctx context.Context, cfg aws.Config, instanceID string) (*EC2Dependencies, error) {
	client := ec2.NewFromConfig(cfg)

	// Describe the instance
	resp, err := client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Reservations) == 0 || len(resp.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("instance not found: %s", instanceID)
	}

	instance := resp.Reservations[0].Instances[0]

	// Extract dependencies
	deps := &EC2Dependencies{
		InstanceID:     *instance.InstanceId,
		InstanceState:  string(instance.State.Name),
		VPCID:          *instance.VpcId,
		SubnetID:       *instance.SubnetId,
		SecurityGroups: make([]string, len(instance.SecurityGroups)),
	}

	// Populate security groups
	for i, sg := range instance.SecurityGroups {
		deps.SecurityGroups[i] = *sg.GroupId
	}

	// Fetch IAM role (if any)
	if instance.IamInstanceProfile != nil && instance.IamInstanceProfile.Arn != nil {
		deps.IAMRole = *instance.IamInstanceProfile.Arn
	} else {
		deps.IAMRole = "None"
	}

	return deps, nil
}