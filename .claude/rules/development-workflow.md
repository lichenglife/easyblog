# 开发工作流

## 1. 核心原则

### 1.1 规划优先 (Plan First)

在任何代码修改之前，必须先进行规划：

1. **理解需求** - 仔细阅读用户请求，明确目标和约束
2. **探索代码库** - 使用 Glob/Grep 搜索相关代码，理解现有实现
3. **制定计划** - 对于复杂任务，使用 `EnterPlanMode` 制定详细实现方案
4. **获取确认** - 在实施重大变更前，向用户展示计划并获取批准

### 1.2 任务拆分 (Task Decomposition)

复杂任务必须拆分为可管理的子任务：

- 使用 `TaskCreate` 创建结构化任务列表
- 每个子任务应该是独立、可测试的
- 明确任务依赖关系 (`blockedBy` / `blocks`)
- 按 ID 顺序处理任务（ID 小的优先）

---

## 2. 开发流程

### 2.1 功能开发流程

```yaml
1. 创建特性分支：git checkout -b feature/feature-name
2. 编写代码：遵循编码标准
3. 运行测试：make test
4. 运行 lint: make lint
5. 检查覆盖：make coverage
6. 提交代码：conventional commits
7. 创建 PR: /pr-create
```

### 2.2 Bug 修复流程

```yaml
1. 创建修复分支：git checkout -b fix/bug-description
2. 编写失败测试
3. 验证测试失败
4. 实施修复
5. 验证测试通过
6. 检查覆盖
7. 提交并创建 PR
```

### 2.3 代码审查流程

```yaml
1. 自动检查：CI 通过
2. 人工审查：团队成员审查
3. 修改完善：根据反馈修改
4. 合并代码：合并到 main 分支
```

---

## 3. 验证标准

所有代码变更必须通过以下验证：

| 验证类型 | 命令 | 标准 |
|---------|------|------|
| 编译检查 | `make build` | 无编译错误 |
| 代码格式 | `make fmt` | gofmt 检查通过 |
| 静态分析 | `make vet` | go vet 检查通过 |
| 代码审查 | `make lint` | golangci-lint 通过 |
| 单元测试 | `make test` | 所有测试通过 |
| 覆盖率 | `make coverage` | ≥ 70% |
| CI 检查 | `/ci-check` | 所有 CI 检查通过 |

---

## 4. 提交前检查清单

- [ ] 代码已格式化 (`make fmt`)
- [ ] 通过 go vet (`make vet`)
- [ ] 通过 linter (`make lint`)
- [ ] 所有测试通过 (`make test`)
- [ ] 测试覆盖率 ≥ 70%
- [ ] 无敏感信息泄露
- [ ] 无 TODO 注释遗漏
- [ ] 文档已更新
- [ ] 变更已自测

---

## 5. 任务管理

### 5.1 任务创建

```go
TaskCreate(
    subject: "实现 XX 功能",
    description: "详细描述功能需求和验收标准",
    activeForm: "正在实现 XX 功能"
)
```

### 5.2 任务状态

| 状态 | 说明 |
|------|------|
| `pending` | 待处理 |
| `in_progress` | 进行中 |
| `completed` | 已完成 |
| `deleted` | 已删除 |

### 5.3 任务依赖

```
任务 1 → 任务 2 → 任务 3
  ↓        ↓        ↓
blocks:  blocks:  blocks: []
 [2]      [3]
```

---

## 6. Bug 修复流程

```
1. 重现问题
   ↓
2. 定位根因（日志、堆栈、代码分析）
   ↓
3. 编写失败测试（验证 BUG 存在）
   ↓
4. 实施修复
   ↓
5. 验证测试通过
   ↓
6. 回归测试（确保无新 BUG）
   ↓
7. 提交修复
```

**修复原则：**
- 优先修复根因，而非表面症状
- 添加测试用例防止回归
- 小步修改，频繁验证
- 如遇到阻塞，及时向用户反馈

---

## 7. Git 规范

### 7.1 分支命名

| 类型 | 命名格式 | 示例 |
|------|---------|------|
| 功能分支 | `feature/xxx` | `feature/kafka-consumer` |
| 修复分支 | `fix/xxx` | `fix/memory-leak` |
| 实验分支 | `experiment/xxx` | `experiment/new-serializer` |

### 7.2 Conventional Commits

```
feat:     新功能
fix:      Bug 修复
docs:     文档更新
style:    代码格式（不影响功能）
refactor: 重构
test:     测试相关
chore:    构建/工具/配置
```

**提交示例：**

```
feat: add SkyWalking v3 protocol support
fix: resolve Kafka consumer rebalance issue
docs: update architecture diagram
refactor: extract converter logic to separate package
test: add unit tests for data cleaning processor
```

---

## 8. 常用命令

### 8.1 开发命令

```bash
make build           # 编译
make run-collector   # 运行 collector
make run-exporter    # 运行 exporter
make test            # 测试
make coverage        # 覆盖率
make fmt             # 格式化
make lint            # Lint
make clean           # 清理
```

### 8.2 Docker 命令

```bash
docker-compose up -d              # 启动
docker-compose down               # 停止
docker-compose logs -f collector  # 查看日志
```

### 8.3 Kubernetes 命令

```bash
kubectl apply -f config/k8s/      # 部署
kubectl get pods -n etap          # 查看状态
kubectl logs -n etap <pod-name>   # 查看日志
```
