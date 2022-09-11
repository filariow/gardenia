#!/bin/sh

set -e

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
