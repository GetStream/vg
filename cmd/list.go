// Copyright © 2017 Stream
//

package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/GetStream/vg/internal/utils"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all existing workspaces",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		files, err := ioutil.ReadDir(utils.VirtualgoRoot())
		if err != nil {
			return err
		}

		for _, file := range files {
			fmt.Println(file.Name())
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
