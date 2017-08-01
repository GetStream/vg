// Copyright © 2017 Stream
//

package cmd

import (
	"github.com/GetStream/vg/internal/workspace"
	"github.com/spf13/cobra"
)

// activateRealCmd represents the activateReal command
var activateRealCmd = &cobra.Command{
	Use:    "activateReal",
	Hidden: true,
	Short:  "This is used to assure the local installs are actually installed",
	Long:   ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		ws, err := workspace.Current()
		if err != nil {
			return err
		}

		return ws.InstallSavedLocalPackages()

	},
}

func init() {
	RootCmd.AddCommand(activateRealCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// activateRealCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// activateRealCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
