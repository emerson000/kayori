#!/bin/bash

docker compose build
docker compose down
docker compose -f docker-compose-prod.yml up --scale drone=2 -d