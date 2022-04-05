## redo日志
为了避免修改一个字节就将脏也刷到磁盘，将一个修改操作按照一定的格式记录下来，当服务器发生崩溃后能使用redo日志复原结果，这就是redo日志。
## 日志格式
redo日志通用结构,根据类型的不同，日志的格式也有所不同
![[Pasted image 20220329155422.png]]
- type：日志的类型
- space ID：表空间ID
- page number：页号
- data：具体内容

type为 `MLOG_8BYTE` 类型时，表示在页面的某个偏移量处写入8字节的数据
![[Pasted image 20220329155954.png]]

type为 `MLOG_WRITE_STRING` 时，因为是写入字节序列不能确定大小所以需要个 len 字段
![[Pasted image 20220329160416.png]]

## 日志的写入过程
将对底层页面进行一次原子访问的过程称为，Mini-Transaction（MTR）。
一个事务可以包含若干个语句，每一条语句又包含若干个MTR，每一个MTR又包含若干条redo日志。
![[Pasted image 20220329162318.png]]

redo日志并不是直接写入磁盘的，会使用一个 redo log buffer 的连续内存空间，这片内存空间被划分为若干个连续的 redo log block。默认值为 16MB
日志的写入并不是以每生成一个redo日志就写入到 redo log 日志中的，而是一个MTR为一组，当MTR结束时才将这一组redo日志写入

### 刷盘时机
- log buffer 空间不足时，当redo日志占了log buffer总容量的50%时，执行一次刷盘
- 事务提交时
- 后台线程以每秒一次的频率将redo日志进行刷盘
- checkpoint时

## checkpoint
redo日志是为了系统崩溃后恢复脏页使用，如果对应的脏页已经刷到磁盘中，即使系统崩溃，在重启后也不会使用该redo日志恢复该页面了，所以该redo日志就没有存在的必要的，它所占用的磁盘空间就可以被后续的redo日志覆盖