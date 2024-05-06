package ui

import (
	"strings"

	"github.com/LaBatata101/goinsta/internal/snapshot"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type reviewModel struct {
	snapshots       []snapshot.Snapshot
	currSnapIndex   int
	paginator       paginator.Model
	viewport        viewport.Model
	isViewportReady bool
	drawScrollBar   bool
	summary         *snapshot.Summary
	windowHeight    int
	windowWidth     int
}

func ReviewSnapshotsModel(snapPaths []string, summary *snapshot.Summary) reviewModel {
	var snapshots []snapshot.Snapshot
	for _, snapshotPath := range snapPaths {
		// Don't need to handle error here, since, we have valid snap paths at this point.
		snap, _ := snapshot.Read(snapshotPath)
		snapshots = append(snapshots, snap)
	}

	p := paginator.New()
	p.Type = paginator.Arabic
	p.PerPage = 1
	p.ArabicFormat = greenText2.Bold(true).Render("  Reviewing: ") + "[" + YellowText.Bold(true).Render("%d/%d") + "]"
	p.KeyMap.NextPage = key.NewBinding(key.WithDisabled())
	p.KeyMap.PrevPage = key.NewBinding(key.WithDisabled())
	p.SetTotalPages(len(snapshots))

	return reviewModel{
		snapshots:     snapshots,
		currSnapIndex: 0,
		paginator:     p,
		summary:       summary,
	}
}

const useHighPerformanceRenderer = false

func (m reviewModel) currSnapshot() *snapshot.Snapshot {
	return &m.snapshots[m.currSnapIndex]
}

func (m reviewModel) Init() tea.Cmd {
	return nil
}

func (m reviewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd                  tea.Cmd
		cmds                 []tea.Cmd
		verticalMarginHeight int
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			m.summary.AddAccepted(m.snapshots[m.currSnapIndex])
			m.snapshots[m.currSnapIndex].Accept()
			m.currSnapIndex++
			if m.currSnapIndex < len(m.snapshots) {
				m.paginator.NextPage()
			}
		case "r":
			m.summary.AddRejected(m.snapshots[m.currSnapIndex])
			m.snapshots[m.currSnapIndex].Reject()
			m.currSnapIndex++
			if m.currSnapIndex < len(m.snapshots) {
				m.paginator.NextPage()
			}
		case "s":
			m.summary.AddSkipped(m.snapshots[m.currSnapIndex])
			m.currSnapIndex++
			if m.currSnapIndex < len(m.snapshots) {
				m.paginator.NextPage()
			}
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.windowHeight, m.windowWidth = msg.Height, msg.Width

		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight = headerHeight + footerHeight

		if !m.isViewportReady {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.isViewportReady = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		if useHighPerformanceRenderer {
			// Render (or re-render) the whole viewport. Necessary both to
			// initialize the viewport and when the window is resized.
			//
			// This is needed for high-performance rendering only.
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}

	if m.currSnapIndex >= len(m.snapshots) {
		return m, tea.Quit
	}

	snap := m.currSnapshot()
	contentHeight := lipgloss.Height(snap.Diff())
	if contentHeight > m.viewport.Height {
		if !m.drawScrollBar {
			m.drawScrollBar = true
			m.viewport.Width -= lipgloss.Width(scrollBarBlock)
		}
	} else {
		m.drawScrollBar = false
	}
	m.viewport.SetContent(diffView(m.viewport.Width, snap))

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	m.paginator, cmd = m.paginator.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m reviewModel) View() string {
	var b strings.Builder

	if m.currSnapIndex >= len(m.snapshots) {
		return ""
	}

	scrollBar := ""
	if m.drawScrollBar {
		scrollBar = m.scrollBar()
	}

	// FIX: viewport is not moved 0.05% to the right
	b.WriteString(lipgloss.JoinVertical(0.05, m.headerView(),
		// The `viewport` contains the diff view
		lipgloss.JoinHorizontal(0, m.viewport.View(), scrollBar)))
	b.WriteString(m.footerView())

	return b.String()
}

const scrollBarBlock = "█▍"

func (m reviewModel) scrollBar() string {
	percent := m.viewport.ScrollPercent()
	pos := int(float64(m.viewport.Height) * percent)
	return strings.Repeat("\n", max(0, pos-lipgloss.Height(scrollBarBlock))) + scrollBarBlock
}

func (m reviewModel) headerView() string {
	return lipgloss.JoinVertical(0, summaryHeader(m.windowWidth, m.currSnapshot()),
		diffHeader(m.windowWidth, m.currSnapshot()))
}

func (m reviewModel) footerView() string {
	var b strings.Builder
	b.WriteString("\n" + strings.Repeat("─", m.windowWidth))
	b.WriteString("\n" + m.paginator.View())
	b.WriteString("\n\n")
	b.WriteString("  " + GreenText.Render("a") + " accept " + grayText.Render("keep the new snapshot") + "\n")
	b.WriteString("  " + RedText.Render("r") + " reject " + grayText.Render("reject the new snapshot") + "\n")
	b.WriteString("  " + YellowText.Render("s") + " skip   " + grayText.Render("keep both for now") + "\n")
	b.WriteString("  " + RedText.Bold(true).Render("q quit   ") + grayText.Render("stop reviewing") + "\n")
	return b.String()
}
