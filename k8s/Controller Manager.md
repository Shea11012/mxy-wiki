---
tags: [k8s]
date created: 2022-05-11 04:15
date modified: 2022-05-11 04:18
title: Controller Manage
---

# Controller Manage

Controller Manager 是确保集群启动的关键，作用是确保 kubernetes 遵循声明式系统规范，确保系统的真实状态与用户定义的期望状态一致。

Controller Manager 是多个控制器组合，每个 controller 事实上都是一个 control loop，负责监听其管控的对象，当对象发生变更时完成配置。当配置失败通过会触发自动重试，整个集群会在控制器不断重试的机制下确保最终一致性。