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


# download a model
download_model () {
  VAN_GOGH_MODEL_GZIP="$PROJECT_DIR/models/model_van-gogh_ckpt.tar.gz"
  VAN_GOGH_MODEL_DIR="$PROJECT_DIR/models/model_van-gogh"

  # Create models folder if it doesn't exist
  mkdir -p $PROJECT_DIR/models

  # If van gogh model dir doesn't exist unzip it from tar file
  if [ ! -d "$VAN_GOGH_MODEL_DIR" ]; then

    # Fetch tar file if it doesn't exist
    if [ ! -f "$VAN_GOGH_MODEL_GZIP" ]; then
      echo "Downloading model..."
      curl -o $VAN_GOGH_MODEL_GZIP \
        https://saas-interview.s3.amazonaws.com/model_van-gogh_ckpt.tar.gz
    fi
  
    cd $PROJECT_DIR/models && tar -zxvf model_van-gogh_ckpt.tar.gz
  fi
}

# download a test image
download_data () {
  TEST_IMG="$PROJECT_DIR/data/test.jpg"

  # Create data directories
  mkdir -p $PROJECT_DIR/data/stylized
  mkdir -p $PROJECT_DIR/data/uploaded

  # Fetch test image if it doesn't exist
  if [ ! -f "$TEST_IMG" ]; then
    echo "Downloading test image..."
    curl -o $TEST_IMG \
      https://saas-interview.s3.amazonaws.com/67ffd833-ea51-49da-bda2-d835a8e81ceb-response.jpg
  fi
}