# tryssh

[![Go Report Card](https://goreportcard.com/badge/github.com/Driver-C/tryssh)](https://goreportcard.com/report/github.com/Driver-C/tryssh)

[English](README.md) | 简体中文

`tryssh`是一个具有密码猜测功能的命令行SSH终端工具。

它可以使用SSH协议交互登录服务器，或者将本地文件上传到服务器或将远程文件下载到本地。

当然，它也可以管理用于尝试登陆服务器的用户名、端口号、密码以及已经成功登陆服务器的缓存信息。

> 注意！不要将`tryssh`用于生产场景！

## 我为什么需要 tryssh ?

* 我只喜欢使用命令行工具，不想使用图形化工具
* 我有很多登陆信息相似的服务器，但是我不想每次登陆都输入登陆信息
* 我经常跨操作系统使用SSH终端，但是没有找到让我在多种操作系统上使用习惯不变的工具
* `tryssh` 我没用过，看着还不错，想试试

## 当前开发状态

目前`tryssh`处于功能完善阶段，基本功能已有，但是在功能的细节上做得不好还需要改进，比如安全。

目前仅有 *Driver-C* 一人参与开发，而且需要利用业务时间来完成，所以开发进度不会很快。

如果遇到任何使用问题，任何建议请提交`issue`，会尽快回复。

目前项目仅保留`master`分支用于发布稳定版本，`tag`也从master分支创建。

## 待做清单

排名不区分优先级，以下内容在完成后删除对应条目

1. 传输文件支持通配符
2. 完成单元测试代码
3. 安全相关功能，配置文件加密、隐藏明文显示的敏感信息、密码输入应改为交互式等

## 快速开始

```bash
# 创建一个名为 testuser 的备选用户
tryssh create users testuser

# 创建备选端口号 22
tryssh create ports 22

# 创建一个备选密码
tryssh create passwords 123456

# 用以上创建的信息尝试登陆 192.168.1.1
tryssh ssh 192.168.1.1
```

## 怎么查看其他帮助

在`tryssh`的帮助信息中已经写好了所有子命令的使用帮助，可以通过下列命令查看

```bash
tryssh -h

# 查看子命令 ssh 的帮助
tryssh ssh -h
```

## 功能详解

```
$ tryssh -h
 command line ssh terminal tool.

Usage:
  tryssh [command]

Available Commands:
  alias       Set, unset, and list aliases, aliases can be used to log in to servers
  create      Create alternative username, port number, password, and login cache information
  delete      Delete alternative username, port number, password, and login cache information
  get         Get alternative username, port number, password, and login cache information
  help        Help about any command
  prune       Check if all current caches are available and clear the ones that are not available
  scp         Upload/Download file to/from the server through SSH protocol
  ssh         Connect to the server through SSH protocol
  version     Print the client version information for the current context

Flags:
  -h, --help   help for tryssh

Use "tryssh [command] --help" for more information about a command.
```

### create 命令

tryssh 的`create`命令用于创建用来猜密码登陆的各类配置，比如用户名、端口号和密码，也可以直接创建已知用户名、端口号和密码的缓存。

#### create 帮助信息
```
$ tryssh create -h
Create alternative username, port number, password, and login cache information

Usage:
  tryssh create [command]

Available Commands:
  caches      Create an alternative cache
  keys        Create a alternative key file path
  passwords   Create an alternative password
  ports       Create an alternative port
  users       Create an alternative username

Flags:
  -h, --help   help for create

Use "tryssh create [command] --help" for more information about a command.
```

#### create 使用举例

```
# 创建一个名为 testuser 的备选用户
tryssh create users testuser

# 创建备选端口号 22
tryssh create ports 22

# 创建一个备选密码
tryssh create passwords 123456
```

### delete 命令

tryssh 的`delete`命令用于删除用来猜密码登陆的各类配置，比如用户名、端口号和密码，也可以直接删除缓存。

#### delete 帮助信息

```
$ tryssh delete -h
Delete alternative username, port number, password, and login cache information

Usage:
  tryssh delete [command]

Available Commands:
  caches      Delete an alternative cache
  passwords   Delete an alternative password
  ports       Delete an alternative port
  users       Delete an alternative username

Flags:
  -h, --help   help for delete

Use "tryssh delete [command] --help" for more information about a command.
```

#### delete 使用举例

```
# 删除一个名为 testuser 的备选用户
tryssh delete users testuser

# 删除备选端口号 22
tryssh delete ports 22

# 删除一个备选密码
tryssh delete passwords 123456

# 删除服务器192.168.1.1的登陆缓存
tryssh delete caches 192.168.1.1
```

### get 命令

tryssh 的`get`命令用于查看用来猜密码登陆的各类配置，比如用户名、端口号、密码以及登陆缓存。

#### get 帮助信息

```
$ tryssh get -h
Get alternative username, port number, password, and login cache information

Usage:
  tryssh get [command]

Available Commands:
  caches      Get alternative caches by ipAddress
  keys        Delete a alternative key file path
  passwords   Get alternative passwords
  ports       Get alternative ports
  users       Get alternative usernames

Flags:
  -h, --help   help for get

Use "tryssh get [command] --help" for more information about a command.
```

#### get 使用举例

```
# 查看用于猜密码的候选用户
tryssh get users

# 查看用于猜密码的候选端口号
tryssh get ports

# 查看当前已有的登陆缓存
tryssh get caches
```

### prune 命令

tryssh的`prune`命令用于测试当前已有缓存是否依然可用，如果不可用可以选择执行删除缓存，也可以不询问直接删除缓存。

#### prune 帮助信息

```
$ tryssh prune -h
Check if all current caches are available and clear the ones that are not available

Usage:
  tryssh prune [flags]

Flags:
  -a, --auto               Automatically perform concurrent cache optimization without asking for confirmation to delete
  -c, --concurrency int    Number of multiple requests to perform at a time (default 8)
  -h, --help               help for prune
  -t, --timeout duration   SSH timeout when attempting to log in. It can be "1s" or "1m" or other duration (default 2s)
```

#### prune 使用举例

```
# 交互式进行缓存可用性测试
tryssh prune

# 非交互进行缓存可用性测试
tryssh prune -a

# 非交互进行缓存可用性测试，同时设置并发数为10(默认为8)，连接超时时间为5秒(默认为2秒)
tryssh prune -c 10 -t 5s -a
```

> 交互式模式下设置并发数是无效的

### alias 命令

tryssh 的`alias`命令是用于给 *已有* 的缓存设置别名用的，方便在登陆或者传输文件时直接使用别名来操作

#### alias 帮助信息

```
$ tryssh alias -h
Set, unset, and list aliases, aliases can be used to log in to servers

Usage:
  tryssh alias [command]

Available Commands:
  list        List all alias
  set         Set an alias for the specified server address
  unset       Unset the alias

Flags:
  -h, --help   help for alias

Use "tryssh alias [command] --help" for more information about a command.
```

#### alias 使用举例

```
# 查看当前所有别名
tryssh alias list

# 给192.168.1.1服务器设置一个名为"host1"的别名
tryssh alias set host1 -t 192.168.1.1

# 取消名为"host1"的别名
tryssh alias unset host1
```

### ssh 命令

tryssh 的`ssh`命令用于猜密码登陆服务器，在成功获取正确登陆信息后会缓存这些信息以便下次直接使用缓存登陆，不用重新猜密码。

#### ssh 帮助信息

```
chenjingyu@MacBook ~ % tryssh ssh -h
Connect to the server through SSH protocol

Usage:
  tryssh ssh <ipAddress> [flags]

Flags:
  -c, --concurrency int    Number of multiple requests to perform at a time (default 8)
  -h, --help               help for ssh
  -t, --timeout duration   SSH timeout when attempting to log in. It can be "1s" or "1m" or other duration (default 1s)
  -u, --user string        Specify a username to attempt to login to the server,
                           if the specified username does not exist, try logging in using that username
```

#### ssh 使用举例

```
# 登陆192.168.1.1服务器，如果没有缓存则尝试猜密码登陆
tryssh ssh 192.168.1.1

# 登陆别名为host1的服务器
tryssh ssh host1

# 登陆192.168.1.1服务器，如果没有缓存则尝试猜密码登陆，同时设置并发数为20，超时时间为500毫秒，指定登陆的用户为root
tryssh ssh 192.168.1.1 -c 20 -t 500ms -u root
```

### scp 命令

tryssh 的`scp`命令用于上传或者下载文件或者目录，`scp`命令支持使用别名

#### scp 帮助信息

```
chenjingyu@MacBook ~ % tryssh scp -h
Upload/Download file to/from the server through SSH protocol

Usage:
  tryssh scp <source> <destination> [flags]

Examples:
# Download test.txt file from 192.168.1.1 and place it under ./
tryssh scp 192.168.1.1:/root/test.txt ./
# Upload test.txt file to 192.168.1.1 and place it under /root/
tryssh scp ./test.txt 192.168.1.1:/root/
# Download test.txt file from 192.168.1.1 and rename it to test2.txt and place it under ./
tryssh scp 192.168.1.1:/root/test.txt ./test2.txt

# Download testDir directory from 192.168.1.1 and place it under ~/Downloads/
tryssh scp -r 192.168.1.1:/root/testDir ~/Downloads/
# Upload testDir directory to 192.168.1.1 and rename it to testDir2 and place it under /root/
tryssh scp -r ~/Downloads/testDir 192.168.1.1:/root/testDir2

Flags:
  -c, --concurrency int    Number of multiple requests to perform at a time (default 8)
  -h, --help               help for scp
  -r, --recursive          Recursively copy entire directories
  -t, --timeout duration   SSH timeout when attempting to log in. It can be "1s" or "1m" or other duration (default 1s)
  -u, --user string        Specify a username to attempt to login to the server,
                           if the specified username does not exist, try logging in using that username
```

#### scp 使用举例

> scp的使用例子在帮助信息里已经有阐述，下面只是做翻译

```
# 从192.168.1.1服务器上下载test.txt文件放到本地的./目录
tryssh scp 192.168.1.1:/root/test.txt ./

# 从本地上传test.txt文件到192.168.1.1的/root/目录下
tryssh scp ./test.txt 192.168.1.1:/root/

# 从192.168.1.1服务器上下载test.txt文件到本地./目录下并改名为test2.txt
tryssh scp 192.168.1.1:/root/test.txt ./test2.txt

# 从192.168.1.1服务器上下载testDir目录到本地的~/Downloads/下
tryssh scp -r 192.168.1.1:/root/testDir ~/Downloads/

# 上传本地的testDir目录到192.168.1.1服务器/root/下并改名为testDir2
tryssh scp -r ~/Downloads/testDir 192.168.1.1:/root/testDir2
```
