import QtQuick
import QtQuick.Shapes
import qs.Commons

// Drawn mark — Omarchy 4's bar font does not carry the Material shield
// glyphs we use in the TUI / Waybar output, and they fall back to a person.
Item {
  id: root

  property real iconSize: Style.font.icon
  property color color: Color.foreground
  property bool connected: false

  width: iconSize
  height: iconSize
  implicitWidth: iconSize
  implicitHeight: iconSize

  readonly property real inset: iconSize * 0.10
  readonly property real stroke: Math.max(1.2, iconSize * 0.12)
  readonly property real left: inset + (width - inset * 2) * 0.16
  readonly property real right: inset + (width - inset * 2) * 0.84
  readonly property real top: inset + (height - inset * 2) * 0.06
  readonly property real midY: inset + (height - inset * 2) * 0.48
  readonly property real tipY: inset + (height - inset * 2) * 0.92
  readonly property real midX: width / 2

  Shape {
    anchors.fill: parent
    antialiasing: true
    layer.enabled: true
    layer.samples: 4

    ShapePath {
      // Outline when down, filled when up — readable at bar size without
      // depending on a nerd-font glyph the Quattro bar font does not have.
      fillColor: root.connected ? root.color : "transparent"
      strokeColor: root.color
      strokeWidth: root.connected ? 0 : root.stroke
      joinStyle: ShapePath.RoundJoin
      startX: root.left
      startY: root.top
      PathLine { x: root.right; y: root.top }
      PathLine { x: root.right; y: root.midY }
      PathLine { x: root.midX; y: root.tipY }
      PathLine { x: root.left; y: root.midY }
      PathLine { x: root.left; y: root.top }
    }
  }
}
