#!/bin/bash

# 项目根目录执行
# 统计所有 .go 文件代码行数（排除空行、注释行）
echo "正在统计 Go 项目有效代码行数..."
echo "----------------------------------------"

# 查找所有 go 文件，排除空行 + // 注释行
find . -type f -name "*.go" \
    -not -path "./vendor/*" \
    -not -path "./.git/*" |
    xargs wc -l |
    grep -v '^[[:space:]]*$' |
    grep -v '^[[:space:]]*//' |
    awk '{sum += $1} END {print "✅ 项目总代码行数（.go）: " sum " 行"}'

echo -e "\n----------------------------------------"
echo "文件总数：$(find . -type f -name "*.go" | wc -l) 个 Go 文件"
