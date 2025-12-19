package cmd

import (
	"warehouse-go/warehouse-service/app"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the user service",
	Run: func(cmd *cobra.Command, args []string) {
		app.RunServer()
	},
}

func init()	{
	rootCmd.AddCommand(startCmd)
}
