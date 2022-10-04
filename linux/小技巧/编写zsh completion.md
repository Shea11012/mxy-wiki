---
date created: 2021-11-28 13:34
date modified: 2022-02-11 04:58
title: zsh completion
---
# zsh completion
## 配置 .zshrc 设置
```shell
fpath=($HOME/.zsh-completions $fpath)

autoload -Uz compinit
compinit
````

## 示例
在 `$HOME/.zsh-completions` 目录内创建一个 `_quick-start` 文件
```shell
#compdef _quick-start quick-start

_quick-start() {
	local subcmds
	if (( CURRENT == 2 ));then
		subcmds=("restart_publicfns_redis_cluster:restart public fns redis cluster")
		_describe -t subcmds 'subcmds' subcmds
	fi
}

_quick-start

```
`exec zsh` 重新加载 shell

## 工具方法
### compdef
>[!note]
> compdef function command command command
>> [!example]
>> compdef _hello hello1 hello2 hello3


### 整体补全
| 名称 | 描述|
| --- | --- |
| `_alternative` | 生成补全候选列表 |
| `_arguments` | 指定如何对一个命令进行选项补全和参数补全 |
| `_describe` | 创建仅由单词、描述信息组成的简单补全 |
| `_regex_arguments` | 创建一个方法，使用正则表达式匹配命令行参数 |

### 单个单词补全
| 名称 | 描述 |
| --- | --- |
| `_values` | 任意关键词及其参数进行补全 |
| `_combination` | 用于补全值的组合，比如 hostname 和 username 的组合|
| `_multi_parts` | 对单词的多个部分进行补全，其中每个部分使用某个字符分隔。如对部分文件路径进行补全： `/u/i/sy` 补全为 `/usr/include/sys` |
| `_sep_parts` | 类似 `_multi_parts`,允许不同的补全部分使用不同的分隔符 |
|  `_sequence` | 封装另一个补全方法，从而对另一个补全方法生成的匹配项进行补全 |

### 为指定类型的对象进行补全
| 名称 | 描述 |
| --- | --- |
| `_path_files` | 补全文件路径 |
| `_files` | 使用除了 `-g` 和 `-/` 以外的所有选项调用 |
| `_net_interfaces` | 补全网络接口名称 |
| `_users` | 补全用户名 |
| `_groups` | 补全用户组名称 |
| `_options` | 补全 shell 名称|
| `_parameters` | 补全 shell 参数、变量的名称|

### Action
`_arguments`，`_regex_arguments`，`_alternative` 等最后一个参数都是 Action。表示如何补全相应的参数
| 名称 | 描述 |
| --- | --- |
| `()` | 参数是必须的，但没有匹配项。相当于占位符|
| `ITEM1 ITEM2` | 可能的匹配列表 |
|  `->string` |  将 `$state` 设置为 string 并继续后续逻辑 |
| `FUNCTION` | 将要调用的方法名称，该方法能够生成匹配项或调用其他 Action |
| `{EVAL-STRING}` | 将字符串作为 shell 代码执行，可以生成匹配项。也可以调用工具方法并传入参数 |
| `=ACTION` | 在不改变补全位置节点的情况下，向补全命令中插入一个伪造单词 |





## link
[Writing Zsh Completion for Padrino](https://wikimatze.de/writing-zsh-completion-for-padrino/)
[writing zsh completion scripts](https://tylerthrailkill.com/2019-01-13/writing-zsh-completion-scripts/)
[Zsh 自动补全脚本入门 | 楚权的世界](http://chuquan.me/2020/10/02/zsh-completion-tutorial/)