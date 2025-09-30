# Go-React-Admin Backend Makefile
# 项目根目录: backend/

.PHONY: help dev build run test fmt vet lint clean tidy migrate

# 默认目标
.DEFAULT_GOAL := help

# 项目变量
BINARY_NAME=server
MAIN_PATH=./cmd/server/main.go
BUILD_DIR=./bin
CONFIG_PATH=./configs/config.yml

# Go 命令
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

## help: 显示帮助信息
help:
	@echo "Go-React-Admin Backend - 可用命令:"
	@echo ""
	@echo "  make dev          - 启动开发服务器(热重载,需要 air)"
	@echo "  make run          - 直接运行服务器"
	@echo "  make build        - 编译二进制文件到 bin/ 目录"
	@echo "  make test         - 运行所有测试"
	@echo "  make fmt          - 格式化代码(go fmt)"
	@echo "  make vet          - 代码检查(go vet)"
	@echo "  make lint         - 代码静态分析(需要 staticcheck)"
	@echo "  make tidy         - 整理依赖(go mod tidy)"
	@echo "  make clean        - 清理编译文件"
	@echo "  make migrate      - 运行数据库迁移"
	@echo "  make check        - 运行所有代码质量检查(fmt+vet+test)"
	@echo ""

## dev: 启动开发服务器(热重载)
dev:
	@echo "🚀 启动开发服务器(热重载)..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "❌ 未找到 air 工具,请先安装:"; \
		echo "   go install github.com/air-verse/air@latest"; \
		exit 1; \
	fi

## run: 直接运行服务器
run:
	@echo "🚀 启动服务器..."
	@if [ ! -f $(CONFIG_PATH) ]; then \
		echo "⚠️  配置文件不存在,正在复制示例配置..."; \
		cp configs/config.example.yml $(CONFIG_PATH); \
		echo "✅ 已创建 $(CONFIG_PATH),请修改数据库配置"; \
	fi
	$(GOCMD) run $(MAIN_PATH)

## build: 编译二进制文件
build:
	@echo "🔨 编译项目..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v $(MAIN_PATH)
	@echo "✅ 编译完成: $(BUILD_DIR)/$(BINARY_NAME)"

## test: 运行所有测试
test:
	@echo "🧪 运行测试..."
	$(GOTEST) -v -cover ./...

## test-verbose: 运行详细测试
test-verbose:
	@echo "🧪 运行详细测试..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "✅ 覆盖率报告: coverage.html"

## fmt: 格式化代码
fmt:
	@echo "🎨 格式化代码..."
	$(GOFMT) ./...
	@echo "✅ 代码格式化完成"

## vet: 代码检查
vet:
	@echo "🔍 运行 go vet..."
	$(GOVET) ./...
	@echo "✅ 代码检查完成"

## lint: 静态分析(需要 staticcheck)
lint:
	@echo "🔍 运行静态分析..."
	@if command -v staticcheck > /dev/null; then \
		staticcheck ./...; \
		echo "✅ 静态分析完成"; \
	else \
		echo "⚠️  未找到 staticcheck,跳过静态分析"; \
		echo "   安装: go install honnef.co/go/tools/cmd/staticcheck@latest"; \
	fi

## tidy: 整理依赖
tidy:
	@echo "📦 整理依赖..."
	$(GOMOD) tidy
	@echo "✅ 依赖整理完成"

## clean: 清理编译文件
clean:
	@echo "🧹 清理编译文件..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	@echo "✅ 清理完成"

## check: 运行所有代码质量检查
check: fmt vet test
	@echo "✅ 所有检查通过!"

## migrate: 运行数据库迁移(需要先启动MySQL)
migrate:
	@echo "🔄 运行数据库迁移..."
	@if [ -f ../docs/schema.sql ]; then \
		mysql -u root -p go_react_admin < ../docs/schema.sql; \
		echo "✅ 数据库迁移完成"; \
	else \
		echo "❌ 未找到 schema.sql 文件"; \
		exit 1; \
	fi

## install-tools: 安装开发工具
install-tools:
	@echo "📦 安装开发工具..."
	go install github.com/air-verse/air@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	@echo "✅ 开发工具安装完成"

## deps: 安装项目依赖
deps:
	@echo "📦 安装项目依赖..."
	$(GOMOD) download
	@echo "✅ 依赖安装完成"

## docker-up: 启动 Docker 环境(MySQL + Redis)
docker-up:
	@echo "🐳 启动 Docker 环境..."
	@if command -v docker-compose > /dev/null; then \
		docker-compose up -d; \
		echo "✅ Docker 环境已启动"; \
	else \
		echo "❌ 未找到 docker-compose"; \
	fi

## docker-down: 停止 Docker 环境
docker-down:
	@echo "🐳 停止 Docker 环境..."
	@if command -v docker-compose > /dev/null; then \
		docker-compose down; \
		echo "✅ Docker 环境已停止"; \
	else \
		echo "❌ 未找到 docker-compose"; \
	fi

## db-backup: 备份数据库
db-backup:
	@echo "💾 备份数据库..."
	@mkdir -p ./backups
	mysqldump -u root -p go_react_admin > ./backups/backup_$$(date +%Y%m%d_%H%M%S).sql
	@echo "✅ 数据库备份完成"

## db-restore: 恢复数据库(需要提供 FILE 参数)
db-restore:
	@if [ -z "$(FILE)" ]; then \
		echo "❌ 请提供备份文件: make db-restore FILE=backups/backup.sql"; \
		exit 1; \
	fi
	@echo "💾 恢复数据库从 $(FILE)..."
	mysql -u root -p go_react_admin < $(FILE)
	@echo "✅ 数据库恢复完成"

## info: 显示项目信息
info:
	@echo "📊 项目信息:"
	@echo "  Go 版本:      $$(go version)"
	@echo "  项目路径:     $$(pwd)"
	@echo "  二进制名称:   $(BINARY_NAME)"
	@echo "  配置文件:     $(CONFIG_PATH)"
	@echo ""
	@echo "📦 依赖统计:"
	@$(GOMOD) graph | wc -l | xargs echo "  依赖包数量: "
	@echo ""
	@echo "📈 代码统计:"
	@find . -name "*.go" -not -path "./vendor/*" | xargs wc -l | tail -1