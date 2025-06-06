#!/bin/sh

set -e

echo "Building 'rosina/bot' image"
docker build \
    -f deploy/docker/bot/Dockerfile \
    -t "rosina/bot:latest" \
    .

echo "Exporting bot image"
docker save --output /tmp/bot-latest.tar rosina/bot:latest

echo "Importing bot image into k3s"
sudo k3s ctr images import /tmp/bot-latest.tar

echo "Appling manifests"
kubectl apply -f "manifests/bot/bot.yaml"
