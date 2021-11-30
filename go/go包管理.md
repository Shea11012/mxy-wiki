go module 生成的目录下会出现两个文件

- `go.mod` 声明模块名称、go的版本、该模块的依赖信息
- `go.sum` 该模块的所有依赖项的校验和。go 在编译前会下载并检查校验和，如果不一致将会拒绝编译。

go module 相关的环境变量

```
GOPROXY="https://proxy.golang.org,direct"
GOSUMDB="gosum.io+ce6e7565+AY5qEHUk/qmHc5btzW45JVoENfazw8LielDsaI+lEbq6"
GOPRIVAT="XXX"
GONOPROXY="XXX"
GONOSUNDB="XXX"
```

- `GOPROXY` 指示镜像递增和获取顺序，每个镜像URL之间用 `,` 分隔，按照优先顺序排列。direct 表示源站。
  - `GOPROXY=off` 表示禁止从互联网下载任何依赖
  - `GOPROXY=direct` 表示依赖一律从源站下载
- `GOSUMDB` 指示校验和服务器的地址和公钥，对于地址 `sum.golang.org`，公钥可以省略（内置）
- `GOPRIVATE` 表示私域前缀。私域前缀下的所有依赖一律从源站下载，而且不做校验。支持多个私域前缀，每个私域用 `,` 分隔。支持通配符 `*`

#### go 包代理协议

go 使用协议约定的格式获取模块信息，并缓存到 **GOPATH/pkg/mod/cache/download** 目录下

1. 获取模块列表

   > 代理服务器需要响应 GET $GOPROXY/<module>/@v/list 请求，并返回当前代理服务中的该模块版本列表
   >
   > 如：curl https://proxy.golang.org/github.com/google/uuid/@v/list
   >
   > v1.0.0
   >
   > v1.1.1
   >
   > v1.1.0

2. 获取模块元素数据

   > 响应 GET $GOPROXY/<module>/@v/<version>.info 请求，返回当前代理服务器特定版本信息
   >
   > {“Version”:”v1.1.1”,”Time”:”2019-02-27T21:05:49Z”}

3. 获取go.mod文件

   > 响应 GET $GOPROXY/<module>/@v/<version>.md 请求，返回当前代理服务器特定版本模块的go.mod 文件
   >
   > 如：curl https://proxy.golang.org/github.com/google/uuid/@v/v1.1.1.mod
   >
   > module github.com/google/uuid # 该模块如果无依赖，则只会包含一个module名称

4. 获取代码压缩包

   > 响应 GET $GOPROXY<module>/@v/<version>.zip 请求，返回当前代理服务器中特定版本模块的压缩包
   >
   > 如：curl https://proxy.golang.org/github.com/google/uuid/@v/v1.1.1.zip # 该压缩包只包含该版本的文件，不包含版本控制信息

5. 获取模块的最新可用版本

   > 请求 GET $GOPROXY/<module>/@latest 获取指定模块的最新可用版本
   >
   > go 命令获取模块最新的可用版本，通常是通过 $GOPROXY/<module>/@v/list 获取版本列表，计算出最新的版本。
   >
   > 只有当list不可用时，才会使用 latest
   >
   > **所以这是代理服务器唯一的可选协议**

可以使用 `go mod download -x -json {package}` 查看下载步骤



#### GOSUMDB 工作机制

GOSUMDB 用于 go 命令校验模块时应该信赖哪个数据库。

完整的GOSUMDB配置

> GOSUMDB=“<checksum database name>+<public key><checksum database service url>”
>
> 其中数据库名字和公钥必须指定，校验和数据库服务url则是可以选的，默认是 https://<checksum database name>

**GOSUMDB 不能确保模块是否包含恶意代码，只能保证构建的一致性**

> go 命令会通过 https://<checksum database service url>/lookup/<module>@<version> 来查询特定模块版本的hash值
>
> 校验和数据库通过存储hash值，给特定的模块版本提供了公证服务

##### 校验流程

1. 模块被下载后，go命令会对下载的模块做hash运算，然后与校验和数据库中的数据进行对比，依此确保下载的模块是否合法
2. 模块hash被写入到go.sum 文件之前，对缓存在本地的模块做hash运算，与校验和数据库中的进行对比，确保本地的模块没有被篡改

