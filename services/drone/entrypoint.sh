#!/bin/sh
set -e

if [ "$DEV_MODE" = "true" ]; then
  echo "Running in development mode"
  cd /app
  go install github.com/air-verse/air@latest
  go mod download
  exec air
else
  echo "Running in production mode: executing the compiled binary"
  exec ./drone
fi