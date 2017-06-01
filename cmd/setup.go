// Copyright © 2017 Stream
//

package cmd

import (
	"os"
	"os/user"
	"path/filepath"

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
		err = appendToFile("~/.bashrc", "\neval \"$(vg eval --shell bash)\"\n")
		if err != nil {
			return err
		}

		err = appendToFile("~/.zshrc", "\neval \"$(vg eval --shell zsh)\"\n")
		if err != nil {
			return err
		}

		err = os.MkdirAll("~/.config/fish/", 0755)
		if err != nil {
			return err
		}

		err = appendToFile("~/.config/fish/config.fish", "\nvg eval --shell fish | source\n")
		return err
	},
}

func appendToFile(fileName string, content string) error {
	if fileName[:2] == "~/" {
		usr, err := user.Current()
		if err != nil {
			return err
		}
		homeDir := usr.HomeDir
		fileName = filepath.Join(homeDir, fileName[2:])
	}

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	_, err = file.WriteString(content)
	return err
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
