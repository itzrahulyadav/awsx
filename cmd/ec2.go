package cmd

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/cobra"
	"github.com/itzrahulyadav/awsx/internal/aws"
	"github.com/itzrahulyadav/awsx/internal/output"
)

// ec2Cmd represents the ec2 subcommand under visualize
var ec2Cmd = &cobra.Command{
	Use:   "ec2 [instance-id]",
	Short: "Visualize dependencies of an EC2 instance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		instanceID := args[0]
		format, _ := cmd.Flags().GetString("format")

		// Load AWS configuration
		cfg, err := config.LoadDefaultConfig(context.Background(),
			config.WithRegion(cmd.Flag("region").Value.String()),
			config.WithSharedConfigProfile(cmd.Flag("profile").Value.String()),
		)
		if err != nil {
			return fmt.Errorf("failed to load AWS config: %v", err)
		}

		// Fetch EC2 instance dependencies
		deps, err := aws.GetEC2Dependencies(context.Background(), cfg, instanceID)
		if err != nil {
			return fmt.Errorf("failed to fetch EC2 dependencies: %v", err)
		}

		// Format and output the result
		switch format {
		case "text":
			output.PrintText(deps)
		case "json":
			output.PrintJSON(deps)
		case "dot":
			output.PrintDOT(deps)
		default:
			return fmt.Errorf("unsupported format: %s", format)
		}

		return nil
	},
}

func init() {
	visualizeCmd.AddCommand(ec2Cmd)
}