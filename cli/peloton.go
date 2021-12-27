package cli

import (
	"fmt"
	"github.com/ice1n36/brain/clients"
	"github.com/spf13/cobra"
	"os"
)

var (
	pelotonGetPerformanceSummaryCmd = &cobra.Command{
		Use:   "peloton-get-perf-summary",
		Short: "brain peloton get performance summary",
		Long:  `peloton command to test getting performance summary`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("BOOYA")
			pelotonClient, err := NewOnePelotonAPIClient(nil)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
)

func init() {
}
