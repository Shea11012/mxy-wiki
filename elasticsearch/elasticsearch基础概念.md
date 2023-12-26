---
tags:
  - elasticsearch
date created: 2021-11-30 21:22
date modified: 2023-10-10 05:36
---

## elasticsearch 是什么

elasticsearch 是一个基于 Apache Lucence 的开源搜索引擎，elasticsearch 使用 Java 开发并使用 Lucence 作为其核心实现所有索引和搜索的功能

- 分布式的实时文件存储，每个字段都被索引并可被搜索
- 分布式的实时分析搜索引擎
- 可以扩展到上百台服务器，处理 PB 级结构化或非结构化数据，支持全文搜索

### elastic 生态圈

kibana：可视化

elasticsearch：存储计算查询

logstash、Beat：数据抓取

### elastic 文件目录结构

| 目录    | 配置文件          | 描述                                         |
| ------- | ----------------- | -------------------------------------------- |
| bin     |                   | 脚本文件，启动 es，安装插件。运行统计数据等。 |
| config  | elasticsearch.yml | 集群配置文件，user，role 相关配置             |
| JDK     |                   | java 运行环境                                 |
| data    | path.data         | 数据文件                                     |
| lib     |                   | java 类库                                     |
| logs    | path.log          | 日志文件                                     |
| modules |                   | 包含所有 ES 模块                               |
| plugins |                   | 包含所有已安装插件                           |

## Document 文档

- es 是面向文档的，文档是所有可搜索数据的最小单位
- 文档会被序列化成 JSON 存在 es 中
- 每个文档都有一个 unique id

#### 文档的元数据

用于标注文档的相关信息

- \index 文档所属的索引名
- \_type 文档所属的类型名
- \_id 文档唯一 ID
- \_source 文档的原始 JSON 数据
- \_version 文档的版本信息
- \_score 相关性打分

index 索引是文档的容器，每个索引都有自己的 mapping 定义，用于定义包含文档的字段名和字段类型

shard 索引中的数据分散在 shard 上

索引的 mapping 和 settings

- mapping 定义文档字段的类型
- settings 定义不同的数据分布

### 集群

es 不同的集群通过不同的名字来区分，默认 elasticsearch ，也可以通过指定 cluster.name=xxx 来设定。

### 节点

节点就是一个 es 实例，每个节点都有一个名字，通过 node.name=xxx 指定

每个节点在启动之后，会分配一个 UID，保存在 data 目录下

#### master-eligible 节点 和 master 节点

>[!note]
> 每个节点启动后，默认就是一个 master eligible 节点。可以通过 node.master: false 禁止
>
> master eligible 节点可以参加选主流程，成为 master 节点
>
> 当第一个节点启动时，它会将自己选举成 master 节点
>
> 每个节点上都保存了集群的状态，且只有 master 节点才能修改集群状态信息
>
> - 所有的节点信息
> - 所有的索引和其相关的 mapping 与 setting 信息
> - 分片的路由信息
>
> 如果任意节点都能修改信息会导致数据的不一致性

#### data node & coordinating node

>[!summary]
> **data node**
> 保存数据的节点，叫做 data node。负责保存分片数据。在数据的扩张上起到关键作用
> **coordinating node**
> 负责接收 client 请求，将请求分发到合适的节点，最终把结果汇集到一起
> 每个节点默认起到 coordinating node 的职责

#### hot & warm node (冷热节点)

>[!summary]
> hot 节点配置更高，io 和 cpu 更高
> warm 节点，存储旧数据
> 通过 hot & warm 可以降低集群部署成本

#### machine learning node

跑机器学习的 job，用来做异常检测

### 配置节点类型

生产环境中，应该为每个节点设置单一的角色

| 节点类型          | 配置参数    | 默认值                                                    |
| ----------------- | ----------- | --------------------------------------------------------- |
| master eligible   | node.master | true                                                      |
| data              | node.data   | true                                                      |
| ingest            | node.ingest | true                                                      |
| coordinating only | -           | 每个节点默认都是 coordinating，需要设置其他类型全部为 false |
| machine learning  | node.ml     | true (需额外 enable x-pack)                               |

### 分片（Primary Shard & Replica Shard）

- 主分片，解决数据水平扩张的问题。通过主分片可以将数据分布到集群内的所有节点之上
  - 主分片在索引创建时指定，后续不允许修改，除非 reindex
  - 主分片不能和自己的副本分片在同一个节点上（节点会宕机）
- 副本，负责容错，承担读请求负载，解决数据高可用的问题。主分片的拷贝
  - 副本分片数，可以动态调整
  - 增加副本数，可以在一定程度上提高服务的可用性（读取的吞吐）

### 倒排索引

- 单词词典（term dictionary）记录所有文档的单词，记录单词到倒排列表的关联关系
- 倒排列表（posting list）记录单词对应的文档结合，由倒排索引项组成
  - 倒排索引项（posting）
    - 文档 ID
    - 词频 TF，该单词在文档出现的次数，用于相关性评分
    - 位置 Position，单词在文档中分词的位置。用于语句搜索
    - 偏移 offset，记录单词的开始结束位置，实现高亮显示

### Mapping & Dynamic Mapping

mapping 类似数据库的 schema，作用：

