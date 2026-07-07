# 服务端日志功能说明 / Server Logging Documentation

## 概述 / Overview

water-repo 服务端现已集成全面的结构化日志系统，使用 Go 标准库 `log/slog` 实现。所有日志自动写入文件，便于问题排查和审计。

The water-repo server now integrates a comprehensive structured logging system using Go's standard library `log/slog`. All logs are automatically written to files for troubleshooting and auditing.

---

## 日志配置 / Log Configuration

### 日志文件位置 / Log File Location
- **路径 / Path**: `~/.local/state/water-repo/log.txt`
- **格式 / Format**: JSON（结构化日志）
- **模式 / Mode**: 追加写入（append mode）

### 日志级别 / Log Levels

| 级别 / Level | 用途 / Usage | 示例场景 / Example Scenarios |
|-------------|-------------|------------------------------|
| **INFO** | 正常操作记录 / Normal operations | 请求接收、操作成功、状态变更 |
| **WARN** | 警告信息 / Warnings | 权限不足、参数缺失、非致命错误 |
| **ERROR** | 错误信息 / Errors | 文件系统错误、网络错误、严重故障 |
| **DEBUG** | 调试信息 / Debug details | 详细的内部状态、中间过程（需启用 debug 模式） |

### 启用 Debug 模式 / Enable Debug Mode

```bash
# 在启动服务器时添加 debug 标志
wt server --debug

# 或在代码中设置 DEBUG = true
```

Debug 模式下：
- 日志同时输出到文件和控制台
- 显示 DEBUG 级别的详细日志
- 便于开发和调试

---

## 日志内容详解 / Log Content Details

### 1. 包管理操作 / Package Management Operations

#### 搜索包 / Search
```json
{
  "time": "2026-05-22T20:50:55.728+08:00",
  "level": "INFO",
  "msg": "Search request received",
  "method": "GET",
  "path": "/search",
  "query_name": "my-app",
  "remote_addr": "192.168.1.100:54321"
}
```

**记录内容 / Logged Information**:
- 请求方法、路径 / Request method, path
- 搜索关键词 / Search query
- 客户端地址 / Client address
- 搜索结果数量 / Results count
- 权限验证状态 / Token validation status

#### 上传包 / Upload
```json
{
  "time": "2026-05-22T20:51:10.123+08:00",
  "level": "INFO",
  "msg": "Package uploaded successfully",
  "name": "my-app.tar.gz",
  "size_bytes": 1048576,
  "size_human": "1.00 MB",
  "tag": "temp",
  "path": "/data/my-app.tar.gz"
}
```

**记录内容 / Logged Information**:
- 包名称、大小 / Package name, size
- 保存路径 / Save path
- 标签 / Tag
- 备份和恢复操作 / Backup and recovery operations
- 错误详情（如有）/ Error details (if any)

#### 下载包 / Install
```json
{
  "time": "2026-05-22T20:51:15.456+08:00",
  "level": "INFO",
  "msg": "Serving package file",
  "package_name": "my-app.tar.gz",
  "file_path": "/data/my-app.tar.gz",
  "tag": "temp",
  "remote_addr": "192.168.1.100:54322"
}
```

#### 删除包 / Delete
```json
{
  "time": "2026-05-22T20:51:20.789+08:00",
  "level": "INFO",
  "msg": "Package deleted successfully",
  "package_name": "old-app.tar.gz",
  "tag": "temp"
}
```

#### 重命名包 / Rename
```json
{
  "time": "2026-05-22T20:51:25.012+08:00",
  "level": "INFO",
  "msg": "Package renamed successfully",
  "old_name": "app-v1.tar.gz",
  "new_name": "app-v2.tar.gz",
  "tag": "static"
}
```

---

### 2. 标签管理 / Tag Management

#### 创建标签 / Add Tag
```json
{
  "time": "2026-05-22T20:52:00.123+08:00",
  "level": "INFO",
  "msg": "Tag added successfully",
  "tag_name": "frontend"
}
```

#### 更新标签 / Update Tag
```json
{
  "time": "2026-05-22T20:52:05.456+08:00",
  "level": "INFO",
  "msg": "Tag updated successfully",
  "package_name": "my-app.tar.gz",
  "old_tag": "temp",
  "new_tag": "frontend"
}
```

#### 删除标签 / Remove Tag
```json
{
  "time": "2026-05-22T20:52:10.789+08:00",
  "level": "INFO",
  "msg": "Tag removed successfully",
  "tag_name": "old-tag",
  "packages_moved_to_temp": 5
}
```

