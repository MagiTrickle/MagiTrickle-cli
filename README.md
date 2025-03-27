![MagiTrickle-cli](./img/compressed.webp)

# MagiTrickle CLI

`magitrickle` is a command-line tool for managing and configuring MagiTrickle through a UNIX socket.

---

## Table of Contents

1. [Overview](#overview)
2. [Installation](#installation)
3. [How It Works](#how-it-works)
4. [CLI Usage](#cli-usage)
   - [System Commands](#system-commands)
   - [Group Commands](#group-commands)
   - [Rule Commands](#rule-commands)
5. [Examples](#examples)
6. [Tips and Troubleshooting](#tips-and-troubleshooting)
7. [License](#license)

---

## Overview

MagiTrickle CLI helps you interact with MagiTrickle's backend via:
- **Groups** – logical containers for traffic or domain matching
- **Rules** – define the matching patterns or IP/domain restrictions
- **System hooks** – advanced capabilities such as netfilterd, interface listing, and saving config

The CLI communicates over a [UNIX socket](http://unix) to send HTTP requests to the MagiTrickle API.  
Commands are structured in a tree-like format under the single executable **`magitrickle`**:

- **`system`**  
  – For system-level operations (hooks, listing interfaces, saving configs).
- **`group`**  
  – Create, list, update, and delete groups.
- **`rule`**  
  – Create, list, update, and delete rules in a specified group.

---

## Installation

*Will be added soon...*

---

## CLI Usage

You can see the full list of commands by typing:
```bash
magitrickle --help
```
Similarly, each subcommand supports `--help` to list detailed usage:
```bash
magitrickle group --help
magitrickle rule --help
magitrickle system --help
```

## Examples

### 1. List Groups
```bash
magitrickle group list
```
Output:
```
Groups:
 - ID: 182f11dd
   Name: Debug
   Interface: singtun-ru1
   Enabled: true
   Color: #791a3e
```

### 2. Create a Group
```bash
magitrickle group create \
    --name="TestGroup" \
    --interface="br1" \
    --enable=true \
    --color="#abc123"
```
Output:
```
Group created successfully
 ID: e89c1f15
 Name: TestGroup
 Interface: br1
 Enabled: true
 Color: #abc123
```

### 3. Add a Rule to a Group
```bash
magitrickle rule create e89c1f15 \
    --name="BlockExampleDomain" \
    --type="domain" \
    --rule="example.com" \
    --enable=true \
    --save
```
Output:
```
Rule created successfully:
 ID: 4c40d238 | Name: BlockExampleDomain | Type: domain | Rule: example.com | Enabled: true
```

### 4. Update a Group (and Save Configuration)
```bash
magitrickle group update e89c1f15 \
    --name="NewGroupName" \
    --enable=false \
    --save
```
Output:
```
Group updated successfully
 ID: e89c1f15
 Name: NewGroupName
 Interface: br1
 Enabled: false
 Color: #abc123
```
Configuration is immediately persisted due to `--save`.

### 5. List Interfaces
```bash
magitrickle system interfaces
```
Example output:
```
Available Interfaces:
  - br0
  - br1
  - eth0
  - wlan0
```

### 6. Save Overall Configuration
```bash
magitrickle system save-config
```
Output:
```
Configuration saved successfully
```

---

## Tips and Troubleshooting

1. **Make sure the MagiTrickle backend is running.**  
   The CLI will fail to connect if the UNIX socket (e.g., `/var/run/magitrickle.sock`) is not accessible or if the `magitrickled` is offline.

2. **Use `--help` often.**  
   Each subcommand has detailed flags and usage info.

3. **Persisting changes with `--save`.**  
   When you create, update, or delete a group/rule, you can optionally add `--save` to immediately persist those changes to the server configuration. Otherwise, you can always run:
   ```bash
   magitrickle system save-config
   ```
   afterwards to commit all outstanding changes.

---