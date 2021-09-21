**linux 中无法使用 host.docker.internal 访问宿主机**
解决：
-  使用 docker 命令行增加 `--add-host=host.docker.internal:host-gateway`
- 使用 docker-compose 文件，给对应的 service 增加以下

```yaml
extra-hosts:
   - "host.docker.internal:host-gateway"
```