---
tags:
  - elasticsearch
  - wsl
date created: 2021-11-30 21:22
date modified: 2023-10-11 02:49
title: elasticsearch 遇到的问题
---

# max_map_count 问题

修改 vm.max_map_count=262144

win10 wsl 修改

```
wsl -d docker-desktop
sysctl -w vm.max_map_count=262144

# 上面的方法不能持久，每次重启就失效
通过 win+R 输入 shell:startup 进入到目录中，在目录底下创建一个 .bat 的脚本
@echo off
wsl -d docker-desktop sysctl -w vm.max_map_count=262144
这样每次电脑开机都会执行这个脚本了
```

# 索引别名

**tips：别名查询也是有一定限制的**

假设有一个名字叫 `category` 的索引，真实索引可以包含一个版本号：`category_v1,category_v2` 

```
PUT /category_v1  创建一个真实索引
PUT /category_v1/_alias/category  设置别名
PUT /category_v1/_alias/computer  设置一个有指定过滤提交的别名
{
	"routing":"12",
	"filter":{
		"term":{
			"user_id":12
		}
	}
}

创建一个 category_v2 的版本
PUT /category_v2

一个别名可以指向多个索引，在添加别名到新索引的同时必须从旧索引中删除它，这个操作需要原子化
POST /_aliases
{
	"acions":[
		{"remove":{"index":"category_v1","alias":"category"}},
		{"add":{"index":"category_v2":"alias":"category"}}
	]
}
```

