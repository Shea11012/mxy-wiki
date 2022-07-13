---
date created: 2021-11-30 21:22
date modified: 2021-11-30 21:22
title: go时间处理
---
### unix 时间戳转换

```go
const tt = 1618900461998  // 这是一个 millisecond 时间戳

layout := "2006-01-02 15:04:05"
// 1 偏移量
now := time.Now().Local()
_,offset := now.Zone()	// 根据当前时间获取到偏移量
fmt.Println(time.Unix(0,tt * int64(time.Millisecond) + int64(offset)).Format(time.RFC3339))

// 2 指定时间差
fmt.Println(time.Unix(0,tt * int64(time.Millisecond) + int64(8 * 3600)).Format(layout))

// 3 指定时区
fmt.Println(time.Unix(0,tt * int64(time.Millisecond)).Add(8 * time.Hour).UTC().Format(layout))

```

