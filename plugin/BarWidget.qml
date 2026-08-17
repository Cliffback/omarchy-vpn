import QtQuick
import Quickshell
import Quickshell.Io
import qs.Commons
import qs.Ui

BarWidget {
  id: root
  moduleName: "limehawk.vpn"

  property string statusText: "VPN: Disconnected"
  property bool connected: false

  // md-shield / md-shield_check — same Material set as bluetooth, volume,
  // and display. The TUI's shield-lock glyphs (U+F0CCC / U+F099D) collapse
  // to a person silhouette at bar size.
  readonly property string glyph: String.fromCodePoint(root.connected ? 0xF0565 : 0xF0498)

  implicitWidth: button.implicitWidth
  implicitHeight: button.implicitHeight

  function refresh() {
    if (!statusProc.running) statusProc.running = true
  }

  function openTui() {
    if (root.bar) root.bar.run("omarchy-launch-or-focus-tui --app-id=TUI.float omarchy-vpn")
  }

  Process {
    id: statusProc
    command: ["omarchy-vpn", "--waybar"]
    stdout: StdioCollector {
      id: statusOut
      waitForEnd: true
    }
    onExited: function(exitCode) {
      if (exitCode !== 0) return
      try {
        var data = JSON.parse(String(statusOut.text || "").trim())
        root.statusText = data.tooltip || ""
        root.connected = data.class === "connected"
      } catch (e) {}
    }
  }

  Timer {
    interval: 5000
    running: true
    repeat: true
    triggeredOnStart: true
    onTriggered: root.refresh()
  }

  BarIconButton {
    id: button
    anchors.fill: parent
    bar: root.bar
    text: root.glyph
    tooltipText: root.statusText
    dimmed: !root.connected
    onPressed: root.openTui()
  }
}
