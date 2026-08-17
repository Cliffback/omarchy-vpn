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

  readonly property color iconColor: {
    var fg = bar ? bar.barForeground : Color.foreground
    return root.connected ? fg : Qt.darker(fg, 1.55)
  }

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
    slotSize: Style.bar.statusSlot
    tooltipText: root.statusText
    iconComponent: Component {
      Item {
        VpnIcon {
          anchors.centerIn: parent
          iconSize: Style.space(11)
          color: root.iconColor
          connected: root.connected
        }
      }
    }
    onPressed: root.openTui()
  }
}
