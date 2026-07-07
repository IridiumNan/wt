# funnel plan for global file share

## command

```bash
wt public <package name> 
```

create a soft link to /tmp/wt/public

```go
cmd := exec.Command("ln", "-s", packageName, /tmp/wt/public/)
```

```bash
ln -s DataDir/package_name /tmp/wt/public
```

and exec sudo tailscale funnel

- if funnel is not on

```bash
mkdir -p /tmp/wt/public
tailscale funnel --bg --set-path /wt/public /tmp/wt/public
```

user can visit this site directly without wt client

stop share file

```bash
wt private <package name> 
```