**如果校验某模块时，校验和数据库没有收录该模块，则会先拉取该模块版本，计算hash，存入数据库中。因此官方的校验和数据库只能收录公开的模块。**

##### 校验和数据库代理

> <proxyURL>/sumdb/<sumdb-name>/supported
>
> go 在访问 GOSUMDB 时，会通过该接口询问模块代理是否能代理校验和数据库



#### 私有模块

私有模块前缀支持通配符，多个模块之间使用逗号分隔

`GOPRIVATE=*.corp.example.com,rsc.io/private`

私有模块不会从代理服务器下载代码，也不会使用校验服务器来检查下载的代码

#### 配置 go get 获取私有仓库

- 使用访问令牌

  如：gitlab 去主页获取访问 token 来实现访问仓库

  ```
  对项目生效
  git config http.extraheader "PRIVATE-TOKEN:{token}"
  
  对所有项目生效
  git config --global http.extraheader "PRIVATE-TOKEN:{token}"
  ```

- 使用 git 方式拉取代码

  ```
  全局替换，该域名下都走 ssh 协议
  git config --global url."git@{url}:".insteadof "https://{url}/"
  ```

  

go 包管理工具：

- `go mod init xxx` 创建一个新的模块，初始化 `go.mod` 并生成相应的描述
- `go build,go test` 会在需要的时候在 `go.mod` 文件中添加新的依赖项
- `go list -m all` 列出当前模块的所有依赖项
- `go get` 修改制定依赖项的版本或者添加一个新的依赖项
- `go mod tidy` 移除模块中没有用到的依赖项

### 将 exported 类型变为其它 package 不可访问

通过定义一个 `internal` 包，将其它包中的 `exported` 函数、类型放在 `internal` 包下这样，就只会被其父级访问：

`/a/b/c/internal/d/e/f` 可以被 `/a/b/c` 导入，不能被 `/a/b/g` 导入

### 访问其它包中的私有方法

对外声明一个导出函数 `external function`，例如：`time.go 中的 time.Sleep` 这个函数就是在 `time` 包内进行函数声明，实现是在 runtime 中的一个 unexported 函数中（`runtime.timeSleep`) 实现的 。

将相关的函数定义在 runtime 包中的好处是，可以访问 runtime 包中 unexported 的类型，同时相关的函数还可以被其它包访问。

使用 go 提供的一个指令 `//go:linkname localname importpath.name`，这个指令为函数或者变量 `localname` 使用 `importpath` 作为目标文件的符号名。**因为这个指令破坏了类型系统和包的模块化，所以它只能在 `import unsafe` 的情况下才能使用**

`importpath.name`格式：`/a/b/c/d/pkg.foo` ，此时在 `/a/b/c/d/pkg` 包中就可以使用这个函数 `foo`



### 构建包

在包声明和包注释前面加注释：

```go
// +build ignore 有这个注释表示不编译该文件
// +build linux darwin 表示只在目标操作系统是 Linux 或 Mac OS 时才编译该文件
```

### 包文档

go 中，包中每个导出的成员和包声明前都应该包含目的和用法说明的注释。

### 内部包

为了将一个包暂时只暴露给一些受限制的用户使用。

go 对于包含了 internal 名字的路径段的包导入路径做了特殊处理，这种包叫 internal 包，一个 internal 包只能被和 internal 目录有同一个父目录的包所导入。例：

- net/http/internal/chunked 内部包只能被 net/http/httputil 或 net/http 包导入，但是不能被 net/url 包导入。

### 构建约束

- 以前缀 `+build` 开始，后跟一个或多个空格
- 位于文件顶部在保声明之前
- 它和包声明至少有一个空行，以防止被视为包文档

**构建约束指定以下规则：*

- 以 ！开头的构建标记是无效的
- 用空格分隔的构建标记在逻辑上是或
- 用逗号分隔的构建标记在逻辑上是和
- 在多行上构建的约束是逻辑和
- `ignore` 在任何不匹配有效架构或操作系统的上会构建



