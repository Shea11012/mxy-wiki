---
date created: 2021-11-30 21:22
date modified: 2021-11-30 21:22
title: git tips
---
1. auto correction in git

自动修复错误的命令，git 会选择运行第一个建议

```
git config --global help.autocorrect 1
```

2. count commits

统计自己的提交记录

```
git rev-list --count {branch}
```

3. optimize repo

用于清理仓库中 git 无法访问或孤立的对象

```
git gc --prune=now --aggressive
```

4. backup of untracked files

对未追踪的文件进行备份

```
git ls-files --others --exclude-standard -z | xargs -0 tar  rvf ~/backup-untracked.zip
```

5. search in git

在 git 中搜索

```
git rev-list --all | xargs git grep -F 'xxx'
```

