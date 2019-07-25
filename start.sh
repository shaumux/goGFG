#!/usr/bin/env bash

docker-compose up -d storage
RETRIES=10
until docker-compose exec storage curl localhost:9200 > /dev/null 2>&1 || [ $RETRIES -eq 0 ]; do
  echo "Waiting for elastic server, $((RETRIES--)) remaining attempts..."
  sleep 5
done
docker-compose build goGFG
docker-compose up -d goGFG