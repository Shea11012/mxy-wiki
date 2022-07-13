---
date created: 2021-09-01 22:08
date modified: 2021-10-29 21:20
title: go 基准测试
---
工具：`go install golang.org/x/perf/cmd/benchstat@latest`

基准测试 (time)：`go test -bench=. -run=^$`
基准测试（time、memory）：`go test -bench=. -run=^$ -benchmem`
> `-run=^$` 不执行单元测试

使用工具比较：
```shell
go test -bench=. -run=^$ -benchmem > old.txt
go test -bench=. -run=^$ -benchmem > new.txt

benchstat old.txt new.txt
```