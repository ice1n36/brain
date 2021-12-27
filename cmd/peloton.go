package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/ice1n36/brain/clients"
	"github.com/spf13/cobra"
	"go.uber.org/config"
)

var (
	username = ""
	password = ""

	pelotonGetPerformanceSummaryCmd = &cobra.Command{
		Use:   "peloton-get-perf-summary",
		Short: "brain peloton get performance summary",
		Long:  `peloton command to test getting performance summary`,
		Run: func(cmd *cobra.Command, args []string) {
			provider, err := config.NewYAML(config.Static(map[string]map[string]string{
				"onepeloton": {
					"username_or_email": username,
					"password":          password,
				},
			}))
			pelotonClient, err := clients.NewOnePelotonAPIClient(provider)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			session, err := pelotonClient.Login(context.TODO())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(session)
		},
	}
)

func init() {
	pelotonGetPerformanceSummaryCmd.Flags().StringVarP(&username, "username", "u", "", "Peloton username")
	pelotonGetPerformanceSummaryCmd.MarkFlagRequired("username")
	pelotonGetPerformanceSummaryCmd.Flags().StringVarP(&password, "password", "p", "", "Peloton password")
	//pelotonGetPerformanceSummaryCmd.MarkFlagRequired("password")
}
