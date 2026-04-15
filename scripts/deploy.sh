#!/bin/bash

# EasyBlog 本地部署脚本
# 用途：快速部署 EasyBlog 到本地环境

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查依赖
check_dependencies() {
    log_info "检查依赖..."

    # 检查 Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi

    # 检查 Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose 未安装，请先安装 Docker Compose"
        exit 1
    fi

    # 检查 Go
    if ! command -v go &> /dev/null; then
        log_error "Go 未安装，请先安装 Go (1.23+)"
        exit 1
    fi

    # 检查 Node.js
    if ! command -v node &> /dev/null; then
        log_warn "Node.js 未安装，前端功能将不可用"
    fi

    log_info "依赖检查通过"
}

# 启动基础设施
start_infrastructure() {
    log_info "启动基础设施 (MySQL, Redis, MinIO)..."
    docker-compose up -d

    # 等待服务启动
    log_info "等待服务启动..."
    sleep 5

    log_info "基础设施启动完成"
}

# 初始化数据库
init_database() {
    log_info "初始化数据库..."

    # 等待 MySQL 就绪
    for i in {1..30}; do
        if docker exec easyblog-mysql mysqladmin ping -uroot -proot123 &> /dev/null; then
            log_info "MySQL 就绪"
            break
        fi
        log_info "等待 MySQL 启动... ($i/30)"
        sleep 1
    done

    # 执行初始化脚本
    docker exec -i easyblog-mysql mysql -uroot -proot123 < scripts/db/init.sql
    log_info "数据库初始化完成"
}

# 构建后端
build_backend() {
    log_info "构建后端..."
    cd cmd && go build -o ../bin/easyblog .
    cd ..
    log_info "后端构建完成"
}

# 启动后端服务
start_backend() {
    log_info "启动后端服务..."

    # 创建日志目录
    mkdir -p logs

    # 后台启动
    nohup ./bin/easyblog server > logs/deploy.log 2>&1 &
    BACKEND_PID=$!

    # 等待服务启动
    sleep 3

    # 检查服务是否正常启动
    if curl -s http://localhost:8080/healthz &> /dev/null; then
        log_info "后端服务启动成功 (PID: $BACKEND_PID)"
    else
        log_error "后端服务启动失败，请查看 logs/deploy.log"
        exit 1
    fi
}

# 构建前端
build_frontend() {
    if command -v node &> /dev/null; then
        log_info "构建前端..."
        cd frontend
        npm install
        npm run build
        cd ..
        log_info "前端构建完成"
    else
        log_warn "跳过前端构建 (Node.js 未安装)"
    fi
}

# 显示部署信息
show_info() {
    echo ""
    echo "=========================================="
    log_info "EasyBlog 部署完成！"
    echo "=========================================="
    echo ""
    echo "服务访问地址:"
    echo "  后端 API:   http://localhost:8080"
    echo "  健康检查：http://localhost:8080/healthz"
    echo "  前端服务：cd frontend && npm run dev"
    echo ""
    echo "默认账号:"
    echo "  管理员：admin / Admin123456"
    echo ""
    echo "日志目录：logs/"
    echo ""
    echo "常用命令:"
    echo "  查看日志：tail -f logs/apiserver.log"
    echo "  停止服务：make stop-local"
    echo "  重启服务：make deploy-local"
    echo "  停止 Docker: make docker-down"
    echo ""
    echo "=========================================="
}

# 主函数
main() {
    echo ""
    echo "=========================================="
    echo "EasyBlog 本地部署脚本"
    echo "=========================================="
    echo ""

    # 检查依赖
    check_dependencies

    # 启动基础设施
    start_infrastructure

    # 初始化数据库
    init_database

    # 构建后端
    build_backend

    # 启动后端服务
    start_backend

    # 构建前端（可选）
    build_frontend

    # 显示部署信息
    show_info
}

# 执行主函数
main
