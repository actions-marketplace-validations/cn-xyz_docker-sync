#!/bin/bash
set -e

export REPO_URL=${REPO_URL}
export REMOTE_URL=${REMOTE_URL}
export HUB_LOGIN_URL=${HUB_LOGIN_URL}
export DOCKER_USER=${DOCKER_USER}
export DOCKER_PASS=${DOCKER_PASS}

echo "## Check Package Version ##################"

skopeo --version

docker-sync