---

### 3. Tailscale Funnel 公网分享 / Public Sharing via Tailscale Funnel

#### 公开包 / Make Public
```json
{
  "time": "2026-05-22T20:53:00.123+08:00",
  "level": "INFO",
  "msg": "Package made public successfully",
  "package_name": "demo-app.tar.gz",
  "public_link": "https://xxx.ts.net/wt/public/demo-app.tar.gz"
}
```

**记录内容 / Logged Information**:
- Funnel 启动状态 / Funnel startup status
- 符号链接创建 / Symlink creation
- 公网链接生成 / Public link generation
- Tailscale 状态查询 / Tailscale status queries

#### 取消公开 / Make Private
```json
{
  "time": "2026-05-22T20:53:05.456+08:00",
  "level": "INFO",
  "msg": "Package made private successfully",
  "package_name": "demo-app.tar.gz"
}
```

#### 查看公开链接 / List Public Links
```json
{
  "time": "2026-05-22T20:53:10.789+08:00",
  "level": "INFO",
  "msg": "Returning public links",
  "links_count": 3
}
```

---

### 4. 权限验证 / Permission Validation

所有需要权限的操作都会记录 Token 验证结果：

```json
{
  "time": "2026-05-22T20:54:00.123+08:00",
  "level": "WARN",
  "msg": "Upload denied: insufficient write permission",
  "remote_addr": "192.168.1.200:12345"
}
```

**记录的权限事件 / Logged Permission Events**:
- 权限不足拒绝 / Access denied due to insufficient permissions
- Token 验证失败 / Token validation failures
- 标签级权限检查 / Tag-level permission checks

---

### 5. 系统操作 / System Operations

#### 配置重载 / Config Reload
```json
{
  "time": "2026-05-22T20:55:00.123+08:00",
  "level": "INFO",
  "msg": "Server config reloaded successfully",
  "config_path": "/home/cai/.config/water-repo/server_config.json"
}
```

#### 元数据同步 / Metadata Sync
```json
{
  "time": "2026-05-22T20:55:05.456+08:00",
  "level": "INFO",
  "msg": "Metadata synchronized successfully",
  "data_dir": "/data/packages"
}
```

---

## 日志查看方法 / How to View Logs

### 1. 使用 wt 命令 / Using wt Command

```bash
# 查看服务器日志
wt server log

# 实时跟踪日志（类似 tail -f）
wt server log | tail -f
```

### 2. 直接查看日志文件 / Direct File Access

```bash
# 查看最新 50 行
tail -n 50 ~/.local/state/water-repo/log.txt

# 实时跟踪
tail -f ~/.local/state/water-repo/log.txt

# 搜索特定关键词
grep "ERROR" ~/.local/state/water-repo/log.txt
grep "upload" ~/.local/state/water-repo/log.txt
```

### 3. 使用 jq 格式化 JSON 日志 / Format JSON Logs with jq

```bash
# 美化输出
cat ~/.local/state/water-repo/log.txt | jq .

# 只看 ERROR 级别
cat ~/.local/state/water-repo/log.txt | jq 'select(.level == "ERROR")'

# 只看特定消息
cat ~/.local/state/water-repo/log.txt | jq 'select(.msg | contains("upload"))'
```

---

## 日志示例场景 / Log Example Scenarios

### 场景 1: 用户上传包 / Scenario 1: User Uploads Package

```
INFO Upload request received method=POST path=/upload remote_addr=192.168.1.100:54321
INFO Package exists, will replace package_name=app.tar.gz existing_path=/data/app.tar.gz
INFO Removing backup file backup_path=/data/app.tar.gz.bak
INFO Package uploaded successfully name=app.tar.gz size_bytes=2097152 size_human=2.00 MB tag=temp path=/data/app.tar.gz
```

### 场景 2: 权限不足被拒绝 / Scenario 2: Access Denied Due to Insufficient Permissions

```
INFO Install request received method=GET path=/install package_name=secret-app.tar.gz remote_addr=192.168.1.200:12345
WARN Install request denied: insufficient install permission package_name=secret-app.tar.gz tag=private remote_addr=192.168.1.200:12345
```

### 场景 3: Tailscale Funnel 启动 / Scenario 3: Tailscale Funnel Startup

