package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type desktopKind int

const (
	desktopNone desktopKind = iota
	desktopOmarchy3
	desktopOmarchy4
)

func (k desktopKind) String() string {
	switch k {
	case desktopOmarchy3:
		return "omarchy-3"
	case desktopOmarchy4:
		return "omarchy-4"
	default:
		return "none"
	}
}

// lookPath is exec.LookPath, overridable in tests.
var lookPath = exec.LookPath

func detectDesktop(home string) desktopKind {
	if isOmarchy4(home) {
		return desktopOmarchy4
	}
	if fileExists(filepath.Join(home, ".config", "waybar", "config.jsonc")) {
		return desktopOmarchy3
	}
	return desktopNone
}

func isOmarchy4(home string) bool {
	if _, err := lookPath("omarchy-shell"); err != nil {
		return false
	}
	return fileExists(filepath.Join(home, ".config", "omarchy", "shell.json"))
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func setupDesktop(home string) error {
	switch detectDesktop(home) {
	case desktopOmarchy4:
		if err := unpatchWaybarFiles(); err != nil {
			return err
		}
		return setupShell(home)
	case desktopOmarchy3:
		return setupWaybar()
	default:
		return nil
	}
}

func removeDesktop(home string) error {
	if err := removeShell(home); err != nil {
		return err
	}
	return unpatchWaybarFiles()
}

func runDesktopFlag(flag string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	switch flag {
	case "--setup":
		return setupDesktop(home)
	case "--remove":
		return removeDesktop(home)
	case "--setup-waybar":
		return setupWaybar()
	case "--remove-waybar":
		return removeWaybar()
	default:
		return fmt.Errorf("unknown flag %s", flag)
	}
}
