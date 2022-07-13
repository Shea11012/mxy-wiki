---
date created: 2021-11-30 21:22
date modified: 2022-03-09 19:00
title: elasticsearch API总结
---
`shutdown` API

`curl -X``POST 'http://localhost:9200/_shutdown'`

创建文档

```
	  索引名    类型名   ID
PUT /megacorp/employee/1
{
"first_name" : "John",
"last_name" : "Smith",
"age" : 25,
"about" : "I love to go rock climbing",
"interests": [ "sports", "music" ]
}
```

检索文档

`GET /megacorp/employee/1`  DELETE 删除 PUT 更新

返回信息如下：

```
{
    "_index": "megacorp",
    "_type": "employee",
    "_id": "1",
    "_version": 1,
    "_seq_no": 0,
    "_primary_term": 1,
    "found": true,
    "_source": {
        "first_name": "John",
        "last_name": "Smith",
        "age": 25,
        "about": "I love to go rock climbing",
        "interests": [
            "sports",
            "music"
        ]
    }
}
```

简单搜索

```
GET /megacorp/employee/_search  // 搜索全部
GET /megacorp/employee/_search?q=last_name:Smith  // 根据q的键值对查询
```

DSL 查询

DSL（Domain Specific Language 特定领域语言）以 JSON 请求体的形式出现

```
GET /megacorp/employee/_search
{
"query" : {
	"match" : {
		"last_name" : "Smith"
		}
	}
}
```

filter 过滤器

```
GET /megacorp/employee/_search
{
	"query":{
		"filtered":{
			"filter":{
				"range":{
					"age":{"gt":30}
				}
			}
		},
		"query":{
			"match":{
				"last_name":"test"
			}
		}
	}
}
```

全文搜索

```
GET /megacorp/employee/_search
{
    "query" : {
        "match" : {
            "about" : "rock climbing"
         }
     }
}
默认情况下，Elasticsearch 根据结果相关性评分来对结果集进行排序(_score)
```

短语搜索

```
GET /megacorp/employee/_search
{
	"query":{
		"match_phrase":{
			"about":"rock climbing"
		}
	},
	"highlight":{ // 高亮匹配到的关键字，使用 em 包裹住
		"fields":{
			"about":{}
		}
	}
}
```

检索文档的一部分

`GET /index/type/id?_source=xxx,xxx`

如果 `_source` 不带参数则只会得到 `_source` 字段而不会获取其他元数据

检查文档是否存在

`HEAD /index/type/id`

HEAD 方式只会返回 HTTP 头，不会返回响应体

文档局部更新

`POST /index/type/id/_update`

update 请求表单接受一个局部文档参数 `doc` 它会合并到现有文档中，存在的字段被覆盖，新字段添加

```
POST /website/blg/1/_update
{
    "doc":{
        "tag":["testing"],
        "views":0
    }
}
retry_on_conflict=5 参数表示失败重试次数
```

检索多个文档

```
POST /_mget
{
    "docs":[
        {
            "_index":"website",
            "_type":"blog",
            "_id":2,
        },
        {
            "_index":"website",
            "_type":"pageviews",
            "_id":1,
            "_source":"view" // 指定检索字段
        }
    ]
}
如果文档具有相同的 _index 和 _type 可以通过 ids 数组来代替完整的 docs 数组
POST /index/type/_mget
{
    "ids":["2","1"]
}
```

更新时的批量操作



