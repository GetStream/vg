// Copyright © 2017 Stream
//

package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// cdpackagesCmd represents the cdpackages command
var cdpackagesCmd = &cobra.Command{
	Use:   "cdpackages",
	Short: "Change the working directory to the src directory of the active workspace",
	Long: `Example:

	$ vg activate test
	$ vg cdpackages
	$ echo $PWD
	/home/user/.virtualgo/test/src
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New(noEvalError)
	},
}

func init() {
	RootCmd.AddCommand(cdpackagesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cdpackagesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cdpackagesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
