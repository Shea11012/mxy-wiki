---
title: "Linux 文件系统基础与 ext 机制"
created: 2026-07-22
updated: 2026-07-22
status: needs-review
type: synthesis
tags:
  - linux
  - linux/filesystem
confidence: INFERRED
---
# Linux 文件系统基础与 ext 机制

> 本页由重复旧笔记整合而来；正文按原页面分节保留，尚未重新核验技术时效性。

## 文件系统创建与格式化

## 文件系统
### 创建分区
```shell
cfdisk /dev/sdb
```
### 创建文件系统
- mkefs 创建 ext 文件系统
- mke2fs 创建 ext2 文件系统
- mkfs.ext3 创建 ext3 文件系统
- mkfs.ext4 创建 ext4 文件系统
- mkreiserfs 创建 ReiserFS 文件系统
- jfs_mkfs 创建 JFS 文件系统
- mkfs.xfs 创建 XFS 文件系统
- mkfs.zfs 创建 ZFS 文件系统
- mkfs.btrfs 创建 btrfs 文件系统

## ext 文件系统内部机制

### inode 

文件系统上索引技术具体化为索引节点（index node），简称 inode ，inode 中存储了 inode 号、文件类型、权限、文件所有者、大小、时间戳等元数据信息，最重要还存储了指向属于该文件 block 指针。

一般 inode 大小为 128 字节或 256 字节

### bmap 块位图

位图只使用 0 和 1 标识对应 block 是空闲还是被占用，0 和 1 在位图中的位置和 block 的位置一一对应，第一位标识第一块，第二位标识第二块，依次下去直到标记完所有的 block。

位图中 1 个字节 8 个位，可以标识 8 个 block，如果一个 block 为 1KB 容量为 1G 的文件系统，block 数量有 1024*1024，使用位图的话 block 只有 1024\*1024/8=131072 字节 = 128k 。

**bmap 的优化针对的是写优化，因为只有写才需要找到空闲的 block 并分配空闲 block。**

### imap 表

imap 是标识 inode 号是否被分配的位图为 inodemap

### 块组

优化文件系统占用的 block 划分成块组，解决 bmap、inode table 和 imap 太大的问题。

在物理层面的划分是将磁盘按柱面划分为多个分区，即多个文件系统。

在逻辑层面上的划分是将文件系统划分成块组。

每个文件系统包含多个块组，每个块组包含多个元数据区和数据区：元数据区存储 bmap、indoe table、imap 等数据；数据区存储文件数据的区域

### 符号链接存储方式

软链接：符号链接的大小是其指向目标路径占用的字符个数

硬链接：每创建一个文件的硬链接，实质上是多一个指向该 inode 记录的 inode 指针，并且硬链接数加 1。

删除文件的实质是删除该文件所在目录 data block 中的对应的 inode 指针，所以也是减少硬链接次数，当硬链接次数为 1 时再删除文件就是真的删除文件。