```
INFO Public request received method=POST path=/public package_name=demo.tar.gz remote_addr=192.168.1.100:54322
INFO Exposing package via Tailscale Funnel package_name=demo.tar.gz tag=temp
INFO Starting Tailscale Funnel data_dir=/data sub_path=/wt/public
INFO Tailscale Funnel started successfully link_dir=/data/public
INFO Package exposed successfully package_name=demo.tar.gz public_link=https://xxx.ts.net/wt/public/demo.tar.gz
```

### 场景 4: 错误处理 / Scenario 4: Error Handling

```
ERROR Upload failed: write error error="disk full" bytes_written=1024 package_name=large-file.tar.gz
INFO Recovered old package after write failure pkg_name=large-file.tar.gz
```

---

## 最佳实践 / Best Practices

### 1. 定期清理日志 / Regular Log Cleanup

```bash
# 保留最近 7 天的日志
find ~/.local/state/water-repo -name "log.txt" -mtime +7 -delete

# 或压缩旧日志
gzip ~/.local/state/water-repo/log.txt.1
```

### 2. 监控错误日志 / Monitor Error Logs

```bash
# 设置告警（示例）
tail -f ~/.local/state/water-repo/log.txt | grep "ERROR" | while read line; do
    echo "ALERT: $line" | mail -s "wt server error" admin@example.com
done
```

### 3. 日志轮转 / Log Rotation

对于长期运行的服务器，建议配置日志轮转：

```bash
# 使用 logrotate（Linux）
sudo tee /etc/logrotate.d/water-repo <<EOF
~/.local/state/water-repo/log.txt {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0644 cai cai
}
EOF
```

---

## 故障排查 / Troubleshooting

### 问题 1: 日志文件不存在 / Issue 1: Log File Not Found

**原因 / Cause**: 服务器未正常启动或权限问题

**解决 / Solution**:
```bash
# 检查目录权限
ls -la ~/.local/state/water-repo/

# 手动创建目录
mkdir -p ~/.local/state/water-repo
chmod 755 ~/.local/state/water-repo
```

### 问题 2: 日志过多 / Issue 2: Excessive Log Volume

**原因 / Cause**: Debug 模式启用或高频请求

**解决 / Solution**:
```bash
# 关闭 debug 模式重新启动
wt server  # 不加 --debug 标志

# 或过滤查看关键信息
grep -E "(ERROR|WARN)" ~/.local/state/water-repo/log.txt
```

### 问题 3: Tailscale Funnel 相关错误 / Issue 3: Tailscale Funnel Errors

**日志示例 / Log Example**:
```json
{"level":"ERROR","msg":"Failed to start Tailscale Funnel","error":"tailscale not found"}
```

**解决 / Solution**:
```bash
# 安装 Tailscale
curl -fsSL https://tailscale.com/install.sh | sh

# 登录 Tailscale
tailscale up

# 验证状态
tailscale status
```

---

## 技术实现细节 / Technical Implementation Details

### 使用的 slog Handler / slog Handler Used

- **生产环境 / Production**: `slog.NewJSONHandler()` → JSON 格式写入文件
- **Debug 模式 / Debug Mode**: `io.MultiWriter(file, stdout)` → 同时输出到文件和控制台

### 日志属性 / Log Attributes

每个日志条目包含的标准字段：
- `time`: ISO 8601 时间戳
- `level`: 日志级别（INFO/WARN/ERROR/DEBUG）
- `msg`: 人类可读的消息
- 上下文属性：`remote_addr`, `package_name`, `tag`, `error`, etc.

### 性能考虑 / Performance Considerations

- 异步写入：slog 默认同步写入，但对于本项目的规模（小团队）影响可忽略
- 日志级别过滤：Debug 日志在生产环境被过滤，减少 I/O
- 结构化日志：JSON 格式便于后续分析和工具处理

---

## 总结 / Summary

water-repo 的服务端日志系统提供了：

✅ **全面的操作记录** - 所有 API 请求和操作都有日志  
✅ **结构化格式** - JSON 格式便于解析和分析  
✅ **多级别支持** - INFO/WARN/ERROR/DEBUG 满足不同需求  
✅ **自动持久化** - 日志自动写入文件，无需额外配置  
✅ **易于查看** - 通过 `wt server log` 或直接查看文件  
✅ **Tailscale Funnel 集成** - 完整记录公网分享操作  

这为个人和小团队提供了足够的可观测性，便于问题排查和审计。

---

**版本 / Version**: v0.1.3+  
**最后更新 / Last Updated**: 2026-05-22
