---
tags: [k8s]
date created: 2022-05-11 07:55
date modified: 2022-05-11 08:01
title: Scheduler
---

# Scheduler
特殊的 Controller，Scheduler 的职责在于监控当前集群所有未调度的 Pod，并且获取当前集群所有节点的健康状况和资源使用情况，为待调度 Pod 选择最佳计算节点，完成调度。
调度阶段分为：
- Predict：过滤不能满足业务需求的节点，如资源不足、端口冲突等
- Priority：按既定要素将满足调度需求的节点评分，选择最佳节点
- Bind：将计算节点与 Pod 绑定，完成调度
