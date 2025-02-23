#!/usr/bin/env bash
# You’d likely parameterize REGISTRY or ORG for your Docker registry
REGISTRY="kayori"
TAG="latest"

# Build images
docker build -t $REGISTRY/frontend:$TAG ./services/frontend
docker build -t $REGISTRY/backend:$TAG ./services/backend
docker build -t $REGISTRY/drone:$TAG ./services/drone

# Push images
docker push $REGISTRY/frontend:$TAG
docker push $REGISTRY/backend:$TAG
docker push $REGISTRY/drone:$TAG
