# Deployment
deployment 最常用于部署无状态服务。deployment能够更新pod和replicaSet

## 创建一个deployment
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

## 更新deployment
```shell
kubectl set deployment/nginx-deployment nginx=nginx:1.9.1
```

## 回滚deployment
```shell
# 检查deployment更新历史
kubectl rollout history deployment/nginx-deployment

# 回滚到前一个版本
kubectl rollout undo deployment/nginx-deployment

# 回滚到指定版本
kubectl rollout undo deployment/nginx-deployment --to-revision=2
```

## 伸缩deployment
```shell
# 伸缩deployment
kubectl scale deployment/nginx-deployment --replicas=10
```

## 暂停和继续deployment
```shell
# 暂停
kubectl rollout pause deployment/nginx-deployment

# 继续
kubectl rollout resume deployment/nginx-deployment
```