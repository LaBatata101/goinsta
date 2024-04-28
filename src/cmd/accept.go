package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(acceptCmd)
}

var acceptCmd = &cobra.Command{
	Use:   "accept",
	Short: "Accept all snapshots",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Accept")
	},
}
