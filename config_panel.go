package main

import (
	"charm.land/lipgloss/v2"
)

func (m model) renderConfigPanel(width, height int) string {
	inner := panelInnerSize(width, height)
	var lines []string
	lines = append(lines, dimStyle.Render("Tunnels"))
	lines = append(lines, "")

	if m.listLen() == 0 {
		lines = append(lines, dimStyle.Render("No configs."))
	} else {
		idx := 0
		if m.netbirdAvail {
			lines = append(lines, m.renderIndexItem(indexRow{
				name: netbirdRowName, kind: "NetBird",
				connected: m.netbirdStatus.Connected(), index: idx, width: inner.width,
			})...)
			idx++
		}
		if m.warpAvail {
			lines = append(lines, m.renderIndexItem(indexRow{
				name: warpRowName, kind: "WARP",
				connected: m.warpStatus.Connected(), index: idx, width: inner.width,
			})...)
			idx++
		}
		for _, cfg := range m.configs {
			lines = append(lines, m.renderIndexItem(indexRow{
				name: cfg, kind: "WireGuard",
				connected: m.isActive(cfg), index: idx, width: inner.width,
			})...)
			idx++
		}
	}

	return renderPanel(width, height, lines)
}

type indexRow struct {
	name      string
	kind      string
	connected bool
	index     int
	width     int
}

func (m model) renderIndexItem(row indexRow) []string {
	selected := row.index == m.cursor

	if selected && m.modal == modalRenaming {
		return []string{m.renameInput.View(), ""}
	}
	if selected && m.modal == modalDeleting {
		return []string{
			errorStyle.Render("Delete "+row.name+"?") + " " + dimStyle.Render("[y/n]"),
			"",
		}
	}

	status := statusWord(row.connected)
	if selected && m.modal == modalConnecting && m.connectName == row.name {
		status = m.spinner.View()
	}
	name := truncate(row.name, row.width)
	meta := row.kind + " · " + status

	nameLine := name
	metaLine := dimStyle.Render(truncate(meta, row.width))
	if selected {
		nameLine = selectedItemStyle.Render(name)
		if row.connected {
			metaLine = connectedStyle.Render(truncate(meta, row.width))
		}
	} else if row.connected {
		metaLine = connectedStyle.Render(truncate(meta, row.width))
	}

	return []string{nameLine, metaLine, ""}
}

func formatIndexRow(name, status, kind string, width int) string {
	return truncate(name, width) + "\n" + truncate(kind+" · "+status, width)
}

func truncate(s string, width int) string {
	if width <= 0 {
		return s
	}
	if lipgloss.Width(s) <= width {
		return s
	}
	if width <= 1 {
		return "…"
	}
	runes := []rune(s)
	for lipgloss.Width(string(runes)+"…") > width && len(runes) > 0 {
		runes = runes[:len(runes)-1]
	}
	return string(runes) + "…"
}
