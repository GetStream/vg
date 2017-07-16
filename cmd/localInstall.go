// Copyright © 2017 Stream
//

package cmd

import (
	"errors"
	"path/filepath"

	"github.com/GetStream/vg/utils"
	"github.com/spf13/cobra"
)

// localInstallCmd represents the localInstall command
var localInstallCmd = &cobra.Command{
	Use:   "localInstall <package> [path]",
	Short: "Installs a package from your filesystem inside the workspace",
	Long: `To install a package from your global GOPATH inside the workspace:

	vg localInstall github.com/pkg/errors
	
If you want to install a pacage from a specific path:

	vg localInstall github.com/pkg/errors

It is important to note that these installs are persistent and are not undone
by running 'vg ensure'. To remove a locally installed package use 'vg uninstall'.
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("No package specified")
		}
		if len(args) > 2 {
			return errors.New("Too many arguments suplied")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		settings, err := utils.CurrentSettings()
		if err != nil {
			return err
		}

		path := ""
		pkg := args[0]
		if len(args) == 2 {
			path = args[1]
		} else {

			path = filepath.Join(
				utils.OriginalGopath(),
				"src",
				utils.PkgToDir(pkg),
			)
		}
		settings.LocalInstalls[pkg] = utils.LocalInstall{
			Path: path,
		}

		err = utils.SaveCurrentSettings(settings)
		if err != nil {
			return err
		}

		return utils.LinkCurrentLocalInstalls()
	},
}

func init() {
	RootCmd.AddCommand(localInstallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// localInstallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// localInstallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
