// Copyright © 2017 Stream
//

package cmd

import (
	"fmt"
	"os"

	"github.com/GetStream/vg/internal/workspace"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show info about your current workspace",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := os.Getenv("VIRTUALGO")
		if name == "" {
			fmt.Println("No virtualgo workspace is active")
			return nil
		}

		ws := workspace.New(name)

		settings, err := ws.Settings()
		if err != nil {
			return err
		}

		fmt.Println("Active workspace:           ", name)
		fmt.Println("Workspace path:             ", ws.Path())
		fmt.Println("Fallback to global packages:", settings.GlobalFallback)
		return nil

	},
}

func init() {
	RootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
