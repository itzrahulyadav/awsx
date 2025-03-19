package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// visualizeCmd represents the visualize command
var visualizeCmd = &cobra.Command{
	Use:   "visualize [service] [resource-id]",
	Short: "Visualize AWS resource dependencies",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("service not supported: %s (use 'ec2' for now)", args[0])
	},
}

func init() {
	rootCmd.AddCommand(visualizeCmd)
	visualizeCmd.PersistentFlags().String("format", "text", "Output format (text|json|dot)")
}