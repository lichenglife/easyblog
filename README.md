# EasyBlog - 简单博客系统

[![Go Version](https://img.shields.io/badge/go-1.23.5-blue.svg)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/vue-3.4-green.svg)](https://vuejs.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

一个基于 Go + Vue3 的简单博客系统，支持文章管理、评论、点赞、标签分类等功能。

## 📋 目录

- [功能特性](#功能特性)
- [技术栈](#技术栈)
- [快速开始](#快速开始)
- [项目结构](#项目结构)
- [开发指南](#开发指南)
- [部署](#部署)

## 功能特性

### 用户模块
- [x] 用户注册/登录 (JWT 认证)
- [x] 用户名密码登录
- [ ] GitHub OAuth 登录
- [x] 用户信息管理
- [x] 头像上传

### 内容模块
- [x] Markdown 编辑器
- [x] 文章发布/编辑/删除
- [x] 文章草稿
- [x] 文章置顶
- [x] 封面图设置
- [x] 标签管理 (每篇文章最多 5 个标签)
- [x] 分类管理 (支持 3 级分类)
- [x] 图片上传 (本地存储)

### 互动模块
- [x] 评论系统
- [x] 评论回复 (楼中楼)
- [x] 点赞功能

### 管理模块
- [x] 用户管理 (审核/禁用/角色分配)
- [x] 文章管理 (审核/推荐/置顶)
- [x] 内容审核
- [x] 系统配置

## 技术栈

### 后端
- **语言**: Go 1.23.5
- **Web 框架**: Gin v1.10.x
- **ORM**: GORM v1.26.x
- **数据库**: MySQL 8.0+
- **缓存**: Redis 7.x
- **日志**: Zap
- **配置**: Viper
- **CLI**: Cobra
- **认证**: JWT (golang-jwt/jwt/v5)

### 前端
- **框架**: Vue 3.4.x
- **语言**: TypeScript 5.x
- **构建工具**: Vite 5.x
- **状态管理**: Pinia 2.x
- **路由**: Vue Router 4.x
- **HTTP**: Axios 1.x
- **UI 组件**: Element Plus 2.x
- **样式**: TailwindCSS 3.x
- **测试**: Vitest + Playwright

### 运维
- **容器**: Docker + Docker Compose
- **对象存储**: MinIO (兼容 S3 协议)
- **监控**: Prometheus + Grafana

## 快速开始

### 环境要求

- Go 1.23+
- Node.js 18+ (可选，用于前端构建)
- Docker + Docker Compose

### 1. 克隆项目

```bash
git clone https://github.com/lichenglife/easyblog.git
cd easyblog
```

### 2. 一键部署（推荐）

```bash
# 使用部署脚本
./scripts/deploy.sh
```

### 3. 分步部署

```bash
# 启动基础设施
make docker-up

# 初始化数据库
make db-init

# 构建并启动后端
make build
make run

# 启动前端（新终端）
make frontend-dev
```

### 4. 访问服务

- 后端 API: http://localhost:8080
- 健康检查：http://localhost:8080/healthz
- 前端服务：http://localhost:5173

### 5. 默认账号

- 管理员账号：`admin` / `Admin123456`

## 项目结构

```
easyblog/
├── cmd/                    # 命令行入口
│   └── apiserver/
│       ├── main.go         # 主入口
│       └── app/            # 应用初始化
├── internal/
│   ├── config/            # 配置加载
│   ├── handler/           # HTTP 处理器
│   ├── biz/               # 业务逻辑
│   ├── store/             # 数据访问层
│   ├── model/             # 数据模型
│   ├── middleware/        # 中间件
│   └── pkg/               # 内部工具包
│       ├── db/            # 数据库连接
│       ├── cache/         # 缓存连接
│       ├── jwt/           # JWT 工具
│       └── oss/           # 对象存储
├── pkg/                    # 可复用工具包
├── configs/               # 配置文件
│   └── config.yaml
├── scripts/
│   └── db/
│       └── init.sql       # 数据库初始化脚本
├── frontend/              # 前端项目
│   ├── src/
│   │   ├── api/          # API 调用
│   │   ├── stores/       # Pinia stores
│   │   ├── router/       # 路由配置
│   │   └── pages/        # 页面组件
│   └── package.json
├── data/                   # 数据目录 (本地存储)
│   ├── mysql/            # MySQL 数据
│   ├── redis/            # Redis 数据
│   └── minio/            # MinIO 数据
├── docker-compose.yml     # Docker 配置
├── Makefile               # 构建脚本
└── CLAUDE.md             # 开发指南
```

## 开发指南

### 后端开发

```bash
# 运行测试
make test

# 代码检查
make lint

# 编译
make build
```

### 前端开发

```bash
# 单元测试
cd frontend && npm run test:unit

# E2E 测试
cd frontend && npm run test:e2e

# 构建
cd frontend && npm run build
```

## API 文档

启动服务后访问：`http://localhost:8080/swagger/index.html`

## 部署

### Docker 部署

```bash
docker-compose -f docker-compose.prod.yml up -d
```

### Kubernetes 部署

详见 `deploy/k8s/` 目录配置

## 许可证

MIT License
