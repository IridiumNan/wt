# Tag 功能实现计划

> 版本：v0.1.1  
> 基于：DESIGN.md §3（标签系统）、§9.1（细粒度权限）、functions.md（命令汇总）

---

## 一、权限验证逻辑（核心）

### 1.1 检查链路

```
请求到达 → TokenOk(w, r, method, tag)
         → isAccess(auth)
         → ① 先查全局 Token 列表
              ├─ 匹配 → 放行（全局优先，全局是超管兜底）
              └─ 不匹配 ↓
         → ② 再查标签 Token 列表
              ├─ 匹配 → 放行（标签级扩展权限）
              └─ 不匹配 → 拒绝
```

**设计依据**（`functions.md`）：

> "查看指定标签的包 → 根据标签查权限 → 有全局则直接跳过 → 没有的话验证标签权限"

- 全局 Token：超管级兜底，拥有全局 Token 的用户对所有标签都有对应权限
- 标签 Token：扩展权限，只对特定标签生效

### 1.2 各命令权限规则

| 命令 | 权限 | 校验逻辑 |
|------|------|---------|
| `wt ls [tag]` | Read | 查指定 tag 的 Read（全局→标签） |
| `wt search <name>` | Read | 全局搜索 → `GetAvailableTags(Read)` 过滤不可读的包 |
| `wt info <name>` | Read | 先查包的 tag → 校验该 tag 的 Read |
| `wt install <name>` | Install | 先查包的 tag → 校验该 tag 的 Install |
| `wt upload <file> [name]` | Write | 固定写入 `temp` → 校验 `temp` 的 Write |
| `wt mv <old> <new>` | Write | 查包的 tag → 校验该 tag 的 Write |
| `wt rm <name>` | Write | 查包的 tag → 校验该 tag 的 Write |
| `wt sync` | Write | `GetAvailableTags(Write)` → 只同步有权限的标签 |
| `wt tag add <tag>` | Write | 全局 Write（管理操作） |
| `wt tag rm <tag>` | Write | 全局 Write（管理操作） |
| `wt tag <pkg> <new-tag>` | Write | 旧 tag 的 Write + 新 tag 的 Write |
| `wt clear <tag>` | Write | 该 tag 的 Write |

### 1.3 `isAccess` 改造

```go
// 当前：只查全局
// 改造后：全局优先 → 标签兜底
func isAccess(auth model.Auth) (valid bool) {
    clientToken := auth.Token
    if clientToken == "" {
        return false
    }
    // 第一步：全局 Token 检查
    if tokenInList(clientToken, config.GetTokenList(auth.WtMethod)) {
        return true
    }
    // 第二步：标签 Token 检查（tag 为空则跳过）
    if auth.Tag != "" {
        tagTokens, err := config.GetTagTokenList(auth.Tag, auth.WtMethod)
        if err == nil && tokenInList(clientToken, tagTokens) {
            return true
        }
    }
    return false
}
```

### 1.4 `GetAvailableTags` 函数

```go
// 返回客户端对指定操作有权限的所有标签
func GetAvailableTags(clientToken string, wtMethod model.WTMethod) []string {
    // 内置标签 + 自定义标签
    allTags := []string{config.DefaultTagTemp, config.DefaultTagStatic}
    allTags = append(allTags, config.GetAllTags()...)

    var accessible []string
    for _, tag := range allTags {
        // 全局 Token 检查
        if tokenInList(clientToken, config.GetTokenList(wtMethod)) {
            accessible = append(accessible, tag)
            continue
        }
        // 标签 Token 检查
        tagTokens, err := config.GetTagTokenList(tag, wtMethod)
        if err == nil && tokenInList(clientToken, tagTokens) {
            accessible = append(accessible, tag)
        }
    }
    return accessible
}
```

### 1.5 标签创建时的 Token 继承

`wt tag add <tag-name>` 创建标签时，服务端自动将**所有全局 Token** 拷贝到新标签的各权限列表中：

```
新标签.ReadTokens    ← serverConfig.ReadToken
新标签.InstallTokens ← serverConfig.InstallToken
新标签.WriteTokens   ← serverConfig.WriteToken
```

这意味着默认情况下，拥有全局权限的用户对新标签也有同等权限。

---

## 二、实现阶段

### 阶段一：标签 CRUD 命令（核心）

#### 2.1 需要新增的命令

| 命令 | 说明 |
|------|------|
| `wt tag add <tag-name>` | 创建自定义标签，自动继承全局 Token |
| `wt tag rm <tag-name>` | 删除标签（内置标签不可删），包回退到 `temp` |
| `wt tag <pkg-name> <target-tag>` | 修改包的归属标签 |

