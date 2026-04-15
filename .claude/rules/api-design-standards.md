# API 设计规范

## 1. gRPC 服务设计

### 1.1 服务定义

使用 Protocol Buffers 定义清晰的 gRPC 服务接口：

```protobuf
service SkyWalkingReceiver {
  // Collect traces from SkyWalking agents
  rpc Collect(CollectRequest) returns (CollectResponse);

  // Collect metrics
  rpc CollectMetrics(MetricsRequest) returns (MetricsResponse);

  // Collect logs
  rpc CollectLogs(LogsRequest) returns (LogsResponse);
}
```

### 1.2 消息命名

- 请求消息：`XxxRequest`
- 响应消息：`XxxResponse`
- 事件消息：`XxxEvent`

---

## 2. HTTP API 设计

### 2.1 RESTful 风格

```
GET    /api/v1/traces      # 查询 traces
POST   /api/v1/traces      # 创建 trace
GET    /api/v1/traces/{id} # 获取单个 trace
DELETE /api/v1/traces/{id} # 删除 trace
```

### 2.2 响应格式

```go
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}
```

### 2.3 错误响应

```go
type ErrorResponse struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}
```

---

## 3. 内部接口设计

### 3.1 组件接口

```go
// Receiver 数据接收器接口
type Receiver interface {
    Start(ctx context.Context) error
    Shutdown(ctx context.Context) error
}

// Processor 数据处理器接口
type Processor interface {
    ProcessTraces(ctx context.Context, td ptrace.Traces) (ptrace.Traces, error)
    ProcessMetrics(ctx context.Context, md pmetric.Metrics) (pmetric.Metrics, error)
    ProcessLogs(ctx context.Context, ld plog.Logs) (plog.Logs, error)
}

// Exporter 数据导出器接口
type Exporter interface {
    Export(ctx context.Context, data interface{}) error
    Shutdown(ctx context.Context) error
}
```

### 3.2 工厂模式

```go
// ReceiverFactory 接收器工厂
type ReceiverFactory interface {
    CreateDefaultConfig() Config
    CreateReceiver(cfg Config, set Settings, nextConsumer Consumer) (Receiver, error)
}
```

---

## 4. 参数验证

### 4.1 输入验证

```go
func ValidateConfig(cfg Config) error {
    if cfg.Endpoint == "" {
        return errors.New("endpoint is required")
    }
    if cfg.Port < 1 || cfg.Port > 65535 {
        return errors.New("port must be between 1 and 65535")
    }
    if cfg.Timeout <= 0 {
        return errors.New("timeout must be positive")
    }
    return nil
}
```

### 4.2 边界检查

```go
func ProcessBatch(data []Span, maxSize int) error {
    if len(data) == 0 {
        return errors.New("batch cannot be empty")
    }
    if len(data) > maxSize {
        return fmt.Errorf("batch size %d exceeds maximum %d", len(data), maxSize)
    }
    return nil
}
```

---

## 5. 版本控制

### 5.1 API 版本

```
/api/v1/...
/api/v2/...
```

### 5.2 向后兼容

- 不删除已有字段
- 新增字段设为可选
- 使用 `deprecated` 标记废弃字段

---

## 6. 文档要求

### 6.1 API 文档

每个公开 API 必须有文档说明：

```go
// Start starts the SkyWalking receiver.
//
// The receiver will listen for incoming SkyWalking v3 protocol data
// on the configured endpoint and port.
//
// Parameters:
//   - ctx: Context for cancellation
//
// Returns:
//   - error: Error if the receiver fails to start
func (r *Receiver) Start(ctx context.Context) error
```
