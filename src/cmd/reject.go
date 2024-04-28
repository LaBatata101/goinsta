package cmd

import (
	"fmt"
	"goinsta/snapshot"
	"goinsta/ui"
	"log"

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

		acceptedSnaps, err := snapshot.RejectAll(snapshots)
		if err != nil {
			log.Fatal("An error ocurred while accepting snapshots: ", err)
		}

		fmt.Println(ui.RedText.Render("Rejected") + ":")
		for _, snap := range acceptedSnaps {
			fmt.Printf("  %s (%s)\n", snap.Source, snap.Name)
		}
	},
}
