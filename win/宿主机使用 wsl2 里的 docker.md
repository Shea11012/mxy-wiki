---
tags: ['docker']
date created: 2023-01-10 19:22
date modified: 2023-02-19 02:23
---

## 安装 docker

[docker 安装和卸载](../container/docker/docker%20安装和卸载.md)

为 windows 构建 docker cli
```bash
git clone https://github.com/docker/cli.git  
cd cli
```

如果遇到没有 buildx 插件
>下载 buildx https://github.com/docker/buildx/releases
> 将下载的执行文件重命名为 buildx，移至 `$HOME/.docker/cli-plugins`

构建 docker cli
```bash
docker buildx bake --set binary.platform=windows/amd64

cp build/docker-windows-amd64.exe /mnt/c/tools
```

## 进入 win shell

需要将 tools 加入到 path 中
```shell
cd c:\tools  
ren docker-windows-amd64.exe docker.exe
```

方案 1：因为每次 wsl 重启 ip 地址都会变化，创建一个获取 wslip 的 powershell 脚本
```ps1
$wslip = wsl -- ip -o -4 -json addr list eth0 | ConvertFrom-Json | %{ $_.addr_info.local } | ?{ $_ }

if($wslip -eq $null){
    echo "wslip is empty"
    return
}

Write-Host "Setting Docker context 'wsl' to host=tcp://$($wslip):2376"
$command=docker context inspect wsl | Out-Null

if($? -eq $false) {
    echo "create docker context wsl"
    docker context create wsl --docker "host=tcp://$($wslip):2376"
} else {
    docker context update wsl --docker "host=tcp://$($wslip):2376"
}

docker context use wsl
```

方案 2：固定 wsl 的 ip 地址
创建 `~/.wslconfig` 文件
```config
[wsl2]
networkingMode=bridge
vmSwitch=my-switch
ipv6=true
```

创建 switch
- 打开 hyper-v 管理器
- 连接到服务器
- 点击虚拟机交换机管理器
- 选择外部类型，创建虚拟交换机，命名为 `my-switch`
- 创建完成就可以固定 wslip，也可以使用上面的脚本直接配置

