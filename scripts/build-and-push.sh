#!/usr/bin/env bash
# You’d likely parameterize REGISTRY or ORG for your Docker registry
REGISTRY="my-docker-registry"
TAG="latest"

# Build images
docker build -t $REGISTRY/frontend:$TAG ./services/frontend
docker build -t $REGISTRY/backend-api:$TAG ./services/backend-api
docker build -t $REGISTRY/drones:$TAG ./services/drones
docker build -t $REGISTRY/drone-controller:$TAG ./services/drone-controller
docker build -t $REGISTRY/streaming-service:$TAG ./services/streaming-service
docker build -t $REGISTRY/database:$TAG ./services/database

# Push images
docker push $REGISTRY/frontend:$TAG
docker push $REGISTRY/backend-api:$TAG
docker push $REGISTRY/drones:$TAG
docker push $REGISTRY/drone-controller:$TAG
docker push $REGISTRY/streaming-service:$TAG
docker push $REGISTRY/database:$TAG
