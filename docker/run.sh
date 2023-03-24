#!/bin/bash

# This script is used to run docker container pulling from an existing image
## how to run from root: ./docker/run.sh [tag]
args=("$@")
TAG_ARG=${args[0]}
TAG=${TAG_ARG:-latest} docker-compose -f docker-compose.yml up
