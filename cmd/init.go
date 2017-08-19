// Copyright © 2017 Stream
//

package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [workspaceName]",
	Short: "Create and enable a workspace and link it to the current directory",
	Long: `This is normally the command that you need to start using virtualgo
for a project. If you want more control you can use 'vg activate' and 'vg link'
seperately.

The simplest way to use it vg is just to call:
	
	vg init

This will create a workspace named after the current directory, which is
usually a good name. If you want to use a different name, just specify it:

	vg init myCoolWorkspace

To create a workspace in full isolation mode, you should run:

	vg init --full-isolation

A fully isolated workspace has some upsides, but some go commands will not work
as expected anymore. If bindfs is installed it works much better, but relative
arguments still might not work as expected. See the README on Github for full
details on advantages and disadvantages.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("You haven't eval-ed `vg eval` yet.")
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.PersistentFlags().Bool("global-fallback", false, `Fallback to global packages when they are not present in workspace. 
			  This is the default mode if both --full-isolation and --global-fallback are not provided.`)
	initCmd.PersistentFlags().Bool("full-isolation", false, "Create a fully isolated workspace, see project README for advantages and disadvantages")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
