// Copyright © 2017 Stream
//

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/GetStream/vg/internal/utils"
	"github.com/GetStream/vg/internal/workspace"
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
		ws, err := workspace.Current()
		if err != nil {
			return err
		}

		for _, pkg := range args {
			// pkgComponents := strings.Split(pkg, hello
			pkgDir := utils.PkgToDir(pkg)
			fmt.Printf("Uninstalling %q from workspace\n", pkg)
			err := os.RemoveAll(filepath.Join(ws.Src(), pkgDir))
			if err != nil {
				return errors.Wrapf(err, "Couldn't remove package src '%s'", ws.Name())
			}

			pkgInstalledDirs, err := filepath.Glob(filepath.Join(ws.Pkg(), "*", pkgDir))
			if err != nil {
				return errors.Wrapf(err, "Something went wrong when globbing files for '%s'", pkg)
			}

			for _, path := range pkgInstalledDirs {
				fmt.Println("Removing", path)

				err = os.RemoveAll(path)
				if err != nil {
					return errors.Wrapf(err, "Couldn't remove installed package files for '%s'", pkg)
				}
			}

			pkgInstalledFiles, err := filepath.Glob(filepath.Join(ws.Pkg(), "*", pkgDir+".a"))
			if err != nil {
				return errors.Wrapf(err, "Something went wrong when globbing files for '%s'", pkg)
			}

			for _, path := range pkgInstalledFiles {
				fmt.Println("Removing", path)

				err = os.RemoveAll(path)
				if err != nil {
					return errors.Wrapf(err, "Couldn't remove installed package files for '%s'", pkg)
				}
			}

			settings, err := ws.Settings()
			if err != nil {
				return err
			}

			if _, ok := settings.LocalInstalls[pkg]; ok {
				fmt.Printf("Removing %q from persistent local installs\n", pkg)
				delete(settings.LocalInstalls, pkg)

				err = ws.SaveSettings(settings)
				if err != nil {
					return err
				}
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
