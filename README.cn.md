# water-repo

[![License](https://img.shields.io/badge/license-GPL--3.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.25.9+-blue.svg)](https://golang.org/)
[![Version](https://img.shields.io/badge/version-v0.0.3-orange.svg)]()

超轻量级个人/小团队仓库管理工具,基于 Go 开发,纯 CLI 交互(无图形化界面),天生适配服务器环境;CS 架构设计,单二进制文件开箱即用,零依赖零配置,没有任何多余功能。

---

## 🚀 核心特性

- **轻量级**：单二进制文件,零依赖
- **简单易用**：直观的命令行操作
- **安全可靠**：基于 Token 的权限控制系统
- **灵活组织**：标签化管理包资源
- **跨平台支持**：支持 Linux、macOS 和 Windows

---

## 📦 安装方式

### 下载二进制文件

#### Linux

```bash
wget https://repo.waterman.xin/apps/water-repo/wt-lastest-linux-amd64
chmod +x wt-lastest-linux-amd64
mv wt-lastest-linux-amd64 ~/.local/bin/wt  # 或任何在 PATH 中的目录
```

#### macOS (Apple Silicon)

```bash
curl -LO https://repo.waterman.xin/apps/water-repo/wt-lastest-darwin-arm64
chmod +x wt-lastest-darwin-arm64
mv wt-lastest-darwin-arm64 ~/.local/bin/wt  # 或任何在 PATH 中的目录
```

#### Windows

从 [wt-lastest-windows-amd64.exe](https://repo.waterman.xin/apps/water-repo/wt-lastest-windows-amd64.exe) 下载

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
| 显示帮助 | `wt help` | 显示帮助信息 |

---

## 🔐 权限管理

采用极简的 **三级 Token 权限体系**,无需用户注册登录,通过配置文件中的 Token 验证操作权限,未配置对应有效 Token 会直接拒绝操作。

| 权限等级 | 可执行操作 | 适用场景 |
|---------|-----------|---------|
| Read | `search`、`list`、`info` | 仅允许查看仓库内容 |
| Install | `install` | 允许下载使用仓库中的包 |
| Write | `upload`、`mv`、`rm` | 允许管理和修改仓库内容 |

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
| 新增自定义标签 | `wt tag add <标签名>` | 创建新的分类标签 |
| 删除自定义标签 | `wt tag rm <标签名>` | 删除指定标签；原归属该标签的所有包会自动回退到 `temp` 标签 |
| 批量清理标签下的所有包 | `wt clear <标签名>` | **高危操作！** 永久删除该标签下的所有包,不可恢复 |

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
    ]
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

## 🚧 开发中功能 (v0.1.1+)

以下功能正在开发中,将在后续版本中陆续上线,现有配置文件和核心命令保持完全向后兼容。

### 核心待实现（下一版本）

1. **按标签的细粒度权限控制**
   - 基于现有标签系统扩展,无需引入复杂的角色和用户体系
   - 支持为每个标签单独配置 Read/Install/Write 权限 Token
   - 标签权限优先级高于全局权限,实现"不同人管理不同分类的包"
   - 示例：前端组仅拥有 `frontend` 标签的写入权限,后端组仅拥有 `backend` 标签的写入权限

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

### 规划中功能

1. **`wt install-tag` 批量下载**：一条命令下载指定标签下的所有包
2. **`wt public` 一键公开分享**：将指定包公开给整个局域网,任何人无需配置 Token 即可下载
3. **下载进度显示**：命令行实时显示下载速度和进度条
4. **断点续传**：支持大文件断点续传,中断后无需重新下载
5. **`wt config` 配置管理命令**：无需手动编辑 JSON 文件,通过命令行修改配置

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
