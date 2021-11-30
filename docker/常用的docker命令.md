| 命令                                              | 说明                             |
| ------------------------------------------------- | -------------------------------- |
| docker exec -it -u 用户 容器 bash                 | 指定用户登录指定容器             |
| docker cp 容器名:文件路径  宿主机对应的文件路径   | 从容器拷贝到宿主机，同样可以反之 |
| docker rmi $(docker images -f "dangling=true" -q) | 删除标签为 none 的镜像           |
| docker rm $(docker ps -q -s)                      | 删除停止运行的容器               |

