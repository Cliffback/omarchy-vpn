package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"charm.land/lipgloss/v2"
)

// ConfigInfo holds static config details parsed from a .conf file.
type ConfigInfo struct {
	Address    string
	DNS        string
	Endpoint   string
	PeerKey    string
	AllowedIPs []string // one entry per AllowedIPs line, possibly comma-separated
}

type sheetField struct {
	Label string
	Value string
}

// ParseConfigFile reads a WireGuard .conf file and extracts display fields.
func ParseConfigFile(name string) ConfigInfo {
	path := fmt.Sprintf("/etc/wireguard/%s.conf", name)
	out, err := exec.Command("sudo", "cat", path).Output()
	if err != nil {
		data, err2 := os.ReadFile(path)
		if err2 != nil {
			return ConfigInfo{}
		}
		out = data
	}
	data := out

	var info ConfigInfo
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		switch key {
		case "Address":
			info.Address = val
		case "DNS":
			info.DNS = val
		case "Endpoint":
			info.Endpoint = val
		case "PublicKey":
			info.PeerKey = val
		case "AllowedIPs":
			info.AllowedIPs = append(info.AllowedIPs, val)
		}
	}
	return info
}

func (m model) renderStatusPanel(width, height int) string {
	if m.netbirdSelected() {
		return m.renderNetbirdStatusPanel(width, height)
	}
	if m.warpSelected() {
		return m.renderWarpStatusPanel(width, height)
	}

	innerWidth := width - 4
	name := m.selectedConfig()
	if name == "" {
		return renderSheet(width, height, []string{dimStyle.Render("  No configs.")})
	}

	info := ParseConfigFile(name)
	s := m.vpnStatus[name]
	up := m.isActive(name)
	peer := info.PeerKey
	if len(peer) > 24 {
		peer = peer[:24] + "…"
	}

	fields := []sheetField{
		{"Status", statusWord(up)},
		{"Address", info.Address},
		{"DNS", info.DNS},
		{"Endpoint", firstNonEmpty(s.Endpoint, info.Endpoint)},
		{"Peer", peer},
	}
	if up {
		fields = append(fields,
			sheetField{"Download", s.TransferRx},
			sheetField{"Upload", s.TransferTx},
			sheetField{"Handshake", s.Handshake},
		)
	}
	fields = append(fields, sheetField{"Path", fmt.Sprintf("/etc/wireguard/%s.conf", name)})

	return renderSheet(width, height, sheetLines(fields, innerWidth))
}

func (m model) renderNetbirdStatusPanel(width, height int) string {
	innerWidth := width - 4
	s := m.netbirdStatus

	switch {
	case s.Connected():
		fields := []sheetField{
			{"Status", "Connected"},
			{"Address", s.IP},
			{"FQDN", s.FQDN},
			{"Peers", fmt.Sprintf("%d/%d connected", s.PeersConnected, s.PeersTotal)},
			{"Download", formatBytes(s.TransferRx)},
			{"Upload", formatBytes(s.TransferTx)},
			{"Management", s.MgmtURL},
		}
		return renderSheet(width, height, sheetLines(fields, innerWidth))
	case s.DaemonStatus == "":
		return renderSheet(width, height, []string{
			"  " + dimStyle.Render(statusWord(false)),
			"  " + dimStyle.Render("Daemon down"),
			"  " + dimStyle.Render("sudo systemctl enable --now netbird"),
		})
	case s.NeedsLogin():
		return renderSheet(width, height, []string{
			"  " + warnStyle.Render("Login required"),
			"  " + dimStyle.Render("Run netbird up in a terminal to log in."),
		})
	default:
		fields := []sheetField{
			{"Status", firstNonEmpty(s.DaemonStatus, "Disconnected")},
			{"Management", s.MgmtURL},
		}
		return renderSheet(width, height, sheetLines(fields, innerWidth))
	}
}

func (m model) renderWarpStatusPanel(width, height int) string {
	innerWidth := width - 4
	s := m.warpStatus

	switch {
	case s.Connected():
		fields := []sheetField{
			{"Status", "Connected"},
			{"Provider", "Cloudflare WARP"},
		}
		return renderSheet(width, height, sheetLines(fields, innerWidth))
	case s.DaemonDown:
		return renderSheet(width, height, []string{
			"  " + dimStyle.Render("Daemon down"),
			"  " + dimStyle.Render("sudo systemctl enable --now warp-svc"),
		})
	case s.NeedsRegistration():
		return renderSheet(width, height, []string{
			"  " + warnStyle.Render("Login required"),
			"  " + dimStyle.Render("Run warp-cli registration new in a terminal."),
		})
	case s.Connecting():
		fields := []sheetField{{"Status", "Connecting"}}
		return renderSheet(width, height, sheetLines(fields, innerWidth))
	default:
		fields := []sheetField{
			{"Status", "Disconnected"},
			{"Provider", "Cloudflare WARP"},
		}
		return renderSheet(width, height, sheetLines(fields, innerWidth))
	}
}

func dash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "--"
	}
	return s
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func renderSheetField(label, value string) string {
	return "  " + labelStyle.Render(label) + " " + valueStyle.Render(dash(value))
}

func sheetLines(fields []sheetField, innerWidth int) []string {
	lines := []string{""}
	if innerWidth >= 48 && len(fields) > 1 {
		col := innerWidth / 2
		for i := 0; i < len(fields); i += 2 {
			left := renderSheetField(fields[i].Label, fields[i].Value)
			if i+1 < len(fields) {
				right := renderSheetField(fields[i+1].Label, fields[i+1].Value)
				pad := col - lipgloss.Width(left)
				if pad < 1 {
					pad = 1
				}
				lines = append(lines, left+strings.Repeat(" ", pad)+strings.TrimLeft(right, " "))
			} else {
				lines = append(lines, left)
			}
		}
		return lines
	}
	for _, f := range fields {
		lines = append(lines, renderSheetField(f.Label, f.Value))
	}
	return lines
}

type panelSize struct{ width, height int }

func panelInnerSize(width, height int) panelSize {
	return panelSize{width: max(8, width-4), height: max(3, height-2)}
}

func renderPanel(width, height int, lines []string) string {
	inner := panelInnerSize(width, height)
	padded := make([]string, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			padded = append(padded, "")
			continue
		}
		if strings.HasPrefix(line, "  ") {
			padded = append(padded, line)
		} else {
			padded = append(padded, "  "+line)
		}
	}
	for len(padded) < inner.height {
		padded = append(padded, "")
	}
	if len(padded) > inner.height {
		padded = padded[:inner.height]
	}
	return panelBorderStyle.
		Width(width - 2).
		Height(inner.height).
		Render(strings.Join(padded, "\n"))
}

func renderSheet(width, height int, lines []string) string {
	body := []string{dimStyle.Render("Details"), ""}
	body = append(body, lines...)
	return renderPanel(width, height, body)
}
