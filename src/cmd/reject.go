package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rejectCmd)
}

var rejectCmd = &cobra.Command{
	Use:   "reject",
	Short: "Reject all snapshots",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Reject")
	},
}
