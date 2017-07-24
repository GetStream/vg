// Copyright © 2017 Stream
//

package cmd

import (
	"path/filepath"

	"github.com/GetStream/vg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// localInstallCmd represents the localInstall command
var localInstallCmd = &cobra.Command{
	Use:   "localInstall <package> [path]",
	Short: "Installs a package from your filesystem inside the workspace",
	Long: `To install a package from your global GOPATH inside the workspace:

	vg localInstall github.com/pkg/errors
	
If you want to install a pacage from a specific path:

	vg localInstall github.com/pkg/errors ~/some/path/errors

It is important to note that by default these installs are undone when running
'vg ensure' or 'vg moveVendor'. To make sure the local installs are still
present after running these commands you can use the '--persistent' flag.

	vg localInstall github.com/pkg/errors --persistent

To remove a persistently installed local package use 'vg uninstall <pkg>'.
After that a 'vg ensure' will install like normal again.
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

		persist, err := cmd.Flags().GetBool("persistent")
		if err != nil {
			return errors.WithStack(err)
		}

		if !persist {
			return utils.InstallCurrentLocalPackage(pkg, path)
		}

		settings, err := utils.CurrentSettings()
		if err != nil {
			return err
		}

		settings.LocalInstalls[pkg] = utils.LocalInstall{
			Path: path,
		}

		err = utils.SaveCurrentSettings(settings)
		if err != nil {
			return err
		}

		return utils.InstallCurrentPersistentLocalPackages()
	},
}

func init() {
	RootCmd.AddCommand(localInstallCmd)
	localInstallCmd.Flags().BoolP("persistent", "p", false, "Persist local install")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// localInstallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// localInstallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
