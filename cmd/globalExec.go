// Copyright © 2017 Stream
//

package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// globalExecCmd represents the globalExec command
var globalExecCmd = &cobra.Command{
	Use:   "globalExec <cmd> [argsuments for cmd...]",
	Short: "Execute a command globally (outside the active workspace)",
	Long: `For instance to update the globally installed dep binary while a
workspace is active:
	
	vg globalExec go get -u github.com/golang/dep/cmd/dep
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New(noEvalError)
	},
}

func init() {
	RootCmd.AddCommand(globalExecCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// globalExecCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// globalExecCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
