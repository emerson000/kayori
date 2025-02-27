#!/bin/bash

docker compose down
git reset --hard HEAD
git pull
docker compose -f docker-compose-prod.yml up -d --build