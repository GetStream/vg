// Copyright © 2017 Stream
//

package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// deactivateCmd represents the deactivate command
var deactivateCmd = &cobra.Command{
	Use:   "deactivate",
	Short: "Deactivate the current virtualgo workspace",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("You haven't eval-ed `vg eval` yet.")
	},
}

func init() {
	RootCmd.AddCommand(deactivateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deactivateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deactivateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
