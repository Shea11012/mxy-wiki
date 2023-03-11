---
tags: []
date created: 2023-03-12 05:20
date modified: 2023-03-12 05:27
---

## commit message 格式

```
<type>{scope}: <subject>
空一行
<body>
空一行
<footer>
```

header 不能省略，body 和 footer 可以省略，任何一行都不能超过 72 个字符，这是为了避免自动换行

### header

type，subject 是必须的，scope 可以省略

#### type 类型

- feat：新功能（feature）
- fix：修补 bug
- docs：文档
- style：格式
- refactor：重构（既不是新增功能，也不是修改 bug）
- test：增加测试
- chore：构建过程或辅助工具变动

#### scope

scope 用于说明 commit 影响范围

#### subject

subject 是 commit 的简短描述

