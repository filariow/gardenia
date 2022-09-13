#!/bin/sh

set -e

echo "Building valved"
go build \
    -trimpath \
    -ldflags="-s -w" \
    -o "bin/valved" \
    "cmd/valved/main.go"

echo "Installing valved into /usr/local/valved"
sudo install "bin/valved" "/usr/local/valved"

echo "Starting valved.service"
sudo install "deploy/linux/valved-remote.service" "/usr/lib/systemd/system/valved.service"
sudo systemctl daemon-reload
sudo systemctl enable --now "valved.service"
sudo systemctl status "valved.service"

