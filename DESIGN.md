# water-repo 设计文档 / Design Document

> **文档定位**：本文档记录 water-repo 的核心设计决策与背后的「为什么」——面向想要理解项目设计哲学的用户和贡献者。具体的使用说明和配置示例请参阅 [README.md](README.md)。
>
> **Document purpose**: This document captures the core design decisions of water-repo and the "why" behind them — for users and contributors who want to understand the design philosophy. For usage instructions and configuration examples, see [README.md](README.md).

---

## 第一部分：中文

### 1. 设计初衷

water-repo 的诞生源于几个非常具体的痛点：

**手动复制 URL 的摩擦**。之前用 AI 辅助写的仓库管理脚本，每次下载都要手动复制 URL 到服务器上。对于一个重度键盘用户来说，这种鼠标操作是持续的打断。

**目录结构是冗余的**。存储的是已经打包好的压缩包、镜像、配置文件——这些本质上都是单一扁平形态的制品。为它们维护嵌套的目录结构，像是在给不需要分类的东西强行分类。

**小型团队的权限困境**。想过给团队开放仓库，但企业级 ACL / RBAC 系统对于三五个人的小组来说是过度设计。需要一种「看一眼就懂」的权限模型，而不是一套需要培训才能上手的管理后台。

**学 Go 的实践项目**。刚好在学 Golang，这个场景足够简单又足够完整，是理想的练手项目。

这些痛点的交集，催生了 water-repo 的核心定位：

> 一个**单二进制、零配置、纯 CLI** 的个人/小团队包管理工具。像 `apt` 一样简单，但服务于你自己的制品分发。

---

### 2. 核心设计原则

water-repo 的每一个设计决策，都可以追溯到以下原则。当新功能提议与原则冲突时，原则优先。

| 原则 | 含义 | 反例 |
|------|------|------|
| **零摩擦上手** | 单二进制文件，不需要数据库、不需要配置文件即可启动 | 需要先装 PostgreSQL、配 YAML |
| **极简权限** | 三级 Token（Read / Install / Write），没有用户系统、没有角色 | RBAC with roles, groups, policies |
| **扁平分类** | Tag 替代目录结构，一个包只属于一个 Tag | 嵌套目录、多级命名空间 |
| **制品，不是源码** | 只管理打包好的成品，不关心构建过程 | CI/CD 集成、源码版本管理 |
| **CLI 优先** | 所有操作通过命令行完成，没有 GUI、没有 Web 面板 | Web UI dashboard |
| **去中心化就绪** | 一个二进制同时包含 server 和 client，任何节点都可以既是服务端也是客户端 | 中心服务器 + 代理架构 |

这些原则共同构成了 water-repo 的边界：**它不做什么，和它做什么同样重要。**

---

### 3. Tag 标签系统 —— 为什么不是目录？

#### 3.1 问题

传统的包管理方式会用目录来组织文件：

```
repo/
├── temp/
│   ├── random-script.sh
│   └── test-config.json
├── static/
│   └── nginx.conf
└── iso/
    └── ubuntu-24.04.iso
```

这个结构在两种情况下是合理的：(1) 包的数量非常大，需要层级导航；(2) 同一个包有多个关联文件需要放在一起。

但在 water-repo 的场景下，这两个条件都不成立：

- 包的规模是 **几十到几百个**，一个扁平列表完全可以承载
- 每个包是 **单一文件**（压缩包/镜像/配置），不需要关联文件组
- 用户的操作路径是 `search → info → install`，不需要浏览目录树

#### 3.2 方案

用 **Tag（标签）** 替代目录：每个包有且仅有一个 Tag，系统内置 `temp`（临时分享）和 `static`（长期存储）两个不可删除的标签，用户可以创建自定义标签（如 `iso`、`frontend`、`backend`）。

```
操作：wt ls temp
结果：直接列出 temp 标签下所有包，不需要知道它们在哪个「文件夹」
```

#### 3.3 设计考量

