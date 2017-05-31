// Copyright © 2017 Jelte Fennema <license-tech@jeltef.nl>
//

package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// evalCmd represents the eval command
var evalCmd = &cobra.Command{
	Use:   "eval",
	Short: "Expose the script that should be eval-ed in the current shell",
	Long: `For POSIX compatible shells put the following in your .bashrc/.zshrc/.whateverrc:

	eval "$(vg eval)"
	
Or for fish, put this in your config.fish:

	vg eval --shell fish | source
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		shell := cmd.Flag("shell").Value.String()
		data, err := Asset("data/" + shell)
		if err != nil {
			return errors.New("This shell is not supported at the moment")
		}
		fmt.Print(string(data))
		return nil

	},
}

func init() {
	RootCmd.AddCommand(evalCmd)
	evalCmd.Flags().StringP("shell", "s", "sh", "The shell for which the a file to output")
}
