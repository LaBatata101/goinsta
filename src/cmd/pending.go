package cmd

import (
	"fmt"
	"goinsta/snapshot"
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pendingCmd)
}

var pendingCmd = &cobra.Command{
	Use:   "pending-snapshots",
	Short: "List all pending snapshots",
	Run: func(cmd *cobra.Command, args []string) {
		snapshots, err := snapshot.GetNewSnapshotPaths()
		if err != nil {
			log.Fatal("An error ocurred while getting .snap.new snapshots: ", err)
		}

		for _, snap := range snapshots {
			fmt.Println(snap)
		}
	},
}
