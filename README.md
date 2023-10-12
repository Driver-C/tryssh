# tryssh

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

Rankings do not differentiate priority levels.

1. File transfer supports wildcards
2. Completing unit test code
3. Security-related features, such as encrypting configuration files, hiding sensitive information from plain text display, and switching to interactive password input
4. Support for key-based authentication
5. One-click functionality to check the availability of the current cache

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
