package connectivity

import (
	"context"
	"fmt"
	_ "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"awsx/internal/awsclient"
	"awsx/internal/util"
	"strings"
)

// ConnectivityResult holds the result of the connectivity check
type ConnectivityResult struct {
	IsConnected bool
	Reason      string
}


func Check(client *awsclient.Client, srcType, srcID, dstType, dstID, vpcID string) (ConnectivityResult, error) {
	ctx := context.TODO()
	result := ConnectivityResult{IsConnected: true}

	// Resolve source and destination to subnet IDs
	srcSubnetID, srcSGs, err := resolveResource(client, ctx, srcType, srcID, vpcID)
	if err != nil {
		return ConnectivityResult{}, fmt.Errorf("source resource resolution failed: %w", err)
	}
	dstSubnetID, dstSGs, err := resolveResource(client, ctx, dstType, dstID, vpcID)
	if err != nil {
		return ConnectivityResult{}, fmt.Errorf("destination resource resolution failed: %w", err)
	}

	// If either resource is a subnet, skip SG check for that resource
	if srcSubnetID == "" || dstSubnetID == "" {
		return ConnectivityResult{IsConnected: false, Reason: fmt.Sprintf("resource %s (%s) or %s (%s) not found in VPC %s", srcType, srcID, dstType, dstID, vpcID)}, nil
	}

	// Get destination subnet details for CIDR
	dstSubnet, err := client.GetSubnetDetails(ctx, dstSubnetID)
	if err != nil {
		return ConnectivityResult{}, fmt.Errorf("destination subnet check failed: %w", err)
	}

	// Check route tables for connectivity between subnets
	if !checkRouteTable(client, ctx, srcSubnetID, *dstSubnet.CidrBlock) {
		return ConnectivityResult{IsConnected: false, Reason: fmt.Sprintf("no route exists from source subnet %s to destination subnet %s", srcSubnetID, dstSubnetID)}, nil
	}

	// Check security groups if both resources have associated SGs
	if len(srcSGs) > 0 && len(dstSGs) > 0 {
		if !checkSecurityGroups(client, ctx, srcSGs, dstSGs) {
			return ConnectivityResult{IsConnected: false, Reason: "security group rules block TCP port 80 traffic"}, nil
		}
	}

	// Check NACLs
	if !checkNACLs(client, ctx, srcSubnetID, dstSubnetID) {
		return ConnectivityResult{IsConnected: false, Reason: "network ACL rules block traffic"}, nil
	}

	return result, nil
}

func resolveResource(client *awsclient.Client, ctx context.Context, resourceType, resourceID, vpcID string) (string, []types.GroupIdentifier, error) {
	if resourceType == "subnet" {
		subnet, err := client.GetSubnetDetails(ctx, resourceID)
		if err != nil {
			return "", nil, err
		}
		if *subnet.VpcId != vpcID {
			return "", nil, fmt.Errorf("%s %s is not in VPC %s", resourceType, resourceID, vpcID)
		}
		return resourceID, nil, nil // Subnets don't have SGs
	}

	eniID, err := client.GetNetworkInterfaceByResource(ctx, resourceType, resourceID)
	if err != nil {
		return "", nil, err
	}
	eni, err := client.GetNetworkInterfaceDetails(ctx, eniID)
	if err != nil {
		return "", nil, err
	}
	if *eni.VpcId != vpcID {
		return "", nil, fmt.Errorf("%s %s (eni %s) is not in VPC %s", resourceType, resourceID, eniID, vpcID)
	}
	return *eni.SubnetId, eni.Groups, nil
}

func checkRouteTable(client *awsclient.Client, ctx context.Context, srcSubnetID, dstCidr string) bool {
	routeTables, err := client.GetRouteTables(ctx, srcSubnetID)
	if err != nil || len(routeTables) == 0 {
		return false
	}

	for _, rt := range routeTables {
		for _, route := range rt.Routes {
			if route.DestinationCidrBlock != nil && *route.DestinationCidrBlock == dstCidr && route.State == types.RouteStateActive {
				return true
			}
		}
	}
	return false
}

func checkSecurityGroups(client *awsclient.Client, ctx context.Context, srcSGs, dstSGs []types.GroupIdentifier) bool {
	dstSGIds := make([]string, 0, len(dstSGs))
	for _, sg := range dstSGs {
		if sg.GroupId != nil {
			dstSGIds = append(dstSGIds, *sg.GroupId)
		}
	}
	if len(dstSGIds) == 0 {
		return false
	}

	sgs, err := client.GetSecurityGroups(ctx, dstSGIds)
	if err != nil {
		return false
	}

	for _, sg := range sgs {
		for _, rule := range sg.IpPermissions {
			if rule.FromPort == 80 && rule.ToPort == 80 && util.ContainsProtocol(rule.IpProtocol, "tcp") {
				for _, ipRange := range rule.IpRanges {
					if ipRange.CidrIp != nil && strings.Contains(*ipRange.CidrIp, "0.0.0.0/0") {
						return true // Allow traffic from any source
					}
				}
				for _, sgRef := range rule.UserIdGroupPairs {
					for _, srcSG := range srcSGs {
						if sgRef.GroupId != nil && srcSG.GroupId != nil && *sgRef.GroupId == *srcSG.GroupId {
							return true // Allow traffic from source SG
						}
					}
				}
			}
		}
	}
	return false
}

func checkNACLs(client *awsclient.Client, ctx context.Context, srcSubnetID, dstSubnetID string) bool {
	// Get NACLs for destination subnet
	dstSubnet, err := client.GetSubnetDetails(ctx, dstSubnetID)
	if err != nil {
		return false
	}
	nacls, err := client.GetNetworkACLs(ctx, *dstSubnet.NetworkAclId)
	if err != nil {
		return false
	}

	for _, nacl := range nacls {
		for _, entry := range nacl.Entries {
			if entry.Egress != nil && !*entry.Egress && entry.RuleAction == types.RuleActionAllow {
				if entry.Protocol == nil || *entry.Protocol == "-1" || *entry.Protocol == "6" { // All protocols or TCP
					if entry.CidrBlock != nil && strings.Contains(*entry.CidrBlock, "0.0.0.0/0") {
						return true // Allow inbound traffic
					}
				}
			}
		}
	}
	return false
}