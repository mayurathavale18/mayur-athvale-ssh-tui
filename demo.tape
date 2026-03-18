Output docs/demo.gif

Set Shell bash
Set FontSize 14
Set Width 1200
Set Height 600
Set Padding 20
Set Theme "Catppuccin Mocha"
Set TypingSpeed 75ms

# Start the server locally in background
Hide
Type "SSH_PORT=2223 go run ./cmd/server/ &"
Enter
Sleep 3s
Show

# Connect locally
Type "ssh -t -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -p 2223 localhost"
Enter
Sleep 3s

# About tab with ASCII art
Sleep 3s

# Experience tab
Type "2"
Sleep 3s

# Scroll
Type "j"
Sleep 500ms
Type "j"
Sleep 500ms
Type "j"
Sleep 2s

# Projects
Type "3"
Sleep 3s

# Skills
Type "4"
Sleep 3s

# Contact
Type "5"
Sleep 3s

# Stats
Type "6"
Sleep 3s

# Back to About
Type "1"
Sleep 2s

# Quit
Type "q"
Sleep 1s

Hide
Type "kill %1 2>/dev/null"
Enter
