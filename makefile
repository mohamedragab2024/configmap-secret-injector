# Variables
IMAGE_NAME := configmap-secret-injector
IMAGE_TAG := latest
DOCKER_REGISTRY := docker.io/regoo707

# Build the binary
.PHONY: build
build:
	go build -o bin/$(IMAGE_NAME) cmd/main.go

# Build Docker image
.PHONY: docker-build
docker-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

# Push Docker image
.PHONY: docker-push
docker-push:
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)
	docker push $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

# All-in-one command
.PHONY: all
all: build docker-build docker-push

.PHONY: run
run:
	go run cmd/main.go
