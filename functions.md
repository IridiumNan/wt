# 命令汇总

> 版本：v0.1.2

## 已实现命令

- 查看指定标签的包 -> 根据标签查权限 -> 有全局则直接跳过 -> 没有的话验证标签权限

```bash
wt ls [tag]
```

- 搜索包 -> 有权限则返回搜索结果， 没有权限访问的从结果中去除

```bash
wt search <package name>
```

- 获取包信息 -> 有权限则返回包详情

```bash
wt info <package name>
```

- 获取包信息 -> 下载包 -> 有权限则开放下载， 没有权限则拒绝

```bash
wt install <package name>
```

- 上传包 -> 直接上传到 temp 标签

```bash
wt upload <local path> <package name>
```

- 重命名包 -> 获取包信息 -> 有权限则开放， 没有权限则拒绝

```bash
wt mv <old name> <new name>
```

- 删除包 -> 有权限则开放， 没有权限则拒绝

```bash
wt rm <package name>
```

- 同步磁盘和metadata -> 先列出所有可以操作的标签 -> 对这些标签的数据进行磁盘同步

```bash
wt sync
```

- 客户端配置 -> 不查权限

```bash
wt config ...
```

- 添加标签 -> 如果标签存在则拒绝， 不存在则创建并将所有全局 Token 拷贝到标签中

```bash
wt tag add <tag name>
```

- 删除标签 -> 内置标签（temp/static）不可删除，删除后包回退到 temp

```bash
wt tag rm <tag name>
```

- 给某个包更换标签 -> 验证对原来标签的权限和对新标签的权限 -> 如果都有权限则允许

```bash
wt tag <package name> <tag name>
```

- 列出当前客户端可见的标签列表

```bash
wt tag list
wt tag ls
```

- 服务端给某个标签添加新的权限 Token -> 标签级细粒度权限

```bash
wt server tag <tag name> read_token add <new client token>
wt server tag <tag name> install_token add <new client token>
wt server tag <tag name> write_token add <new client token>
```

## 规划中命令（v0.2+）

- 清空标签下所有包（高危操作，需确认）

```bash
wt clear <tag name>
```

- 批量下载标签下所有包

```bash
wt install-tag <tag name>
```

## 权限链路

```
请求到达 → TokenOk(w, r, method, tag)
         → isAccess(auth)
         → ① 先查全局 Token 列表
              ├─ 匹配 → 放行（全局优先，超管兜底）
              └─ 不匹配 ↓
         → ② 再查标签 Token 列表
              ├─ 匹配 → 放行（标签级扩展权限）
              └─ 不匹配 → 拒绝
```

## 相关函数

- `GetAvailableTags()` 获取所有可以访问(下载, 写入)的标签
- `isAccess()` 全局优先 → 标签兜底的权限校验
- `AddTagToken()` 给指定标签添加指定类型的 Token
