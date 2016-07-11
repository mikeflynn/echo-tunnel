#!/usr/bin/osascript

on run (args)
  display dialog (first item of args) with title (second item of args) with icon file (third item of args) buttons {(fourth item of args), (fifth item of args)} giving up after 30
end run