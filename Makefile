# Makefile for the project
.DEFAULT_GOAL := help

# 定义Makefile all伪目标，执行`make`时，会执行all伪目标
all: build test cover clean lint tidy format help

# ==============================================================================
# Includes

# 确保 `include common.mk` 位于第一行，common.mk 中定义了一些变量，后面的子 makefile 有依赖
#include scripts/make-rules/common.mk
#include scripts/make-rules/all.mk

# ==============================================================================
# Usage 
define USAGE_OPTIONS

选项：
   BINS      要构建的二进制文件。默认为cmd中的所有文件
   VERSION   编译到二进制文件中的版本信息
   V         设置为1为启用详细的构建信息输出。默认为0。
endef
export USAGE_OPTIONS

## ----------------------------------------------------
## Binaries
## ----------------------------------------------------

##@ build
build: ## 编译所有二进制文件
	go build -o bin/ ./...

## ----------------------------------------------------
## Testing
## ----------------------------------------------------

test: ## 执行单元测试
	go test ./...

cover: ## 执行单元测试覆盖率
	go test -coverprofile=coverage.out ./...

coverhtml:  ## 查看单元测试覆盖率
   go tool cover -html coverage.out
## ----------------------------------------------------
## Cleanup
## ----------------------------------------------------
clean: # 清理构建产物以及临时文件目录
	@echo "======= clean up all build artifacts ========="
	rm -rf bin/

## ----------------------------------------------------
## Formatting
## ----------------------------------------------------
format: ## 格式化源码文件
	go fmt ./...

## ----------------------------------------------------
## lint / Verify
## ----------------------------------------------------
lint: ## 执行lint检查
	golangci-lint   run   ./...

tidy: ## go mod tidy
	go mod tidy

## ----------------------------------------------------
## swagger 生成swagger 文档
## ----------------------------------------------------
swagger: ## 生成swagger 文档
	swag  init  -g  .\cmd\apiserver\main.go

help: Makefile ## 显示帮助信息
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<TARGETS> <OPTIONS>\033[0m\n\n\033[35mTargets:\033[0m\n"} /^[0-9A-Za-z._-]+:.*?##/ { printf "  \033[36m%-45s\033[0m %s\n", $$1, $$2 } /^\$\([0-9A-Za-z_-]+\):.*?##/ { gsub("_","-", $$1); printf "  \033[36m%-45s\033[0m %s\n", tolower(substr($$1, 3, length($$1)-7)), $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' Makefile #$(MAKEFILE_LIST)
	@echo -e "$$USAGE_OPTIONS"

# 伪目标
.PHONY: all build test cover clean lint tidy format swagger  help
