// Copyright © 2017 Stream
//

package cmd

import (
	"fmt"

	"github.com/GetStream/vg/internal/workspace"
	"github.com/spf13/cobra"
)

// querySettingsCmd represents the querySettings command
var querySettingsCmd = &cobra.Command{
	Use:    "querySettings",
	Hidden: true,
	Short:  "This can be used to query certain settings for a workspace",
	Long: `WARNING: This command is unstable it is only exposed because it is
used internally by virtualgo`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ws, err := workspace.Current()
		if err != nil {
			return err
		}

		settings, err := ws.Settings()
		if err != nil {
			return err
		}

		fmt.Println(settings.GlobalFallback)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(querySettingsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// querySettingsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// querySettingsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
