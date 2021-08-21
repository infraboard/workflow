API_PROJECT_NAME := "workflow-api"
API_MAIN_FILE_PAHT := "api/main.go"
SCH_PROJECT_NAME := "workflow-scheduler"
SCH_MAIN_FILE_PAHT := "scheduler/main.go"
NODE_PROJECT_NAME := "workflow-node"
NODE_MAIN_FILE_PAHT := "node/main.go"
PKG := "github.com/infraboard/workflow"
IMAGE_PREFIX := "github.com/infraboard/workflow"

GO_MODDIR := $(shell go env GOMODCACHE)
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
	@go fmt ./...
	@sh ./script/build.sh local dist/${API_PROJECT_NAME} ${API_MAIN_FILE_PAHT} ${IMAGE_PREFIX} ${PKG}
	@sh ./script/build.sh local dist/${SCH_PROJECT_NAME} ${SCH_MAIN_FILE_PAHT} ${IMAGE_PREFIX} ${PKG}
	@sh ./script/build.sh local dist/${NODE_PROJECT_NAME} ${NODE_MAIN_FILE_PAHT} ${IMAGE_PREFIX} ${PKG}

linux: ## Linux build
	@sh ./script/build.sh linux dist/${API_PROJECT_NAME} ${API_MAIN_FILE_PAHT} ${IMAGE_PREFIX} ${PKG}
	@sh ./script/build.sh linux dist/${SCH_PROJECT_NAME} ${SCH_MAIN_FILE_PAHT} ${IMAGE_PREFIX} ${PKG}
	@sh ./script/build.sh linux dist/${NODE_PROJECT_NAME} ${NODE_MAIN_FILE_PAHT} ${IMAGE_PREFIX} ${PKG}
	
run-api: dep ## Run Server
	@sh ./script/build.sh local dist/${API_PROJECT_NAME} ${API_MAIN_FILE_PAHT} ${IMAGE_PREFIX} ${PKG}
	@./dist/${API_PROJECT_NAME} start

build-sch: dep ## Build the binary file
	@go fmt ./...
	@sh ./script/build.sh local dist/${SCH_PROJECT_NAME} ${SCH_MAIN_FILE_PAHT} ${IMAGE_PREFIX} ${PKG}

linux-sch: ## Linux build
	@sh ./script/build.sh linux dist/${SCH_PROJECT_NAME} ${SCH_MAIN_FILE_PAHT} ${IMAGE_PREFIX} ${PKG}
	
run-sch: dep build-sch ## Run schedule
	@./dist/${SCH_PROJECT_NAME} start

build-node: dep ## Build the binary file
	@go fmt ./...
	@sh ./script/build.sh local dist/${NODE_PROJECT_NAME} ${NODE_MAIN_FILE_PAHT} ${IMAGE_PREFIX} ${PKG}

run-node: dep build-node ## Run node
	@./dist/${NODE_PROJECT_NAME} start


clean: ## Remove previous build
	@go clean .
	@rm -f dist/*

install: ## Install depence go package
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
	@go install github.com/infraboard/protoc-gen-go-ext@v0.0.3
	@go install github.com/infraboard/mcube/cmd/protoc-gen-go-http@latest

codegen: ## Init Service
	@protoc -I=.  -I${GO_MODDIR} --go-ext_out=. --go-ext_opt=module=${PKG} --go-grpc_out=. --go-grpc_opt=module=${PKG} --go-http_out=. --go-http_opt=module=${PKG} api/pkg/*/pb/*.proto
	@go generate ./...

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'