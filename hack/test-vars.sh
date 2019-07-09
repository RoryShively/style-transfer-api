#!/usr/bin/env bash

# # Create variables for project directories
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
  DIR="$( cd -P "$( dirname "$SOURCE" )" >/dev/null 2>&1 && pwd )"
  SOURCE="$(readlink "$SOURCE")"
  [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done
PROJECT_DIR="$( cd -P "$( dirname "$SOURCE" )/.." >/dev/null 2>&1 && pwd )"

MINIKUBE_HOST=$(minikube ip)
NODE_PORT=$(kubectl get svc api -o jsonpath="{.spec.ports[?(@.port==3100)].nodePort}")

export TEST_IMAGE_PATH="$PROJECT_DIR/data/test.jpg"
export API_URL="http://$MINIKUBE_HOST:$NODE_PORT"