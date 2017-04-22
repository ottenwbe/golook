package cmd

import (
	"fmt"
	"github.com/ottenwbe/golook/broker/api"
	"github.com/spf13/cobra"
)

func init() {
	api.ConfigApi()
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: fmt.Sprint("Print information about the API."),
	Long:  fmt.Sprint("Print the support versions of the API."),
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println(fmt.Sprintf("API Versions: %s", api.API_VERSION))
	},
}

func init() {
	api.ConfigApi()

	RootCmd.AddCommand(apiCmd)
}
