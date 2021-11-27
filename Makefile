API_PROJECT_NAME := "workflow-api"
API_MAIN_FILE_PAHT := "api/main.go"
SCH_PROJECT_NAME := "workflow-scheduler"
SCH_MAIN_FILE_PAHT := "scheduler/main.go"
NODE_PROJECT_NAME := "workflow-node"
NODE_MAIN_FILE_PAHT := "node/main.go"
PKG := "github.com/infraboard/workflow"
IMAGE_PREFIX := "github.com/infraboard/workflow"

BUILD_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
BUILD_COMMIT := ${shell git rev-parse HEAD}
BUILD_TIME := ${shell date '+%Y-%m-%d %H:%M:%S'}
BUILD_GO_VERSION := $(shell go version | grep -o  'go[0-9].[0-9].*')
VERSION_PATH := "${PKG}/version"

PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/ | grep -v redis)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: all dep lint vet test test-coverage build clean

all: build

dep: ## Get the dependencies
	@go mod tidy

lint: ## Lint Golang files
	@golint -set_exit_status ${PKG_LIST}

vet: ## Run go vet
	@go vet ${PKG_LIST}

test: ## Run unittests
	@go test -short ${PKG_LIST}
	
test-coverage: ## Run tests with coverage
	@go test -short -coverprofile cover.out -covermode=atomic ${PKG_LIST} 
	@cat cover.out >> coverage.txt

build: dep ## Build the binary file
	@go build -a -o dist/${API_PROJECT_NAME} -ldflags "-s -w" -ldflags "-X '${VERSION_PATH}.GIT_BRANCH=${BUILD_BRANCH}' -X '${VERSION_PATH}.GIT_COMMIT=${BUILD_COMMIT}' -X '${VERSION_PATH}.BUILD_TIME=${BUILD_TIME}' -X '${VERSION_PATH}.GO_VERSION=${BUILD_GO_VERSION}'" ${API_MAIN_FILE_PAHT}
	@go build -a -o dist/${SCH_PROJECT_NAME} -ldflags "-s -w" -ldflags "-X '${VERSION_PATH}.GIT_BRANCH=${BUILD_BRANCH}' -X '${VERSION_PATH}.GIT_COMMIT=${BUILD_COMMIT}' -X '${VERSION_PATH}.BUILD_TIME=${BUILD_TIME}' -X '${VERSION_PATH}.GO_VERSION=${BUILD_GO_VERSION}'" ${SCH_MAIN_FILE_PAHT}
	@go build -a -o dist/${NODE_PROJECT_NAME} -ldflags "-s -w" -ldflags "-X '${VERSION_PATH}.GIT_BRANCH=${BUILD_BRANCH}' -X '${VERSION_PATH}.GIT_COMMIT=${BUILD_COMMIT}' -X '${VERSION_PATH}.BUILD_TIME=${BUILD_TIME}' -X '${VERSION_PATH}.GO_VERSION=${BUILD_GO_VERSION}'" ${NODE_MAIN_FILE_PAHT}

linux: ## Linux build
	@GOOS=linux GOARCH=amd64 go build -a -o dist/${API_PROJECT_NAME} -ldflags "-s -w" -ldflags "-X '${VERSION_PATH}.GIT_BRANCH=${BUILD_BRANCH}' -X '${VERSION_PATH}.GIT_COMMIT=${BUILD_COMMIT}' -X '${VERSION_PATH}.BUILD_TIME=${BUILD_TIME}' -X '${VERSION_PATH}.GO_VERSION=${BUILD_GO_VERSION}'" ${API_MAIN_FILE_PAHT}
	@GOOS=linux GOARCH=amd64 go build -a -o dist/${SCH_PROJECT_NAME} -ldflags "-s -w" -ldflags "-X '${VERSION_PATH}.GIT_BRANCH=${BUILD_BRANCH}' -X '${VERSION_PATH}.GIT_COMMIT=${BUILD_COMMIT}' -X '${VERSION_PATH}.BUILD_TIME=${BUILD_TIME}' -X '${VERSION_PATH}.GO_VERSION=${BUILD_GO_VERSION}'" ${SCH_MAIN_FILE_PAHT}
	@GOOS=linux GOARCH=amd64 go build -a -o dist/${NODE_PROJECT_NAME} -ldflags "-s -w" -ldflags "-X '${VERSION_PATH}.GIT_BRANCH=${BUILD_BRANCH}' -X '${VERSION_PATH}.GIT_COMMIT=${BUILD_COMMIT}' -X '${VERSION_PATH}.BUILD_TIME=${BUILD_TIME}' -X '${VERSION_PATH}.GO_VERSION=${BUILD_GO_VERSION}'" ${NODE_MAIN_FILE_PAHT}
	
