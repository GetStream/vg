// Copyright © 2017 Stream
//

package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/GetStream/vg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// initSettingsCmd represents the initSettings command
var initSettingsCmd = &cobra.Command{
	Use:   "initSettings [workspaceName]",
	Short: "This command initializes the settings file for a certain workspace",
	Long:  ``,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("Too much arguments specified")
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) (err error) {
		workspace := ""
		cwd, err := os.Getwd()
		if err != nil {
			return errors.WithStack(err)
		}

		if len(args) == 1 {
			workspace = args[0]
		} else {
			workspace = filepath.Base(cwd)

		}
		fmt.Println(workspace)

		settingsPath := utils.SettingsPath(workspace)
		if err != nil {
			return errors.WithStack(err)
		}

		dir := filepath.Dir(settingsPath)

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			return errors.WithStack(err)
		}

		// Check if it's a new workspace. Only continue if this is the case or
		// if force is set.
		_, err = os.Stat(dir)
		if err != nil {
			if !os.IsNotExist(err) {
				return errors.WithStack(err)
			}
		} else if !force {
			return nil
		}

		settings := utils.NewWorkspaceSettings()
		settings.GlobalFallback, err = cmd.Flags().GetBool("global-fallback")

		if err != nil {
			return errors.WithStack(err)
		}

		srcpath := filepath.Join(utils.CurrentGopath(), "src") + string(filepath.Separator)

		if strings.HasPrefix(cwd, srcpath) && !settings.GlobalFallback {
			// If current directory is inside the current gopath
			// add it to the packages that need to be symlinked
			pkgDir := strings.TrimPrefix(cwd, srcpath)

			// Make sure pkg is slash seperated
			pkgComponents := filepath.SplitList(pkgDir)
			pkg := path.Join(pkgComponents...)

			settings.LocalInstalls[pkg] = utils.LocalInstall{
				Path: cwd,
			}

		}

		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return errors.WithStack(err)
		}

		err = utils.SaveSettings(workspace, settings)
		if err != nil {
			return err
		}

		return utils.LinkLocalInstalls(workspace, settings)
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
