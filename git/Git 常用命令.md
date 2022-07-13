---
tags: ["git"]
date created: 2021-06-05 00:30
date modified: 2022-04-22 06:37
title: Git 常用命令
---
#git 

### clone
命令 | 说明 
--- | ---
   git clone -b {branch} | {url} 克隆一个指定分支

### commit

命令 | 说明 
--- | ---
git commit -am {msg} |将所有已跟踪文件先 add 到暂存区，再执行 commit 命令
git commit —amed |修改最后一次提交的注释信息
   git reset —hard {commit-id} |移动 head 指向，将快照回滚到暂存区，并将暂存区文件还原到工作目录
   git reset —soft commit-id |只移动 head 指向，但不会将快照回滚到暂存区

### log
命令 | 说明 
--- | ---
  git log —decorate —all —graph —online |查看分支合并图
   git reflog |查看所有执行过的指令

### checkout
命令 | 说明 
--- | ---
   git rm —cached {filename} |只删除暂存区文件，取消追踪
   git rm -r —cached {filename} |删除远程文件，保留本地文件
   git checkout — filename |恢复被删除的文件
   git checkout -b {branch} |创建并切换到 branch 分支
   git checkout -b {branch-name} origin/branch-name | 创建本地和远程对应的分支
   git checkout —track origin/branch-name |跟踪指定上游分支

### branch
命令 | 说明 
--- | ---
git branch {branch-name}|创建分支
  git checkout {branch} |切换分支
  git branch |列出所有分支，当前分支会有 \*
  git branch -r | 查看远程分支
  git branch —merged |显示所有已经合并到当前分支的列表
  git branch —no-merged |显示所有没有合并到当前分支的分支列表
  git branch -d name |删除分支
  git branch -f branch commit-id |-f 强制改变分支的指向，将分支指向指定的 commit-id

### merge
命令 | 说明 
--- | ---
git merge {branch-name} | 合并 branch-name 分支到当前分支

### stash
命令 | 说明 
--- | ---
  git stash| 把当前工作区储藏
  git stash list |查看储藏列表
  git stash apply name |恢复储藏
  git stash drop name |删除储藏

### tag
命令 | 说明 
--- | ---
  git tag|看所有标签
  git tag -a tag-name -m 'msg'|含附注的标签，tag-name：版本号，将 -a 改为 -s 可以使用 GPG 私钥签名
  git tag -v tag-name| 验证标签
  git push origin tag-name| 分享标签，git push 默认不推送标签，推送所有本地标签使用 —tags
  git show tag-name |查看标签信息
  git tag -d tagname |删除本地标签
  git tag tag-name commit-id |给指定的提交添加标签

### push && fetch
命令 | 说明 
--- | ---
  git push {remote\_name} source:destination | 将 source 分支推送到远程的 destination 分支，remote\_name 一般为 origin
  git push {remote\_name} branch^:master | 将 branch 分支的父级推送到 master 上
  git fetch {remote\_name} branch^:master | 从远程分支的 branch 上获取父级提交更新到 master 分支上
  git push {remote\_name} :master| 表示删除 master
  git fetch {remote\_name} :slave | 创建一个 slave 本地分支

### merge && rebase && cherry-pick
命令 | 说明 
--- | ---
  git merge {branch}  | 合并指定分支到当前分支
  git rebase {branch} | 合并分支，创造线性的提交历史，将 branch 分支合并到当前分支

### extra
命令 | 说明 
--- | ---
  git stripspace |去掉行尾空白符，多个空行压缩成一行，必要时在文件末尾增加一个空行
  git show :/{query} |查询之前所有提交信息，找到条件相匹配的最近一条
  git remote add name url|添加一个远程仓库
  git update-index —assume-unchanged {path} |忽略指定的文件或目录，这样对于该文件的修改不会被 git 追踪到
  git checkout HEAD|改变 head 指向，可使用 HEAD^ 或 HEAD~num
  git rebase —interactive HEAD^|整理提交记录
  git reset HEAD^ |撤销更改
  git revert HEAD^ |撤销更改（创建一个副本，此副本与父节点一致）
  git cherry-pick commit-id commit-id commit-id... |整理提交，可以选择需要其他分支的哪些提交
  git config —global help.autocorrect 1 |自动修复错误的命令，git 会选择允许第一个建议
  git rev-list —count (branch} |统计自己的提交记录
  git gc —prune=now —aggressive |用于清理仓库中 git 无法访问或孤立的对象
  git ls-files —others —exclude-standard -z \| xargs -0 tar rvf ~/backup-untracked.zip |对未追踪的文件进行备份
  git rev-list —all \| xargs git grep -F 'xxx' | 在 git 中搜索