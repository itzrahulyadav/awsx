package output

import (
	"encoding/json"
	"fmt"
	_ "strings"

	"github.com/itzrahulyadav/awsx/internal/aws"
)

// PrintText prints the dependencies in a human-readable text format
func PrintText(deps *aws.EC2Dependencies) {
	fmt.Printf("EC2 Instance: %s (%s)\n", deps.InstanceID, deps.InstanceState)
	fmt.Printf("  VPC: %s\n", deps.VPCID)
	fmt.Printf("  Subnet: %s\n", deps.SubnetID)
	fmt.Printf("  Security Groups:\n")
	for _, sg := range deps.SecurityGroups {
		fmt.Printf("    - %s\n", sg)
	}
	fmt.Printf("  IAM Role: %s\n", deps.IAMRole)
}

// PrintJSON prints the dependencies in JSON format
func PrintJSON(deps *aws.EC2Dependencies) error {
	data, err := json.MarshalIndent(deps, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

// PrintDOT prints the dependencies in Graphviz DOT format
func PrintDOT(deps *aws.EC2Dependencies) {
	fmt.Println("digraph ec2_dependencies {")
	fmt.Println(`  rankdir=BT;`)
	fmt.Printf("  instance [label=\"EC2 Instance\\n%s\\n(%s)\"];\n", deps.InstanceID, deps.InstanceState)
	fmt.Printf("  vpc [label=\"VPC\\n%s\"];\n", deps.VPCID)
	fmt.Printf("  subnet [label=\"Subnet\\n%s\"];\n", deps.SubnetID)

	// Relationships
	fmt.Println("  instance -> subnet;")
	fmt.Println("  subnet -> vpc;")

	// Security Groups
	for i, sg := range deps.SecurityGroups {
		fmt.Printf("  sg%d [label=\"Security Group\\n%s\"];\n", i, sg)
		fmt.Printf("  instance -> sg%d;\n", i)
	}

	// IAM Role
	if deps.IAMRole != "None" {
		fmt.Printf("  iam [label=\"IAM Role\\n%s\"];\n", deps.IAMRole)
		fmt.Println("  instance -> iam;")
	}

	fmt.Println("}")
}