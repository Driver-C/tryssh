## Context

tryssh 的配置文件中存储 SSH 密码等敏感字段。当前已有 AES-GCM 加密基础设施（`pkg/utils/crypto.go`），配置加载/保存时自动处理加解密。但启用加密的唯一入口是环境变量 `TRYSSH_MASTER_KEY`，普通 CLI 用户无法方便地启用加密。

最近的代码修复将 `encryptConfigForSave` 改为使用 `GetCachedMasterKey()`（不提示密码），导致纯明文配置的用户没有途径触发首次加密。

## Goals / Non-Goals

**Goals:**
- 提供显式的 `tryssh encrypt` CLI 命令，让用户交互式输入 master password 并加密当前配置
- 加密后配置文件中的密码字段以 `enc:` 前缀存储
- 后续 `tryssh` 任何操作加载配置时，检测到 `enc:` 内容自动提示解密
- 保持 `TRYSSH_MASTER_KEY` 环境变量作为无交互入口（CI/脚本场景）

**Non-Goals:**
- 不实现 `tryssh decrypt` 命令（用户可通过删除 `enc:` 前缀或重新创建配置来回到明文）
- 不修改现有加解密算法或密钥派生方式
- 不实现密钥轮换（更换 master password）

## Decisions

1. **新增 `cmd/encrypt/` 包**：遵循项目现有结构（`cmd/ssh/`、`cmd/scp/` 等），`encrypt` 作为独立子命令
2. **命令流程**：提示输入密码 → 确认密码（二次输入）→ 加载配置 → 加密 → 写回 → 清除内存中的密钥
3. **利用现有函数**：`GetMasterKey()` 处理交互输入和密钥派生，`encryptConfigForSave()` 处理加密，`UpdateConfig()` 处理写回。只需在 `encrypt` 命令中将密钥注入缓存，其余复用
4. **密码确认**：通过二次输入确认密码，避免因输入错误导致配置无法解密。对比两个输入是否一致
5. **幂等性**：已加密的字段（`enc:` 前缀）不会被重复加密，直接跳过

## Risks / Trade-offs

- **用户忘记密码** → 配置文件中的密码无法恢复。命令执行后输出提示信息提醒用户妥善保管密码
- **无密码强度校验** → 仅要求最少 4 个字符。当前 `deriveKey` 已有此限制，保持一致
- **不在 encrypt 命令中实现解密** → 用户如需回到明文需手动编辑配置文件或删除重建。降低实现复杂度
