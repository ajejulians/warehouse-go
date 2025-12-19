package cmd

import (

	"github.com/spf13/cobra"
	"warehouse-go/user-service/app"	
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
