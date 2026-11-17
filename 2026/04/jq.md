---
created: 2026-04-13
tags:
  - linux
  - tool
---
## 核心操作
提取字段： `jq '.name' data.json`
提取嵌套字段： `jq '.user.profile.age' data.json`
遍历数组：`jq '.[]' data.json`
遍历数组并提取字段：`jq '.[].name' data.json`
按索引访问数组：`jq '.[1]' data.json`
数组切片：`jq '.[1:3]' data.json`

## 高级过滤和条件选择
条件过滤：`jq '.[] | select(.age > 25)' data.json`
管道： `jq '.[] | .name' data.json`
多条件过滤：`jq '.[] | select( .age > 25 and .city == "Shanghai")' data.json`
条件判断： `jq 'if age > 25 then "Adult" else "Minor" end' data.json`

## 数据转换和输出控制
数据映射：`jq 'map(.name)' data.json`
添加、修改字段：`jq 'map(. + {"isActive": true})' data.json`
删除字段：`jq 'del(.user.password)' data.json`
字段重命名：`jq '.users[] | {username: .name}' data.json`
构造新对象：`jq '{name,city}' data.json`
原始字符串输出： `jq -r '.name' data.json`

## 常用命令行选项和技巧
- Slurp(-s) 将输入的所有JSON对象，合并为一个大数组
- 传递外部变量：--arg / --argjson
```bash
name='Alice'
jq --arg n "$name" '.name == $n' data.json
```
- 递归下降(..)，可以递归搜索JSON树，查找所有满足条件的元素
```bash
jq '.. | select(type == "string")' data.json
```
- 处理缺失字段(//) 为不存在的字段提供默认值
```bash
jq '.user.profile.sex // "N/A' data.json #若age不存在，则返回 N/A
```