## 0. 密码提示逻辑修复（已完成）

- [x] 0.1 `pkg/utils/crypto.go` 新增 `GetCachedMasterKey()` 函数：只在密钥已缓存或环境变量存在时返回，不触发交互提示
- [x] 0.2 `pkg/config/loader.go` `decryptConfig` 增加预扫描：先检测是否存在 `enc:` 前缀加密内容，纯明文配置跳过密码提示，有加密内容但无密钥时报错
- [x] 0.3 `pkg/config/loader.go` `encryptConfigForSave` 改用 `GetCachedMasterKey()` 替代 `GetMasterKey()`，保存时不再触发交互提示
- [x] 0.4 补充测试：`TestDecryptConfig_EncryptedButNoKey`、`TestGetCachedMasterKey_*` 系列测试
- [x] 0.5 验证全部测试和 lint 通过

## 1. 命令注册

- [x] 1.1 创建 `cmd/encrypt/encrypt.go`，实现 `NewEncryptCommand()` 返回 cobra.Command，包含 Use/Short/Long 描述
- [x] 1.2 在 `cmd/cmd.go` 中注册 `encrypt` 子命令，更新 `cmd/cmd_test.go` 子命令列表和计数

## 2. 核心逻辑

- [x] 2.1 实现交互式密码输入流程：`readPasswordFn` 可注入函数封装 `term.ReadPassword`，提示 "Enter master password" 和 "Confirm master password"，不回显
- [x] 2.2 实现密码验证：空输入报错退出、两次不一致报错退出、长度不足 4 报错退出，错误通过 `fatalFn` 可注入
- [x] 2.3 实现环境变量优先逻辑：`TRYSSH_MASTER_KEY` 已设置时跳过交互输入直接使用
- [x] 2.4 实现加密主流程：加载配置 → 通过 `GetMasterKey()` 获取密钥 → `encryptConfigForSave()` 加密 → `UpdateConfig()` 写回 → `ClearMasterKey()` 清理
- [x] 2.5 处理无密码可加密的边界情况：输出 "No passwords to encrypt" 提示

## 3. 测试

- [x] 3.1 命令结构测试：`TestNewEncryptCommand_Structure`（Use/Short/Long）、`TestNewEncryptCommand_Run`（通过 cobra Run 闭包执行）
- [x] 3.2 `countPlaintextPasswords` 测试：全明文、混合、空、仅空密码、全加密 — 5 个场景
- [x] 3.3 环境变量路径测试：`TestRunEncrypt_EnvVar_NoPasswords`（无密码）、`TestRunEncrypt_EnvVar_WithPasswords`（有密码 + 验证文件中 enc: 前缀）
- [x] 3.4 交互式路径测试（注入 `readPasswordFn`）：空密码、过短、不匹配、首次读取错误、确认读取错误、成功加密 — 6 个场景
- [x] 3.5 `executeEncrypt` 错误路径：`TestExecuteEncrypt_LoadConfigError`（不可读路径）、`TestExecuteEncrypt_MixedPlaintextAndEncrypted`（混合加密验证）
- [x] 3.6 覆盖率验证：`cmd/encrypt` 从 19.6% 提升到 85.0%，总体 85.3%，lint 0 issues

## 4. 文档

- [x] 4.1 更新 `README.md`：命令列表添加 encrypt、安全特性更新为 `tryssh encrypt` 主入口 + 环境变量补充、encrypt 命令详解和示例
- [x] 4.2 更新 `README_zh.md`：同步中文文档（命令列表、安全特性、encrypt 命令详解）
- [x] 4.3 确认总体覆盖率 ≥ 80%（当前 85.3%）
