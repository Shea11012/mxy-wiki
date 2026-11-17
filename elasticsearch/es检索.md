---
updated: 2025-08-13
tags:
  - elasticsearch
---

![[../assets/2e7d1be9b53e4bef8f538611dd46bafa_MD5.jpg]]

# 全文检索

## match 和 match_phrase

match 适合应用于高召回率和结果精准度要求较低的场景
match_phrase 要求查询的词条顺序和文档中的词条顺序一致

## multi_match

可以在多个字段中搜索指定关键词
```
POST my_index/_search
{
  "query": {
    "multi_match": {
      "query": "哈基米",
      "fields": [ "title^3", "message"]
    }
  }
}
```
在 title 和 message 字段中搜索，使用 `^3` 提升 title 权重

## match_phrase_prefix

查询词语需要按顺序匹配文档中的内容，同时允许最后一个词语只匹配前缀

## query_string

允许使用 lucene 查询语法编写复杂查询表达式。应用场景包括高级搜索、数据分析和报表

# 组合检索

## bool 组合检索

包含了四种子查询类型
- must：查询结果必须满足指定条件
- must_not：查询结果必须不满足指定条件
- filter：过滤条件，使用 filter 可以借助缓存机制提高查询性能
- should：查询结果可以满足部分条件，满足条件的最小数量有 minimum_should_match 参数控制

>[!note] query 和 filter 区别
>query 用于评估文档相关性，并对结果评分，通常用于搜索场景
>filter 用于筛选文档，不会对文档评分，通常用于过滤场景
>在性能方面，filter 比 query 更有效，它不需要进行评分，且可以被缓存

当同时需要评分和筛选时，可以组合使用 query 和 filter
```
POST my_index/_search
{
  "query": {
     "bool": {
        "must": [
           {
             "match": {
                "title": "search"
             }
           },
           {
             "match": {
                "content": "elasticsearch"
             }
           }
        ],
        "filter": [
           {
              "term": {
                 "status": "published"
              }
           },
           {
              "range": {
                 "publish_date: {
                    "gte": "2024-01-01"   
                 }
              }
           }
        ]
     }
  }
}
```

