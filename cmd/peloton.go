package cmd

import (
	"fmt"
	"os"

	"github.com/ice1n36/brain/service"
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
			pelotonService, err := service.NewOnePelotonService(provider)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			distanceByYear, err := pelotonService.GetTotalDistanceTravelledInYear()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			totalDistance := 0.0
			for year, distance := range distanceByYear {
				fmt.Printf("Distanced Travelled in %d: %.2f miles\n", year, distance)
				totalDistance = totalDistance + distance
			}
			fmt.Printf("You've travelled a grand total of %.2f miles\n", totalDistance)
		},
	}
)

func init() {
	pelotonGetPerformanceSummaryCmd.Flags().StringVarP(&username, "username", "u", "", "Peloton username")
	pelotonGetPerformanceSummaryCmd.MarkFlagRequired("username")
	pelotonGetPerformanceSummaryCmd.Flags().StringVarP(&password, "password", "p", "", "Peloton password")
	pelotonGetPerformanceSummaryCmd.MarkFlagRequired("password")
}
