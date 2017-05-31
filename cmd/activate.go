// Copyright © 2017 Stream
//

package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// activateCmd represents the activate command
var activateCmd = &cobra.Command{
	Use:   "activate [workspaceName]",
	Short: "Activate a specific virtualgo workspace",
	Long: `
The most simple way to use it is to activate an workspace named after the
current directory, by just calling:

	vg activate

If you want to activate a specific workspace you can specify it as an
optional argument like this:

	vg activate my-personal-env
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("You haven't eval-ed `vg eval` yet.")
	},
}

func init() {
	RootCmd.AddCommand(activateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// activateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// activateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
