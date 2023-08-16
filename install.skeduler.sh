#!/bin/sh

set -e

echo "Building 'skeduler' image"
docker build \
    -f deploy/docker/skeduler/Dockerfile \
    -t "rosina/skeduler:latest" \
    .

echo "Building 'rosinacli' image"
docker build \
    -f deploy/docker/rosinacli/Dockerfile \
    -t "rosina/rosinacli:latest" \
    .

echo "Building 'rosina/bot' image"
docker build \
    -f deploy/docker/bot/Dockerfile \
    -t "rosina/bot:latest" \
    .

echo "Exporting skeduler image"
docker save --output /tmp/skeduler-latest.tar rosina/skeduler:latest

echo "Exporting rosina image"
docker save --output /tmp/rosinacli-latest.tar rosina/rosinacli:latest

echo "Exporting bot image"
docker save --output /tmp/bot-latest.tar rosina/bot:latest

echo "Importing skeduler image into k3s"
sudo k3s ctr images import /tmp/skeduler-latest.tar

echo "Importing rosina image into k3s"
sudo k3s ctr images import /tmp/rosinacli-latest.tar

echo "Importing bot image into k3s"
sudo k3s ctr images import /tmp/bot-latest.tar

echo "Appling manifests"
kubectl apply -f "manifests/skeduler.yaml" -f "manifests/bot.yaml"
