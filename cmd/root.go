package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)


var rootCmd = &cobra.Command{
	Use: "awsx",
	Short: "awsx - Enhanced AWS CLI",
	Long: "awsx is a CLI tool that makes it easy to use aws cli",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("region", "", "AWS region to use (overrides default)")
	rootCmd.PersistentFlags().String("profile", "", "AWS profile to use (overrides default)")
}