## ADDED Requirements

### Requirement: tryssh encrypt command registration
系统 SHALL 注册 `encrypt` 子命令到根命令，用户可通过 `tryssh encrypt` 调用。

#### Scenario: encrypt 命令出现在帮助列表中
- **WHEN** 用户运行 `tryssh --help`
- **THEN** 输出中包含 `encrypt` 子命令及其描述

#### Scenario: encrypt 命令独立执行
- **WHEN** 用户运行 `tryssh encrypt`
- **THEN** 命令开始交互式密码输入流程

### Requirement: 交互式密码输入与确认
系统 SHALL 提示用户输入 master password 两次（输入 + 确认），两次输入一致时才继续。密码输入 SHALL 不回显到终端。

#### Scenario: 两次输入一致
- **WHEN** 用户输入密码 "mypassword" 并确认输入 "mypassword"
- **THEN** 系统使用该密码派生加密密钥，继续执行加密流程

#### Scenario: 两次输入不一致
- **WHEN** 用户输入密码 "password1" 并确认输入 "password2"
- **THEN** 系统输出错误信息 "Passwords do not match" 并退出，不修改配置文件

#### Scenario: 密码少于 4 个字符
- **WHEN** 用户输入密码 "abc"
- **THEN** 系统输出错误信息并退出，不修改配置文件

#### Scenario: 首次输入为空
- **WHEN** 用户直接按回车跳过密码输入
- **THEN** 系统输出错误信息并退出，不修改配置文件

### Requirement: 配置加密与写回
系统 SHALL 加载当前配置文件，加密所有明文密码字段（含 `main.passwords` 和 `serverList[].password`），并将加密后的配置写回磁盘。已加密的字段 SHALL 跳过不重复处理。

#### Scenario: 纯明文配置加密
- **WHEN** 配置文件包含明文密码 "mysecret"，用户执行 `tryssh encrypt` 并正确输入密码
- **THEN** 配置文件中该密码变为 `enc:` 前缀的加密字符串

#### Scenario: 混合明文和已加密配置
- **WHEN** 配置中一个密码是明文 "plain"，另一个已是 `enc:xxx`
- **THEN** 明文密码被加密，已有的 `enc:xxx` 保持不变

#### Scenario: 配置中没有密码字段
- **WHEN** 配置中 passwords 和 serverList 均为空
- **THEN** 命令成功完成，提示 "No passwords to encrypt"

### Requirement: 内存密钥清理
加密完成并写回配置后，系统 SHALL 清除内存中缓存的 master key。

#### Scenario: 加密完成后密钥被清除
- **WHEN** 加密流程成功完成
- **THEN** 缓存的 master key 被清零释放

### Requirement: 环境变量优先
如果环境变量 `TRYSSH_MASTER_KEY` 已设置，系统 SHALL 使用环境变量中的密码，不进行交互式输入。

#### Scenario: 环境变量已设置
- **WHEN** `TRYSSH_MASTER_KEY` 设置为 "envpassword"，用户运行 `tryssh encrypt`
- **THEN** 系统直接使用环境变量中的密码加密，不提示交互输入
