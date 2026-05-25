## Context

tryssh 是一个 Go 编写的 SSH 终端工具，支持密码猜测、凭证缓存、SCP 文件传输和服务器别名管理。当前技术栈：

- Go 1.23.0（本机为 1.25.10）
- Cobra CLI、Logrus、SFTP、go-cartesian-product
- YAML 配置存储于 `~/.tryssh/tryssh.db`
- 零测试覆盖率
- 单人维护，dev 分支开发

架构分层：`cmd/`（Cobra 命令）→ `pkg/control/`（业务逻辑）→ `pkg/launcher/`（SSH/SFTP 连接）→ `pkg/config/`（配置管理）+ `pkg/utils/`（工具函数）

## Goals / Non-Goals

**Goals:**
- 将 Go 版本升级到 1.25，所有依赖更新到最新兼容版本
- 重构设计不佳的部分，提升代码质量和可维护性
- 建立完整的单元测试体系，覆盖率 100%
- 通过代码审计发现并修复安全和质量问题
- 完善项目文档（README、架构说明、构建文档）
- 更新 CI/CD workflow 适配新工具链

**Non-Goals:**
- 不改变项目的核心功能和行为
- 不增加新功能特性
- 不更换主要依赖库（如 Cobra、Logrus）
- 不改变配置文件格式和存储位置
- 不改变命令行接口（CLI flags/commands）

## Decisions

### D1: Go 版本升级策略

**决策**: 直接升级到 Go 1.25.10，更新 go.mod 中的 `go` 指令和 `toolchain` 指令。

**理由**: 本机已安装 Go 1.25.10，项目代码量不大，升级风险可控。

**替代方案**: 逐步升级（1.23→1.24→1.25）— 不必要，增加工作量。

### D2: 依赖更新策略

**决策**: 使用 `go get -u` 更新所有依赖到最新 minor/patch 版本，运行 `go mod tidy` 清理。

**理由**: 保持依赖新鲜度，减少已知漏洞。主要依赖（Cobra、SFTP、crypto）向后兼容。

### D3: 重构重点区域

**决策**: 针对以下区域进行重构：

1. **错误处理**: 当前大量使用 `log.Fatalf`，应改为返回 error 并在顶层统一处理
2. **接口抽象**: `pkg/launcher/` 中的连接逻辑缺少接口定义，不利于测试 mock
3. **并发控制**: `pkg/control/control.go` 的并发模型可简化
4. **配置管理**: `pkg/config/` 职责过重，应拆分为配置加载、凭证组合、缓存管理

**理由**: 这些区域是测试覆盖率的最大障碍——紧耦合和全局状态使得 mock 困难。

**替代方案**: 仅添加测试不重构 — 会导致测试代码被迫依赖文件系统和网络。

### D4: 测试架构

**决策**:
- 使用标准 `testing` 包 + `testify` 断言库
- 接口层 mock 使用 Go 1.25 的 mock 机制或手写 mock（避免引入重框架）
- 使用 `t.Setenv` 和临时目录处理文件依赖
- SSH 连接使用接口 mock，不依赖真实网络
- 覆盖率目标 100%，通过 `go test -coverprofile` 验证
- 在 CI 中强制覆盖率门禁

**理由**: 标准库 + testify 是 Go 社区最主流的测试方案，手写 mock 比框架更适合小项目。

**替代方案**: 使用 mockgen 等代码生成工具 — 对本项目规模来说过重。

### D5: 代码审计方案

**决策**:
- 使用 `golangci-lint` 作为主要静态分析工具，启用所有推荐 linter
- 使用 `govulncheck` 检查已知漏洞
- 手动审计安全敏感代码（SSH 密钥处理、密码存储）

**理由**: golangci-lint 是 Go 社区标准工具，配置灵活，覆盖面广。

### D6: 文档改进范围

**决策**:
- 更新 README.md 和 README_zh.md 反映最新构建和测试方法
- 更新 docs/build.md 添加测试说明
- 添加 CONTRIBUTING.md 说明开发流程
- 代码内不添加冗余注释，仅在必要处保留

**替代方案**: 生成 API 文档（godoc）— 对 CLI 工具价值有限。

### D7: CI/CD 更新

**决策**:
- 更新 goreleaser workflow 中的 Go 版本
- 添加 PR 检查 workflow：lint + test + 覆盖率检查
- 使用 `golangci-lint-action` 和 Go 官方测试 action

**理由**: 现有 workflow 仅处理发布，缺少质量门禁。PR 检查防止问题合入。

## Risks / Trade-offs

- **重构范围**: 大范围重构可能引入回归 → 通过 100% 测试覆盖率来捕获回归
- **Go 1.25 兼容性**: 部分依赖可能不兼容 Go 1.25 → 逐个验证，必要时 pin 版本
- **100% 覆盖率**: 某些路径（如 signal 处理）难以测试 → 使用构建标签或集成测试补充
- **破坏性变更**: 重构可能导致外部行为变化 → 严格保持 CLI 接口不变
