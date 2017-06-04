// Copyright © 2017 Stream
//

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall <package> [otherPackages]",
	Short: "Uninstall a package from the active workspace",
	Long: `To remove github.com/pkg/errors:
	
	vg uninstall github.com/pkg/errors`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("No package specified")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		workspace := os.Getenv("VIRTUALGO")
		if workspace == "" {
			return errors.New("A virtualgo workspace should be activated first by using `vg activate [workspaceName]`")
		}

		for _, pkg := range args {
			fmt.Printf("Uninstalling %s from workspace\n", pkg)
			err := os.RemoveAll(filepath.Join(virtualgoDir, workspace, "src", pkg))
			if err != nil {
				return errors.Wrapf(err, "Couldn't remove package src '%s'", workspace)
			}

			err = os.RemoveAll(filepath.Join(virtualgoDir, workspace, "pkg", pkg))
			if err != nil {
				return errors.Wrapf(err, "Couldn't remove installed package files for '%s'", workspace)
			}

			err = os.RemoveAll(filepath.Join(virtualgoDir, workspace, "pkg", pkg+".a"))
			if err != nil {
				return errors.Wrapf(err, "Couldn't remove installed package files for '%s'", workspace)
			}
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(uninstallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uninstallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uninstallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
