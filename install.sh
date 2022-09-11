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
sudo install "deploy/linux/valved.service" "/usr/lib/systemd/system/valved.service"
sudo systemctl daemon-reload
sudo systemctl enable --now "valved.service"

echo "Building 'skeduler' image"
docker build \
    -f deploy/docker/skeduler/Dockerfile \
    -t "rosina/skeduler:latest" \
    .

echo "Building 'rosina' image"
docker build \
    -f deploy/docker/rosina/Dockerfile \
    -t "rosina/rosina:latest" \
    .

echo "Appling manifests"
kubectl apply -f "manifests/skeduler.yaml"
