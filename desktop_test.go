package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestDetectDesktop_Omarchy4_WhenShellAndConfigExist(t *testing.T) {
	home := t.TempDir()
	mustMkdirWrite(t, filepath.Join(home, ".config", "omarchy", "shell.json"), `{ "version": 1 }`)
	stubLookPath(t, map[string]string{"omarchy-shell": "/usr/bin/omarchy-shell"})

	if got := detectDesktop(home); got != desktopOmarchy4 {
		t.Fatalf("detectDesktop() = %v, want Omarchy 4", got)
	}
}

func TestDetectDesktop_Omarchy3_WhenWaybarConfigExists(t *testing.T) {
	home := t.TempDir()
	mustMkdirWrite(t, filepath.Join(home, ".config", "waybar", "config.jsonc"), `{ "modules-right": ["network"] }`)
	stubLookPath(t, nil)

	if got := detectDesktop(home); got != desktopOmarchy3 {
		t.Fatalf("detectDesktop() = %v, want Omarchy 3", got)
	}
}

func TestDetectDesktop_None_WhenNeitherDesktopExists(t *testing.T) {
	home := t.TempDir()
	stubLookPath(t, nil)

	if got := detectDesktop(home); got != desktopNone {
		t.Fatalf("detectDesktop() = %v, want none", got)
	}
}

func TestDetectDesktop_Prefers4_WhenWaybarLeftoversRemain(t *testing.T) {
	home := t.TempDir()
	mustMkdirWrite(t, filepath.Join(home, ".config", "omarchy", "shell.json"), `{ "version": 1 }`)
	mustMkdirWrite(t, filepath.Join(home, ".config", "waybar", "config.jsonc"), `{ "modules-right": ["custom/vpn", "network"] }`)
	stubLookPath(t, map[string]string{"omarchy-shell": "/usr/bin/omarchy-shell"})

	if got := detectDesktop(home); got != desktopOmarchy4 {
		t.Fatalf("detectDesktop() = %v, want Omarchy 4 when both are present", got)
	}
}

const sampleWaybarConfig = `{
    "modules-right": [
    "network"
    ]
}`

