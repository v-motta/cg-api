#
# Vinicius Motta - 04-2024
#

COMMIT := $(shell git log -1 --pretty=format:"%H")

PROJECT_NAME := "cost-guardian-api"
DOCKER_IMAGE := "$(PROJECT_NAME)-i"
RUN_SCRIPT := "run-container.sh"

#
# Build image
#
build_image:
	@echo
	@echo "Building image..."
	docker build --rm \
		--build-arg COMMIT=$(COMMIT) \
		-t $(DOCKER_IMAGE) .
	docker images | grep $(PROJECT_NAME)
	@echo "Image built."

#
# Run container
#
run_container:
	@echo
	@echo "Running container..."
	./$(RUN_SCRIPT)
	@echo "Container running."