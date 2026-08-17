# README gif. Fake data only — always --demo. No window-bar stoplights.

Output demo.gif

Require ./omarchy-vpn

Set Shell bash
Set FontFamily "CaskaydiaMono Nerd Font"
Set FontSize 14
Set Width 1200
Set Height 720
Set Padding 20
Set Theme "Catppuccin Mocha"
Set TypingSpeed 0ms
Set Framerate 24

Hide
Type "./omarchy-vpn --demo"
Enter
Sleep 800ms
Show

# Land on connected homelab
Sleep 1.4s

# Details follow the cursor
Down
Sleep 1.2s
Down
Sleep 1.2s
Up
Sleep 1s

# Toggle office up, then down
Enter
Sleep 1.6s
Enter
Sleep 1.2s

# Help overlay
Type "?"
Sleep 2s
Type "j"
Sleep 800ms

Type "q"
Sleep 400ms