- **一个包一个 Tag**：简化了数据模型和操作语义。移动包就是改一个字段，不需要处理「一个文件属于多个目录」的复杂情况。
- **内置标签不可删除**：`temp` 和 `static` 作为兜底分类，保证系统在任何状态下都有可用的标签。
- **删除标签时包回退到 `temp`**：防止包变成「无标签」的孤儿状态，保持数据完整性。
- **不引入嵌套标签**：这会让 Tag 退化成目录，违背扁平化原则。

---

### 4. Token 权限系统 —— 为什么不是用户系统？

#### 4.1 问题

传统的权限管理思路是：先创建用户 → 给用户分配角色 → 给角色配置权限。这在 water-repo 的场景下有两个问题：

1. **管理负担**：小型团队（3-5 人）不需要用户生命周期管理（入职/离职/密码重置）
2. **认知负担**：用户想知道「我能下载这个包吗？」需要理解用户→角色→权限→资源的映射链

water-repo 的使用场景是：分享一个配置文件给同事，或者让团队成员下载构建好的镜像。这些场景需要的权限模型是「有这个 Token 就能做这件事」，而不是「你是张三，张三属于开发者组，开发者组有读写权限」。

#### 4.2 方案

**三级 Token 权限**：

| 级别 | 允许的操作 | 典型场景 |
|------|-----------|---------|
| Read Token | `search`, `list`, `info` | 浏览仓库，看看有什么 |
| Install Token | 以上 + `install` | 下载使用包 |
| Write Token | 以上 + `upload`, `mv`, `rm`, `tag` | 管理仓库内容 |

- 服务端可以配置**多个 Token**（每个级别可以有一个 Token 列表），客户端只配置**单个 Token**
- 没有用户注册、没有登录态、没有 Session
- 每次请求携带 Token，服务端做简单的列表匹配

#### 4.3 设计考量

- **Token 而非用户名/密码**：Token 可以直接分发、直接撤销（从列表删除即可），不需要维护用户数据库。
- **三级而非二级**：Read 和 Install 分开是因为「能看到」和「能下载」在信任模型上是不同的。你可能想让所有人浏览仓库目录，但只让信任的人下载。
- **服务端多 Token、客户端单 Token**：服务端可以给不同人发不同的 Token（方便单独撤销），客户端只需要知道自己该用哪个。
- **与 `wt public` 的关系**（计划中）：对于真正「任何人可下载」的场景，`wt public` 绕过 Token 检查，不需要让所有人配置 Token。

---

### 5. 架构设计

#### 5.1 整体架构

```
┌──────────────────────────────────────────┐
│                  wt 二进制                 │
│                                           │
│  ┌─────────────┐    ┌──────────────────┐ │
│  │   server 模式 │    │    client 模式    │ │
│  │  (wt server) │    │  (wt search/...) │ │
│  │             │    │                  │ │
│  │  HTTP API   │◄──►│  HTTP Client     │ │
│  │  + Store    │    │  + CLI           │ │
│  └─────────────┘    └──────────────────┘ │
└──────────────────────────────────────────┘
```

同一个 `wt` 二进制，通过子命令切换角色：

- `wt server` → 启动 HTTP 服务，暴露 REST API
- `wt search / install / upload / ...` → 作为客户端调用远程或本地的 server

#### 5.2 模块划分

```
cmd/wt/main.go          # 唯一入口，解析子命令分发
internal/
  model/                # 共享数据结构（Package, Config, Auth）
  config/               # 配置文件加载与路径管理
  store/                # 内存存储 + 磁盘同步
  server/               # HTTP 路由、Handler、中间件
  client/               # HTTP 客户端封装
  command/              # 各子命令的业务逻辑编排
  presets/              # 交互式配置的预设问题
pkg/
  fileutil/             # 文件操作工具
  httphelper/           # HTTP 工具
  loghelper/            # 日志工具
```

#### 5.3 关键设计决策

