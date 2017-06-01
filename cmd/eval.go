// Copyright © 2017 Stream
//

package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// evalCmd represents the eval command
var evalCmd = &cobra.Command{
	Use:   "eval",
	Short: "Expose the script that should be eval-ed in the current shell",
	Long: `For bash put this in your .bashrc:

	eval "$(vg eval)"

Or for zsh, put his in your .zshrc:

	eval "$(vg eval --shell zsh)"
	
Or for fish, put this in your config.fish:

	vg eval --shell fish | source
	`,
	RunE: func(cmd *cobra.Command, args []string) error {

		shell := cmd.Flag("shell").Value.String()

		if shell != "fish" {
			data, err := Asset("data/sh")
			if err != nil {
				return errors.New("Basic POSIX shell commands were not found, this is a bug")
			}
			fmt.Print(string(data))
		}

		data, err := Asset("data/" + shell)
		if err != nil {
			return errors.New("This shell is not supported at the moment")
		}
		fmt.Print(string(data))

		if shell == "bash" {
			err := RootCmd.GenBashCompletion(os.Stdout)
			if err != nil {
				return errors.New("Bash completions could not be generated")
			}
		}
		return nil

	},
}

func init() {
	RootCmd.AddCommand(evalCmd)
	evalCmd.Flags().StringP("shell", "s", "bash", "The shell for which the a file to output")
}
