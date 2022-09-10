#!/bin/sh

set -e

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

echo "Ensuring kind cluster 'rosina-cluster' is present"
kind get clusters | grep '^rosina-cluster$' || kind create cluster --config "manifests/kind/config.yaml"

echo "Loading skeduler image into kind cluster"
kind load docker-image --name rosina-cluster "rosina/skeduler:latest"

echo "Loading 'rosina' image into kind cluster"
kind load docker-image --name rosina-cluster "rosina/rosina:latest"

echo "Appling manifests"
kubectl apply -f "manifests/skeduler.yaml"

