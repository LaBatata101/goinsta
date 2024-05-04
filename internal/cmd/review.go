package cmd

import (
	"fmt"
	"log"

	"github.com/LaBatata101/goinsta/internal/snapshot"
	"github.com/LaBatata101/goinsta/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(reviewCmd)
}

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Interactively review snapshots",
	Run: func(cmd *cobra.Command, args []string) {
		snapshots, err := snapshot.GetNewSnapshotPaths()
		if err != nil {
			log.Fatal("An error ocurred while getting .snap.new snapshots: ", err)
		}

		if len(snapshots) == 0 {
			fmt.Println("no snapshots to review")
			return
		}

		rc := snapshot.Summary{}
		model := ui.ReviewSnapshotsModel(snapshots, &rc)
		p := tea.NewProgram(model)
		p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}

		ui.PrintSummary(&rc)
	},
}
