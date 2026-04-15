---
name: 代码质量要求
description: 项目代码质量标准和测试要求
type: feedback
---

**所有代码必须有单元测试**

- Handler 层：覆盖率 ≥ 80%
- Biz 层：覆盖率 ≥ 80%
- Store 层：覆盖率 ≥ 80%
- Model 层：覆盖率 ≥ 70%
- Middleware：覆盖率 ≥ 80%
- 工具函数：覆盖率 ≥ 90%

**测试规范**：
- 测试文件命名：`xxx_test.go`
- 测试函数命名：`TestXxx_Function_Scenario`
- 多场景使用表驱动测试
- 使用 mock 隔离外部依赖

**Why**: 用户明确要求所有代码必须有单元测试，这是项目的硬性约束。

**How to apply**: 实现任何功能时，必须同时编写对应的单元测试文件，否则视为未完成。
