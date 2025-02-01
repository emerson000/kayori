#!/bin/sh
set -e

if [ "$DEV_MODE" = "true" ]; then
  echo "Running in development mode: executing 'go run'"
  cd /app
  go install github.com/air-verse/air@latest
  go mod download
  exec air
else
  echo "Running in production mode: executing the compiled binary"
  exec ./server
fi