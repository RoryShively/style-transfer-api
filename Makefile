SHELL := /bin/bash

DATA_DIR := data/
MODEL_DIR := models/

TEST_IMG := $(DATA_DIR)/test.img
VAN_GOGH_MODEL := $(MODEL_DIR)/model_van-gogh/


minikube-mount-pid = $(word 1,$(shell ps | grep -v grep \
                                         | grep 'minikube mount'))

ifndef VERBOSE
.SILENT:
endif

.PHONY: start clean unit_tests e2e_tests all_tests


$(DATA_DIR):
	mkdir -p $(DATA_DIR)/stylized
	mkdir -p $(DATA_DIR)/uploaded

$(TEST_IMG): $(DATA_DIR)
	source "./hack/setup.sh" && \
		download_data

$(MODEL_DIR):
	mkdir -p $(MODEL_DIR)

$(VAN_GOGH_MODEL): $(MODEL_DIR)
	source "./hack/setup.sh" && \
		download_model

##-- RULES --##

start: $(TEST_IMG) $(VAN_GOGH_MODEL)
	source "./hack/minikube.sh" && \
		minikube_up && \
		mount_directories && \
		skaffold_deploy

	# minikube start --memory=8192 --cpus=4 \
	# minikube mount data:/data & \
	# minikube mount models:/models & \
	# skaffold run -p minikube

clean:
	source "./hack/minikube.sh" && \
		kill_mounts \
	minikube delete


unit_tests:
	echo "Running unit tests..." && \
	cd src/api && \
	go test -tags=unit -v
	
e2e_tests: $(TEST_IMG) $(VAN_GOGH_MODEL)
	echo "Running e2e tests..." && \
	source "./hack/test-vars.sh" && \
		cd src/api && \
		go test -tags=e2e -v

all_tests: unit_tests e2e_tests
	# echo "Running all tests..."
	# unit_tests
	# e2e_tests