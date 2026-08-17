package main

import (
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m model) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("")
	}

	// File picker
	if m.modal == modalImporting {
		v := tea.NewView(m.filePicker.View())
		v.AltScreen = true
		return v
	}

	// Help overlay replaces everything
	if m.modal == modalHelp {
		helpView := m.help.View(m.keys)
		overlay := helpOverlayStyle.Render(
			helpTitleStyle.Render("Keyboard shortcuts") + "\n\n" +
				helpView + "\n\n" +
				dimStyle.Render("omarchy-vpn "+displayVersion()+" · any key closes"),
		)
		v := tea.NewView(lipgloss.Place(
			m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			overlay,
		))
		v.AltScreen = true
		return v
	}

	// Layout: header (2) + gap (1) + panels + footer (1)
	titleBar := m.renderTitleBar()
	titleHeight := 3
	bottomHeight := 1
	panelHeight := m.height - titleHeight - bottomHeight
	if panelHeight < 5 {
		panelHeight = 5
	}

	// Panel widths: 40/60 split
	leftWidth := m.width * 2 / 5
	rightWidth := m.width - leftWidth

	if leftWidth < 24 {
		leftWidth = 24
		rightWidth = m.width - leftWidth
	}

	// Render panels
	left := m.renderConfigPanel(leftWidth, panelHeight)
	right := m.renderStatusPanel(rightWidth, panelHeight)

	// Join panels horizontally
	panels := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	// Bottom bar
	var bottom string
	if m.message != "" && time.Now().Before(m.messageExp) {
		bottom = " " + m.message
	} else {
		m.message = ""
		bottom = " " + dimStyle.Render("?")
	}

	v := tea.NewView(lipgloss.JoinVertical(lipgloss.Left,
		titleBar,
		panels,
		bottom,
	))
	v.AltScreen = true
	return v
}

func (m model) sessionConnected() bool {
	return len(m.activeVPNs) > 0 || m.netbirdStatus.Connected() || m.warpStatus.Connected()
}

func (m model) sessionNames() []string {
	names := append([]string{}, m.activeVPNs...)
	if m.netbirdStatus.Connected() {
		names = append(names, "NetBird")
	}
	if m.warpStatus.Connected() {
		names = append(names, "Cloudflare WARP")
	}
	return names
}

func statusWord(connected bool) string {
	if connected {
		return "Connected"
	}
	return "Disconnected"
}

func (m model) renderTitleBar() string {
	line1 := " " + itemStyle.Render("omarchy-vpn")

	word := statusWord(m.sessionConnected())
	var status string
	if m.sessionConnected() {
		status = connectedStyle.Render(word)
	} else {
		status = dimStyle.Render(word)
	}

	names := strings.Join(m.sessionNames(), "  ")
	left := " " + status
	if names != "" {
		left += "    " + valueStyle.Render(names)
	}
	right := dimStyle.Render(displayVersion())
	pad := m.width - lipgloss.Width(left) - lipgloss.Width(right) - 1
	if pad < 2 {
		pad = 2
	}
	line2 := left + strings.Repeat(" ", pad) + right

	return line1 + "\n" + line2 + "\n"
}
