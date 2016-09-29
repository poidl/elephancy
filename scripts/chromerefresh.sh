#!/bin/bash
sleep 1
WID=$(xdotool search --onlyvisible --class chromium|head -1)
xdotool windowactivate ${WID}
xdotool key ctrl+2
xdotool key ctrl+F5
