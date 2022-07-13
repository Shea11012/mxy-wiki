---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: Laravel DingoApi开发复习
---
Api 指定访问版本推荐使用 Accept 头来指定

```
Accept:application/API_STANDARDS_TREE.API_SUBTYPE.v1+json
```

API_STANDARDS_TREE 三个值可选：

- x 本地开发或私有环境
- prs 未对外发布的，提供给公司 app、单页应用、桌面应用等
- vnd 对外发布，开放给所有用户

