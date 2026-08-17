import QtQuick
import Quickshell
import Quickshell.Io
import qs.Commons
import qs.Ui

BarWidget {
  id: root
  moduleName: "limehawk.vpn"

  property string glyph: "󰳌"
  property string statusText: "VPN: Disconnected"
  property bool connected: false

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
        root.glyph = data.text || root.glyph
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
    slotSize: Style.bar.statusSlot
    fontSize: Style.font.caption
    tooltipText: root.statusText
    dimmed: !root.connected
    onPressed: root.openTui()
  }
}
