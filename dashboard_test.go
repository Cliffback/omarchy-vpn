package main

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestRenderTitleBar_Disconnected(t *testing.T) {
	initColors()
	m := model{width: 80}
	got := m.renderTitleBar()
	if !strings.Contains(got, "omarchy-vpn") {
		t.Fatalf("missing name:\n%s", got)
	}
	if !strings.Contains(got, "Disconnected") {
		t.Fatalf("missing status word:\n%s", got)
	}
	if strings.Contains(got, "●") || strings.Contains(got, "○") || strings.ContainsAny(got, "󰳌󰦝") {
		t.Fatalf("old title chrome survived:\n%s", got)
	}
}

func TestRenderTitleBar_ConnectedListsNames(t *testing.T) {
	initColors()
	m := model{width: 80, activeVPNs: []string{"homelab"}}
	got := m.renderTitleBar()
	if !strings.Contains(got, "Connected") {
		t.Fatalf("missing Connected:\n%s", got)
	}
	if !strings.Contains(got, "homelab") {
		t.Fatalf("missing tunnel name:\n%s", got)
	}
	if strings.Contains(got, "●") {
		t.Fatalf("dot badge survived:\n%s", got)
	}
}

func TestTitledBox_ExactSize(t *testing.T) {
	initColors()
	got := titledBox("Tunnels", []string{"hello"}, 36, 12)
	lines := strings.Split(got, "\n")
	if len(lines) != 12 {
		t.Fatalf("height %d want 12\n%s", len(lines), got)
	}
	for i, line := range lines {
		if w := lipgloss.Width(line); w != 36 {
			t.Fatalf("line %d width %d want 36: %q", i, w, line)
		}
	}
	if !strings.Contains(got, "Tunnels") {
		t.Fatalf("missing title:\n%s", got)
	}
}

func TestConfigPanel_GlowRows(t *testing.T) {
	initColors()
	m := model{warpAvail: true, configs: []string{"homelab"}}
	got := m.renderConfigPanel(36, 16)
	if !strings.Contains(got, "Cloudflare WARP") || !strings.Contains(got, "WARP · Disconnected") {
		t.Fatalf("missing glow row:\n%s", got)
	}
	if !strings.Contains(got, "WireGuard · Disconnected") {
		t.Fatalf("missing wireguard meta:\n%s", got)
	}
	for i, line := range strings.Split(got, "\n") {
		if w := lipgloss.Width(line); w != 36 {
			t.Fatalf("line %d width %d want 36: %q", i, w, line)
		}
	}
}

func TestTitledBox_InnerPad(t *testing.T) {
	initColors()
	got := titledBox("Tunnels", []string{"hello"}, 36, 12)
	lines := strings.Split(got, "\n")
	// title in the top edge, then one blank row, then content
	if !strings.Contains(lines[0], "Tunnels") {
		t.Fatalf("title not in top edge:\n%s", got)
	}
	if strings.Contains(lines[1], "hello") {
		t.Fatalf("content glued to top border:\n%s", got)
	}
	if !strings.Contains(got, "hello") {
		t.Fatalf("missing content:\n%s", got)
	}
}

func TestConfigPanel_OneBlankBetweenItems(t *testing.T) {
	initColors()
	m := model{warpAvail: true, configs: []string{"homelab"}}
	got := m.renderConfigPanel(40, 16)
	idxWARP := strings.Index(got, "Cloudflare WARP")
	idxMeta := strings.Index(got, "WARP · Disconnected")
	idxNext := strings.Index(got, "homelab")
	if idxWARP < 0 || idxMeta < 0 || idxNext < 0 || !(idxWARP < idxMeta && idxMeta < idxNext) {
		t.Fatalf("unexpected order:\n%s", got)
	}
	between := got[idxMeta+len("WARP · Disconnected") : idxNext]
	// one visual blank line, not two
	if strings.Count(between, "\n") != 2 {
		t.Fatalf("want one blank between items, got %d newlines in %q\n%s", strings.Count(between, "\n"), between, got)
	}
}

