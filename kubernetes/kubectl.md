---
tags: 
date created: 2021-12-03 20:20
date modified: 2023-12-01 15:11
title: kubectl
---

# 基础命令

| 命令    | 介绍                            |
| ------- | ------------------------------- |
| create  | 通过文件名或标准输入创建资源    |
| expose  | 将一个资源公开为一个新的 service |
| run     | 在集群中运行一个特定镜像        |
| set     | 在对象上设置特定的功能          |
| get     | 显示一个或多个资源              |
| explain | 文档参考资料                    |
| edit    | 使用默认编辑器编辑一个资源      |
| delete  | 删除资源                        |

# 部署命令

| 命令           | 介绍                                             |
| -------------- | ------------------------------------------------ |
| rollout        | 管理资源发布                                     |
| rolling-update | 对给定的控制器滚动更新                           |
| scale          | 扩、缩容 pod 数量，deployment、replicaset、rc、job |
| autoscale      | 创建一个自动扩或缩容的 pod 数量|

# 集群管理命令

| 命令         | 介绍              |
| ------------ | ----------------- |
| certificate  | 修改证书资源      |
| cluster-info | 显示集群信息      |
| top          | 显示资源          |
| cordon       | 标记节点不可调度  |
| uncordon     | 标记节点可被调度  |
| drain        | 驱逐节点上的应用  |
| taint        | 修改节点 taint 标记 |

# 故障和调试命令

| 命令         | 介绍                                          |
| ------------ | --------------------------------------------- |
| describe     | 显示特定资源或资源组的详细信息                |
| logs         | 显示 pod 内容器日志                           |
| attach       | 追加一个进程到 pod 的容器内                   |
| exec         | 在容器内执行命令                              |
| port-forward | 端口转发                                      |
| proxy        | 在 localhost 和 api server 之间创建一个 proxy |
| cp           | 拷贝文件或目录到容器                          |
| auth         | 检查授权                                      |

# 其它命令

| 命令         | 介绍                        |
| ------------ | --------------------------- |
| apply        | 通过文件名应用配置          |
| patch        | 更新资源字段                |
| replace      | 替换资源                    |
| convert      | 不同 api 版本之间转换配置文件 |
| label        | 更新资源标签                |
| annotate     | 更新资源注释                |
| completion   | kubectl 补全                 |
| api-versions | api 版本                     |
| config       | kubeconfig                  |

