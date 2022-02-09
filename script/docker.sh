#!/bin/bash

docker build --no-cache -f Dockerfile -t ghcr.io/distbuild/worker:latest .
