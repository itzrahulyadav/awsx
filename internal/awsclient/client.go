package awsclient

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// Client wraps the AWS EC2 client
type Client struct {
	ec2Client *ec2.Client
}

// New initializes a new AWS client
func New() (*Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}
	return &Client{ec2Client: ec2.NewFromConfig(cfg)}, nil
}

// GetNetworkInterfaceDetails fetches details for a network interface
func (c *Client) GetNetworkInterfaceDetails(ctx context.Context, eniID string) (*types.NetworkInterface, error) {
	resp, err := c.ec2Client.DescribeNetworkInterfaces(ctx, &ec2.DescribeNetworkInterfacesInput{
		NetworkInterfaceIds: []string{eniID},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe network interface %s: %w", eniID, err)
	}
	if len(resp.NetworkInterfaces) == 0 {
		return nil, fmt.Errorf("network interface %s not found", eniID)
	}
	return &resp.NetworkInterfaces[0], nil
}

// GetSubnetDetails fetches details for a subnet
func (c *Client) GetSubnetDetails(ctx context.Context, subnetID string) (*types.Subnet, error) {
	resp, err := c.ec2Client.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{
		SubnetIds: []string{subnetID},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe subnet %s: %w", subnetID, err)
	}
	if len(resp.Subnets) == 0 {
		return nil, fmt.Errorf("subnet %s not found", subnetID)
	}
	return &resp.Subnets[0], nil
}

// GetRouteTables fetches route tables for a subnet
func (c *Client) GetRouteTables(ctx context.Context, subnetID string) ([]types.RouteTable, error) {
	resp, err := c.ec2Client.DescribeRouteTables(ctx, &ec2.DescribeRouteTablesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("association.subnet-id"),
				Values: []string{subnetID},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe route tables for subnet %s: %w", subnetID, err)
	}
	return resp.RouteTables, nil
}

// GetSecurityGroups fetches security group details
func (c *Client) GetSecurityGroups(ctx context.Context, groupIDs []string) ([]types.SecurityGroup, error) {
	resp, err := c.ec2Client.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{
		GroupIds: groupIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe security groups: %w", err)
	}
	return resp.SecurityGroups, nil
}

// GetNetworkACLs fetches NACL details
func (c *Client) GetNetworkACLs(ctx context.Context, naclID string) ([]types.NetworkAcl, error) {
	resp, err := c.ec2Client.DescribeNetworkAcls(ctx, &ec2.DescribeNetworkAclsInput{
		NetworkAclIds: []string{naclID},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe NACL %s: %w", naclID, err)
	}
	return resp.NetworkAcls, nil
}

// GetNetworkInterfaceByResource resolves a resource to its network interface
func (c *Client) GetNetworkInterfaceByResource(ctx context.Context, resourceType, resourceID string) (string, error) {
	switch resourceType {
	case "ec2":
		resp, err := c.ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
			InstanceIds: []string{resourceID},
		})
		if err != nil {
			return "", fmt.Errorf("failed to describe EC2 instance %s: %w", resourceID, err)
		}
		if len(resp.Reservations) == 0 || len(resp.Reservations[0].Instances) == 0 {
			return "", fmt.Errorf("EC2 instance %s not found", resourceID)
		}
		if len(resp.Reservations[0].Instances[0].NetworkInterfaces) == 0 {
			return "", fmt.Errorf("EC2 instance %s has no network interfaces", resourceID)
		}
		return *resp.Reservations[0].Instances[0].NetworkInterfaces[0].NetworkInterfaceId, nil
	case "eni":
		return resourceID, nil // Already an ENI
	case "subnet":
		return "", nil // Subnet-based check, no ENI needed
	case "lambda", "rds":
		// Simplified: assume resource has an ENI (real implementation would query Lambda/RDS APIs)
		resp, err := c.ec2Client.DescribeNetworkInterfaces(ctx, &ec2.DescribeNetworkInterfacesInput{
			Filters: []types.Filter{
				{
					Name:   aws.String("description"),
					Values: []string{fmt.Sprintf("*%s*", resourceID)},
				},
			},
		})
		if err != nil {
			return "", fmt.Errorf("failed to find network interface for %s %s: %w", resourceType, resourceID, err)
		}
		if len(resp.NetworkInterfaces) == 0 {
			return "", fmt.Errorf("no network interface found for %s %s", resourceType, resourceID)
		}
		return *resp.NetworkInterfaces[0].NetworkInterfaceId, nil
	default:
		return "", fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}