---
date created: 2021-11-30 21:22
date modified: 2022-03-02 21:50
title: protobuf 使用
---
```protobuf
message SearchRequest {
	required string query = 1; // 字段规则 字段类型 = 标识号
	optional int32 page_num = 2;
	optional int32 result_per_page = 3;
}
```

#### 标量数值类型



| proto 类型 | go 类型 | 备注                                                        |
| ---------- | ------- | ----------------------------------------------------------- |
| double     | float64 |                                                             |
| float      | float32 |                                                             |
| int32      | int32   | 可变长编码，编码为负数时比较低效，如果包含负数，使用 sint32 |
| int64      | int64   | 可变长编码，包含负数，使用 sint64                           |
| uint32     | uint32  | 可变长编码                                                  |
| uint64     | uint64  | 可变长编码                                                  |
| sint32     | int32   | 可变长编码，有符号整型值，编码时比 int32 高效               |
| sint64     | int64   | 同上                                                        |
| fixed32    | uint32  | 总是 4 个字节，如果数值总是比 228 大的话，会比 int32 高效   |
| fixed64    | uint64  | 总是 8 个字节，如果数值比 256 大的话，会比 int64 高效       |
| sfixed32   | int32   | 总是 4 个字节                                               |
| sfixed64   | int64   | 总是 8 bytes                                                |
| bool       | bool    |                                                             |
| string     | string  | 必须是 utf-8 或者 7 bit ASCII 文本并且不能大于 $2^{32}$     |
| bytes      | []byte  | 包含任何字节序列不能大于 $2^{32}$                           |



#### 标识号

在消息定义中， 每个字段都有一个唯一的 **数字标识符**。这些标识符是用来在消息二进制格式中识别各个字段的，一旦开始使用就不能再改变。

1-15 之内的标识符在编码时只会占用一个字节，16-2047 之内的标识符占用 2 个字节。所以频繁使用的消息元素应该在 1-15 之内。

最小的标识符从 1 开始，最大在 $(2^{29})-1$ ，其中 19000 - 19999 的标识符，protobuf 协议进行了保留，使用预留标识符编译时会报警。

#### 指定字段规则

- required：表示该字段值必须设置
- optional：表示该字段可以有 0 个或者 1 个值
- repeated：该字段可以重复任意多次，重复的值顺序会被保留

#### enum 枚举类型

```protobuf
enum EnumAllowingAlias {
	option allow_alias = true
	UNKNOW = 0;
	STARTED = 1;
	RUNNING = 1;
}
```

`allow_alias` ，允许字段编号重复，`RUNNING 是 STARTED 别名 `

枚举常量必须是一个 32 byte 的整数，不推荐采用负数。**第一个枚举值必须是 0，且必须定义**

#### reserved 保留字段

```protobuf
		enum Foo {
			reserved 2,15,9 to 11,40 to max;
			reserved "FOO","BAR";
		}
```

#### 引入其它 proto 文件

```protobuf
import 'other.proto'
```

#### Nested Types 嵌套类型

```protobuf
message SearchResponse {
	message Result {
		string url = 1;
		string title = 2;
		repeated string snippets = 3;
	}
	repeated Result results = 1;
}
```

#### 更新一个 message 类型

- 不要改变任何已存在的字段标识符
- 字段可以被移除，可以添加前缀 `OBSOLETE_` 或者使用 `reserved`

#### Any 类型

any 以 bytes 呈现序列化的消息，并且包含一个 URL 作为这个类型的唯一标识和元数据

为了使用 any 类型，需要引入 `google/protobuf/any.proto`

```protobuf
import "google/protobuf/any.proto"

message ErrorStatus {
	string message = 1;
	repeated google.protobuf.Any details = 2;
}
```

#### Oneof 

oneof 可以在同时定义一组字段且共享同一块内存，一旦其中一个字段值被设置会自动删除其余字段

**oneof 内不能使用 repeated**

```protobuf
message SampleMessage {
	oneof test_oneof {
		string name = 4;
		SubMessage sub_message = 9;
	}
}
```

#### Maps 

**map** 字段内不能使用 repeated

```protobuf
map<key_type,value_type> map_field = N;
```

#### 定义服务

在使用 RPC 服务时，可以在 proto 文件内定义一个接口，将会自动生成该接口和代码

```protobuf
syntax = "proto3"

service SearchService {
	rpc Search(SearchRequest) returns (SearchResponse);
}
```

#### JSON mapping 

