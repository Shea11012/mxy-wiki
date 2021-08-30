# pacman 常用命令
## 安装包
- pacman -S 包名：可以同时安装多个以空格分隔
- pacman -Sy 包名：在同步包数据库后再执行安装
- pacman -Sv 包名：显示一些操作信息后执行安装
- pacman -U： 安装本地包，扩展名为 pkg.tar.gz
- pacman -U url：安装一个远程包，如：`pacman -U http://www.example.com/repo/example.pkg.tar.xz`
- pacman -Sw 包名：只下载包，不安装

## 删除包
- pacman -R 包名：只删除包，保留其全部已经安装的依赖关系
- pacman -Rs 包名：删除包，同时删除其所有没有被其他已安装软件包使用的依赖关系
- pacman -Rsc 包名：删除包时，删除所有依赖
- pacman -Rd 包名：删除包时不检查依赖

## 搜索包
- pacman -Ss 关键字：在仓库搜索含关键字的包
- pacman -Qs 关键字：搜索已安装的包
- pacman -Qi 包名：产看包的详细信息
- pacman -QI 包名：列出包文件

## 其他
- pacman -Sc：清理未安装的包文件，包文件位于 `/var/cache/pacman/pkg`
- pacman -Scc：清理所有缓存文件
- pacman -Syy：仅同步软件包数据库
- pacman -Syyu：同步数据库且更新已安装包