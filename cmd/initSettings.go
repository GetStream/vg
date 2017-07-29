// Copyright © 2017 Stream
//

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/GetStream/vg/internal/utils"
	"github.com/GetStream/vg/internal/workspace"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// initSettingsCmd represents the initSettings command
var initSettingsCmd = &cobra.Command{
	Use:    "initSettings [workspaceName]",
	Hidden: true,
	Short:  "This command initializes the settings file for a certain workspace",
	Long:   ``,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("Too much arguments specified")
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) (err error) {
		name := ""
		cwd, err := os.Getwd()
		if err != nil {
			return errors.WithStack(err)
		}

		if len(args) == 1 {
			name = args[0]
		} else {
			name = filepath.Base(cwd)

		}
		fmt.Println(name)
		ws := workspace.New(name)

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			return errors.WithStack(err)
		}

		// Check if it's a new workspace. Only continue if this is the case or
		// if force is set.
		_, err = os.Stat(ws.Path())
		if err != nil {
			if !os.IsNotExist(err) {
				return errors.WithStack(err)
			}
		} else if !force {
			return nil
		}

		settings := workspace.NewSettings()
		settings.GlobalFallback, err = cmd.Flags().GetBool("global-fallback")

		if err != nil {
			return errors.WithStack(err)
		}

		err = os.MkdirAll(ws.Path(), 0755)
		if err != nil {
			return errors.WithStack(err)
		}

		err = ws.SaveSettings(settings)
		if err != nil {
			return err
		}

		originalSrcPath := filepath.Join(utils.OriginalGopath(), "src") + string(filepath.Separator)
		if settings.GlobalFallback || !strings.HasPrefix(cwd, originalSrcPath) {
			// Finished no need to do a local install of the current
			// directory

			return nil
		}

		// If current directory is inside the current gopath
		// add it to the packages that need to be symlinked
		pkgDir := strings.TrimPrefix(cwd, originalSrcPath)

		// Make sure pkg is slash seperated
		pkg := utils.DirToPkg(pkgDir)

		return ws.InstallLocalPackagePersistently(pkg, cwd)
	},
}

func init() {
	RootCmd.AddCommand(initSettingsCmd)
	initSettingsCmd.PersistentFlags().Bool("global-fallback", false, "Fallback to global packages when they are not present in workspace")
	initSettingsCmd.PersistentFlags().BoolP("force", "f", false, "")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initSettingsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initSettingsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
