# 多服务器管理计划

- 查看可用的服务器

```bash
wt list-servers

# 输出样例
# alias             host
# home-server    http://100.120.81.67:12212
# cloud-server   http://13.51.31.41:12212
# ...
```

- 切换服务器

```bash
wt change-server <alias>/<host>

# 输出样例
# change the server as <alias> => <host>        successfully

# server <alias> not found

# exec the list-servers command automatically
# available servers
# ... 
```

- 添加服务器

```bash
wt add-server <alias> <host>

# 输出样例
# add the server <alias> => <host>      successfully
```

- 删除服务器

```bash
wt del-server <alias>/<host>

# 输出样例
# delete the server <alias> => <host>  successfully

# server <alias> not found
# exec the list-servers command automatically
# available servers
# ... 

```