**为什么是 CS 架构而不是纯本地？**
因为核心场景是「分发」——一个人上传，其他人下载。纯本地工具无法满足这个需求。但 CS 不意味着中心化：由于 server 和 client 在同一个二进制里，任何节点都可以作为服务端，天然支持去中心化部署。

**为什么用标准库 `net/http` 而不是框架？**
zero-dependency 原则。项目规模小，标准库的路由和中间件完全够用，引入 Gin/Chi 等框架只会增加二进制体积和理解成本。

**为什么没有用 `cobra` 做 CLI？**
同样的原因。初期用 Go 标准库的 `flag` 包足够，后期如果需要更复杂的子命令嵌套，可以再考虑迁移——但迁移成本很低，因为 `internal/command` 已经隔离了命令逻辑。

**为什么数据目录默认是当前目录？**
零配置原则的直接体现。用户可以在任意目录下直接 `wt server` 启动服务，数据就落在这个目录里——不需要创建特定路径、不需要迁移文件、不需要指定挂载点。这对于「在服务器上临时起一个仓库共享几个文件」的场景尤为重要：没有前置步骤，即起即用。如果需要持久化到特定位置，`-d` 参数提供显式控制。

**为什么包名就是文件名？**
极简原则的体现。上传 `my-app.tar.gz`，下载下来就是 `my-app.tar.gz`——没有中间层的命名映射，没有元数据与文件实体的分离。这带来三个好处：(1) 用户对下载产物有明确的预期；(2) 不需要额外的索引层来维护「逻辑名→物理文件」的映射；(3) 可以直接在文件系统里 `ls` 看到仓库的实际状态，和 `meta.json` 互相印证。在任何应用场景下，拒绝引入不必要的命名抽象都是最佳选择。

---

### 6. 去中心化架构哲学

#### 6.1 传统中心化 vs. water-repo 的去中心化

传统的仓库管理工具（如 Nexus、Artifactory、Docker Registry）都是中心化的：一个中央服务器存储所有制品，所有客户端从这个唯一的源拉取。这种架构在大型组织中运转良好，但对小型团队和临时场景来说，有不可忽视的成本：

- **单点故障**：中央服务器宕机，所有人都无法下载
- **网络瓶颈**：大文件下载受限于服务器带宽
- **维护负担**：需要专人管理服务器、配置备份、监控
- **隐性门槛**：不能「想分享就分享」，必须先部署中央服务

water-repo 的答案是：**让每个节点既是客户端也是服务端**。同一个 `wt` 二进制，你可以用它上传包到一个远端，也可以在自己机器上起一个 server 让别人来拉。没有「中央」的概念——或者更准确地说，任何节点都可以成为临时的中心。

#### 6.2 去中心化工作场景

想象一个典型的内网场景：

> 你和公司的其他人都连在同一个局域网里，每个人都装了 `wt` 并启动了 server。你想下载一个常用的系统镜像——你出于隐私或安全考虑，不希望从外部网络下载。
>
> 通过 `wt discover`，你的客户端自动发现了局域网内所有运行中的 wt 服务端，并把它们加入你的镜像源列表。你执行 `wt search ubuntu-24.04`，客户端并行查询所有已知的镜像源，合并返回结果。你看到同事 B 的机器上有这个包，而且已经通过 `wt public` 公开了。你执行 `wt install ubuntu-24.04`，客户端自动选择最快的源（局域网传输速度远快于外网），几秒内下载完成。
>
> 在这个场景中，没有中心服务器。公司的所有设备都是仓库的一部分，各取所需。

#### 6.3 设计层面的含义

去中心化不是一个「功能」，而是架构层面的选择，它影响了多个模块的设计：

- **server 和 client 在同一二进制**：任何安装了 `wt` 的节点都具备服务能力，不需要额外安装「服务端版本」
- **镜像源是列表而非单点**：客户端连接配置是一个 URL 列表，而非单一地址——这为 Discover 和 Mirror 功能预留了空间
- **Public 机制独立于 Token**：「公开」意味着绕过认证，这让包可以在没有信任关系的节点之间自由流动，是去中心化协作的关键
- **零配置启动**：不需要声明「我是中心服务器」还是「我是边缘节点」，启动即服务

