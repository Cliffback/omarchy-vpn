package main

import (
	"encoding/json"
	"os"
	"regexp"
	"testing"
)

// PKGBUILD pkgver is the only version we ship. The plugin manifest must
// match it so the bar widget and the TUI don't drift.
func TestShippedVersionSourcesMatch(t *testing.T) {
	pkg, err := os.ReadFile("PKGBUILD")
	if err != nil {
		t.Fatal(err)
	}
	m := regexp.MustCompile(`(?m)^pkgver=(.+)$`).FindSubmatch(pkg)
	if m == nil {
		t.Fatal("PKGBUILD has no pkgver=")
	}
	pkgver := string(m[1])
	if pkgver == "" || pkgver == "dev" {
		t.Fatalf("pkgver %q is not a release", pkgver)
	}

	raw, err := os.ReadFile("plugin/manifest.json")
	if err != nil {
		t.Fatal(err)
	}
	var man struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(raw, &man); err != nil {
		t.Fatal(err)
	}
	if man.Version != pkgver {
		t.Fatalf("plugin version %q != PKGBUILD pkgver %q", man.Version, pkgver)
	}
}
