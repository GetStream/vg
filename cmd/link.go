// Copyright © 2017 Stream
//

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// linkCmd represents the link command
var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Link the current virtualgo workspace to the this directory",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := os.Getenv("VIRTUALGO")
		if name == "" {
			return errors.New("A virtualgo workspace should be activated first by using `vg activate [workspaceName]`")
		}

		curdir, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "Couldn't get current working directory")
		}

		fmt.Printf("Linking workspace '%s' to %s\n", name, curdir)

		err = ioutil.WriteFile(".virtualgo", []byte(name), 0644)
		if err != nil {
			return errors.Wrap(err, "Something went wrong when writing the file")
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(linkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// linkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// linkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
