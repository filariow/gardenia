#!/bin/sh

set -e

echo "Building valved"
go build \
    -trimpath \
    -ldflags="-s -w" \
    -o "bin/valved" \
    "cmd/valved/main.go"

echo "Building rosina"
go build \
    -trimpath \
    -ldflags="-s -w" \
    -o "bin/rosina" \
    "cmd/rosina/main.go"

echo "Building flowmeter"
go build \
    -trimpath \
    -ldflags="-s -w" \
    -o "bin/flowmeter" \
    "cmd/flowmeter/main.go"

echo "Installing valved into /usr/local/bin/valved"
sudo cp "/usr/local/bin/valved" "/usr/local/bin/valved.bkp"
sudo install "bin/valved" "/usr/local/bin/valved"

echo "Installing rosina into /usr/local/bin/rosina"
sudo cp "/usr/local/bin/rosina" "/usr/local/bin/rosina.bkp"
sudo install "bin/rosina" "/usr/local/bin/rosina"

echo "Installing flowmeter into /usr/local/bin/flowmeter"
sudo cp "/usr/local/bin/flowmeter" "/usr/local/bin/flowmeter.bkp"
sudo install "bin/flowmeter" "/usr/local/bin/flowmeter"

echo "Starting valved.service"
sudo install "deploy/linux/valved-remote.service" "/usr/lib/systemd/system/valved.service"
sudo systemctl daemon-reload
sudo systemctl enable --now "valved.service"
sudo systemctl status "valved.service"

echo "Starting rosina.service"
sudo install "deploy/linux/rosina-remote.service" "/usr/lib/systemd/system/rosina.service"
sudo systemctl daemon-reload
sudo systemctl enable --now "rosina.service"
sudo systemctl status "rosina.service"

echo "Starting flowmeter.service"
sudo install "deploy/linux/flowmeter-remote.service" "/usr/lib/systemd/system/flowmeter.service"
sudo systemctl daemon-reload
sudo systemctl enable --now "flowmeter.service"
sudo systemctl status "flowmeter.service"