#### 2.2 Store 层改动

**文件：`internal/store/memory.go`**

新增函数：

```go
// AddTag 创建新标签，初始化空的包列表
func AddTag(tagName string) error

// RemoveTag 删除标签，所有包回退到 temp
func RemoveTag(tagName string) error

// ChangePackageTag 修改包的标签（复用已有 UpdateTag）
// 已存在 UpdateTag(pkgName, newTag)，检查是否需要加额外校验
```

校验项：
- `AddTag`：标签名不能为空，不能与已有标签重复
- `RemoveTag`：不能删除 `temp` 和 `static`
- `ChangePackageTag`：包必须存在，目标标签必须存在

#### 2.3 Server API 层改动

**文件：`internal/server/handler.go`** — 新增 handler：

```go
// POST /tag/add?name=<tag-name>
// 权限：全局 Write Token
func tagAddHandler(w http.ResponseWriter, r *http.Request)

// DELETE /tag/rm?name=<tag-name>  
// 权限：全局 Write Token
// 行为：删除标签，包回退到 temp
func tagRmHandler(w http.ResponseWriter, r *http.Request)

// PUT /tag/change?pkg=<pkg-name>&tag=<new-tag>
// 权限：旧 tag + 新 tag 的 Write Token
func tagChangeHandler(w http.ResponseWriter, r *http.Request)
```

**文件：`internal/server/router.go`** — 注册路由：

```go
mux.HandleFunc("/tag/add", tagAddHandler)
mux.HandleFunc("/tag/rm", tagRmHandler)
mux.HandleFunc("/tag/change", tagChangeHandler)
```

**文件：`internal/server/middleware.go`** — 改造：

```go
// 改造 isAccess：全局优先 → 标签兜底（见 §1.3）
func isAccess(auth model.Auth) bool

// 改造 TokenOk：正确传递 tag 参数
func TokenOk(w http.ResponseWriter, r *http.Request, wtMethod model.WTMethod, tag string) bool

// 新增：获取客户端的可用标签列表
func GetAvailableTags(clientToken string, wtMethod model.WTMethod) []string
```

**文件：`internal/server/handler.go`** — 修改已有 handler 的 TokenOk 调用：

| Handler | 当前调用 | 修改为 |
|---------|---------|--------|
| `searchHandler` | `TokenOk(w, r, WTRead, "")` | `TokenOk(w, r, WTRead, "")` — 保持不变，search 后过滤 |
| `infoHandler` | 复用 searchHandler | 同上 |
| `installHandler` | `TokenOk(w, r, WTInstall, "")` | 先查包的 tag，再 `TokenOk(w, r, WTInstall, pkgTag)` |
| `uploadHandler` | `TokenOk(w, r, WTWrite, "")` | `TokenOk(w, r, WTWrite, "temp")` |
| `mvHandler` | `TokenOk(w, r, WTWrite, "")` | 先查包的 tag，再 `TokenOk(w, r, WTWrite, pkgTag)` |
| `rmHandler` | `TokenOk(w, r, WTWrite, "")` | 先查包的 tag，再 `TokenOk(w, r, WTWrite, pkgTag)` |
| `listHandler` | `TokenOk(w, r, WTRead, tag)` | 保持不变（已经传了 tag） |
| `syncHandler` | `TokenOk(w, r, WTWrite, tag)` | 改为 `GetAvailableTags` 驱动同步 |

#### 2.4 Client 层改动

**文件：`internal/client/client.go`** — 新增函数：

```go
// TagAddRequest 创建新标签
func TagAddRequest(tagName string) error

// TagRmRequest 删除标签
func TagRmRequest(tagName string) error

// TagChangeRequest 修改包的标签
func TagChangeRequest(pkgName string, newTag string) error
```

#### 2.5 Command 层改动

**文件：`internal/command/client.go`** — 新增 `"tag"` 命令分发：

```
wt tag add <tag-name>        → handleTwoTarget("tag", args) 中判断 secondTarget == "add"
wt tag rm <tag-name>         → 同上 secondTarget == "rm"  
wt tag <pkg-name> <tag-name> → handleTwoTarget("tag", args) 调用 TagChangeRequest
```

建议在 `handleTwoTarget` 的 switch 中新增 `"tag"` case，内部根据 `firstTarget` 的值分流：
- `firstTarget == "add"` → `TagAddRequest(secondTarget)`  
- `firstTarget == "rm"` → `TagRmRequest(secondTarget)`
- 其他 → `TagChangeRequest(firstTarget, secondTarget)`（修改包标签）

#### 2.6 帮助手册更新

**文件：`internal/command/Default.txt`** — 在 COMMANDS 段添加：

