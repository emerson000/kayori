#!/bin/sh
set -e

if [ "$DEV_MODE" = "true" ]; then
  echo "Running in development mode"
  exec npm run dev
else
  echo "Running in production mode"
  exec npm run start
fi