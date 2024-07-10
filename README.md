# tryssh

[![Go Report Card](https://goreportcard.com/badge/github.com/Driver-C/tryssh)](https://goreportcard.com/report/github.com/Driver-C/tryssh)

English | [简体中文](README_zh.md)

`tryssh` is a command line SSH terminal tool with password guessing function. 

It can use the SSH protocol to interactively log in to the server, or upload local files to the server or download remote files to the local location.

Of course, it can also manage the usernames, port numbers, passwords, and cached information of successfully logged-in servers for login attempts.

> Attention! Do not use `tryssh` in a production environment!

## Why do I need "tryssh"?

* I prefer command-line tools and do not want to use graphical tools
* I have many servers with similar login information, but I don't want to input login details every time I log in
* I frequently use SSH terminal across multiple operating systems, but I haven't found a tool that allows me to maintain the same workflow across different OSes
* I haven't used `tryssh` before, but it looks good, and I want to give it a try

## Current development status

Currently, `tryssh` is in the stage of feature completion. The core functionalities are already implemented, but there is room for improvement in terms of details, particularly in areas such as security.

Currently, only one person *Driver-C* is involved in the development, and the progress is limited by the need to allocate time from other work responsibilities. Therefore, the development progress is not expected to be fast.

If you encounter any usage issues or have any suggestions, please submit an `issue`. We will respond as soon as possible.

Currently, the project only maintains the `master` branch for releasing stable versions, and `tags` are created from the `master` branch as well.

## TODO list

Rankings do not differentiate priority levels. Delete the corresponding entry after completion of the following content.

1. File transfer supports wildcards
2. Completing unit test code
3. Security-related features, such as encrypting configuration files, hiding sensitive information from plain text display, and switching to interactive password input

## Quick Start

```bash
# Create an alternative user named "testuser"
tryssh create users testuser

# Create an alternative port number 22
tryssh create ports 22

# Create an alternative password
tryssh create passwords 123456

# Attempt to log in to 192.168.1.1 using the information created above
tryssh ssh 192.168.1.1
```

## How to view other help documentation?

All usage help for the subcommands has been documented in the `tryssh` help information. You can view it by using the following command:

```bash
tryssh -h

# View the help documentation for the subcommand "ssh"
tryssh ssh -h
```

## Detailed Function Explanation

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

### Command: create

The `"create"` command of `tryssh` is used to create various configurations for password guessing login, such as usernames, port numbers, and passwords. It can also directly create caches with known usernames, port numbers, and passwords.

#### Help information

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

#### Example

```
# Create an alternative user named testuser
tryssh create users testuser

# Create an alternative port: 22
tryssh create ports 22

# Create an alternative passwords: 123456
tryssh create passwords 123456
```

### Command: delete

The `"delete"` command of `tryssh` is used to delete various configurations for password guessing login, such as usernames, port numbers, and passwords. It can also directly delete caches.

#### Help information

```
$ tryssh delete -h
Delete alternative username, port number, password, and login cache information

Usage:
  tryssh delete [command]

Available Commands:
  caches      Delete an alternative cache
  keys        Delete a alternative key file path
  passwords   Delete an alternative password
  ports       Delete an alternative port
  users       Delete an alternative username

Flags:
  -h, --help   help for delete

Use "tryssh delete [command] --help" for more information about a command.
```

#### Example

```
# Delete an alternative user named testuser
tryssh delete users testuser

# Delete an alternative port: 22
tryssh delete ports 22

# Delete an alternative passwords: 123456
tryssh delete passwords 123456

# Delete the cache information about 192.168.1.1
tryssh delete caches 192.168.1.1
```

### Command: get

The `"get"` command of `tryssh` is used to view various configurations for password guessing login, such as usernames, port numbers, passwords, and login caches.

#### Help information

```
$ tryssh get -h
Get alternative username, port number, password, and login cache information

Usage:
  tryssh get [command]

Available Commands:
  caches      Get alternative caches by ipAddress
  passwords   Get alternative passwords
  ports       Get alternative ports
  users       Get alternative usernames

Flags:
  -h, --help   help for get

Use "tryssh get [command] --help" for more information about a command.
```

#### Example

```
# View candidate users for password guessing
tryssh get users

# View candidate ports for password guessing
tryssh get ports

# View the currently existing login caches
tryssh get caches
```

### Command: prune

The `"prune"` command of `tryssh` is used to test whether the existing caches are still usable. If they are not, you can choose to delete the cache, or directly delete it without confirmation.

#### Help information

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

#### Example

```
# Interactively conduct cache availability testing
tryssh prune

# Conduct non-interactive cache availability testing
tryssh prune -a

# Conduct non-interactive cache availability testing, while setting the concurrency to 10 (default is 8) and the connection timeout to 5 seconds (default is 2 seconds).
tryssh prune -c 10 -t 5s -a
```

> The setting for concurrency is invalid in interactive mode.

### Command: alias

The `"alias"` command of `tryssh` is used to assign aliases to existing caches, making it convenient to use these aliases directly for login or file transfer operations.

#### Help information

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

#### Example

```
# View all current aliases
tryssh alias list

# Set an alias named 'host1' for the server with the IP address 192.168.1.1
tryssh alias set host1 -t 192.168.1.1

# Remove the alias named 'host1'
tryssh alias unset host1
```

### Command: ssh

The `"ssh"` command of `tryssh` is used for password guessing login to a server. Upon successfully obtaining the correct login information, it will cache these details for direct login in future attempts, without the need to guess the password again.

#### Help information

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

#### Example

```
# Login to the server at 192.168.1.1. If there is no cache available, attempt to guess the password for login
tryssh ssh 192.168.1.1

# Login to the server with the alias 'host1'
tryssh ssh host1

# Login to the server at 192.168.1.1. If there is no cache available, attempt to guess the password for login. Set the concurrency to 20, timeout to 500 milliseconds, and specify the user as 'root'.
tryssh ssh 192.168.1.1 -c 20 -t 500ms -u root
```

### Command: scp

The `"scp"` command of `tryssh` is used to upload or download files or directories. The scp command supports the use of aliases.

#### Help information

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

#### Example

> Same as the information in the help section.

```
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
```