```
  tag        Tag management (add/rm/change package tag)
```

**文件：`internal/command/Simple.txt`** — 添加：

```
TAG MANAGEMENT:
  wt tag add <name>             Create a new tag
  wt tag rm <name>              Remove a tag (packages revert to temp)
  wt tag <pkg> <tag>            Change package tag
```

**文件：`internal/command/Advance.txt`** — 添加详细说明：

```
10. TAG MANAGEMENT
    Usage: wt tag <subcommand> [arguments]
    Subcommands:
      add <tag-name>        Create a new custom tag
      rm <tag-name>         Remove a custom tag
      <pkg> <tag>           Move package to a different tag
    Examples:
      wt tag add frontend          # Create a new tag
      wt tag rm frontend           # Remove the tag
      wt tag my-app frontend       # Move my-app to frontend tag
```

---

### 阶段二：标签级细粒度权限

#### 2.7 需要新增的命令

| 命令 | 说明 |
|------|------|
| `wt server tag <tag> read_token add <token>` | 给标签添加读 Token |
| `wt server tag <tag> install_token add <token>` | 给标签添加安装 Token |
| `wt server tag <tag> write_token add <token>` | 给标签添加写 Token |
| `wt tag list` | 列出当前客户端可见的标签 |

#### 2.8 Server 命令层改动

**文件：`internal/command/server.go`** — 在 `execCommand` 中新增 `"tag"` 分支：

```
wt server tag <tag-name> read_token add <token>
                                 ↓
args = ["server", "tag", "<tag>", "read_token", "add", "<token>"]
                                   ↑         ↑          ↑        ↑
                              cmdFlag   FirstTarget SecondTarget ThirdTarget
```

当前 `execCommand` 的 `len(args) == 5` 分支处理 `config`，需要扩展处理 `tag`。

#### 2.9 Config 层改动

**文件：`internal/config/server.go`** — 在 `AlterServerConfig` 或新增函数中处理 tag token：

```go
// AddTagToken 给指定标签添加指定类型的 Token
func AddTagToken(tag string, wtMethod model.WTMethod, token string) error
```

#### 2.10 Client 层改动

**文件：`internal/client/client.go`** — 新增：

```go
// TagListRequest 获取客户端有权访问的标签列表
func TagListRequest() error
```

对应的 API 端点改造（`/taglist` handler）：

```go
// GET /taglist
// 返回当前客户端可访问的标签列表
func tagListHandler(w http.ResponseWriter, r *http.Request) {
    // 从 Header 中提取 token
    // 调用 GetAvailableTags(token, WTRead)
    // 返回标签列表
}
```

---

### 阶段三：批量操作

#### 2.11 需要新增的命令

| 命令 | 说明 |
|------|------|
| `wt clear <tag-name>` | 清空标签下所有包（高危，需确认） |
| `wt install-tag <tag>` | 批量下载标签下所有包（规划中，v0.2+） |

#### 2.12 Store 层

`DeletePackageByTag` 已存在，无需新增。

#### 2.13 Server API 层

```go
// DELETE /clear?tag=<tag-name>
// 权限：该 tag 的 Write Token
// 行为：永久删除该标签下所有包，需额外确认机制
func clearHandler(w http.ResponseWriter, r *http.Request)
```

#### 2.14 路由

```go
mux.HandleFunc("/clear", clearHandler)
```

---

## 三、数据流全景图

```
┌─────────────────────────────────────────────────────────────┐
│  wt tag add frontend                                        │
│    │                                                        │
│    ▼                                                        │
│  command/client.go → handleTwoTarget("tag", ["add","frontend"])│
│    │                                                        │
│    ▼                                                        │
│  client.TagAddRequest("frontend")                           │
│    │  Header: X-Write-Token                                 │
│    ▼                                                        │
│  POST /tag/add?name=frontend                                │
│    │                                                        │
│    ▼                                                        │
│  server.tagAddHandler                                       │
│    │  TokenOk(w, r, WTWrite, "")  ← 全局 Write 检查          │
│    │  store.AddTag("frontend")                              │
│    │  config.AddTagTokens("frontend", globalTokens...)      │
│    │  config.writeServerConfig()  ← 持久化到磁盘             │
│    ▼                                                        │
│  Response: {"code":200, "status":"success"}                  │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  wt ls frontend                                             │
│    │                                                        │
│    ▼                                                        │
│  client.ListRequest("frontend")                             │
│    │  Header: X-Read-Token                                  │
│    ▼                                                        │
│  GET /list?tag=frontend                                     │
│    │                                                        │
│    ▼                                                        │
│  server.listHandler                                         │
│    │  TokenOk(w, r, WTRead, "frontend")                     │
│    │    → isAccess: 全局 Read? → 标签 Read?                   │
│    │  store.ListPackagesByTag("frontend")                   │
│    ▼                                                        │
│  Response: ["pkg1", "pkg2", ...]                            │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  wt tag my-app frontend                                     │
│    │                                                        │
│    ▼                                                        │
│  client.TagChangeRequest("my-app", "frontend")               │
│    │  Header: X-Write-Token                                 │
│    ▼                                                        │
│  PUT /tag/change?pkg=my-app&tag=frontend                    │
│    │                                                        │
│    ▼                                                        │
│  server.tagChangeHandler                                    │
│    │  查 my-app 当前 tag（如 "temp"）                         │
│    │  TokenOk(w, r, WTWrite, "temp")     ← 旧 tag Write      │
│    │  TokenOk(w, r, WTWrite, "frontend") ← 新 tag Write      │
│    │  store.UpdateTag("my-app", "frontend")                  │
│    ▼                                                        │
│  Response: {"code":200, "status":"success"}                  │
└─────────────────────────────────────────────────────────────┘
```

