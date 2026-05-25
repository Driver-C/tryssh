## Why

当前配置加密功能只能通过环境变量 `TRYSSH_MASTER_KEY` 启用，普通用户难以发现这个入口。最近的代码修改修复了"纯明文配置也会提示输入密码"的问题，但也导致纯明文用户没有任何途径触发加密。需要一个显式的 CLI 命令让用户主动启用配置加密。

## What Changes

- 修复密码提示逻辑：纯明文配置不再弹出密码提示，有加密内容但无密钥时报错而非静默跳过（已实现）
- 新增 `GetCachedMasterKey()` 函数：只返回已缓存或环境变量中的密钥，不交互提示，用于保存时按需加密（已实现）
- 新增 `tryssh encrypt` 子命令，交互式输入 master password 后加密当前配置中的所有密码字段
- 加密完成后提示用户记住密码，后续加载配置时会自动提示输入
- 环境变量 `TRYSSH_MASTER_KEY` 保留作为补充入口（CI/自动化场景）
- 不引入 **BREAKING** 变更：已有加密配置的加载/解密行为不变

## Capabilities

### New Capabilities
- `encrypt-command`: `tryssh encrypt` CLI 子命令，负责交互式输入 master password、加密配置文件中的明文密码字段、并将加密后的配置写回磁盘

### Modified Capabilities
（无 — 配置加载/解密/保存的行为已在之前的修复中完成，本次只新增入口命令）

## Impact

- `pkg/utils/crypto.go`：新增 `GetCachedMasterKey()` 函数（已修改）
- `pkg/config/loader.go`：`decryptConfig` 预扫描加密内容、`encryptConfigForSave` 改用 `GetCachedMasterKey`（已修改）
- `cmd/encrypt/` 目录：新增子命令注册
- 用户文档需补充加密功能说明
