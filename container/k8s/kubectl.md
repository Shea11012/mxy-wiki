---
date created: 2021-12-03 20:20
date modified: 2022-05-11 08:11
title: kubectl
---
## 常用命令
查看 pods

```bash
kubectl get pods -n <namespace>
```

展示资源的详细信息和相关event
```bash
kubectl describe pod <podname>
```

进入容器
```bash
# 默认选中pod中的第一个容器或者kubectl.kubernetes.io/default-container注解下的容器
kubectl exec -it <podname> bash

# 选中指定的容器
kubectl exec -it <podname> -c <container-name> bash
```

查看日志
```bash
kubectl logs <podname>
```

`kubectl get deployments` ：列出 deployments 信息

`kubectl get rs`：查看复制的数量

`kubectl scale deployname --replicas=4`：给指定的 deployname 创建 4 个复制

`kubectl get pods -o wide` ：查看 pods 的变化

`kubectl describe deployname` ：查看 deployment 的事件日志

`kubectl set image deployname deployname=deployname:version`：更新并设置应用版本

`kubectl rollout status deployname` ：确认更新状态

`kubectl rollout undo deployname`：回滚更新

