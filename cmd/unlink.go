// Copyright © 2017 Stream
//

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// unlinkCmd represents the unlink command
var unlinkCmd = &cobra.Command{
	Use:   "unlink",
	Short: "Unlinks the current directory from the workspace it's linked to",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_ = os.Remove(".virtualgo")
	},
}

func init() {
	RootCmd.AddCommand(unlinkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// unlinkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// unlinkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
