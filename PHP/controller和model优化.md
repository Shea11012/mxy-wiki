---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: controller和model优化
---
对于中大型项目的 controller 和 model 的优化

使用 service 将外部行为注入到 service，在 service 使用外部行为，将 service 注入到 controller

使用 `Repository` 辅助 model ，将相关的数据库逻辑封装在不同的 repository ，方便中大型项目的维护。

方式：

1. 将 model 依赖注入到 repository
2. 将数据库逻辑写在 repository
3. 将 repository 依赖注入到 service

优点：

1. 解决了 controller 臃肿的问题
2. 符合 SOLID 的单一职责原则：数据库逻辑写在 repository ，没写在 controller
3. 符合 SOLID 的依赖反转原则：controller 并非直接相依于 repository，而是将 repository 依赖注入进 controller

参考 [Service](https://oomusou.io/laravel/service/) [Repository](https://oomusou.io/laravel/repository/)

