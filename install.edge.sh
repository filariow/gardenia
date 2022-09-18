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

echo "Installing valved into /usr/local/valved"
sudo install "bin/valved" "/usr/local/valved"

echo "Installing rosina into /usr/local/rosina"
sudo install "bin/rosina" "/usr/local/rosina"

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

