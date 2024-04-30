package cmd

import (
	"fmt"
	"log"

	"github.com/LaBatata101/goinsta/internal/snapshot"
	"github.com/LaBatata101/goinsta/internal/ui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(acceptCmd)
}

var acceptCmd = &cobra.Command{
	Use:   "accept",
	Short: "Accept all snapshots",
	Run: func(cmd *cobra.Command, args []string) {
		snapshots, err := snapshot.GetNewSnapshotPaths()
		if err != nil {
			log.Fatal("An error ocurred while getting .snap.new snapshots: ", err)
		}

		if len(snapshots) == 0 {
			fmt.Println("no snapshots to review")
			return
		}

		acceptedSnaps, err := snapshot.AcceptAll(snapshots)
		if err != nil {
			log.Fatal("An error ocurred while accepting snapshots: ", err)
		}
		ui.PrintAccepted(acceptedSnaps)
	},
}
