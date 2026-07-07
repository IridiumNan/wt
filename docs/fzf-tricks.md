# fzf tricks

- search and get package info

```bash
wt search <pkg> | fzf | xargs -I {} wt info {}
```

- search and remove a package

```bash
wt search <pkg> | fzf | xargs -I {} wt rm {}
```

- search and download a package

```bash
wt search <pkg> | fzf | xargs -I {} wt install {}
```

- find a file and upload

```bash
wt upload $(fzf)
```

- find tag and list packages

```bash
wt ls $(wt tag ls)
```

- tag a package

```bash
wt tag $(wt ls | fzf) $(wt tag ls | fzf)
```

- search and public a package

```bash
wt search <pkg> | fzf | xargs -I {} wt public
```
