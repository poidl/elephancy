#!/bin/bash
export GOROOT=$HOME/programs/goroot
export PATH=$PATH:$GOROOT/bin
if hoit=$(pgrep mypage)
then
  echo "killing process mypage, removing executable"
  kill $hoit;
  rm mypage
fi
sleep 0.2
if go build
then
  echo "starting mypage server";
  ./chromerefresh.sh &
  ./mypage;
  sleep 0.2
  # WID=$(xdotool search --onlyvisible --class chromium|head -1)
  # xdotool windowactivate ${WID}
  # xdotool key ctrl+1
  # xdotool key ctrl+F5
  exit 0
else
  exit 1
fi
