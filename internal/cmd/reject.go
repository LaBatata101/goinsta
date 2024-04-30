package cmd

import (
	"fmt"
	"log"

	"github.com/LaBatata101/goinsta/internal/snapshot"
	"github.com/LaBatata101/goinsta/internal/ui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rejectCmd)
}

var rejectCmd = &cobra.Command{
	Use:   "reject",
	Short: "Reject all snapshots",
	Run: func(cmd *cobra.Command, args []string) {
		snapshots, err := snapshot.GetNewSnapshotPaths()
		if err != nil {
			log.Fatal("An error ocurred while getting .snap.new snapshots: ", err)
		}

		if len(snapshots) == 0 {
			fmt.Println("no snapshots to review")
			return
		}

		rejectedSnaps, err := snapshot.RejectAll(snapshots)
		if err != nil {
			log.Fatal("An error ocurred while accepting snapshots: ", err)
		}
		ui.PrintReject(rejectedSnaps)
	},
}