- 定义索引中的字段名称
- 定义字段的数据类型，例如字符串，数字，布尔
- 字段，倒排索引的相关配置

mapping 会把 json 文档映射成 Lucene 所需要的扁平格式

一个 mapping 属于一个索引的 Type

- 每个文档都属于一个 Type
- 一个 Type 有一个 mapping 定义
- 7.0 开始 不需要在 mapping 中指定 Type 信息

字段的数据类型：

- 简单类型
  - Text、keyword
  - Date
  - Integer、float
  - Boolean
  - IPv4 & IPv6
- 复杂类型 对象和嵌套对象
  - 对象类型、嵌套类型
- 特殊类型
  - geo_point & geo_shape 、percolator

#### dynamic mapping 类型自动识别

| JSON 类型 | es 类型                                                       |
| -------- | ------------------------------------------------------------ |
| 字符串   | 匹配日期格式，设置成 date<br />配置数字设置为 float 或 long，该选项默认关闭<br />设置为 text，并且增加 keyword 字段 |
| 布尔值   | boolean                                                      |
| 浮点数   | float                                                        |
| 整数     | long                                                         |
| 对象     | object                                                       |
| 数组     | 由第一个非空数值类型所决定                                   |
| 空值     | 忽略                                                         |

#### 更改 mapping 类型

1. 新增字段
   - dynamic 设为 true，一旦有新增字段写入，mapping 同时也被更新
   - dynamic 设为 false，mapping 不会被更新，新增字段的数据无法被索引，但是信息会出现在 \_source 中
   - dynamic 设为 strict，文档写入失败
2. 对已有字段，一旦已有数据写入，就不再支持修改字段定义。如果希望改变字段类型，必须 reindex api，重建索引

#### 动态字段映射

```http
PUT test-dynamic-mapping
{
	"mappings": {
		// 整体不自动增加字段
		"dynamic": false
		"properties": {
			"person": {
				// 内嵌person对象可以自动增加字段
				"dynamic": true,
				"properties": {
					"name": {
						"type": "keyword"
					}
				}
			},
			"company": {
				// company对象发现新字段会报错
				"dynamic": "strict",
				"properties": {
					"company_id": {
						"type": "keyword"
					}
				}
			}
		}
	}
}
```

#### 动态模板

##### match_mapping_type

按照数据类型匹配，当相对json中某种数据类型字段做特殊配置时使用
```http
PUT test-dynamic-mapping
{
	"mappings": {
		"dynamic_templates": [
			{
				"test_float": {
					// 匹配json中的 long
					"match_mapping_type": "long",
					"mapping": {
						// 将字段设置为 integer
						"type": "integer"
					}
				}
			}
		]
	}
}
```

##### match、unmatch

match、unmatch使用 `*` 匹配0个或多个字符
设置 `"match_pattern": "regex"`，match、unmatch可以使用正则表达式

```http
PUT test-dynamic-mapping
{
	"mappings": {
		"dynamic_templates": [
			{
				"match_test": {
					"match_pattern": "regex",
					"match": "long_*", // 属性名以long_开头
					"unmatch": ".*_test$*", // 属性名不以_test结尾
					"mapping": {
						"type": "long" //匹配到的字段设为long
					}
				}
			}
		]
	}
}
```

##### path_match、path_unmatch

用于匹配内嵌对象时使用
```http
PUT test-dynamic-mapping
{
	"mappings": {
		"dynamic_templates": [
			{
				"test_path_match": {
					"match_pattern": "long_", // 以long_开头
					"path_match": "person.*", // 内嵌对象person所有字段
					"path_unmatch": "*.age", // 排除age字段
					"mapping": {
						"type": "text" // 匹配字段设置为text
					}
				}
			}
		]
	}
}
```

##### 模板变量

`{name}`和 `{dynamic_type}`会根据动态匹配的内容进行替换
```http
PUT test-dynamic-mapping
{
	"mappings": {
		"dynamic_templates": [
			{
				"named_analyzers": {
					// 匹配所有string类型
					"match_mapping_type": "string",
					// 匹配任意属性名
					"match": "*",
					"mapping": {
						"type": "text,
						// 将匹配到的字段名设为解析器
						"analyzer": "{name}"
					}
				}
			},
			{
				"no_doc_values": {
				// 匹配所有类型，但是上面有了一个string匹配，所以实际上是匹配除了string之外的所有其他字段
					"match_mapping_type": "*",
					"mapping": {
						// type设为匹配到的类型
						"type": "{dynamic_type}",
						"doc_values": false
					}
				}
			}
		]
	}
}

PUT test-dynamic-mapping/_doc
{
	"english": "Some English text", //会用english作为解析器
	"count": 5 // 会设置为动态类型 long
}
```

>[!tip]
>1. 所有null值及空数组属于无效值，不会被任何动态模板匹配，在索引文档时，只有字段第一次有有效值时，才会与各动态模板匹配，找出匹配的模板创建新字段
>2. 规则匹配时，按照动态模板配置的顺序依次对比，使用最先匹配成功的模板。当一个字段符合2个以上的模板时，只会使用最靠前的那一个。每个动态模板至少应包含`match`、`path_match`、`match_mapping_type`中的一个，`unmatch`和`path_unmatch`不能单独使用
>3. `dynamic_templates`可以在运行时修改，每次修改都是整体替换，而非追加