func TestUseStackedLayout(t *testing.T) {
	cases := []struct {
		w, h  int
		stack bool
	}{
		{80, 24, false},
		{70, 20, false},
		{100, 30, false},
		{40, 80, true},
		{50, 100, true},
		{60, 24, true},
		{70, 90, true},
		{80, 100, true},
	}
	for _, tc := range cases {
		if got := useStackedLayout(tc.w, tc.h); got != tc.stack {
			t.Errorf("%dx%d stacked=%v want %v", tc.w, tc.h, got, tc.stack)
		}
	}
}

func titleLine(view, name string) int {
	for i, line := range strings.Split(view, "\n") {
		if strings.Contains(line, name) {
			return i
		}
	}
	return -1
}

func TestDashboard_ColumnsOnWide(t *testing.T) {
	initColors()
	h := newHelp()
	h.SetWidth(78)
	m := model{
		width:     80,
		height:    24,
		warpAvail: true,
		configs:   []string{"homelab"},
		keys:      newKeyMap(),
		help:      h,
	}
	got := m.View().Content
	tun, det := titleLine(got, "Tunnels"), titleLine(got, "Details")
	if tun < 0 || det < 0 {
		t.Fatalf("missing pane titles:\n%s", got)
	}
	if tun != det {
		t.Fatalf("wide layout should be columns, titles on lines %d and %d\n%s", tun, det, got)
	}
}

func TestDashboard_StacksOnPhone(t *testing.T) {
	initColors()
	h := newHelp()
	h.SetWidth(38)
	m := model{
		width:      40,
		height:     80,
		warpAvail:  true,
		warpStatus: WarpStatus{DaemonDown: true},
		configs:    []string{"homelab"},
		keys:       newKeyMap(),
		help:       h,
	}
	got := m.View().Content
	tun, det := titleLine(got, "Tunnels"), titleLine(got, "Details")
	if tun < 0 || det < 0 {
		t.Fatalf("missing pane titles:\n%s", got)
	}
	if det <= tun {
		t.Fatalf("phone layout should stack Details below Tunnels, titles on lines %d and %d\n%s", tun, det, got)
	}
	if !strings.Contains(got, "WireGuard · Disconnected") {
		t.Fatalf("meta truncated on phone:\n%s", got)
	}
	if !strings.Contains(got, "Daemon down") {
		t.Fatalf("details missing on phone:\n%s", got)
	}
	if !strings.Contains(got, "import") || !strings.Contains(got, "quit") {
		t.Fatalf("footer truncated on phone:\n%s", got)
	}
}

func TestDashboard_MetaNotTruncatedAt80(t *testing.T) {
	initColors()
	h := newHelp()
	h.SetWidth(78)
	m := model{
		width:     80,
		height:    24,
		warpAvail: true,
		configs:   []string{"homelab"},
		keys:      newKeyMap(),
		help:      h,
	}
	got := m.View().Content
	if !strings.Contains(got, "WireGuard · Disconnected") {
		t.Fatalf("meta truncated at 80 cols:\n%s", got)
	}
	if strings.Contains(got, "Disconnec…") {
		t.Fatalf("ellipsis on status word:\n%s", got)
	}
}

func TestSheetField_EmptyIsDash(t *testing.T) {
	initColors()
	got := renderSheetField("DNS", "")
	if !strings.Contains(got, "DNS") || !strings.Contains(got, "--") {
		t.Fatalf("expected dashed empty: %q", got)
	}
}

func TestSheetField_NoNerdIcon(t *testing.T) {
	initColors()
	got := renderSheetField("Endpoint", "1.2.3.4:51820")
	if strings.ContainsAny(got, "󰖟󰩟󰇖") {
		t.Fatalf("nerd icon on sheet: %q", got)
	}
}
