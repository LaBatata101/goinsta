package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(reviewCmd)
}

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Interactively review snapshots",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Review")
	},
}
