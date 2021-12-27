package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "brain-cli",
		Short: "brain-cli is the command line interface for the brain service",
		Long:  `A command line tool used for testing the brain service on command line`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
)

func init() {
	rootCmd.AddCommand(pelotonGetPerformanceSummaryCmd)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
