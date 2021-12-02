## k8s 命令

`kubectl get deployments` ：列出 deployments 信息

`kubectl get rs`：查看复制的数量

`kubectl scale deployname --replicas=4`：给指定的 deployname 创建 4 个复制

`kubectl get pods -o wide` ：查看 pods 的变化

`kubectl describe deployname` ：查看 deployment 的事件日志

`kubectl set image deployname deployname=deployname:version`：更新并设置应用版本

`kubectl rollout status deployname` ：确认更新状态

`kubectl rollout undo deployname`：回滚更新

