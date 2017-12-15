// Copyright © 2017 Stream
//

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/GetStream/vg/internal/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Enables virtualgo in your shell",
	Long: `After running this you have to restart your shell or run:

	source ~/.bashrc                   # for bash
	source ~/.zshrc                    # for zsh
	source ~/.config/fish/config.fish  # for fish
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		vgExistCheck := "command -v vg >/dev/null 2>&1"
		shellInfos := []struct {
			shell       string
			configFile  string
			evalCommand string
		}{
			{
				shell:       "bash",
				configFile:  "~/.bashrc",
				evalCommand: vgExistCheck + ` && eval "$(vg eval --shell bash)"`,
			},
			{
				shell:       "zsh",
				configFile:  "~/.zshrc",
				evalCommand: vgExistCheck + ` && eval "$(vg eval --shell zsh)"`,
			},
			{
				shell:       "fish",
				configFile:  "~/.config/fish/config.fish",
				evalCommand: vgExistCheck + `; and vg eval --shell fish | source`,
			},
		}

		for _, info := range shellInfos {
			// bash
			shellExists, err := utils.CommandExists(info.shell)
			if err != nil {
				return err
			}

			if !shellExists {
				fmt.Printf("Skipping setup for %q shell, because it is not installed\n", info.shell)
				continue
			}

			evalCommandExists, err := lineExists(info.configFile, info.evalCommand)
			if err != nil {
				return err
			}
			if evalCommandExists {
				fmt.Printf("Skipping setup for %q shell, because setup has been performed already\n", info.shell)
				continue
			}

			fmt.Printf("Editing %q to setup %q shell\n", info.configFile, info.shell)
			err = appendToFile(info.configFile, "\n"+info.evalCommand+"\n")
			if err != nil {
				return err
			}
		}
		return err
	},
}

func appendToFile(filename string, content string) error {
	filename = utils.ReplaceHomeDir(filename)
	fileDir := filepath.Dir(filename)
	err := os.MkdirAll(fileDir, 0755)
	if err != nil {
		return errors.WithStack(err)
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	_, err = file.WriteString(content)
	return err
}

func lineExists(filename string, line string) (bool, error) {
	filename = utils.ReplaceHomeDir(filename)
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// If the file doesn't exist the line doesn't exist as well
			return false, nil
		}

		return false, errors.WithStack(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == line {
			return true, nil
		}
	}
	return false, nil

}

func init() {
	RootCmd.AddCommand(setupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
