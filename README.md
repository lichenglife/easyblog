# Easyblog

Easyblog is a simple blog application that allows users to create, read, update, and delete blog posts.

## project architecture

项目采用了分层架构，分为以下几层：

1. **表示层**：API接口层和前端页面层。
2. **业务层**：业务逻辑层，包括用户管理、文章管理、评论管理等。
3. **数据访问层**：数据访问层，包括数据库操作和缓存操作。
4. **基础设施层**：基础设施层，包括日志、监控、缓存等。

### 系统组件

- **HTTP服务**： 使用gin框架作为HTTP服务，提供API接口。
- **Grpc服务**： 使用gRPC框架作为RPC服务，提供高性能的接口服务
- **认证服务**： 使用JWT作为认证服务，提供用户认证和授权。
- **用户服务**： 提供用户管理和认证功能。
- **博客服务**： 提供博客管理和发布功能。
- **缓存**： 使用Redis作为缓存，提高系统性能。
- **数据库**： 使用MySQL作为数据库，存储用户和博客数据。
- **前端**：使用Vue3作为前端框架，提供用户界面。

## 项目目录结构
```
easyblog/
├── api/                  # API定义
│   ├── openapi/          # OpenAPI规范
│   ├── proto/            # gRPC协议文件
│   └── swagger/          # Swagger文档
├── build/                # 构建相关文件
├── cmd/                  # 命令行入口
│   ├── apiserver/        # API服务器
│   └── tools/            # 辅助工具
├── configs/              # 配置文件
├── deployments/          # 部署文件
├── docs/                 # 文档
│   └── architecture/     # 架构文档
├── examples/             # 示例代码
├── frontend/             # 前端项目 (Vue 3)
├── init/                 # 初始化脚本
├── internal/             # 内部代码
│   ├── apiserver/        # API服务器
│   ├── authz/            # 认证授权
│   ├── pkg/              # 内部包
│   └── service/          # 业务服务
├── pkg/                  # 公共包
│   ├── auth/             # 认证
│   ├── config/           # 配置
│   ├── db/               # 数据库
│   ├── errors/           # 错误处理
│   ├── log/              # 日志
│   └── util/             # 工具函数
├── scripts/              # 脚本
└── third_party/          # 第三方代码
```

## 核心功能
1. **命令行框架**：使用cobra作为命令行框架，提供命令行工具。
2. **配置管理**：使用viper作为配置管理，支持多种配置格式。
3. **日志管理**：使用zap作为日志管理，支持多种日志格式。
4. **数据库操作**：使用gorm作为数据库操作，支持多种数据库。
5. **缓存操作**：使用redis作为缓存操作，支持多种缓存。
6. **认证授权**：使用jwt作为认证授权，支持多种认证方式。
7. **API文档**：使用swagger作为API文档，支持多种API文档格式。
8. **API服务器**：使用gin作为API服务器，提供API接口。
9. **gRPC服务器**：使用gRPC作为RPC服务器，提供高性能的接口服务。
10. **服务治理**： 使用监控检查、限流、熔断等服务治理功能。
11. **前端**：使用Vue3作为前端框架，提供用户界面。

## 项目启动流程

### 后端

1. 加载命令行参数
2. 初始化配置（从配置文件、环境变量）
3. 初始化日志系统
4. 初始化数据库连接
5. 初始化缓存连接
6. 初始化服务实例
7. 初始化HTTP服务器
8. 注册中间件和路由
9. 启动HTTP服务
10. 监听系统信号，实现优雅关闭

### 前端
1.安装依赖：` cd frontend && npm install`
2.开发模式：`npm run server`
3.打包构建：`npm run build`