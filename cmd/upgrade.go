// Copyright © 2017 Stream
//

package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the virtualgo binary and reload it in the current shell",
	Long: `Upgrade the virtualgo binary and reload it in the current shell

NOTE: This does not always work. 
NOTE: If you encounter issues please restart your terminal.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New(noEvalError)
	},
}

func init() {
	RootCmd.AddCommand(upgradeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upgradeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upgradeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
