# 代码规范

## 1. Go 代码质量标准

### 1.1 包注释

每个包必须有文档注释：

```go
// Package skywalking implements the SkyWalking v3 protocol receiver.
//
// This receiver accepts trace, metrics, and logs data from SkyWalking agents
// and converts them to OTLP format for downstream processing.
package skywalking
```

### 1.2 导出函数注释

所有导出函数必须有 godoc 注释：

```go
// NewReceiver creates a new SkyWalking protocol receiver.
//
// Parameters:
//   - cfg: Receiver configuration
//   - set: Collector settings
//   - nextConsumer: Next pipeline component
//
// Returns:
//   - *Receiver: The receiver instance
//   - error: Error if creation fails
func NewReceiver(cfg Config, set Settings, nextConsumer Consumer) (*Receiver, error)
```

### 1.3 错误处理

必须检查并包装错误：

```go
result, err := someFunction()
if err != nil {
    return fmt.Errorf("failed to process span %s: %w", spanID, err)
}
```

### 1.4 资源管理

必须使用 defer 清理资源：

```go
func processFile(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer f.Close()
    // ... 使用文件
}
```

### 1.5 上下文取消

必须支持上下文取消：

```go
func (p *Processor) Process(ctx context.Context, data interface{}) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        // 处理数据
    }
}
```

### 1.6 并发安全

必须考虑并发场景并使用适当的同步机制：

```go
type SafeBuffer struct {
    mu   sync.Mutex
    data []byte
}

func (b *SafeBuffer) Write(p []byte) (n int, err error) {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.data = append(b.data, p...)
    return len(p), nil
}
```

---

## 2. 命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 包名 | `lower_snake_case` | `skywalking`, `otlp_kafka` |
| 导出函数 | `CamelCase` | `NewReceiver`, `ProcessTraces` |
| 私有函数 | `camelCase` | `processSpan`, `convertMetric` |
| 导出常量 | `UpperCamelCase` | `DefaultBufferSize`, `MaxRetries` |
| 私有常量 | `camelCase` | `defaultTimeout`, `maxConnections` |
| 接口 | `-er` 后缀 | `Processor`, `Converter`, `Writer` |
| 测试函数 | `TestXxx_State_Desc` | `TestProcessor_Process_WithEmptyInput` |

---

## 3. 日志规范

使用 zap 结构化日志：

```go
logger.Info("receiver started",
    zap.String("endpoint", cfg.Endpoint),
    zap.Int("port", cfg.Port),
)

logger.Error("failed to process span",
    zap.String("span_id", spanID),
    zap.Error(err),
)

logger.Debug("processing batch",
    zap.Int("size", batchSize),
    zap.Duration("duration", elapsed),
)
```

### 日志级别使用

| 级别 | 使用场景 |
|------|---------|
| `DEBUG` | 调试信息，开发环境使用 |
| `INFO` | 正常业务日志 |
| `WARN` | 警告信息，不影响业务 |
| `ERROR` | 错误信息，需要处理 |

---

## 4. 配置规范

### 4.1 环境变量前缀

```bash
# Collector 配置
ETAP_COLLECTOR_ENDPOINT=0.0.0.0:11800
ETAP_COLLECTOR_SW_TOKEN=xxx

# Exporter 配置
ETAP_EXPORTER_KAFKA_BROKERS=kafka-1:9092,kafka-2:9092
ETAP_EXPORTER_ES_ENDPOINTS=http://es-1:9200,http://es-2:9200
```

### 4.2 配置文件命名

```
config/
├── collector.yaml          # Collector 配置
├── collector-dev.yaml      # 开发环境配置
├── exporter.yaml           # Exporter 配置
└── k8s/                    # Kubernetes 配置
    ├── deployment.yaml
    ├── service.yaml
    └── configmap.yaml
```

---

## 5. 目录结构规范

待完善
