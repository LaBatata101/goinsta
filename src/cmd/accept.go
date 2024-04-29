package cmd

import (
	"fmt"
	"goinsta/snapshot"
	"goinsta/ui"
	"log"

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

		fmt.Println(ui.GreenText.Render("Accepted") + ":")
		for _, snap := range acceptedSnaps {
			fmt.Printf("  %s (%s)\n", snap.Source, snap.Name)
		}
	},
}
