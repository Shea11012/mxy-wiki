wsl2 因为使用了hyper-v的虚拟机，所以和宿主机不在同一个网段中，在使用宿主机的vpn时，需要先获取宿主机的ip：

```bash
cat /etc/resolv.conf | grep 'nameserver' | awk '{print $2}'
```

为了便于使用还可以将这个设置一个脚本自动加载proxy



#### 改变wsl 的默认登录用户

> wsl -l # 查询版本
>
> \# 如 ubuntu20.04
>
> ubuntu2004 config --default-user shea # 需要保证 shea 这个用户已经被创建在该版本中



#### wsl2 与 trojan 搭配使用时遇到的端口被占用问题

启动Trojan时提示socks5的端口被占用了，使用 `netstat -ano | findstr 51837` 查看占用端口发现，该端口没有被占用，经过一番谷歌，hyper-v 虚拟机占用了该端口。

使用命令查看

```powershell
# 查看系统默认端口占用访问
netsh int ipv4 show dynamicport tcp

# 查看hyper-v 启动后的保留端口范围
netsh interface ipv4 show excludedportrange protocol=tcp  # 根据这个命令的输出，可以看到Trojan的默认socks5的端口是被占用了，所以更改Trojan默认的socks5的端口即可解决
```

#### wsl2 使用 git 时每次登录都需要手动添加key的解决办法

在没有key的时候，可以将win上的key复制到 wsl 上， `cp -r /mnt/c/Users/<username>/.ssh ~/.ssh`

win 上的key权限在linux下不适用，需要修改权限 `chmod 600 ~/.ssh/id_rsa`

> 上面步骤如果发生 permission denid，则需要在 /etc/wsl.conf （没有则新建）添加
>
> [automount]
>		enabled = true
> 		options = "metadata,umask=22,fmask=11"

为了不再每次打开一个新tab时，都需要手动添加key：

```
sudo apt install keychain

# 将下面的内容写到 bashrc 内
eval `keychain --eval --agents ssh id_rsa`
```

#### wsl2 docker exit 139

```
%userprofile%\.wslconfig # 此文件中写入配置

[wsl2]
kernelCommandLine = vsyscall=emulate
```

