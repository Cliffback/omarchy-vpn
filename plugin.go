package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed plugin
var pluginFS embed.FS

const pluginID = "limehawk.vpn"

func enableVPNWidgetJSON(data []byte) ([]byte, error) {
	cfg, err := decodeObject(data)
	if err != nil {
		return nil, err
	}
	right := ensureLayoutSection(cfg, "right")
	if widgetIndex(right, pluginID) >= 0 {
		return encodeJSON(cfg)
	}
	entry := map[string]any{"id": pluginID}
	if i := widgetIndex(right, "omarchy.network"); i >= 0 {
		right = insertAt(right, i, entry)
	} else {
		right = append(right, entry)
	}
	setLayoutSection(cfg, "right", right)
	return encodeJSON(cfg)
}

func disableVPNWidgetJSON(data []byte) ([]byte, error) {
	cfg, err := decodeObject(data)
	if err != nil {
		return nil, err
	}
	for _, section := range []string{"left", "center", "right"} {
		sec, ok := layoutSection(cfg, section)
		if !ok {
			continue
		}
		kept := make([]any, 0, len(sec))
		for _, e := range sec {
			if widgetID(e) != pluginID {
				kept = append(kept, e)
			}
		}
		setLayoutSection(cfg, section, kept)
	}
	return encodeJSON(cfg)
}

func decodeObject(data []byte) (map[string]any, error) {
	var cfg map[string]any
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("shell.json: %w", err)
	}
	if cfg == nil {
		cfg = map[string]any{}
	}
	return cfg, nil
}

func encodeJSON(cfg map[string]any) ([]byte, error) {
	return json.MarshalIndent(cfg, "", "  ")
}

func ensureLayoutSection(cfg map[string]any, section string) []any {
	if sec, ok := layoutSection(cfg, section); ok {
		return sec
	}
	bar, _ := cfg["bar"].(map[string]any)
	if bar == nil {
		bar = map[string]any{}
		cfg["bar"] = bar
	}
	layout, _ := bar["layout"].(map[string]any)
	if layout == nil {
		layout = map[string]any{}
		bar["layout"] = layout
	}
	sec := []any{}
	layout[section] = sec
	return sec
}

func layoutSection(cfg map[string]any, section string) ([]any, bool) {
	bar, _ := cfg["bar"].(map[string]any)
	if bar == nil {
		return nil, false
	}
	layout, _ := bar["layout"].(map[string]any)
	if layout == nil {
		return nil, false
	}
	sec, ok := layout[section].([]any)
	return sec, ok
}

func setLayoutSection(cfg map[string]any, section string, entries []any) {
	ensureLayoutSection(cfg, section)
	bar := cfg["bar"].(map[string]any)
	layout := bar["layout"].(map[string]any)
	layout[section] = entries
}

func widgetID(entry any) string {
	switch v := entry.(type) {
	case map[string]any:
		id, _ := v["id"].(string)
		return id
	case string:
		return v
	default:
		return ""
	}
}

func widgetIndex(entries []any, id string) int {
	for i, e := range entries {
		if widgetID(e) == id {
			return i
		}
	}
	return -1
}

func insertAt(entries []any, i int, entry any) []any {
	out := make([]any, 0, len(entries)+1)
	out = append(out, entries[:i]...)
	out = append(out, entry)
	out = append(out, entries[i:]...)
	return out
}

func pluginDir(home string) string {
	return filepath.Join(home, ".config", "omarchy", "plugins", pluginID)
}

func writePlugin(dest string) error {
	if err := os.RemoveAll(dest); err != nil {
		return err
	}
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return err
	}
	entries, err := fs.ReadDir(pluginFS, "plugin")
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		data, err := pluginFS.ReadFile("plugin/" + e.Name())
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dest, e.Name()), data, 0o644); err != nil {
			return err
		}
	}
	return nil
}

func setupShell(home string) error {
	if err := writePlugin(pluginDir(home)); err != nil {
		return fmt.Errorf("plugin: %w", err)
	}
	if err := patchShellFile(home, enableVPNWidgetJSON); err != nil {
		return err
	}
	fmt.Println("Omarchy shell VPN widget installed.")
	afterDesktopChange()
	return nil
}

func removeShell(home string) error {
	if err := patchShellFile(home, disableVPNWidgetJSON); err != nil {
		return err
	}
	if err := os.RemoveAll(pluginDir(home)); err != nil {
		return err
	}
	fmt.Println("Omarchy shell VPN widget removed.")
	afterDesktopChange()
	return nil
}

func patchShellFile(home string, fn func([]byte) ([]byte, error)) error {
	path := filepath.Join(home, ".config", "omarchy", "shell.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	out, err := fn(data)
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(out, '\n'), 0o644)
}

func rescanPlugins() {
	if path, err := lookPath("omarchy-shell"); err == nil {
		exec.Command(path, "shell", "rescanPlugins").Run()
	}
}

// afterDesktopChange notifies the running shell. Tests replace it so they
// never talk to a live omarchy-shell.
var afterDesktopChange = rescanPlugins