这种架构选择让 water-repo 的使用方式从「搭建服务 → 配置客户端 → 使用」变成了「启动 → 使用」——每一步都是可选的，每一步都可以动态调整。

---

### 7. 数据存储设计 —— 为什么没有数据库？

#### 7.1 方案：内存 + JSON 快照

water-repo 使用 **内存 Map + 定时 JSON 磁盘同步** 的存储方案：

- 所有包的元数据（名称、Tag、大小、修改时间）存储在内存的 `MetaData` 结构体中
- 每次增删改操作后，将内存中的完整元数据**原子写入**磁盘的 JSON 文件
- 启动时从 JSON 文件恢复元数据到内存
- 包的实际文件直接存储在磁盘上，按包名平铺

#### 7.2 为什么不用 SQLite / BoltDB？

| 考虑 | 内存 + JSON | SQLite / BoltDB |
|------|------------|-----------------|
| 依赖 | 零（Go 标准库） | CGO 或额外依赖 |
| 包数量 | 几十到几百，全量加载毫无压力 | 适用于成千上万的记录 |
| 查询模式 | 只有「按名称匹配」和「按 Tag 列出」两种 | 复杂查询 |
| 运维 | JSON 文件可以直接查看、手动修改 | 需要专用工具 |
| 并发 | 单 writer + 读写锁足够 | 内置并发控制 |

对于 water-repo 的规模（几十到几百个包），数据库是杀鸡用牛刀。JSON 文件的透明性反而是一个特性——用户可以 `cat meta.json` 直接看到仓库状态。

#### 7.3 并发控制

使用 `sync.RWMutex`：多个读操作可以并发，写操作独占锁。操作模式以读为主（search/list/info），极少并发写（upload/rm），这个方案天然适合。

---

### 8. 有意的「不做」

好的设计文档也应该记录「决定不做」的事情。以下功能被明确排除在 water-repo 的范围之外：

| 不做 | 原因 |
|------|------|
| **嵌套目录 / 多级命名空间** | Tag 扁平分类已满足需求，目录增加复杂度 |
| **用户注册 / 登录系统** | Token 认证足够，用户管理是负担 |
| **版本管理（SemVer）** | 这不是包管理器（如 npm/apt），版本信息由上传者通过命名表达 |
| **递归上传（自动打包目录）** | 压缩和打包是上传者的责任，water-repo 只管理制品 |
| **GUI / Web 面板** | CLI 优先，不引入前端技术栈 |
| **数据库依赖** | 几十到几百个包用数据库是过度设计 |
| **CI/CD 集成** | 不在分发层解决问题，保持边界清晰 |
| **包依赖解析** | 只做存储和分发，不做构建 |

---

### 9. 未来方向（v0.1.1+）

以下功能在设计上已经预留了空间，与现有架构不冲突：

#### 9.1 按 Tag 的细粒度权限

**现状**：Write Token 可以操作所有包。
**计划**：允许每个 Tag 单独配置 Token，Tag 级权限优先于全局权限。这样前端组只能写入 `frontend` Tag 的包，后端组只能写入 `backend` Tag 的包。

设计考量：这不需要引入「用户」概念，只是把 Token 的作用域从「全局」细化到「Tag」，保持了 Token 模型的简单性。

#### 9.2 动态镜像源 + 局域网发现

