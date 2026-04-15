# 安全规范

## 1. 输入验证

### 1.1 所有外部输入必须验证

```go
func ValidateEndpoint(endpoint string) error {
    if endpoint == "" {
        return errors.New("endpoint cannot be empty")
    }

    // 验证格式
    host, port, err := net.SplitHostPort(endpoint)
    if err != nil {
        return fmt.Errorf("invalid endpoint format: %w", err)
    }

    // 验证端口范围
    portNum, err := strconv.Atoi(port)
    if err != nil || portNum < 1 || portNum > 65535 {
        return errors.New("port must be between 1 and 65535")
    }

    return nil
}
```

### 1.2 防止注入攻击

```go
// 使用参数化查询，避免 SQL 注入
// 使用 zap.SugaredLogger 避免日志注入
logger.Info("user login",
    zap.String("username", sanitizeInput(username)),
    zap.String("ip", clientIP),
)
```

---

## 2. 凭据管理

### 2.1 禁止硬编码凭据

```go
// 错误示例 - 禁止这样做
const apiKey = "sk-1234567890"

// 正确示例 - 从环境变量读取
apiKey := os.Getenv("ETAP_API_KEY")
if apiKey == "" {
    return errors.New("ETAP_API_KEY environment variable is required")
}
```

### 2.2 敏感信息脱敏

```go
func maskPassword(password string) string {
    if len(password) <= 4 {
        return "****"
    }
    return password[:2] + strings.Repeat("*", len(password)-4) + password[len(password)-2:]
}

// 日志中脱敏
logger.Info("connecting to database",
    zap.String("host", cfg.Host),
    zap.String("username", cfg.Username),
    zap.String("password", maskPassword(cfg.Password)),
)
```

---

## 3. 依赖安全

### 3.1 定期更新依赖

```bash
# 检查过时依赖
go list -u -m all

# 更新依赖
go get -u ./...

# 清理依赖
go mod tidy
```

### 3.2 检查安全漏洞

```bash
# 使用 govulncheck
govulncheck ./...

# 使用 nancy
cat go.sum | nancy sleuth
```

---

## 4. 并发安全

### 4.1 使用同步原语

```go
type SafeCache struct {
    mu     sync.RWMutex
    cache  map[string]interface{}
}

func (c *SafeCache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    val, ok := c.cache[key]
    return val, ok
}

func (c *SafeCache) Set(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.cache[key] = value
}
```

### 4.2 避免 Goroutine 泄漏

```go
func ProcessWithContext(ctx context.Context, dataChan <-chan Data) error {
    errChan := make(chan error, 1)

    go func() {
        defer close(errChan)
        for {
            select {
            case <-ctx.Done():
                errChan <- ctx.Err()
                return
            case data, ok := <-dataChan:
                if !ok {
                    errChan <- nil
                    return
                }
                // 处理数据
            }
        }
    }()

    return <-errChan
}
```

---

## 5. 资源管理

### 5.1 及时释放资源

```go
func ReadFile(path string) ([]byte, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer f.Close()

    return io.ReadAll(f)
}
```

### 5.2 限制资源使用

```go
// 限制 HTTP 请求体大小
func withMaxBytes(next http.Handler, maxBytes int64) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
        next.ServeHTTP(w, r)
    })
}
```

---

## 6. 错误处理

### 6.1 不泄露敏感信息

```go
// 错误示例 - 可能泄露敏感信息
return fmt.Errorf("failed to connect to %s:%d with password %s", host, port, password)

// 正确示例
if authErr := connect(host, port, password); authErr != nil {
    return errors.New("authentication failed")
}
```

### 6.2 包装错误保留上下文

```go
result, err := process(data)
if err != nil {
    return fmt.Errorf("processing failed for item %s: %w", itemID, err)
}
```

---

## 7. 安全检查清单

提交前确认：

- [ ] 无硬编码凭据
- [ ] 敏感信息已脱敏
- [ ] 所有输入已验证
- [ ] 无 SQL/命令注入风险
- [ ] 并发访问已加锁
- [ ] 资源已正确释放
- [ ] 依赖无已知漏洞
- [ ] 错误信息不泄露敏感数据

---

## 8. 安全工具

```bash
# 静态分析
golangci-lint run --enable=gosec ./...

# 漏洞扫描
govulncheck ./...

# 依赖检查
go list -m -versions all | grep -E "(vulnerability|security)"
```
