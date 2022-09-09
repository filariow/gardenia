#!/bin/sh

echo "Building valvedmock"
go build \
    -trimpath \
    -ldflags="-s -w" \
    -o "bin/valvedmock" \
    "cmd/valvedmock/main.go"

echo "Installing valvedmock into /usr/local/valved"
sudo install "bin/valvedmock" "/usr/local/valved"

echo "Starting valved.service"
sudo install "deploy/linux/valved.service" "/usr/lib/systemd/system/valved.service"
sudo systemctl daemon-reload
sudo systemctl enable --now "valved.service"

# echo "Appling manifests"
# kubectl apply -f "manifests"

