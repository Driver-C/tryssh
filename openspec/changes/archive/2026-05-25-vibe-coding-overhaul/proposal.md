## Why

项目自创建以来积累了技术债务：Go 版本落后（1.23 vs 本机 1.25）、依赖版本陈旧、零测试覆盖率、缺少安全审计。作为 vibe coding 全面接管的起点，需要一次系统性 overhaul 来建立现代化的工程实践基础。

## What Changes

- **BREAKING**: Go 版本从 1.23 升级到 1.25，更新 go.mod 和所有依赖包到最新兼容版本
- 全面审查项目架构，对设计不佳的部分进行重构（如配置管理、并发控制、错误处理）
- 建立单元测试规范，补充所有包的测试用例，目标覆盖率 100%
- 执行代码审计，修复发现的安全问题、代码缺陷和设计缺陷
- 完善项目文档（README、代码注释、API 文档、架构说明）
- 更新 GitHub Actions workflow 适配新的 Go 版本和测试要求

## Capabilities

### New Capabilities
- `testing-framework`: 单元测试规范、测试工具函数、mock 基础设施、覆盖率报告配置
- `code-audit`: 静态分析规则、安全审计清单、代码质量门禁

### Modified Capabilities
<!-- 无已有 spec 需要修改 -->

## Impact

- **代码**: 所有 `pkg/` 和 `cmd/` 包可能受重构影响
- **依赖**: go.mod 中所有依赖可能需要升级
- **CI/CD**: `.github/workflows/goreleaser.yml` 和 `.goreleaser.yml` 需要更新
- **文档**: README.md、README_zh.md、docs/build.md 需要更新
- **构建**: Makefile 构建目标可能需要调整
