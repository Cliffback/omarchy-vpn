package main

import (
	"encoding/json"
	"strings"
	"testing"
)

const sampleShellJSON = `{
  "bar": {
    "layout": {
      "left": [{"id": "omarchy.menu"}],
      "center": [{"id": "omarchy.clock"}],
      "right": [
        {"id": "omarchy.tray"},
        {"id": "omarchy.network"},
        {"id": "omarchy.audio"}
      ]
    }
  },
  "plugins": [],
  "version": 1
}`

func TestEnableVPNWidget_InsertsBeforeNetwork(t *testing.T) {
	out, err := enableVPNWidgetJSON([]byte(sampleShellJSON))
	if err != nil {
		t.Fatal(err)
	}

	ids := rightIDs(t, out)
	if !containsID(ids, pluginID) {
		t.Fatalf("widget not on the bar: %v", ids)
	}
	if i, j := indexOf(ids, pluginID), indexOf(ids, "omarchy.network"); j >= 0 && i > j {
		t.Fatalf("wanted %s before omarchy.network, got %v", pluginID, ids)
	}
	if i, j := indexOf(ids, pluginID), indexOf(ids, "omarchy.tray"); j >= 0 && i < j {
		t.Fatalf("wanted %s after omarchy.tray, got %v", pluginID, ids)
	}
}

func TestEnableVPNWidget_Idempotent(t *testing.T) {
	once, err := enableVPNWidgetJSON([]byte(sampleShellJSON))
	if err != nil {
		t.Fatal(err)
	}
	twice, err := enableVPNWidgetJSON(once)
	if err != nil {
		t.Fatal(err)
	}

	if got := strings.Count(string(twice), pluginID); got != 1 {
		t.Fatalf("expected one %s entry, got %d\n%s", pluginID, got, twice)
	}
}

func TestEnableVPNWidget_AppendsWhenNoNetwork(t *testing.T) {
	in := `{
	  "bar": {"layout": {"right": [{"id": "omarchy.tray"}]}},
	  "version": 1
	}`
	out, err := enableVPNWidgetJSON([]byte(in))
	if err != nil {
		t.Fatal(err)
	}
	ids := rightIDs(t, out)
	if got := ids[len(ids)-1]; got != pluginID {
		t.Fatalf("expected %s last on the right, got %v", pluginID, ids)
	}
}

func TestDisableVPNWidget_RemovesFromAllSections(t *testing.T) {
	enabled, err := enableVPNWidgetJSON([]byte(sampleShellJSON))
	if err != nil {
		t.Fatal(err)
	}
	out, err := disableVPNWidgetJSON(enabled)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), pluginID) {
		t.Fatalf("widget still present:\n%s", out)
	}
	ids := rightIDs(t, out)
	if !containsID(ids, "omarchy.network") || !containsID(ids, "omarchy.tray") {
		t.Fatalf("removed neighbouring widgets: %v", ids)
	}
}

func TestDisableVPNWidget_MissingIsNoop(t *testing.T) {
	out, err := disableVPNWidgetJSON([]byte(sampleShellJSON))
	if err != nil {
		t.Fatal(err)
	}
	ids := rightIDs(t, out)
	if containsID(ids, pluginID) {
		t.Fatalf("unexpected widget: %v", ids)
	}
}

func rightIDs(t *testing.T, data []byte) []string {
	t.Helper()
	var cfg struct {
		Bar struct {
			Layout struct {
				Right []struct {
					ID string `json:"id"`
				} `json:"right"`
			} `json:"layout"`
		} `json:"bar"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatal(err)
	}
	ids := make([]string, 0, len(cfg.Bar.Layout.Right))
	for _, e := range cfg.Bar.Layout.Right {
		ids = append(ids, e.ID)
	}
	return ids
}

func containsID(ids []string, want string) bool {
	return indexOf(ids, want) >= 0
}

func indexOf(ids []string, want string) int {
	for i, id := range ids {
		if id == want {
			return i
		}
	}
	return -1
}
