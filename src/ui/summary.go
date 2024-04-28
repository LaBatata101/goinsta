package ui

import (
	"bufio"
	"fmt"
	"goinsta/snapshot"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

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
	s2 := fmt.Sprintf("Snapshot: %s", yellowText.Render(snap.Name))
	s3 := fmt.Sprintf("Source: %s:%s", greenText2.Render(snap.Source),
		lipgloss.NewStyle().Bold(true).Render(strconv.Itoa(snap.Loc)))

	return lipgloss.JoinVertical(0, header, s1, s2, s3, strings.Repeat("─", termWidth))
}

func diffView(termWidth int, snap *snapshot.Snapshot) string {
	line := strings.Repeat("─", termWidth)

	var totalLoc int
	var coloredLines []string
	scanner := bufio.NewScanner(strings.NewReader(snap.Diff()))
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "+"):
			coloredLines = append(coloredLines, GreenText.Render(line))
		case strings.HasPrefix(line, "-"):
			coloredLines = append(coloredLines, redText.Render(line))
		default:
			coloredLines = append(coloredLines, line)
		}
		totalLoc++
	}

	var lineNumbers []string
	for i := range totalLoc {
		lineNumbers = append(lineNumbers, lineNumberColor.Render(strconv.Itoa(i+1)))
	}

	lineNumbersText := strings.Join(lineNumbers, " \n")
	diffText := strings.Join(coloredLines, "\n")

	sourceBorder := lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderLeft(true)

	var header string
	if snap.IsNew() {
		header = lipgloss.JoinVertical(0, GreenText.Render("+new results"), line)
	} else {
		header = lipgloss.JoinVertical(0, redText.Render("-old snapshot"), GreenText.Render("+new results"), line)
	}

	return lipgloss.JoinVertical(0.05, header,
		lipgloss.JoinHorizontal(0, lineNumbersText, sourceBorder.Render(diffText)), line)
}
