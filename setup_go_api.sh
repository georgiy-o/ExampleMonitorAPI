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
sudo chmod +x /opt/monitor_api/monitor_api
cd /opt/monitor_api || exit
sudo ./monitor_api --createdb
if sudo ./monitor_api --start 2>&1 | grep -q "Failed to start server"; then
    printf "Cannot start API! Check if port 8030 is open!\n"
else
    printf "API started!\n"
fi
sudo cp ./goapi.sh /opt/monitor_api/goapi.sh || exit
sudo cp ./goapiprogram.service /etc/systemd/system/goapiprogram.service || exit
sudo systemctl daemon-reload
printf "Created systemd service!\n\n"
printf "Service installed and running! Reboot system for complete installation.\n"
