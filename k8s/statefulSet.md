---
date created: 2022-03-12 14:19
date modified: 2022-03-12 14:23
title: statefulSet
---
# statefulSet
- 管理有状态的应用程序
- 管理 pod 时，确保 pod 有一个按顺序增长的 ID
- 每个 pod 都对应一个特有的持久话存储标识

## 使用场景
- 稳定、唯一的网络标识
- 每个 pod 始终对应各自的存储路径
- 按顺序的增加副本、减少副本，并在减少副本时执行清理
- 按顺序自动执行滚动更新