---

## 四、文件改动清单

### 阶段一改动文件

| 文件 | 改动类型 | 内容 |
|------|---------|------|
| `internal/store/memory.go` | 新增 | `AddTag()`, `RemoveTag()` |
| `internal/server/middleware.go` | 改造 | `isAccess()` 全局→标签链路 |
| `internal/server/handler.go` | 新增+改造 | 3 个 tag handler + 已有 handler 传 tag |
| `internal/server/router.go` | 新增 | 3 条 tag 路由 |
| `internal/client/client.go` | 新增 | `TagAddRequest`, `TagRmRequest`, `TagChangeRequest` |
| `internal/command/client.go` | 改造 | `"tag"` 命令分发 |
| `internal/command/Default.txt` | 更新 | 帮助文本 |
| `internal/command/Simple.txt` | 更新 | 快速参考 |
| `internal/command/Advance.txt` | 更新 | 详细文档 |

### 阶段二改动文件

| 文件 | 改动类型 | 内容 |
|------|---------|------|
| `internal/config/server.go` | 新增 | `AddTagToken()` |
| `internal/command/server.go` | 改造 | `"tag"` 子命令分发 |
| `internal/server/handler.go` | 改造 | `tagListHandler` 实现 |
| `internal/client/client.go` | 新增 | `TagListRequest` |

### 阶段三改动文件

| 文件 | 改动类型 | 内容 |
|------|---------|------|
| `internal/server/handler.go` | 新增 | `clearHandler` |
| `internal/server/router.go` | 新增 | `/clear` 路由 |
| `internal/client/client.go` | 新增 | `ClearRequest` |
| `internal/command/client.go` | 改造 | `"clear"` 命令 |

---

## 五、测试要点

### 5.1 权限链路测试

| 场景 | 预期 |
|------|------|
| 全局 Read Token 用户 `ls` 任意标签 | 放行 |
| 仅有 `frontend` Read Token 用户 `ls frontend` | 放行 |
| 仅有 `frontend` Read Token 用户 `ls backend` | 拒绝 |
| 全局 Write Token 用户 `upload` | 放行（写入 temp） |
| 仅有 `temp` Write Token 用户 `upload` | 放行 |
| 仅有 `frontend` Write Token 用户 `upload` | 拒绝（写入 temp 无权限） |

### 5.2 标签生命周期测试

| 场景 | 预期 |
|------|------|
| 创建标签 `wt tag add mytag` | 成功，TagTokenMap 中有该标签且继承了全局 Token |
| 重复创建同名标签 | 拒绝，提示已存在 |
| 删除 `temp` 或 `static` | 拒绝，提示内置标签不可删除 |
| 删除自定义标签 | 成功，其下所有包回退到 `temp` |
| `wt tag <pkg> <new-tag>` | 包从旧标签移除，加入新标签 |

### 5.3 Search 过滤测试

| 场景 | 预期 |
|------|------|
| 全局 Read Token 搜索 | 返回所有匹配的包 |
| 仅有 `frontend` Read Token 搜索 | 只返回 tag 为 `frontend` 的匹配包 |
| 无 Read Token 搜索 | 返回空或拒绝 |

---

## 六、不做的边界

| 不做 | 原因 |
|------|------|
| 嵌套标签 | DESIGN.md §3.3 明确拒绝 |
| 一个包多个标签 | 保持模型简洁 |
| 标签重命名 | 可通过删除+创建实现，避免额外复杂度 |
| 标签描述/元数据 | 保持扁平，按需在 v0.2+ 评估 |