func TestSetupDesktop_Omarchy3_PatchesWaybar(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	stubLookPath(t, nil)
	mustMkdirWrite(t, filepath.Join(home, ".config", "waybar", "config.jsonc"), sampleWaybarConfig)
	mustMkdirWrite(t, filepath.Join(home, ".config", "waybar", "style.css"), sampleDefaultStyle)
	mustMkdirWrite(t, filepath.Join(home, ".config", "hypr", "hyprland.conf"), "source = monitors.conf\n")

	if err := setupDesktop(home); err != nil {
		t.Fatal(err)
	}

	cfg, err := os.ReadFile(filepath.Join(home, ".config", "waybar", "config.jsonc"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(cfg), `"custom/vpn"`) {
		t.Fatalf("waybar config was not patched:\n%s", cfg)
	}
}

func TestSetupDesktop_Omarchy4_InstallsPluginAndEnablesWidget(t *testing.T) {
	home := omarchy4Home(t)

	if err := setupDesktop(home); err != nil {
		t.Fatal(err)
	}

	manifest := filepath.Join(home, ".config", "omarchy", "plugins", pluginID, "manifest.json")
	widget := filepath.Join(home, ".config", "omarchy", "plugins", pluginID, "BarWidget.qml")
	if !fileExists(manifest) {
		t.Fatalf("missing plugin manifest at %s", manifest)
	}
	if !fileExists(widget) {
		t.Fatalf("missing plugin widget at %s", widget)
	}

	shell, err := os.ReadFile(filepath.Join(home, ".config", "omarchy", "shell.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(shell), pluginID) {
		t.Fatalf("shell.json was not enabled:\n%s", shell)
	}
}

func TestSetupDesktop_Omarchy4_CleansWaybarLeftovers(t *testing.T) {
	home := omarchy4Home(t)
	mustMkdirWrite(t, filepath.Join(home, ".config", "waybar", "config.jsonc"), `{
    "modules-right": [
    "custom/vpn",
    "network"
    ],
  "custom/vpn": {
    "exec": "omarchy-vpn --waybar"
  }
}`)
	mustMkdirWrite(t, filepath.Join(home, ".config", "waybar", "style.css"), sampleDefaultStyle+"\n#custom-vpn {\n  margin-right: 19px;\n}\n")
	mustMkdirWrite(t, filepath.Join(home, ".config", "hypr", "hyprland.conf"), "source = monitors.conf\n\n"+hyprlandWindowRule+"\n")

	if err := setupDesktop(home); err != nil {
		t.Fatal(err)
	}

	cfg, _ := os.ReadFile(filepath.Join(home, ".config", "waybar", "config.jsonc"))
	if strings.Contains(string(cfg), `"custom/vpn"`) {
		t.Fatalf("waybar leftover survived setup:\n%s", cfg)
	}
	hypr, _ := os.ReadFile(filepath.Join(home, ".config", "hypr", "hyprland.conf"))
	if strings.Contains(string(hypr), "org.omarchy.omarchy-vpn") {
		t.Fatalf("hyprland leftover survived setup:\n%s", hypr)
	}
	if !fileExists(filepath.Join(home, ".config", "omarchy", "plugins", pluginID, "manifest.json")) {
		t.Fatal("expected plugin to be installed on the 4 path")
	}
}

func TestSetupDesktop_None_IsNoop(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	stubLookPath(t, nil)

	if err := setupDesktop(home); err != nil {
		t.Fatal(err)
	}
}

func TestRemoveDesktop_Omarchy4_RemovesPlugin(t *testing.T) {
	home := omarchy4Home(t)
	if err := setupDesktop(home); err != nil {
		t.Fatal(err)
	}

	if err := removeDesktop(home); err != nil {
		t.Fatal(err)
	}

	if fileExists(filepath.Join(home, ".config", "omarchy", "plugins", pluginID, "manifest.json")) {
		t.Fatal("plugin files were left behind")
	}
	shell, err := os.ReadFile(filepath.Join(home, ".config", "omarchy", "shell.json"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(shell), pluginID) {
		t.Fatalf("widget still enabled:\n%s", shell)
	}
}

func TestRemoveDesktop_MissingFiles_NoError(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	stubLookPath(t, nil)

	if err := removeDesktop(home); err != nil {
		t.Fatal(err)
	}
}

func TestRemoveWaybar_MissingFiles_NoError(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	if err := removeWaybar(); err != nil {
		t.Fatal(err)
	}
}

func TestWritePlugin_PassesOmarchyValidate(t *testing.T) {
	path, err := exec.LookPath("omarchy-plugin-validate")
	if err != nil {
		t.Skip("omarchy-plugin-validate not installed")
	}
	dir := filepath.Join(t.TempDir(), pluginID)
	if err := writePlugin(dir); err != nil {
		t.Fatal(err)
	}
	out, err := exec.Command(path, dir).CombinedOutput()
	if err != nil {
		t.Fatalf("omarchy plugin validate: %v\n%s", err, out)
	}
}

func TestWritePlugin_WritesManifestAndWidget(t *testing.T) {
	dir := filepath.Join(t.TempDir(), pluginID)
	if err := writePlugin(dir); err != nil {
		t.Fatal(err)
	}
	manifest, err := os.ReadFile(filepath.Join(dir, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(manifest), `"id": "limehawk.vpn"`) {
		t.Fatalf("manifest missing plugin id:\n%s", manifest)
	}
	widget, err := os.ReadFile(filepath.Join(dir, "BarWidget.qml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(widget), "omarchy-vpn --waybar") && !strings.Contains(string(widget), `"omarchy-vpn", "--waybar"`) {
		t.Fatalf("widget does not poll --waybar:\n%s", widget)
	}
	if !strings.Contains(string(widget), "TUI.float") {
		t.Fatalf("widget does not launch as TUI.float:\n%s", widget)
	}
	if strings.Contains(string(widget), "󰳌") || strings.Contains(string(widget), "󰦝") {
		t.Fatal("bar widget still uses shield-lock glyphs; at bar size they read as a person")
	}
	if !strings.Contains(string(widget), "0xF0498") || !strings.Contains(string(widget), "0xF0565") {
		t.Fatal("expected md-shield / md-shield_check (0xF0498 / 0xF0565)")
	}
}

func omarchy4Home(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	mustMkdirWrite(t, filepath.Join(home, ".config", "omarchy", "shell.json"), sampleShellJSON)
	stubLookPath(t, map[string]string{"omarchy-shell": "/usr/bin/omarchy-shell"})
	orig := afterDesktopChange
	afterDesktopChange = func() {}
	t.Cleanup(func() { afterDesktopChange = orig })
	return home
}

func stubLookPath(t *testing.T, found map[string]string) {
	t.Helper()
	orig := lookPath
	lookPath = func(file string) (string, error) {
		if path, ok := found[file]; ok {
			return path, nil
		}
		return "", exec.ErrNotFound
	}
	t.Cleanup(func() { lookPath = orig })
}

func mustMkdirWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
