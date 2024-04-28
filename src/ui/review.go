package ui

import (
	"goinsta/snapshot"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbletea"
)

type reviewModel struct {
	snapshots     []snapshot.Snapshot
	currSnapIndex int
	paginator     paginator.Model
}

func NewReviewModel() reviewModel {
	var snapshots []snapshot.Snapshot
	snapPaths, err := snapshot.GetNewSnapshotPaths()
	if err != nil {
		log.Fatal("Failed to get snapshot paths: ", err)
	}

	for _, snapshotPath := range snapPaths {
		// Don't need to handle error here, since, we have valid snap paths at this point.
		snap, _ := snapshot.Read(snapshotPath)
		snapshots = append(snapshots, snap)
	}

	p := paginator.New()
	p.Type = paginator.Arabic
	p.PerPage = 1
	p.ArabicFormat = greenText2.Bold(true).Render("  Reviewing: ") + "[" + yellowText.Bold(true).Render("%d/%d") + "]"
	p.KeyMap.NextPage = key.NewBinding(key.WithDisabled())
	p.KeyMap.PrevPage = key.NewBinding(key.WithDisabled())
	p.SetTotalPages(len(snapshots))

	return reviewModel{
		snapshots,
		0,
		p,
	}
}

func (m reviewModel) Init() tea.Cmd {
	return nil
}

func (m reviewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			m.snapshots[m.currSnapIndex].Accept()
			m.paginator.NextPage()
			m.currSnapIndex++
		case "r":
			m.snapshots[m.currSnapIndex].Reject()
			m.paginator.NextPage()
			m.currSnapIndex++
		case "s":
			m.paginator.NextPage()
			m.currSnapIndex++
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	}

	if m.currSnapIndex >= len(m.snapshots) {
		return m, tea.Quit
	}

	m.paginator, cmd = m.paginator.Update(msg)
	return m, cmd
}

func (m reviewModel) View() string {
	var b strings.Builder

	start, end := m.paginator.GetSliceBounds(len(m.snapshots))
	for _, snapshot := range m.snapshots[start:end] {
		b.WriteString(SnapshotSummary(&snapshot))
	}

	b.WriteString("\n" + m.paginator.View())
	b.WriteString("\n\n")
	b.WriteString("  " + GreenText.Render("a") + " accept " + grayText.Render("keep the new snapshot") + "\n")
	b.WriteString("  " + redText.Render("r") + " reject " + grayText.Render("reject the new snapshot") + "\n")
	b.WriteString("  " + yellowText.Render("s") + " skip   " + grayText.Render("keep both for now") + "\n")
	b.WriteString("  " + redText.Bold(true).Render("q quit   ") + grayText.Render("stop reviewing") + "\n")

	return b.String()
}
