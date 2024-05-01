package ui

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/LaBatata101/goinsta/internal/snapshot"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wrap"
	"golang.org/x/term"
)

func PrintSummary(summary *snapshot.Summary) {
	fmt.Println(BoldText.Render("\nreview finished"))
	if len(summary.Accepted) > 0 {
		PrintAccepted(summary.Accepted)
	}

	if len(summary.Rejected) > 0 {
		PrintReject(summary.Rejected)
	}

	if len(summary.Skipped) > 0 {
		PrintSkipped(summary.Skipped)
	}
}

func PrintAccepted(snaps []snapshot.Snapshot) {
	fmt.Println(GreenText.Render("Accepted") + ":")
	for _, snap := range snaps {
		fmt.Printf("  %s (%s)\n", snap.Source, snap.Name)
	}
}

func PrintReject(snaps []snapshot.Snapshot) {
	fmt.Println(RedText.Render("Rejected") + ":")
	for _, snap := range snaps {
		fmt.Printf("  %s (%s)\n", snap.Source, snap.Name)
	}
}

func PrintSkipped(snaps []snapshot.Snapshot) {
	fmt.Println(YellowText.Render("Skipped") + ":")
	for _, snap := range snaps {
		fmt.Printf("  %s (%s)\n", snap.Source, snap.Name)
	}
}

func RenderSnapshotSummary(snap *snapshot.Snapshot) {
	fmt.Println(SnapshotSummary(snap))
}

func SnapshotSummary(snap *snapshot.Snapshot) string {
	termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		termWidth = 100
		log.Println("Failed to get terminal size, using default size")
	}

	header := Header(termWidth, snap)
	diff := diffView(termWidth, snap)
	return lipgloss.JoinVertical(0, header, diff)
}

func Header(termWidth int, snap *snapshot.Snapshot) string {
	headerText := lipgloss.NewStyle().
		SetString(" Snapshot Summary ").
		Inherit(BoldText).
		Render()
	headerTextWidth, _ := lipgloss.Size(headerText)
	headerLine := strings.Repeat("━", (termWidth-headerTextWidth)/2)

	header := fmt.Sprintf("%s%s%s", headerLine, headerText, headerLine)

	s1 := fmt.Sprintf("Snapshot file: %s", GreenText2Underlined.Render(snap.CleanPath()))
	s2 := fmt.Sprintf("Snapshot: %s", YellowText.Render(snap.Name))
	s3 := fmt.Sprintf("Source: %s:%s", greenText2.Render(snap.Source),
		lipgloss.NewStyle().Bold(true).Render(strconv.Itoa(snap.Loc)))

	return lipgloss.JoinVertical(0, header, s1, s2, s3, strings.Repeat("─", termWidth))
}

func diffView(termWidth int, snap *snapshot.Snapshot) string {
	lineSeparator := strings.Repeat("─", termWidth)

	var loc int
	var coloredLines []string
	var lineNumbersColumn []string
	scanner := bufio.NewScanner(strings.NewReader(snap.Diff()))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		line = wrap.String(line, int(float32(termWidth)*0.98)) // wrap the line to 98% of the terminal width
		switch {
		case strings.HasPrefix(line, "+"):
			coloredLines = append(coloredLines, GreenText.Render(line))
		case strings.HasPrefix(line, "-"):
			coloredLines = append(coloredLines, RedText.Render(line))
		default:
			coloredLines = append(coloredLines, line)
		}
		loc++
		lineNumbersColumn = append(lineNumbersColumn, lineNumberColor.Render(strconv.Itoa(loc)))

		// Add some space for wrapped lines before showing the next number.
		if lineHeight := lipgloss.Height(line); lineHeight > 1 {
			for range lineHeight - 1 {
				lineNumbersColumn = append(lineNumbersColumn, " ")
			}
		}
	}

	lineNumbersText := strings.Join(lineNumbersColumn, " \n")
	diffText := strings.Join(coloredLines, "\n")

	sourceBorder := lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderLeft(true)

	var header string
	if snap.IsNew() {
		header = lipgloss.JoinVertical(0, GreenText.Render("+new results"), lineSeparator)
	} else {
		header = lipgloss.JoinVertical(0, RedText.Render("-old snapshot"), GreenText.Render("+new results"), lineSeparator)
	}

	return lipgloss.JoinVertical(0.05, header,
		lipgloss.JoinHorizontal(0, lineNumbersText, sourceBorder.Render(diffText)), lineSeparator)
}
