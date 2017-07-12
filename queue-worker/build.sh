#!/bin/sh
cp -r ../gateway .

docker build -t functions/queue-worker:latest .

