# 测试规范

## 1. 测试文件命名

### 1.1 文件命名

```
processor.go       → processor_test.go
receiver.go        → receiver_test.go
```

### 1.2 测试函数命名

使用 `TestXxx_State_Desc` 格式：

```go
func TestProcessor_Process_WithValidInput(t *testing.T)
func TestProcessor_Process_WithEmptyInput(t *testing.T)
func TestProcessor_Process_WithContextCanceled(t *testing.T)
```

---

## 2. 测试结构

### 2.1 表格驱动测试

```go
func TestProcessor_Process_WithValidInput(t *testing.T) {
    tests := []struct {
        name    string
        input   ptrace.Traces
        want    ptrace.Traces
        wantErr bool
    }{
        {
            name:  "valid input",
            input: createValidTraces(),
            want:  createExpectedTraces(),
        },
        {
            name:    "empty input",
            input:   ptrace.NewTraces(),
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 测试逻辑
        })
    }
}
```

### 2.2 测试辅助函数

```go
func createValidTraces() ptrace.Traces {
    traces := ptrace.NewTraces()
    rs := traces.ResourceSpans().AppendEmpty()
    // ... 初始化测试数据
    return traces
}
```

---

## 3. 测试覆盖率要求

| 组件类型 | 最低覆盖率 |
|---------|-----------|
| 核心业务逻辑 | ≥ 80% |
| 数据处理器 | ≥ 75% |
| 协议转换器 | ≥ 75% |
| 工具函数 | ≥ 70% |
| 整体项目 | ≥ 70% |

### 3.1 生成覆盖报告

```bash
# 生成覆盖率
make coverage

# 查看覆盖率详情
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 4. 测试类型

### 4.1 单元测试

```go
func TestCleaner_CleanSpan_WithNilInput(t *testing.T) {
    // Arrange
    cleaner := NewCleaner(Config{})

    // Act
    result, err := cleaner.CleanSpan(context.Background(), nil)

    // Assert
    assert.Error(t, err)
    assert.Nil(t, result)
}
```

### 4.2 集成测试

```go
func TestIntegration_KafkaProducerConsumer(t *testing.T) {
    // 跳过 CI 环境
    if testing.Short() {
        t.Skip("skipping integration test in short mode")
    }

    // 设置测试环境
    setupTestKafka(t)
    defer teardownTestKafka(t)

    // 执行测试
    // ...
}
```

### 4.3 基准测试

```go
func BenchmarkProcessor_Process(b *testing.B) {
    processor := NewProcessor(Config{})
    input := createBenchmarkInput()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := processor.Process(context.Background(), input)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

### 4.4 竞态检测

```bash
# 运行竞态检测
go test -race ./...
```

---

## 5. Mock 和 Stub

### 5.1 使用 testify/mock

```go
type MockConsumer struct {
    mock.Mock
}

func (m *MockConsumer) Consume(ctx context.Context, data ptrace.Traces) error {
    args := m.Called(ctx, data)
    return args.Error(0)
}

func TestProcessor_WithMockConsumer(t *testing.T) {
    mockConsumer := new(MockConsumer)
    mockConsumer.On("Consume", mock.Anything, mock.Anything).Return(nil)

    // 测试逻辑

    mockConsumer.AssertExpectations(t)
}
```

---

## 6. 测试最佳实践

### 6.1 测试独立性

每个测试应该是独立的，不依赖其他测试的状态：

```go
func TestXxx(t *testing.T) {
    // 每个测试创建自己的测试数据
    testData := createTestData()
    defer cleanupTestData(testData)

    // 测试逻辑
}
```

### 6.2 使用 t.Cleanup

```go
func TestXxx(t *testing.T) {
    resource := setupResource(t)
    t.Cleanup(func() {
        teardownResource(resource)
    })

    // 测试逻辑
}
```

### 6.3 并行测试

```go
func TestXxx(t *testing.T) {
    t.Parallel()

    // 测试逻辑（确保没有共享状态）
}
```

---

## 7. 验证命令

```bash
# 运行所有测试
make test

# 运行单个测试
go test -v -run TestName ./path/to/package

# 运行并生成覆盖率
go test -coverprofile=coverage.out ./...

# 检查竞态条件
go test -race ./...
```
