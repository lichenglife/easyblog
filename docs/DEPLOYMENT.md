# EasyBlog 部署指南

本文档描述如何将 EasyBlog 部署到本地或生产环境。

---

## 目录

- [环境要求](#环境要求)
- [本地部署](#本地部署)
- [生产环境部署](#生产环境部署)
- [配置说明](#配置说明)
- [常见问题](#常见问题)

---

## 环境要求

### 最低要求

| 组件 | 版本 | 说明 |
|------|------|------|
| Go | 1.23+ | 后端运行环境 |
| Node.js | 18+ | 前端构建环境（可选） |
| Docker | 20+ | 容器化部署 |
| MySQL | 8.0+ | 数据库 |
| Redis | 7+ | 缓存 |

### 推荐配置

- CPU: 2 核心
- 内存：4GB
- 磁盘：20GB

---

## 本地部署

### 方式一：使用部署脚本（推荐）

```bash
# 执行部署脚本
./scripts/deploy.sh
```

脚本会自动完成以下操作：
1. 检查依赖
2. 启动 Docker 服务（MySQL, Redis, MinIO）
3. 初始化数据库
4. 构建后端
5. 启动后端服务
6. 构建前端（如果安装了 Node.js）

### 方式二：使用 Makefile

```bash
# 1. 启动基础设施
make docker-up

# 2. 初始化数据库
make db-init

# 3. 构建后端
make build

# 4. 启动后端服务
make run

# 5. 启动前端（新终端）
make frontend-dev
```

### 访问服务

部署完成后：
- 后端 API: http://localhost:8080
- 健康检查：http://localhost:8080/healthz
- 前端服务：http://localhost:5173 (开发模式)

### 默认账号

- 管理员：`admin` / `Admin123456`

---

## 生产环境部署

### 1. 准备环境变量

```bash
# 设置生产环境变量
export JWT_SECRET="your-secret-key-here"
export DB_HOST="your-db-host"
export DB_PASSWORD="your-db-password"
export REDIS_HOST="your-redis-host"
export REDIS_PASSWORD="your-redis-password"
```

### 2. 构建生产版本

```bash
# 全量构建
make build-prod
```

### 3. 使用生产配置启动

```bash
# 使用生产配置文件
CONFIG_PATH=configs/config.prod.yaml ./bin/easyblog server
```

### 4. 后台运行（使用 systemd）

创建 systemd 服务文件 `/etc/systemd/system/easyblog.service`:

```ini
[Unit]
Description=EasyBlog Service
After=network.target mysql.service redis.service

[Service]
Type=simple
User=easyblog
WorkingDirectory=/opt/easyblog
Environment="CONFIG_PATH=/opt/easyblog/configs/config.prod.yaml"
ExecStart=/opt/easyblog/bin/easyblog server
Restart=on-failure
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable easyblog
sudo systemctl start easyblog
sudo systemctl status easyblog
```

---

## 配置说明

### 环境变量

生产环境推荐使用环境变量管理敏感配置：

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `JWT_SECRET` | JWT 密钥 | - |
| `DB_HOST` | 数据库主机 | 127.0.0.1 |
| `DB_PORT` | 数据库端口 | 3306 |
| `DB_USERNAME` | 数据库用户名 | root |
| `DB_PASSWORD` | 数据库密码 | - |
| `DB_DATABASE` | 数据库名称 | easyblog |
| `REDIS_HOST` | Redis 主机 | localhost |
| `REDIS_PORT` | Redis 端口 | 6379 |
| `REDIS_PASSWORD` | Redis 密码 | - |
| `FRONTEND_URL` | 前端地址 | http://localhost:5173 |
| `OSS_ENABLED` | 启用对象存储 | false |

### 配置文件

| 文件 | 用途 |
|------|------|
| `configs/config.yaml` | 开发环境配置 |
| `configs/config.prod.yaml` | 生产环境配置 |

---

## 日志管理

### 日志位置

- 应用日志：`logs/apiserver.log`
- 错误日志：`logs/error.log`
- 部署日志：`logs/deploy.log`

### 日志轮转

日志自动轮转配置：
- 单个文件最大：100MB
- 保留文件数：7
- 保留天数：30 天
- 压缩：启用

### 查看日志

```bash
# 实时查看日志
make logs

# 查看最近 100 行
tail -n 100 logs/apiserver.log

# 搜索错误日志
grep "ERROR" logs/apiserver.log
```

---

## 常见问题

### 1. Docker 容器启动失败

```bash
# 查看容器日志
docker logs easyblog-mysql
docker logs easyblog-redis

# 重启容器
docker-compose restart
```

### 2. 数据库连接失败

检查数据库配置：
```bash
# 测试数据库连接
docker exec -it easyblog-mysql mysql -uroot -proot123
```

### 3. 端口被占用

修改配置文件中的端口：
```yaml
server:
  http:
    port: 8081  # 修改端口
```

### 4. 前端跨域问题

确保配置文件中 CORS 设置正确：
```yaml
cors:
  allowedOrigins:
    - http://your-frontend-url.com
```

### 5. JWT Token 失效

检查 JWT 配置：
```yaml
jwt:
  expire: 7200  # Token 有效期（秒）
```

---

## 监控与告警

### 健康检查

```bash
curl http://localhost:8080/healthz
```

响应示例：
```json
{"code":0,"message":"success","data":{"status":"OK"}}
```

### 性能监控

推荐使用以下工具监控：
- Prometheus + Grafana：指标监控
- ELK Stack：日志分析
- Jaeger：链路追踪

---

## 备份与恢复

### 数据库备份

```bash
# 备份数据库
docker exec easyblog-mysql mysqldump -uroot -proot123 easyblog > backup.sql

# 恢复数据库
docker exec -i easyblog-mysql mysql -uroot -proot123 < backup.sql
```

### 日志备份

```bash
# 打包日志
tar -czf logs-backup-$(date +%Y%m%d).tar.gz logs/
```

---

## 更新升级

```bash
# 1. 停止服务
make stop-local

# 2. 拉取最新代码
git pull

# 3. 重新构建
make build-prod

# 4. 启动服务
make deploy-local
```

---

## 技术支持

如有问题，请提交 Issue 或联系开发团队。
