#!/bin/bash

if ! dpkg -l | grep -q "golang" || ! go version | grep -q "linux/amd64"; then
    printf "Error! Golang is not installed. Installing from longsleep repository...\n"
    sudo add-apt-repository ppa:longsleep/golang-backports
    sudo apt update
    sudo apt install golang-go
fi

printf "Building project...\n"
sudo mkdir -p /opt/monitor_api
sudo mkdir -p /opt/monitor_api/pub
sudo cp ./index.html /opt/monitor_api/pub/index.html
go mod tidy
go build -o /opt/monitor_api/monitor_api ./
sudo touch /opt/monitor_api/monitors.txt
sudo echo '123,"LG 27GL850"' > /opt/monitor_api/monitors.txt
sudo echo '456,"Dell U2720Q"' >> /opt/monitor_api/monitors.txt
sudo echo '789,"Samsung LU28R550"' >> /opt/monitor_api/monitors.txt

sudo cp ./goapi.sh /opt/monitor_api/goapi.sh || exit
sudo cp ./goapiprogram.service /etc/systemd/system/goapiprogram.service || exit
sudo systemctl daemon-reload
sudo systemctl enable goapiprogram.service
printf "Created systemd service!\n\n"

sudo chmod +x /opt/monitor_api/monitor_api
cd /opt/monitor_api || exit
sudo ./monitor_api --createdb
if pgrep -x "monitor_api" > /dev/null; then
    printf "monitor_api is already running!\n"
else
    sudo ./monitor_api --start 2>&1 &
    printf "API started!\n"
fi

printf "Service installed and running! Reboot system for complete installation.\n"
