.PHONY: mockgen build test docker docker-build docker-push

GIT_COMMIT   ?= $(shell git rev-parse HEAD)
GO_VERSION   ?= $(shell go version | awk {'print $$3'})
DOCKER_IMAGE ?= figmentnetworks/celo-indexer
DOCKER_TAG   ?= latest

# Generate mocks
mockgen:
	@echo "[mockgen] generating mocks"
	@mockgen -destination mock/baseclient/mocks.go github.com/figment-networks/celo-indexer/client Client,RequestCounter
	@mockgen -destination mock/client/mocks.go github.com/figment-networks/celo-indexer/client/figmentclient Client
	@mockgen -destination mock/indexer/mocks.go github.com/figment-networks/celo-indexer/indexer ConfigParser
	@mockgen -destination mock/store/mocks.go github.com/figment-networks/celo-indexer/store AccountActivitySeq,BlockSeq,BlockSummary,Database,Reports,Syncables,SystemEvents,ValidatorAgg,ValidatorSeq,ValidatorSummary,ValidatorGroupAgg,ValidatorGroupSeq,ValidatorGroupSummary

# Build the binary
build:
	go build \
		-ldflags "\
			-X github.com/figment-networks/celo-indexer/cli.gitCommit=${GIT_COMMIT} \
			-X github.com/figment-networks/celo-indexer/cli.goVersion=${GO_VERSION}"

# Run tests
test:
	go test -race -cover ./...

# Build a local docker image for testing
docker:
	docker build -t celo-indexer -f Dockerfile .

# Build a public docker image
docker-build:
	docker build \
		-t ${DOCKER_IMAGE}:${DOCKER_TAG} \
		-f Dockerfile \
		.

# Push docker images
docker-push:
	docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:latest
	docker push ${DOCKER_IMAGE}:${DOCKER_TAG}
	docker push ${DOCKER_IMAGE}:latest