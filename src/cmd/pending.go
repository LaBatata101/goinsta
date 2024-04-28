package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pendingCmd)
}

var pendingCmd = &cobra.Command{
	Use:   "pending-snapshots",
	Short: "List all pending snapshots",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Pending")
	},
}
