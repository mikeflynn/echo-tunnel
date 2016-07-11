#!/usr/bin/osascript

on run (args)
  tell application "System Preferences"
    activate
    reveal anchor "displaysDisplayTab" of pane id "com.apple.preference.displays"
    tell application "System Events"
      delay 1
      set value of slider 1 of group 1 of tab group 1 of window "Built-in Retina Display" of process "System Preferences" to (first item of args as number)

    end tell
    quit
  end tell
end run