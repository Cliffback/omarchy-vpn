package main

import (
	"slices"
	"testing"
)

func TestDemoListIsCanned(t *testing.T) {
	enableDemoMode()
	got := ListConfigs()
	want := []string{"homelab", "office", "travel"}
	if !slices.Equal(got, want) {
		t.Fatalf("ListConfigs() = %v, want canned demo names %v", got, want)
	}
}

func TestDemoDoesNotEnableDaemons(t *testing.T) {
	enableDemoMode()
	if WarpAvailable() || NetBirdAvailable() {
		t.Fatal("demo mode must not probe live WARP or NetBird")
	}
}

func TestDemoConnectDisconnect(t *testing.T) {
	enableDemoMode()
	demoMu.Lock()
	demoActive = []string{"homelab"}
	demoMu.Unlock()
	if err := ConnectVPN("office"); err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(GetActiveVPNs(), "office") {
		t.Fatal("demo connect should mark office active")
	}
	if err := DisconnectVPN("office"); err != nil {
		t.Fatal(err)
	}
	if slices.Contains(GetActiveVPNs(), "office") {
		t.Fatal("demo disconnect should clear office")
	}
}
