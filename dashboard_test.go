package main

import (
	"strings"
	"testing"
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

func TestFormatIndexRow_HasKindAndStatus(t *testing.T) {
	got := formatIndexRow("homelab", "Connected", "WireGuard", 48)
	if !strings.Contains(got, "homelab") || !strings.Contains(got, "Connected") || !strings.Contains(got, "WireGuard") {
		t.Fatalf("row missing columns: %q", got)
	}
	if strings.Contains(got, "▸") || strings.Contains(got, "●") {
		t.Fatalf("old row chrome: %q", got)
	}
}

func TestTruncate_FitsWidth(t *testing.T) {
	got := truncate("Cloudflare WARP", 8)
	if got == "Cloudflare WARP" {
		t.Fatal("expected truncation")
	}
	if !strings.HasSuffix(got, "…") {
		t.Fatalf("expected ellipsis, got %q", got)
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
