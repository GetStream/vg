// Copyright © 2017 Stream
//

package cmd

import (
	"os"

	"github.com/GetStream/vg/internal/utils"
	"github.com/GetStream/vg/internal/workspace"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// moveVendorCmd represents the moveVendor command
var moveVendorCmd = &cobra.Command{
	Use:   "moveVendor",
	Short: "Moves the vendor directory to the workspace",
	Long: `This command can be useful when using virtualgo with projects that
don't use dep yet. For instance to install the dependencies of a glide
based project in your workspace do this:

	glide install
	vg moveVendor
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ws, err := workspace.Current()
		if err != nil {
			return err
		}

		exist, err := utils.VendorExists()
		if err != nil {
			return err
		}
		if !exist {
			return errors.New("No vendor directory is present")
		}

		err = ws.ClearSrc()
		if err != nil {
			return errors.Wrap(err, "Couldn't clear the current src directory inside the workspace")
		}

		err = os.Rename("vendor", ws.Src())
		if err != nil {
			return errors.Wrap(err, "Couldn't move the vendor directory to the active workspace")
		}

		return ws.InstallSavedLocalPackages()

	},
}

func init() {
	RootCmd.AddCommand(moveVendorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// moveVendorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// moveVendorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
