KIND_CLUSTER_NAME ?= "superk"
K8S_NODE_IMAGE ?= v1.15.3

## ----------Targets------------

## help: 
##      Show this help.
help: Makefile
	@sed -n 's/^##//p' $<

## all:
##      Default action - format, check, test and build.
all: fmt checks test build

## kind-create:
##      Create a new kind cluster for local testing.
kind-create:
ifeq (1, $(shell kind get clusters | grep ${KIND_CLUSTER_NAME} | wc -l))
	@echo "Cluster already exists" 
else
	@echo "Creating Cluster"	
	kind create cluster --name ${KIND_CLUSTER_NAME} --image=kindest/node:${K8S_NODE_IMAGE}
endif

## kind-delete:
##      Delete the kind cluster created with kind-create.
kind-delete:
	kind delete cluster --name ${KIND_CLUSTER_NAME}

## kind-recreate:
##      Create a new kind cluster for local testing after deleting it if it already exists. 
##      Use this if you start getting the following error when running kubectl commands:
##      "The connection to the server 127.0.0.1:40495 was refused - did you specify the right host or port?"
kind-recreate: 
ifeq (1, $(shell kind get clusters | grep ${KIND_CLUSTER_NAME} | wc -l))
	@echo "Deleting Cluster" 
	kind delete cluster --name ${KIND_CLUSTER_NAME}
endif
	@echo "Creating Cluster"	
	kind create cluster --name ${KIND_CLUSTER_NAME} --image=kindest/node:${K8S_NODE_IMAGE}

## fmt:
##      Format the code.
fmt:
	find . -name '*.go' | grep -v vendor | xargs gofmt -s -w

## checks:
##      Check/lint code on the code.
checks:
	golangci-lint run

## test:
##      Run unit tests.
test: fmt checks
	mkdir -p codecoverage
	go test ./... -v -coverprofile ./codecoverage/cover.out
	go tool cover -html=./codecoverage/cover.out -o ./codecoverage/cover.html

## build:
##      Build superk binary.
build: fmt checks 
	go build -o superk ./cmd

## run:
##      Build and run the tool.
run: build
	./superk

## debug:
##      Start superk using Delve ready for debugging from VSCode.
debug: build
	dlv exec ./superk --headless --listen localhost:2345 --api-version 2