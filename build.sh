#!/bin/bash
# build and remove dangling images in the system
docker build -t email4tong/dedup:latest .
docker rmi -f $(docker images -f "dangling=true" -q)