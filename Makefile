.DEFAULT_GOAL:=help

#============================================================================

# Load environment variables for local development
include .env
export

#============================================================================

.PHONY: dev
dev:							## Run dev container
	@docker compose ls -q | grep -q "instill-vdp" && true || \
		(echo "Error: Run \"make dev PROFILE=controller\" in vdp repository (https://github.com/instill-ai/vdp) in your local machine first." && exit 1)
	@docker inspect --type container ${SERVICE_NAME} >/dev/null 2>&1 && echo "A container named ${SERVICE_NAME} is already running." || \
		echo "Run dev container ${SERVICE_NAME}. To stop it, run \"make stop\"."
	@docker run -d --rm \
		-e DOCKER_HOST=${SOCAT_HOST}:${SOCAT_PORT} \
		-e CFG_ETCD_PORT=${ETCD_CLIENT_PORT} \
		-v $(PWD):/${SERVICE_NAME} \
		-p ${SERVICE_PORT}:${SERVICE_PORT} \
		--network instill-network \
		--name ${SERVICE_NAME} \
		instill/${SERVICE_NAME}:dev >/dev/null 2>&1

.PHONY: logs
logs:							## Tail container logs with -n 10
	@docker logs ${SERVICE_NAME} --follow --tail=10

.PHONY: stop
stop:							## Stop container
	@docker stop -t 1 ${SERVICE_NAME}

.PHONY: rm
rm:								## Remove container
	@docker rm -f ${SERVICE_NAME}

.PHONY: top
top:							## Display all running service processes
	@docker top ${SERVICE_NAME}

.PHONY: build
build:							## Build dev docker image
	@docker build \
		--build-arg SERVICE_NAME=${SERVICE_NAME} \
		--build-arg GOLANG_VERSION=${GOLANG_VERSION} \
		--build-arg K6_VERSION=${K6_VERSION} \
		-f Dockerfile.dev  -t instill/${SERVICE_NAME}:dev .

.PHONY: go-gen
go-gen:       					## Generate codes
	go generate ./...

.PHONY: unit-test
unit-test:       				## Run unit test
	@go test -v -race -coverpkg=./... -coverprofile=/vdp/coverage.out ./...
	@go tool cover -func=/vdp/coverage.out
	@go tool cover -html=/vdp/coverage.out
	@rm /vdp/coverage.out

.PHONY: integration-test
integration-test:				## Run integration test
	TEST_FOLDER_ABS_PATH=${PWD} k6 run \
		-e API_GATEWAY_PROTOCOL=${API_GATEWAY_PROTOCOL} -e API_GATEWAY_HOST=${API_GATEWAY_HOST} -e API_GATEWAY_PORT=${API_GATEWAY_PORT} \
		integration-test/grpc.js --no-usage-report --quiet

.PHONY: help
help:       	 				## Show this help
	@echo "\nMakefile for local development"
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m (default: help)\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
