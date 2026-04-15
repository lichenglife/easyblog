.PHONY: help build run test clean lint docker-up docker-down db-init

# 默认目标
help:
	@echo "EasyBlog Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build          - 编译后端项目"
	@echo "  make run            - 运行后端项目 (开发模式)"
	@echo "  make run-prod       - 运行后端项目 (生产模式)"
	@echo "  make test           - 运行后端单元测试"
	@echo "  make clean          - 清理编译文件"
	@echo "  make lint           - 运行代码检查"
	@echo "  make docker-up      - 启动 Docker 服务 (MySQL, Redis, MinIO)"
	@echo "  make docker-down    - 停止 Docker 服务"
	@echo "  make db-init        - 初始化数据库"
	@echo "  make frontend       - 安装前端依赖并构建"
	@echo "  make frontend-dev   - 启动前端开发服务器"
	@echo "  make test-all       - 运行所有测试 (后端 + 前端)"
	@echo "  make logs           - 查看后端日志"
	@echo "  make logs-clean     - 清理日志文件"
	@echo "  make deploy-local   - 本地部署 (启动所有服务)"
	@echo "  make stop-local     - 停止本地服务"

# 编译后端
build:
	@echo "Building backend..."
	@cd cmd && go build -o ../bin/easyblog .
	@echo "Build completed."

# 运行后端 (开发模式)
run:
	@echo "Starting easyblog (dev mode)..."
	@mkdir -p logs
	@CONFIG_PATH=configs/config.yaml go run cmd/apiserver/main.go server

# 运行后端 (生产模式)
run-prod:
	@echo "Starting easyblog (prod mode)..."
	@mkdir -p logs
	@CONFIG_PATH=configs/config.prod.yaml ./bin/easyblog server

# 运行后端单元测试
test:
	@echo "Running backend tests..."
	@go test -v -coverprofile=coverage.out ./... 2>&1 | tee test-output.log
	@echo "Test completed. Coverage report: coverage.out"

# 清理编译文件
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf coverage.out
	@rm -rf test-output.log
	@rm -rf logs/*.log
	@echo "Clean completed."

# 代码检查
lint:
	@echo "Running linter..."
	@golangci-lint run ./...
	@echo "Lint completed."

# 启动 Docker 服务
docker-up:
	@echo "Starting Docker services..."
	@docker-compose up -d
	@echo "Docker services started: MySQL, Redis, MinIO"
	@echo "  - MySQL:   localhost:3306 (root/root123)"
	@echo "  - Redis:   localhost:6379"
	@echo "  - MinIO:   localhost:9000 (minioadmin/minioadmin123)"

# 停止 Docker 服务
docker-down:
	@echo "Stopping Docker services..."
	@docker-compose down
	@echo "Docker services stopped."

# 重启 Docker 服务
docker-restart: docker-down docker-up
	@echo "Docker services restarted."

# 初始化数据库
db-init:
	@echo "Initializing database..."
	@docker exec -i easyblog-mysql mysql -uroot -proot123 < scripts/db/init.sql
	@echo "Database initialized."

# 查看数据库日志
db-logs:
	@docker logs -f easyblog-mysql

# 查看 Redis 日志
redis-logs:
	@docker logs -f easyblog-redis

# 查看 MinIO 日志
minio-logs:
	@docker logs -f easyblog-minio

# 安装前端依赖并构建
frontend:
	@echo "Building frontend..."
	@cd frontend && npm install
	@cd frontend && npm run build
	@echo "Frontend build completed."

# 启动前端开发服务器
frontend-dev:
	@cd frontend && npm run dev

# 运行前端测试
frontend-test:
	@cd frontend && npm run test:unit

# 运行所有测试
test-all: test frontend-test
	@echo "All tests completed."

# 查看后端日志
logs:
	@echo "Tailing logs from logs/apiserver.log..."
	@tail -f logs/apiserver.log

# 清理日志文件
logs-clean:
	@echo "Cleaning log files..."
	@rm -rf logs/*.log
	@echo "Log files cleaned."

# 本地部署 - 启动所有服务
deploy-local: docker-up
	@echo "Deploying EasyBlog locally..."
	@mkdir -p logs
	@echo "Starting backend server..."
	@nohup make run > logs/deploy.log 2>&1 &
	@echo "Backend started. PID: $$(pgrep -f 'go run cmd/apiserver/main.go')"
	@echo ""
	@echo "=========================================="
	@echo "EasyBlog 本地部署完成！"
	@echo "=========================================="
	@echo "后端 API:   http://localhost:8080"
	@echo "前端服务：cd frontend && npm run dev"
	@echo "日志目录：logs/"
	@echo ""
	@echo "查看日志：make logs"
	@echo "停止服务：make stop-local"
	@echo "=========================================="

# 停止本地服务
stop-local:
	@echo "Stopping local services..."
	@-pkill -f 'go run cmd/apiserver/main.go'
	@-pkill -f './bin/easyblog'
	@echo "Backend server stopped."
	@echo ""
	@echo "提示：如需停止 Docker 服务，运行 make docker-down"

# 全量部署 (本地测试)
deploy-local-full: docker-up db-init
	@echo "Deploying EasyBlog locally (full)..."
	@mkdir -p logs
	@echo "Starting backend server..."
	@nohup make run > logs/deploy.log 2>&1 &
	@sleep 2
	@echo ""
	@echo "=========================================="
	@echo "EasyBlog 全量部署完成！"
	@echo "=========================================="
	@echo "后端 API:   http://localhost:8080"
	@echo "健康检查：curl http://localhost:8080/healthz"
	@echo "前端服务：cd frontend && npm run dev"
	@echo "日志目录：logs/"
	@echo ""
	@echo "查看日志：make logs"
	@echo "停止服务：make stop-local"
	@echo "=========================================="

# 生产环境构建
build-prod:
	@echo "Building for production..."
	@cd frontend && npm install && npm run build
	@cd cmd && go build -ldflags="-s -w" -o ../bin/easyblog .
	@echo "Production build completed."
	@ls -lh bin/

# 生成 API 文档 (Swagger)
docs:
	@echo "Generating API documentation..."
	@swag init -g cmd/apiserver/main.go -o docs/swagger
	@echo "Documentation generated at docs/swagger/"

# 检查依赖
deps-check:
	@echo "Checking dependencies..."
	@go mod verify
	@go mod tidy
	@echo "Dependencies verified."

# 更新依赖
deps-update:
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@echo "Dependencies updated."

# 格式化代码
fmt:
	@echo "Formatting code..."
	@gofmt -w .
	@echo "Code formatted."

# 准备提交
pre-commit: fmt lint test
	@echo "Pre-commit checks completed."
