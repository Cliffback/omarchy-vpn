# TUI index + sheet — Design

**Date:** 2026-08-17
**Status:** Building

## Goal

Restyle the TUI as a Charm app: rounded panels, a list you scan, a sheet
of facts. Same keys, same backends, same single screen.

## Constraints

- One screen. Left pane is the index; right pane is the sheet.
- Charm only (Bubble Tea / Lip Gloss / Bubbles v2). Rounded borders.
- ANSI colors via `theme.go`. No hex. Accent (ANSI 4) for selection.
  Green is for `Connected` only.
- Do not touch WireGuard / NetBird / WARP backends, waybar JSON, or the
  shell plugin.

## Surfaces

| Surface | Job |
|---------|-----|
| Title bar | Name + one status word |
| Left pane | Find a tunnel |
| Right pane | Facts about the selected row |
| `?` overlay | Shortcuts |

## Header

```
omarchy-vpn
Connected    homelab  NetBird          v0.4.2
```

- Line 1: name only. Not accent. No shield glyph.
- Line 2: `Connected` / `Disconnected`, then active names in mono,
  version at the end, dim.
- No keys in the header.

## Index

Charm list item, two lines, never wraps:

```
Cloudflare WARP
WARP · Disconnected

homelab
WireGuard · Connected
```

- Selection is accent + bold on the name.
- Connected meta line is green.
- Empty: `No configs.`
- Rename / delete stay inline on the selected item.
- Pane title: `Tunnels`. One rounded border.

## Sheet

Label → value. Missing values are `--`. Pane title: `Details`.

WireGuard always shows Status, Address, DNS, Endpoint, Peer, Path.
When up, also Download, Upload, Handshake.

NetBird / WARP use the same label/value shape. Daemon-down and
needs-login are a status word plus one quiet next-step line.

No nerd-font field icons.

## Controls

- Enter connects (the one primary).
- Disconnect is quiet (`d`).
- Footer is a flash, or a single dim `?`. Full bindings live on `?`.

## Files

`dashboard.go`, `config_panel.go`, `status_panel.go`, `styles.go`,
`help.go`, `help_render_test.go`, `dashboard_test.go`.
