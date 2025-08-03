package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"github.com/spf13/cobra"
	"awsx/internal/awsclient"
	"awsx/internal/connectivity"

)

func main() {
	rootCmd := &cobra.Command{
		Use:   "awsx",
		Short: "CLI tool to check connectivity between AWS resources in the same VPC",
	}
	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "Check connectivity between source and destination resources",
		RunE:  runCheck,
	}
	rootCmd.AddCommand(checkCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runCheck(cmd *cobra.Command, args []string) error {
	client, err := awsclient.New()
	if err != nil {
		return fmt.Errorf("failed to initialize AWS client: %w", err)
	}
    
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter Source Resource Type (e.g., ec2, lambda, rds, eni, subnet): ")
	scanner.Scan()
	srcType := strings.ToLower(strings.TrimSpace(scanner.Text()))
	if !isValidResourceType(srcType) {
		return fmt.Errorf("invalid source resource type: must be one of ec2, lambda, rds, eni, subnet")
	}

	fmt.Print("Enter Source Resource ID (e.g., i-1234567890abcdef0, eni-12345678, subnet-12345678): ")
	scanner.Scan()
	srcID := strings.TrimSpace(scanner.Text())
	if !isValidID(srcType, srcID) {
		return fmt.Errorf("invalid source ID for %s: must start with correct prefix and be valid", srcType)
	}

	fmt.Print("Enter Destination Resource Type (e.g., ec2, lambda, rds, eni, subnet): ")
	scanner.Scan()
	dstType := strings.ToLower(strings.TrimSpace(scanner.Text()))
	if !isValidResourceType(dstType) {
		return fmt.Errorf("invalid destination resource type: must be one of ec2, lambda, rds, eni, subnet")
	}

	fmt.Print("Enter Destination Resource ID (e.g., i-0987654321fedcba0, eni-12345678, subnet-12345678): ")
	scanner.Scan()
	dstID := strings.TrimSpace(scanner.Text())
	if !isValidID(dstType, dstID) {
		return fmt.Errorf("invalid destination ID for %s: must start with correct prefix and be valid", dstType)
	}

	fmt.Print("Enter VPC ID (e.g., vpc-12345678): ")
	scanner.Scan()
	vpcID := strings.TrimSpace(scanner.Text())
	if !strings.HasPrefix(vpcID, "vpc-") || len(vpcID) < 8 {
		return fmt.Errorf("invalid VPC ID: must start with 'vpc-' and be a valid ID")
	}

	// Check connectivity
	result, err := connectivity.Check(client, srcType, srcID, dstType, dstID, vpcID)
	if err != nil {
		return fmt.Errorf("connectivity check failed: %w", err)
	}

	// Output result
	if result.IsConnected {
		fmt.Printf("Connectivity: SUCCESS\nFrom %s (%s) to %s (%s) in VPC %s\n", srcType, srcID, dstType, dstID, vpcID)
	} else {
		fmt.Printf("Connectivity: FAILED\nFrom %s (%s) to %s (%s) in VPC %s\nReason: %s\n", srcType, srcID, dstType, dstID, vpcID, result.Reason)
	}
	return nil
}

func isValidResourceType(resourceType string) bool {
	validTypes := map[string]bool{"ec2": true, "lambda": true, "rds": true, "eni": true, "subnet": true}
	return validTypes[resourceType]
}

func isValidID(resourceType, id string) bool {
	switch resourceType {
	case "ec2":
		return strings.HasPrefix(id, "i-") && len(id) >= 10
	case "lambda", "rds":
		return len(id) > 0 // Simplified; can add specific validation if needed
	case "eni":
		return strings.HasPrefix(id, "eni-") && len(id) >= 8
	case "subnet":
		return strings.HasPrefix(id, "subnet-") && len(id) >= 8
	}
	return false
}