---
tags:
  - k8s
date created: 2022-05-11 04:15
date modified: 2023-12-01 18:10
---

# Deployment Controller

- 部署无状态应用
- 管理 pod 和 replicaSet
- 部署和滚动升级

## 创建一个 deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80
```

```shell
# 创建deployment
kubectl apply -f nginx-deployment.yaml

# 查看deployment状态
kubectl get deployments

# 查看deployment发布状态
kubectl rollout status deployment/nginx-deployment

# 查看replica set
kubectl get rs

# 查看pod标签
kubectl get pods --show-labels
```

## 更新 deployment

```shell
kubectl set deployment/nginx-deployment nginx=nginx:1.9.1
```

## 回滚 deployment

```shell
# 检查deployment更新历史
kubectl rollout history deployment/nginx-deployment

# 回滚到前一个版本
kubectl rollout undo deployment/nginx-deployment

# 回滚到指定版本
kubectl rollout undo deployment/nginx-deployment --to-revision=2
```

## 伸缩 deployment

```shell
# 伸缩deployment
kubectl scale deployment/nginx-deployment --replicas=10
```

## 暂停和继续 deployment

```shell
# 暂停
kubectl rollout pause deployment/nginx-deployment

# 继续
kubectl rollout resume deployment/nginx-deployment
```

# StatefulSetController

用来部署有状态应用

statefulset 中的 pod，每个 pod 挂载独立存储，如果一个 pod 出现故障，从其它节点启动一个同样名字的 pod，挂载上原来 pod 的存储继续提供它的状态提供服务。

- 保持 pod 启动顺序和唯一性
- 唯一的网络标识符，持久存储
- 有序的自动滚动更新

应用场景：MySQL、PostgreSQL、etcd 等

# DaemonSet

确保符合条件的 node 上运行一个 pod 副本。当有节点加入/离开时，也会为其创建/删除 pod

应用场景：运行日志收集进程、监控守护进程

# Job

用于执行一次性任务或批处理作业
Job 可以控制 pod 数量，确保一定量的 pod 完成任务后停止

# Scheduler

特殊的 Controller，Scheduler 的职责在于监控当前集群所有未调度的 Pod，并且获取当前集群所有节点的健康状况和资源使用情况，为待调度 Pod 选择最佳计算节点，完成调度。
调度阶段分为：
- Predict：过滤不能满足业务需求的节点，如资源不足、端口冲突等
- Priority：按既定要素将满足调度需求的节点评分，选择最佳节点
- Bind：将计算节点与 Pod 绑定，完成调度