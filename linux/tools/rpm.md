---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: rpm 包命令
---
# rpm 包命令

1. rpm 包安装

`rpm {-i | --install} [install-options] package-file` 

- options:
  - `-i，--install` ：安装
  - `-v`：verbos，输出详细信息
  - `-vv`：verbos，输出更详细信息
- `install-options`
  - `-h`：以 hash marks 格式输出进度条，每个 # 代表 2% 的进度
  - `--test`：测试安装，只做环境检查，并不真正安装
  - `--nodeps`：忽略程序依赖关系
  - `--replacepkgs`：覆盖安装，如果文件修改错误，需要将其找回，使用此方法，需要把修改错误的文件提前删除
  - `--justdb`：不执行安装操作，只更新数据库
  - `--noscripts`：不执行 rpm 自带的所有脚本
  - `--nosignature`：不检查包签名信息，即不检查来源合法性
  - `--nodigest`：不检查包完整性信息

2. RPM 包升级

- `rpm {-U | --upgrade} [install-options] package_file`
- `rpm {-F | --freshen} [install-options] package_file`
- `options`
  - `-U`：升级或安装
  - `-F`：升级但不安装
- `install-options`
  - `--oldpackage`：降级
  - `--force`：强制升级

3. RPM 卸载

- `rpm {-e | --erase} [--allmatches] [--nodeps] [--noscripts]`
- `options`
  - `-e`：删除指定程序
  - `--allmatches`：卸载所有匹配指定名称的程序包的各个版本
  - `--nodeps`：卸载时忽略依赖关系
  - `--test`：测试卸载，dry run 模式

4. RPM 包查询

`rpm {-q | --query} [select-options] [query-options]`

- `options`

  - `-qa,-all`：查询所有已安装的包
  - `-f，--file File`：查询指定的文件是由哪个包安装生成的
  - `-g，--group GROUP`：查询指定包由哪个包组提供
  - `-p，--package PACKAGE_FILE`：对未安装的程序包执行查询操作
  - `--whatprovides CAPABILITY`：查询指定的 capability 由哪个程序包提供
  - `--whatrequires CAPABILITY`：查询指定的 capability 被哪个包所依赖

- `query-options`

  - `--changelog`：查询 rpm 包的 changelog
  - `-l，--list`：列出程序包安装生成的所有文件列表
  - `-i，--info`：查询程序包的 infomation，包括版本号、大小和所属的包组等信息
  - `-c，--configfiles`：查询指定的程序包提供的配置文件
  - `-d，--docfiles`：列出指定的程序包提供的文档
  - `--provides`：列出指定的程序包提供的所有 capability
  - `-R，--requires`：查询指定程序包的依赖关系
  - `--script`：查查看程序包自带的脚本

- 校验程序包后的 format

  - 例：`rpm -V zsh`
  - `S.5....T.  d /usr/share/doc/zsh-4.3.11/README`

  说明：

  > S 文件大小改变
  >
  > M 文件权限改变
  >
  > 5 MD5 校验码改变
  >
  > D 主次设备不匹配
  >
  > L read link 变化
  >
  > U 属主改变
  >
  > G 属组改变
  >
  > T 修改时间改变
  >
  > P CAPABILITY 改变

  改变才会显示相应值，`.` 代表未修改过

5. RPM 包校验

- `options`
  - `-V package_file`：自动检查其数据的完整性及合法性
  - `-K package_file`：收到检查其数据的完整性及合法性
  - `--import key`：手动导入 GPG KEY