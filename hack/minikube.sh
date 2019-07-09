#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# # Create variables for project directories
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
  DIR="$( cd -P "$( dirname "$SOURCE" )" >/dev/null 2>&1 && pwd )"
  SOURCE="$(readlink "$SOURCE")"
  [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done
PROJECT_DIR="$( cd -P "$( dirname "$SOURCE" )/.." >/dev/null 2>&1 && pwd )"


# Start minikube with proper resources
minikube_up () {
  MINIKUBE_STATUS="$( minikube status | awk '/host:[[:space:]]/ { print $2 }' )"
  echo "Minikube status: $MINIKUBE_STATUS"
  if [ "$MINIKUBE_STATUS" != "Running" ]; then
    minikube start --memory=8192 --cpus=4
  fi
  eval $(minikube docker-env)
}


kill_mounts () {
  for pid in $(ps -ef | grep -v grep | grep 'minikube mount' | awk '{print $2}'); do 
    kill -9  $pid
  done
}


# # Mount data and model directories
# # These processes will run in the background and be terminated when script exits
mount_directories () {
  kill_mounts
  
  echo "Mounting directories to minikube host..."
  mkdir -p $PROJECT_DIR/data
  minikube mount $PROJECT_DIR/data:/data &
  # trap "trap - SIGTERM && kill -- $!" SIGINT SIGTERM EXIT

  sleep 1  # To separate output of background tasks

  mkdir -p $PROJECT_DIR/models
  minikube mount $PROJECT_DIR/models:/models &
  # trap "trap - SIGTERM && kill -- $!" SIGINT SIGTERM EXIT
}


# Deploy services with skaffold
skaffold_deploy () {
  cd $PROJECT_DIR
  skaffold run -p minikube
}