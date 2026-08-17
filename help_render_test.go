package main

import (
	"strings"
	"testing"
)

func TestShortHelpContents(t *testing.T) {
	initColors()
	initStyles()
	h := newHelp()
	h.SetWidth(200)
	view := h.ShortHelpView(newKeyMap().ShortHelp())
	if !strings.Contains(view, "?") {
		t.Errorf("short help missing %q in %q", "?", view)
	}
	for _, stale := range []string{"delete config", "rename config", "enter", "│"} {
		if strings.Contains(view, stale) {
			t.Errorf("short help still contains %q: %s", stale, view)
		}
	}
}
