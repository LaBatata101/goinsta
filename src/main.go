// TODO: use cobra for the cli
//  https://github.com/spf13/cobra

package main

import (
	// "goinsta/ui"
	// "log"

	"goinsta/cmd"
	// tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cmd.Execute()
	// p := tea.NewProgram(ui.NewReviewModel())
	// if _, err := p.Run(); err != nil {
	// 	log.Fatal(err)
	// }
}
