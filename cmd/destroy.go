// Copyright © 2017 Stream
//

package cmd

import (
	"os"

	"github.com/GetStream/vg/internal/workspace"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy <workspace> [extraWorkspaces...]",
	Short: "Removes one or multiple workspace and all their contents",
	Long: `To remove workspace 'myWorkspace' and 'someOtherWorkspace':
	
	vg destroy myWorkspace someOtherWorkspace`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("No workspace specified")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, wsName := range args {
			err := os.RemoveAll(workspace.New(wsName).Path())
			if err != nil {
				return errors.Wrapf(err, "Couldn't remove workspace %q", wsName)
			}
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(destroyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// destroyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// destroyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
