---
date created: 2021-11-30 21:22
date modified: 2022-01-31 15:45
title: Fiddler内置命令
---
# Fiddler 内置命令

| 符号 | 说明|
|----| ----|
|？|后边跟一个字符串，Fiddler 会将所有会话存在该字符串匹配的全部高亮显示（匹配的是 Protocol，Host，URL）|
|\>，< | 后边跟一个数值，表示高亮所有尺寸大于或小于改数值的会话（匹配的 Body） |
|=|  后边可以接 HTTP 状态码或 HTTP 方法|
|@| 后边跟的是 Host|

## 断点调试

以下是断点调试命令，会话被中断下来之后，点击页面上方的 GO 按钮放行当前中断下来的会话，但新的匹配内容还是会被断下来，输入命令但不带参数表示取消之前设置的断点。

| 命令 | 说明 |
| --- | ---|
|bpafter |后边跟一个字符串，表示中断所有包含该字符串的会话（这是在收到响应后中断）|
|bps |  后边跟的是 HTTP 状态码，表示中断所有为该状态码的会话|
|bpv 或 bpm| 后边跟的是 HTTP 的方法，表示中断所有为该方法的会话|
|bpu 与 bpafter 类似|（这是在发起请求时中断）|
|cls 或 clear|   清除当前所有会话|
|dump|  将所有会话打包成.zip 压缩包的形式保存到 c 盘根目录下|
|g 或 go |  放行所有中断下来的会话|
|hide |   将 Fiddler 隐藏|
|show |  将 Fiddler 恢复|
|urlreplace |  后边跟两个字符串，表示替换 URL 中的字符串，比如 urlreplace baidu fishc 表示将所有 URL 的 baidu 替换成 fishc（直接输入 urlreplace 不带任何参数表示恢复原来的样子）|
|start |  Fiddler 开始工作|
|stop |   Fiddler 停止工作|
|quit |   关闭 Fiddler|
|select |  后边跟响应的类型（Content-Type），表示选中所有匹配的会话，比如 select image 表示选中所有图片|
|allbut 或 keeponly | 和 select 类似不过 allbut 和 keeponly 会将所有无关的会话删除，比如 keeponly image 表示只保留图片删除其余无关的会话|
|!dns |  后边各一个域名，执行 DNS 查找并在右边的 LOG 栏打印结果|
|!listen |  设置其他监听端口|