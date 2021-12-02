### protoc 命令

```
protoc --proto_path=. --go_out=plugins=grpc,paths=source_relative:. ./proto/*.proto

# 上面的命也可写成
protoc --proto_path=. --go_out=plugins=grpc,paths=source_relative:. ./proto/*.proto

```

- `--proto_path` 或者 `-I` 指定编译源码路径
- `--go_out` 指定插件和go代码生成位置

### 遇到的问题

WARNING: Missing 'go_package' option 

需要在 proto 文件中添加 option go_package 配置，如果当前module是 `github.com/Shea11012/learn-grpc` 中，proto 的文件是写在 `proto` 目录下

```protobuf
package a
option go_package = "github.com/Shea11012/learn-grpc/proto;proto"			
// 根据上面的go_package执行生成命令,生成的pb文件会在当前目录下，包名是 a
protoc --go_out=plugins=grpc,paths=source_relative:. ./proto/*.proto

-------------------------

package b
option go_package = "./pb;pb"
// 以上面的go_package执行命令，pb文件会生成在 ./pb 目录下,包名是 pb
protoc --go_out=. ./proto/*.proto

```

