---
tags: []
date created: 2023-02-18 22:13
date modified: 2023-02-18 22:40
---

## 查询启动时间节点

```bash
journalctl --list-boots
```

## 从指定时间查询日志

```bash
journalctl -S "xxxx-xx-xx xx:xx:xx" -U ""
```

