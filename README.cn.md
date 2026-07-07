# water-repo

[![License](https://img.shields.io/badge/license-GPL--3.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.25.9+-blue.svg)](https://golang.org/)
[![Version](https://img.shields.io/badge/version-v0.2.1-orange.svg)]()

超轻量级个人/小团队仓库管理工具,基于 Go 开发,纯 CLI 交互(无图形化界面),天生适配服务器环境;CS 架构设计,单二进制文件开箱即用,零依赖零配置,没有任何多余功能。

---

## 🚀 核心特性

- **轻量级**：单二进制文件,零依赖
- **简单易用**：直观的命令行操作
- **安全可靠**：基于 Token 的权限控制系统,支持标签级细粒度权限
- **灵活组织**：标签化管理包资源
- **跨平台支持**：支持 Linux、macOS 和 Windows
- **公开分享**：集成 Tailscale Funnel，轻松实现包的公开分享

---

## 📦 安装方式

### 下载二进制文件

#### Linux

```bash
wget https://repo.waterman.xin/apps/water-repo/wt-latest-linux-amd64
chmod +x wt-latest-linux-amd64
mv wt-latest-linux-amd64 ~/.local/bin/wt  # 或任何在 PATH 中的目录
```

#### macOS (Apple Silicon)

```bash
curl -LO https://repo.waterman.xin/apps/water-repo/wt-latest-darwin-arm64
chmod +x wt-latest-darwin-arm64
mv wt-latest-darwin-arm64 ~/.local/bin/wt  # 或任何在 PATH 中的目录
```

#### Windows

从 [wt-latest-windows-amd64.exe](https://repo.waterman.xin/apps/water-repo/wt-latest-windows-amd64.exe) 下载

> **注意**：将下载的文件重命名为 `wt.exe` 并添加到 PATH 环境变量中。

### 从源码构建

```bash
git clone <repository-url>
cd water-repo
go build -o wt ./cmd/wt
```

---

## 🛠️ 快速开始

### 服务端设置

#### 方式一：交互式配置（首次使用推荐）

1. **启动服务器**：

   ```bash
   wt server
   ```

2. **配置服务器**（交互式设置）：

   ```
   >>设置服务器的主机和端口<<
   示例: 0.0.0.0:12212
   默认值: 0.0.0.0:12212
   请输入您的值[留空使用默认值] ->

   >>设置服务器读取超时时间<<
   示例: 20s
   默认值: 20s
   请输入您的值[留空使用默认值] ->

   >>设置服务器安装超时时间<<
   示例: 3h
   默认值: 3h
   请输入您的值[留空使用默认值] ->

   >>设置服务器写入超时时间<<
   示例: 3h
   默认值: 3h
   请输入您的值[留空使用默认值] ->

   >>设置服务器读取令牌<<
   示例: xxxxx
   默认值: 未设置默认令牌
   请输入要添加的令牌数量:1
   设置您的值[留空则随机生成]
   输入第 1 个令牌 ->[your-token]

   >>设置服务器安装令牌<<
   ...
   
   >>设置服务器写入令牌<<
   ...
   ```

#### 方式二：手动配置

1. **启动服务器并指定数据目录**（可选）：

   ```bash
   wt server -d /path/to/data/directory
   ```

2. **使用命令手动配置服务器**：

   ```bash
   # 设置服务器地址和超时时间（替换现有值）
   wt server config server 0.0.0.0:12212
   wt server config read_timeout 20s
   wt server config write_timeout 30s
   wt server config install_timeout 3h
   
   # 添加令牌（可以添加多个令牌）
   wt server config read_token add your-read-token
   wt server config install_token add your-install-token
   wt server config write_token add your-write-token
   
   # 查看当前配置
   wt server config show
   ```

3. **重启服务器**使配置生效：

   ```bash
   # 停止服务器（Ctrl+C）并重新启动
   wt server
   ```

### 客户端设置

#### 方式一：交互式配置（首次使用推荐）

1. **初始化客户端**：

   ```bash
   wt ls
   ```

2. **配置客户端**（交互式设置）：

   ```
   预期配置文件路径:  /home/user/.config/water-repo/client_config.json
   客户端配置文件不存在,开始初始化...

   >>设置服务器地址 -> host:port<<
   示例: http://192.168.1.2:12212
   默认值: http://192.168.1.2:12212
   请输入您的值[留空使用默认值] ->

   >>设置客户端读取超时时间<<
   示例: 20s
   默认值: 20s
   请输入您的值[留空使用默认值] ->

   >>设置客户端安装超时时间<<
   示例: 3h
   默认值: 3h
   请输入您的值[留空使用默认值] ->

   >>设置客户端写入超时时间<<
   示例: 3h
   默认值: 3h
   请输入您的值[留空使用默认值] ->

   >>设置客户端读取令牌<<
   示例: xxxxxxx
   默认值: 随机生成
   请输入您的令牌 [按回车键随机生成] -> 
   生成的随机令牌: [token]

   >>设置客户端安装令牌<<
   ...
   
   >>设置客户端写入令牌<<
   ...
   ```

#### 方式二：手动配置

1. **使用命令手动配置客户端**：

   ```bash
   # 设置服务器地址和超时时间
   wt config server http://192.168.1.2:12212
   wt config read_timeout 20s
   wt config install_timeout 3h
   wt config write_timeout 30s
   
   # 设置令牌（替换现有令牌）
   wt config read_token your-read-token
   wt config install_token your-install-token
   wt config write_token your-write-token
   
   # 查看当前配置
   wt config show
   ```

> **重要提示**：确保客户端令牌与服务器相应权限级别的令牌匹配。

---

## 📋 基础命令

| 操作场景 | 命令示例 | 说明 |
|---------|----------|------|
| 搜索包 | `wt search <包名>` | 根据包名模糊检索所有资源 |
| 查看包详情 | `wt info <包名>` | 获取包大小、上传时间、标签等信息 |
| 下载安装包 | `wt install <包名>` | 从仓库下载指定包到本地 |
| 上传本地包 | `wt upload <文件路径> [包名]` | 将本地文件上传至仓库（如不提供包名则自动生成） |
| 重命名包 | `wt mv <旧名称> <新名称>` | 修改仓库内包的名称 |
| 删除包 | `wt rm <包名>` | 永久移除仓库中的指定包 |
| 列出包 | `wt list [标签]` 或 `wt ls [标签]` | 列出所有包或按标签过滤 |
| 同步元数据 | `wt sync` | 同步本地元数据与服务器 |
| 公开分享 | `wt public <包名>` | 通过 Tailscale Funnel 公开分享包 |
| 取消公开 | `wt private <包名>` | 移除包的公开分享 |
| 查看公开链接 | `wt links` | 显示所有公开分享的包 |
| 查看配置服务器 | `wt list-servers` | 显示配置的客户端服务器 |
| 切换当前服务器 | `wt change-server <名称>` | 切换活动客户端服务器 |
| 添加服务器 | `wt add-server <名称> <服务器地址>` | 添加命名的客户端服务器 |
| 删除服务器 | `wt del-server <名称>` | 删除已配置的客户端服务器 |
| 重载配置 | `wt reload` | 无需重启即可重载服务器配置 |
| 显示帮助 | `wt help` 或 `wt help <command>` | 显示命令帮助信息 |

---

## 🔐 权限管理

采用 **两级 Token 权限体系** —— 全局 Token 作为超管兜底,标签级 Token 提供细粒度权限控制。无需用户注册登录。

| 权限等级 | 可执行操作 | 适用场景 |
|---------|-----------|---------|
| Read | `search`、`list`、`info` | 仅允许查看仓库内容 |
| Install | `install` | 允许下载使用仓库中的包 |
| Write | `upload`、`mv`、`rm`、`tag`、`sync` | 允许管理和修改仓库内容 |

**权限校验顺序**：全局 Token → 标签 Token → 拒绝

- 全局 Token：对所有标签拥有对应级别的权限
- 标签 Token：仅对特定标签拥有权限

---

## 🏷️ 标签管理

通过 **标签（Tag）** 对仓库内的包进行逻辑分类,替代复杂的目录结构和命名空间,所有包必须归属且仅归属一个标签。

系统默认内置两个不可删除的系统标签：

- `temp`：临时文件标签,所有上传的包默认归属此标签
- `static`：持久文件标签,用于存放长期使用的静态资源

### 标签操作命令

| 操作场景 | 命令示例 | 说明 |
|---------|----------|------|
| 列出指定标签的所有包 | `wt list <标签名>` | 展示归属该标签的全部包 |
| 修改包的标签 | `wt tag <包名> <目标标签>` | 将指定包移动到目标标签下 |
| 新增自定义标签 | `wt tag add <标签名>` | 创建新的分类标签（自动继承全局 Token） |
| 删除自定义标签 | `wt tag rm <标签名>` | 删除指定标签；原归属该标签的所有包会自动回退到 `temp` 标签 |
| 列出可见标签 | `wt tag list` | 列出当前客户端有权访问的所有标签 |

### 标签级 Token 管理（服务端）

| 操作场景 | 命令示例 | 说明 |
|---------|----------|------|
| 为标签添加读权限 | `wt server tag <标签> read_token add <token>` | 授予对特定标签的读取权限 |
| 为标签添加下载权限 | `wt server tag <标签> install_token add <token>` | 授予对特定标签的下载权限 |
| 为标签添加写权限 | `wt server tag <标签> write_token add <token>` | 授予对特定标签的写入权限 |

> **注意**：标签 Token 变更后需重启服务器才能生效。

---

## ⚙️ 配置文件

### 客户端配置

路径：`~/.config/water-repo/client_config.json`

```json
{
    "server": "http://服务器IP:端口",
    "read_timeout": "10s",
    "install_timeout": "30m",
    "write_timeout": "30s",
    "read_token": "你的阅读权限Token",
    "install_token": "你的下载权限Token",
    "write_token": "你的写入权限Token"
}
```

### 服务端配置

路径：`~/.config/water-repo/server_config.json`

```json
{
    "server": "127.0.0.1:8080",
    "read_timeout": "15s",
    "write_timeout": "30s",
    "install_timeout": "1h",
    "read_token": [
        "全局阅读Token1",
        "全局阅读Token2"
    ],
    "install_token": [
        "全局下载Token1"
    ],
    "write_token": [
        "管理员写入Token1"
    ],
    "tag_token": {
        "temp": {
            "read_token": [],
            "install_token": [],
            "write_token": []
        },
        "static": {
            "read_token": [],
            "install_token": [],
            "write_token": []
        },
        "frontend": {
            "read_token": ["前端读取Token"],
            "install_token": [],
            "write_token": ["前端写入Token"]
        }
    }
}
```

---

## 📁 文件位置

- **数据存储**：当前目录（或由 `-d` 标志指定）
- **日志文件**：`~/.local/state/water-repo/log.txt`
- **客户端配置**：`~/.config/water-repo/client_config.json`
- **服务端配置**：`~/.config/water-repo/server_config.json`

---

## 🧪 使用示例

### 上传包

```bash
wt upload ./build/my-app.tar.gz my-app-v1
```

### 搜索包

```bash
wt search my-app
```

### 查看包信息

```bash
wt info my-app-v1
```

### 下载包

```bash
wt install my-app-v1
```

### 重命名包

```bash
wt mv my-app-v1 my-app-v2
```

### 列出所有包

```bash
wt ls
```

### 按标签列出包

```bash
wt list temp
```

### 同步元数据

```bash
wt sync
```

### 标签操作

```bash
# 创建新标签
wt tag add frontend

# 将包移动到标签
wt tag my-app frontend

# 列出所有可见标签
wt tag list

# 为标签添加 Token
wt server tag frontend write_token add my-team-token
```

### 使用 Tailscale Funnel 公开分享

Water-Repo 集成了 [Tailscale Funnel](https://tailscale.com/kb/1223/funnel)，可以安全地公开分享包，而无需暴露整个服务器。

#### 前置要求

1. 在服务器上安装并配置 [Tailscale](https://tailscale.com/)
2. 为您的 Tailscale 节点启用 Funnel：
   ```bash
   tailscale funnel 443 on
   ```

#### 命令

```bash
# 公开分享包
wt public my-app-v1
# 返回: https://your-node.tailnet.ts.net/install?name=my-app-v1

# 查看所有公开分享的包
wt links

# 取消包的公开分享
wt private my-app-v1
```

#### 工作原理

- 运行 `wt public <包名>` 时，Water-Repo 会在 Tailscale Funnel 暴露的特殊目录中创建符号链接
- 包可以通过公共 HTTPS URL 访问，无需认证令牌
- 使用 `wt private <包名>` 移除公开链接
- 使用 `wt links` 审计当前所有分享的包

> **安全提示**：公开分享的包对任何拥有链接的人都可访问。仅分享您打算公开的包。

---

## 🛠️ 进阶用法

### 配置管理

#### 客户端配置

```bash
# 显示当前配置
wt config show

# 修改配置（替换现有值）
wt config server http://192.168.1.2:12212
wt config read_timeout 20s
wt config install_timeout 3h
wt config write_timeout 30s
wt config read_token <新令牌>
wt config install_token <新令牌>
wt config write_token <新令牌>
```

#### 服务端配置

```bash
# 显示服务器配置
wt server config show

# 修改服务器地址和超时时间（替换现有值）
wt server config server 0.0.0.0:8080
wt server config read_timeout 30s
wt server config write_timeout 30s
wt server config install_timeout 3h

# 添加令牌（追加到令牌列表）
wt server config read_token add <新令牌>
wt server config install_token add <新令牌>
wt server config write_token add <新令牌>

# 标签级 Token 管理
wt server tag frontend read_token add <新令牌>
wt server tag backend write_token add <新令牌>
```

> **注意**：
>
> - 客户端配置命令会替换现有值
> - 服务端令牌命令使用 `add` 关键字追加到令牌列表（支持多个令牌）
> - 服务器配置更改需要重启服务器才能生效

### 服务器日志

```bash
# 查看服务器日志
wt server log
```

---

## 💡 FZF 集成技巧

[fzf](https://github.com/junegunn/fzf) 是一个强大的命令行模糊查找工具，与 Water-Repo 配合使用效果极佳。以下是一些实用的一行命令：

### 搜索与交互

```bash
# 搜索并获取包信息
wt search <pkg> | fzf | xargs -I {} wt info {}

# 搜索并删除包
wt search <pkg> | fzf | xargs -I {} wt rm {}

# 搜索并下载包
wt search <pkg> | fzf | xargs -I {} wt install {}

# 搜索并公开分享包
wt search <pkg> | fzf | xargs -I {} wt public {}
```

### 文件操作

```bash
# 选择文件并上传（交互式文件选择）
wt upload $(fzf)
```

### 标签管理

```bash
# 选择标签并列出包
wt ls $(wt tag ls | fzf)

# 为包添加标签（交互式选择包和标签）
wt tag $(wt ls | fzf) $(wt tag ls | fzf)
```

> **提示**：从 [https://github.com/junegunn/fzf](https://github.com/junegunn/fzf) 安装 fzf 以获得增强的交互式工作流。

---

## 🚧 开发中功能 (v0.3.0+)

以下功能正在开发中,将在后续版本中陆续上线,现有配置文件和核心命令保持完全向后兼容。

### 核心待实现（下一版本）

1. **批量操作**
   - `wt clear <tag>`：清空标签下所有包（需确认）
   - `wt install-tag <tag>`：批量下载标签下所有包

2. **动态镜像源管理**
   - 支持添加多个远程 wt 服务端作为镜像源
   - 搜索时自动并行查询所有镜像源,合并返回结果
   - 下载时自动选择最快的可用源,失败自动切换到其他源
   - 提供 `wt mirror add/remove/list` 命令管理镜像源

3. **局域网自动发现**
   - 服务端启动后自动广播自身存在
   - 客户端自动发现同一局域网内所有运行中的 wt 服务端
   - 发现的节点自动加入镜像源列表,无需手动配置
   - 支持一键关闭自动发现功能

### 最近完成 (v0.2.0)

1. ✅ **Tailscale Funnel 集成**：通过 `wt public/private/links` 实现公开包分享
2. ✅ **服务器配置重载**：通过 `wt reload` 动态更新配置，无需重启
3. ✅ **增强的帮助系统**：全面的文档覆盖所有命令级别

---

## 🤝 贡献指南

欢迎贡献代码！请随时提交 Pull Request。

---

## 📄 许可证

本项目采用 GNU General Public License v3.0 许可证 - 详见 [LICENSE](LICENSE) 文件。

---

## 🙏 致谢

- 使用 [Go](https://golang.org/) 构建
- 受简单轻量级包管理解决方案需求的启发
