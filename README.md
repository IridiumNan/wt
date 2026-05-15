# water-repo

## 概述

- 一款基于go的超轻量级个人仓库管理器
- 纯命令行， 没有UI界面， 服务器友好
- CS 架构， 二进制文件开箱即用
- 简便高效,没有任何花里胡哨

---

## 基本用法

- 搜索仓库中的包

```bash
wt search <package name>
```

- 下载仓库中的包

```bash
wt install <package name>
```

- 上传本地的文件

```bash
wt upload <path to your local package>
```

- 替换仓库中的包(用于更新)

```bash
wt replace <package name> <path to your new package>
```

- 重命名仓库的包

```bash
wt mv <package name> <new package name>
```

- 删除仓库中的包

```bash
wt rm <package name>
```

---

## 权限管理

> 为了实现多人使用时的精细化控制，使用token进行不同权限的验证

这里分成三种权限

1. Read -> search
2. Install -> install
3. Write -> upload replace mv rm tag

> 权限通过配置文件中的token进行控制， 如果客户端的配置文件中没有正确的token, 则会导致相应的操作失败

---

## 配置文件

- Client 端的配置文件是 ~/.config/water-repo/client_config.json

```json
{
    "server": "服务器ip:端口"

    "read_timeout": "读取的超时时间",
    "install_timeout": "下载的超时时间",
    "write_timeout": "写入超时时间",

    "read_token": "具体的阅读权限token",
    "install_token": "具体的下载token",
    "write_token": "具体的写入权限token"

}
```

- Server 端的配置文件是 ~/.config/water-repo/server_config.json

```json
{
    "server": "服务器ip:端口",
    
    "read_timeout": "读取超时时间",
    "write_timeout": "写入超时时间",
    "install_timeout": "下载超时时间",

    "read_token": [
        "阅读token1",
        "阅读token2",
        ...
    ],
    "install_token": [
        "下载token1",
        ...
    ],
    "write_token": [
        "写入token1",
        ...
    ]
}
```

---

## 进阶用法 TODO

### tag 管理

> 为了管理不同类型的包，加入了 `tag` 这个属性， 可以根据tag来进行整理和分离

- 默认带有 `static` 和 `temp` 两种tag来区别临时文件和持久文件, 这两种标签不可删除

- 使用下面的方法列出 tag 为 target tag 的所有包

```bash
wt list <target tag> 
```

- 使用upload上传的包的tag默认是`temp`, 可以通过以下方法修改标签

```bash
wt tag <package name> <target tag>
```

- 添加新的标签

```bash
wt tag add <tag name>
```

- 删除标签

```bash
wt tag rm <tag name>
```

> 标签被删除之后， 原来属于这个标签的所有包都会变成 `temp` 标签

- 清理所有 tag 为 target tag 的包 **慎用!!!**

```bash
wt clear <target tag>
```

---