**现状**：客户端连接单一 server。
**计划**（去中心化的具体实现，设计哲学见[第 6 节](#6-去中心化架构哲学)）：

1. `wt mirror add <url>` 添加多个远程 wt 服务端为镜像源
2. 搜索时并行查询所有镜像源，合并去重
3. 下载时自动选择最快源，失败自动切换
4. `wt discover` 自动发现局域网内的 wt 服务端，加入镜像源列表

设计考量：这利用了 water-repo 的 CS 同体特性——每个 wt server 都可以作为别人的 mirror。局域网内的设备互相发现，形成去中心化的仓库网络。

#### 9.3 `wt public` 一键公开分享

**场景**：临时分享一个包给很多人，不想让每个人配置 Token。
**方案**：`wt public <package-name>` 将该包标记为公开，绕过 Token 检查即可下载。

设计考量：这与 Token 系统不矛盾——Token 是「默认需要认证」，Public 是「显式声明公开」，后者是前者的例外而非替代。

#### 9.4 体验增强（批量下载、进度条、断点续传）

- `wt install-tag <tag>`：一条命令下载整个 Tag 下的所有包
- 下载进度条：对大文件的 CLI 进度反馈
- 断点续传：大文件下载中断后从断点继续

---

## Part 2: English

### 1. Design Motivation

water-repo was born from a few very specific pain points:

**The friction of manually copying URLs**. A previous AI-assisted repository script required manually copying URLs to the server for every download. For a heavy keyboard user, this mouse interaction was a constant interruption.

**Directory structures are redundant**. The artifacts being stored — compressed archives, images, configuration files — are inherently flat, single-file entities. Maintaining nested directories for them feels like forcing categories onto things that don't need them.

**The permission dilemma of small teams**. Opening the repository to a team is desirable, but enterprise ACL/RBAC systems are overkill for groups of 3-5 people. The ideal permission model is one you can understand at a glance, not one that requires training.

**A project to learn Go**. The scope is simple enough to be manageable, yet complete enough to be a meaningful learning exercise.

These intersecting pain points define water-repo's core positioning:

> A **single-binary, zero-config, pure-CLI** package management tool for individuals and small teams. As simple as `apt`, but built for your own artifact distribution.

---

### 2. Core Design Principles

Every design decision in water-repo traces back to these principles. When a proposed feature conflicts, the principle wins.

| Principle | Meaning | Counter-example |
|-----------|---------|-----------------|
| **Zero-friction onboarding** | Single binary, no database, works without a config file | Requires PostgreSQL and YAML setup first |
| **Minimal permissions** | Three-tier Token (Read / Install / Write), no user system, no roles | RBAC with roles, groups, policies |
| **Flat categorization** | Tags replace directory structures; one package, one tag | Nested directories, multi-level namespaces |
| **Artifacts, not source** | Manages only finished, packaged outputs; doesn't care about build processes | CI/CD integration, source version management |
| **CLI-first** | All operations via command line; no GUI, no web dashboard | Web UI dashboard |
| **Decentralization-ready** | One binary contains both server and client; any node can serve or consume | Central server + agent architecture |

These principles define water-repo's boundary: **what it doesn't do is as important as what it does.**

---

### 3. Tag System — Why Not Directories?

#### 3.1 The Problem

Traditional package management uses directories to organize files:

```
repo/
├── temp/
│   ├── random-script.sh
│   └── test-config.json
├── static/
│   └── nginx.conf
└── iso/
    └── ubuntu-24.04.iso
```

This structure makes sense when (1) the number of packages is very large, requiring hierarchical navigation; and (2) a single package has multiple associated files.

Neither condition holds for water-repo:

- Package count is in the **tens to low hundreds** — a flat list works perfectly
- Each package is a **single file** (archive/image/config), no associated file groups needed
- The user workflow is `search → info → install`, which doesn't involve browsing directory trees

#### 3.2 The Solution

**Tags instead of directories**: every package has exactly one tag. The system provides two built-in, non-deletable tags — `temp` (temporary sharing) and `static` (long-term storage). Users can create custom tags (e.g., `iso`, `frontend`, `backend`).

```
Command: wt ls temp
Result: directly lists all packages under the "temp" tag
        no need to know which "folder" they're in
```

#### 3.3 Design Rationale

- **One package, one tag**: simplifies the data model and operational semantics. Moving a package is a single field change — no "file belongs to multiple directories" complexity.
- **Built-in tags are non-deletable**: `temp` and `static` serve as fallback categories, ensuring the system always has usable tags.
- **Packages revert to `temp` on tag deletion**: prevents orphan packages (tagless state), maintaining data integrity.
- **No nested tags**: nested tags would degenerate back into directories, violating the flatness principle.

---

### 4. Token Permission System — Why Not a User System?

#### 4.1 The Problem

The traditional approach: create users → assign roles → configure permissions. This has two problems for water-repo's use case:

1. **Management burden**: small teams (3-5 people) don't need user lifecycle management (onboarding/offboarding/password resets)
2. **Cognitive burden**: answering "Can I download this package?" requires tracing a user→role→permission→resource mapping chain

water-repo's usage: sharing a config file with a colleague, or letting team members download built images. These scenarios require a permission model of "you have this Token, you can do this thing" — not "you are Zhang San, Zhang San belongs to the developer group, the developer group has read-write access."

#### 4.2 The Solution

**Three-tier Token permissions**:

| Level | Allowed Operations | Typical Scenario |
|-------|-------------------|------------------|
| Read Token | `search`, `list`, `info` | Browse the repository |
| Install Token | Above + `install` | Download packages |
| Write Token | Above + `upload`, `mv`, `rm`, `tag` | Manage repository content |

- Server can configure **multiple tokens** per level (a token list); client configures **a single token**
- No user registration, no login state, no sessions
- Each request carries a token; the server does a simple list membership check

#### 4.3 Design Rationale

- **Tokens, not username/password**: tokens can be distributed directly and revoked instantly (remove from the list). No user database to maintain.
- **Three tiers, not two**: Read and Install are separated because "can see" and "can download" are different trust levels. You may want everyone to browse the catalog, but only trusted people to download.
- **Server multi-token, client single-token**: the server can issue different tokens to different people (for individual revocation); the client only needs to know its own.
- **Relationship with `wt public`** (planned): for truly "anyone can download" scenarios, `wt public` bypasses token checks entirely — no need to configure tokens for everyone.

---

### 5. Architecture Design

#### 5.1 High-Level Architecture

```
┌──────────────────────────────────────────┐
│               wt binary                   │
│                                           │
│  ┌─────────────┐    ┌──────────────────┐ │
│  │  server mode │    │   client mode    │ │
│  │  (wt server) │    │ (wt search/...)  │ │
│  │             │    │                  │ │
│  │  HTTP API   │◄──►│  HTTP Client     │ │
│  │  + Store    │    │  + CLI           │ │
│  └─────────────┘    └──────────────────┘ │
└──────────────────────────────────────────┘
```

The same `wt` binary switches roles via subcommands:

- `wt server` → starts an HTTP server, exposes REST API
- `wt search / install / upload / ...` → acts as a client, calling a remote or local server

#### 5.2 Module Breakdown

```
cmd/wt/main.go          # Single entry point; parses subcommands and dispatches
internal/
  model/                # Shared data structures (Package, Config, Auth)
  config/               # Config file loading and path management
  store/                # In-memory storage + disk sync
  server/               # HTTP routing, handlers, middleware
  client/               # HTTP client wrapper
  command/              # Business logic orchestration for each subcommand
  presets/              # Interactive config question presets
pkg/
  fileutil/             # File operation utilities
  httphelper/           # HTTP utilities
  loghelper/            # Logging utilities
```

#### 5.3 Key Architectural Decisions

**Why CS architecture instead of purely local?**
Because the core use case is *distribution* — one person uploads, others download. A purely local tool can't serve this. But CS doesn't mean centralized: since server and client live in the same binary, any node can be a server, naturally supporting decentralized deployment.

**Why `net/http` standard library instead of a framework?**
Zero-dependency principle. The project is small enough that the standard library's routing and middleware are sufficient. Adding Gin/Chi would only increase binary size and cognitive overhead.

**Why not `cobra` for CLI?**
Same reason. The `flag` package from the standard library is adequate for the current subcommand complexity. If more complex nesting is needed later, migration cost is low because `internal/command` already isolates command logic.

**Why does the data directory default to the current directory?**
A direct expression of the zero-config principle. Users run `wt server` in any directory and data lands right there — no need to create specific paths, migrate files, or specify mount points. This is especially important for the "spin up a temporary repo on a server to share a few files" scenario: zero pre-steps, instant-on. For persistent storage at a specific location, the `-d` flag provides explicit control.

**Why does the package name equal the file name?**
A direct expression of the minimalism principle. Upload `my-app.tar.gz`, download `my-app.tar.gz` — no intermediate naming layer, no separation between metadata and the file entity. This yields three benefits: (1) users have a clear expectation of what they'll get; (2) no extra index layer to maintain "logical-name → physical-file" mappings; (3) you can `ls` the data directory and see the actual repository state, cross-validating against `meta.json`. In every scenario, rejecting unnecessary naming abstraction is the optimal choice.

---

### 6. Decentralized Architecture Philosophy

#### 6.1 Traditional Centralization vs. water-repo's Decentralization

Traditional repository tools (Nexus, Artifactory, Docker Registry) are centralized: one central server stores all artifacts, and all clients pull from this single source. This architecture works well in large organizations, but for small teams and ad-hoc scenarios, it carries non-trivial costs:

- **Single point of failure**: central server goes down, nobody can download
- **Network bottleneck**: large file downloads constrained by server bandwidth
- **Maintenance burden**: someone must manage the server, backup configs, monitor health
- **Implicit barrier**: you can't "just share" — you must first deploy a central service

water-repo's answer: **let every node be both client and server**. The same `wt` binary can upload packages to a remote, and simultaneously serve packages to others from your own machine. There is no "center" — or more precisely, any node can become a temporary center.

#### 6.2 A Decentralized Work Scenario

Imagine a typical LAN scenario:

> You and your coworkers are on the same company intranet. Everyone has `wt` installed and server running. You need a common system image — but for privacy or security reasons, you don't want to download from external networks.
>
> Via `wt discover`, your client automatically finds all running wt servers on the LAN and adds them to your mirror list. You run `wt search ubuntu-24.04` — the client queries all known mirrors in parallel, merging results. You see coworker B's machine has this package, already made public via `wt public`. You run `wt install ubuntu-24.04` — the client auto-selects the fastest source (LAN speeds far exceed external), and the download completes in seconds.
>
> In this scenario, there is no central server. Every device in the company is part of the repository, each taking what it needs.

#### 6.3 Design Implications

Decentralization is not a "feature" — it's an architectural choice that shapes multiple module designs:

- **Server and client in the same binary**: any node with `wt` installed can serve, no separate "server edition" required
- **Mirror sources are a list, not a single endpoint**: the client connection config is a URL list, not a single address — this reserves space for Discover and Mirror functionality
- **Public mechanism is independent of Token**: "public" means bypassing authentication, enabling packages to flow freely between nodes with no trust relationship — the key to decentralized collaboration
- **Zero-config startup**: no need to declare "I am the central server" or "I am an edge node" — start and you're serving

This architectural choice transforms the water-repo usage pattern from "deploy service → configure client → use" into "start → use" — every step is optional, every step can be adjusted dynamically.

---

### 7. Data Storage — Why No Database?

#### 7.1 Solution: In-Memory Map + JSON Snapshots

water-repo uses an **in-memory Map + atomic JSON disk sync** storage strategy:

- All package metadata (name, tag, size, modification time) is held in an in-memory `MetaData` struct
- After every create/update/delete operation, the full metadata is **atomically written** to a JSON file on disk
- On startup, metadata is restored from the JSON file into memory
- Actual package files are stored directly on disk, flat by package name

#### 7.2 Why Not SQLite / BoltDB?

| Consideration | In-Memory + JSON | SQLite / BoltDB |
|---------------|------------------|-----------------|
| Dependencies | Zero (Go stdlib) | CGO or external dependency |
| Package count | Tens to low hundreds — full load is trivial | Suited for thousands+ records |
| Query patterns | Only two: "match by name" and "list by tag" | Complex queries |
| Operations | JSON file can be viewed and edited directly | Requires specialized tools |
| Concurrency | Single writer + RWMutex is sufficient | Built-in concurrency control |

For water-repo's scale (tens to low hundreds of packages), a database is using a sledgehammer to crack a nut. The transparency of JSON is actually a feature — users can `cat meta.json` to directly see the repository state.

#### 7.3 Concurrency Control

Uses `sync.RWMutex`: multiple reads can proceed concurrently; writes acquire an exclusive lock. The access pattern is read-heavy (search/list/info) with rare concurrent writes (upload/rm), making this a natural fit.

---

### 8. Deliberate Omissions

A good design document also records what was *decided against*. The following features are explicitly out of water-repo's scope:

| Won't Do | Reason |
|----------|--------|
| **Nested directories / multi-level namespaces** | Flat Tag categorization is sufficient; directories add complexity |
| **User registration / login system** | Token authentication is enough; user management is a burden |
| **Version management (SemVer)** | This is not a package manager (like npm/apt); version info is expressed by the uploader through naming |
| **Recursive upload (auto-pack directories)** | Compression and packaging are the uploader's responsibility; water-repo manages only artifacts |
| **GUI / Web dashboard** | CLI-first; no frontend tech stack |
| **Database dependency** | Overkill for tens to low hundreds of packages |
| **CI/CD integration** | Problem solved at a different layer; keep boundaries clear |
| **Package dependency resolution** | Storage and distribution only, no build |

---

### 9. Future Direction (v0.1.1+)

The following features have architectural headroom reserved; they don't conflict with the current design:

#### 9.1 Tag-Level Fine-Grained Permissions

**Current**: Write Token grants access to all packages.
**Planned**: Allow configuring separate tokens per Tag; Tag-level permissions take priority over global ones. A frontend team would only have write access to the `frontend` tag's packages.

Design note: this doesn't introduce a "user" concept — it merely scopes tokens from "global" to "per-tag", preserving the Token model's simplicity.

#### 9.2 Dynamic Mirror Sources + LAN Discovery

**Current**: Client connects to a single server.
**Planned** (concrete implementation of decentralization; see [Section 6](#6-decentralized-architecture-philosophy) for design philosophy):

1. `wt mirror add <url>` — add multiple remote wt servers as mirror sources
2. Search queries all mirrors in parallel, merging and deduplicating results
3. Downloads automatically select the fastest available source, with failover
4. `wt discover` — auto-discovers wt servers on the LAN, adds them as mirrors

Design note: this leverages water-repo's CS-same-binary nature — every wt server can be someone else's mirror. Devices on the LAN discover each other, forming a decentralized repository network.

#### 9.3 `wt public` — One-Click Public Sharing

**Scenario**: temporarily share a package with many people without configuring tokens for everyone.
**Solution**: `wt public <package-name>` marks the package as public, bypassing token checks for downloads.

Design note: this doesn't contradict the Token system — Token means "authentication required by default"; Public means "explicitly declared open." The latter is an exception, not a replacement.

#### 9.4 UX Enhancements (batch download, progress bar, resume)

- `wt install-tag <tag>`: download all packages under a tag in one command
- Download progress bar: CLI feedback for large file transfers
- Resume broken downloads: continue interrupted large-file downloads from the breakpoint

---

## 版本历史 / Version History

| 版本 / Version | 日期 / Date | 变更 / Changes |
|---------------|-------------|----------------|
| v0.1 | 2026-05 | 初始版本，记录 v0.0.1 的设计决策 / Initial version, documents v0.0.1 design decisions |
