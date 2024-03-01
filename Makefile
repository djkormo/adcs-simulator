dpl ?= deploy.env
include $(dpl)
export $(shell sed 's/=.*//' $(dpl))

COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
PROJECT?=github.com/djkormo/adcs-simulator

# DOCKER TASKS
# Build the container
build: ## Build the container
	docker build -t $(APP_NAME) . --progress=plain --build-arg VERSION=${VERSION} --build-arg  COMMIT=${COMMIT} --build-arg BUILD_TIME=${BUILD_TIME} --build-arg PROJECT=${PROJECT}

build-nc: ## Build the container without caching
	docker build --no-cache -t $(APP_NAME) . --progress=plain --build-arg VERSION=${VERSION} --build-arg  COMMIT=${COMMIT} --build-arg BUILD_TIME=${BUILD_TIME} --build-arg PROJECT=${PROJECT}

run: ## Run container on port configured in `config.env`
	docker run --name $(APP_NAME) -it  $(DOCKER_REPO)/$(APP_NAME):latest  bash

stop: ## Stop and remove a running container
	docker rm -f $(APP_NAME)

release: build-nc publish ## Make a release by building and publishing the `{version}` ans `latest` tagged containers to ECR

# Docker publish
publish: repo-login publish-latest publish-version ## Publish the `{version}` ans `latest` tagged containers to ECR

publish-latest: tag-latest ## Publish the `latest` taged container to ECR
	@echo 'publish latest to $(DOCKER_REPO)'
	docker push $(DOCKER_REPO)/$(APP_NAME):latest

publish-version: tag-version ## Publish the `{version}` taged container to ECR
	@echo 'publish $(VERSION) to $(DOCKER_REPO)'
	docker push $(DOCKER_REPO)/$(APP_NAME):$(VERSION)

# Docker tagging
tag: tag-latest tag-version ## Generate container tags for the `{version}` ans `latest` tags

tag-latest: ## Generate container `{version}` tag
	@echo 'create tag latest'
	docker tag $(APP_NAME) $(DOCKER_REPO)/$(APP_NAME):latest

tag-version: ## Generate container `latest` tag
	@echo 'create tag $(VERSION)'
	docker tag $(APP_NAME) $(DOCKER_REPO)/$(APP_NAME):$(VERSION)



inspect: ## Generate container `latest` tag
	@echo 'inspect $(APP_NAME)'
	docker history $(APP_NAME)
	docker inspect $(APP_NAME)



.PHONY: operator-build 
operator-build:
	go clean 
	go fmt ## Build manager binary.
	go vet
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build \
		-ldflags "-s -w -X ${PROJECT}/version.Release=${VERSION} \
		-X ${PROJECT}/version.Commit=${COMMIT} -X ${PROJECT}/version.BuildTime=${BUILD_TIME}" \
		-o ${APP_NAME}

.PHONY: docker-build
docker-build: ## Build docker image with the manager.
	DOCKER_BUILDKIT=1 docker build -t $(APP_NAME) . --progress=plain