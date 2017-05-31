// Copyright © 2017 Jelte Fennema <license-tech@jeltef.nl>
//

package cmd

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect the current virtualgo workspace to the this directory",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := os.Getenv("VIRTUALGO")
		if name == "" {
			return errors.New("A virtualgo workspace should be activated first by using `vg activate [workspaceName]`")
		}
		err := ioutil.WriteFile(".virtualgo", []byte(name), 0644)
		if err != nil {
			return errors.Wrap(err, "Something went wrong when writing the file")
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(connectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// connectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// connectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
