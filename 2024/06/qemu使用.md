---
created: 2024-06-28
tags: [qemu]
updated: 2024-06-28
---

# 参数

`-enable-kvm`：启动 kvm 虚拟机，这种模式虚拟机运行速度更快
`-name`: 指定虚拟机 name
`-hda`：自动设置虚拟硬盘和镜像
`-drive`：自定义虚拟硬盘
`-cdrom`：指定 iso 镜像
`-m`：内存大小
`-vga`：设置 [[https://wiki.archlinux.org/title/QEMU#Graphic_card|vga类型]]
`-boot`：指定启动顺序。c(虚拟驱动硬盘），d(虚拟 CD-ROM），n(虚拟网络)
`-smp`：指定内核数
`-cpu host`：模拟主机处理器

启动一个 arch 虚拟机
```shell
qemu-system-x86_64 -enable-kvm -machine type=pc,accel=kvm -vga virtio -display gtk,gl=on -cpu host -m 4096 -hda Arch-Linux-x86_64-basic-20240601.239419.qcow2
```