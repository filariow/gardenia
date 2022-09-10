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
sudo systemctl start "deploy/linux/valved.service"

echo "Appling manifests"
kubectl apply -f "manifests"