run-api: dep ## Run Server
	@go run ${API_MAIN_FILE_PAHT} start

build-api: dep ## Build the binary file
	@GOOS=linux GOARCH=amd64 go build -a -o dist/${API_PROJECT_NAME} -ldflags "-s -w" -ldflags "-X '${VERSION_PATH}.GIT_BRANCH=${BUILD_BRANCH}' -X '${VERSION_PATH}.GIT_COMMIT=${BUILD_COMMIT}' -X '${VERSION_PATH}.BUILD_TIME=${BUILD_TIME}' -X '${VERSION_PATH}.GO_VERSION=${BUILD_GO_VERSION}'" ${API_MAIN_FILE_PAHT}
	
run-sch: dep build-sch ## Run schedule
	@go run ${SCH_MAIN_FILE_PAHT} start

build-sch: dep ## Build the binary file
	@GOOS=linux GOARCH=amd64 go build -a -o dist/${SCH_PROJECT_NAME} -ldflags "-s -w" -ldflags "-X '${VERSION_PATH}.GIT_BRANCH=${BUILD_BRANCH}' -X '${VERSION_PATH}.GIT_COMMIT=${BUILD_COMMIT}' -X '${VERSION_PATH}.BUILD_TIME=${BUILD_TIME}' -X '${VERSION_PATH}.GO_VERSION=${BUILD_GO_VERSION}'" ${SCH_MAIN_FILE_PAHT}

run-node: dep build-node ## Run node
	@go run ${NODE_MAIN_FILE_PAHT} start

build-node: dep ## Build the binary file
	@GOOS=linux GOARCH=amd64 go build -a -o dist/${NODE_PROJECT_NAME} -ldflags "-s -w" -ldflags "-X '${VERSION_PATH}.GIT_BRANCH=${BUILD_BRANCH}' -X '${VERSION_PATH}.GIT_COMMIT=${BUILD_COMMIT}' -X '${VERSION_PATH}.BUILD_TIME=${BUILD_TIME}' -X '${VERSION_PATH}.GO_VERSION=${BUILD_GO_VERSION}'" ${NODE_MAIN_FILE_PAHT}

clean: ## Remove previous build
	@go clean .
	@rm -f dist/*

install: ## Install depence go package
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
	@go install github.com/infraboard/protoc-gen-go-ext@v0.0.3
	@go install github.com/infraboard/mcube/cmd/protoc-gen-go-http@latest

gen: ## Init Service
	@protoc -I=. -I=/usr/local/include --go_out=. --go_opt=module=${PKG} --go-grpc_out=. --go-grpc_opt=module=${PKG} api/app/*/pb/*.proto
	@protoc-go-inject-tag -input=api/app/application/*.pb.go
	@protoc-go-inject-tag -input=api/app/deploy/*.pb.go
	@protoc-go-inject-tag -input=api/app/pipeline/*.pb.go
	@protoc-go-inject-tag -input=api/app/action/*.pb.go
	@protoc-go-inject-tag -input=api/app/scm/*.pb.go
	@protoc-go-inject-tag -input=api/app/template/*.pb.go
	@go generate ./...

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'