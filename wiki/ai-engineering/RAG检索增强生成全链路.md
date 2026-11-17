---
title: "RAG检索增强生成全链路"
created: 2026-07-15
updated: 2026-07-16
status: stable
type: concept
tags: [ai/rag, ai/embedding, ai/retrieval, ai/vector-db, ai/eval]
source: "workspace/agents/AI-Engineer-from-scrach/02-RAG/"
confidence: SYNTHESIZED
related:
  - "AI工程师从零开始（总览）"
  - "Agent范式详解"
  - "Agent记忆系统"
---

# RAG检索增强生成全链路

> RAG 的本质不是“向量库 + Prompt”，而是在生成时从受控知识源取得证据，并让答案可以被验证、更新和治理。

## 一、端到端架构

```text
离线/准实时数据面
Source → Parse/OCR → Normalize → Segment → Enrich → Index → Version

在线查询面
Query → Intent/Auth → Rewrite/Decompose → Retrieve → Fuse → Rerank
      → Select/Compress Context → Generate → Cite/Abstain → Verify/Log
```

生产系统还需要删除传播、权限、增量更新、质量评测和观测，不能只关注检索调用。

## 二、数据摄取与解析

| 阶段 | 关键问题 |
|------|----------|
| 连接器 | 数据源版本、增量同步、删除和权限变更能否传播 |
| 解析 | 文本、表格、图片、公式、代码和阅读顺序是否保留 |
| 规范化 | 编码、页眉页脚、重复内容、语言、元数据是否清理 |
| 分段 | chunk 是否保持语义单元及父子结构 |
| 增强 | 标题、实体、日期、ACL、来源 URL、页码、版本 |
| 索引 | lexical、dense、structured、graph 是否按需求组合 |

### PDF 不是普通文本

PDF 保存的是页面绘制指令，不保证自然阅读顺序。解析时应区分：

- 原生文本 PDF 与扫描件 OCR
- 单栏、多栏、表格、脚注和跨页结构
- 图片/图表的视觉理解与文本说明
- 页码、坐标和章节层级等引用锚点

解析质量应单独评测；解析错了，后续 embedding 和 reranker 无法补救。

## 三、Chunking

常见策略：

- 固定 token + overlap：简单基线
- 按标题/段落/代码块/表格的结构切分
- 语义切分：代价较高，边界不一定稳定
- Parent-child：小块召回，大块提供上下文
- Late chunking / contextualized chunk：借助长上下文表征减少局部歧义

不要把 chunk 大小当固定真理。它取决于语料、查询粒度、embedding 模型、reranker 和答案所需证据跨度，应通过 retrieval eval 调参。

## 四、检索不是只有向量相似度

### 4.1 多路召回

| 通路 | 擅长 |
|------|------|
| BM25/全文 | 专名、ID、代码、原句、稀有词 |
| Dense embedding | 同义改写、概念相似 |
| Sparse neural | 兼顾词项与语义扩展 |
| Structured/SQL | 精确字段、范围、聚合、权限 |
| Graph | 显式关系和多跳实体路径 |

Hybrid Search 常用 RRF 或学习融合，不应直接相加未经校准的异构分数。

### 4.2 Query Transformation

- rewrite：消歧、补上下文、标准化
- multi-query：生成多个检索视角
- decomposition：多跳问题拆成子问题
- HyDE：生成假设文档辅助 embedding
- metadata filter：先按租户、时间、产品、语言等缩小范围

查询改写也可能改变用户意图，应保留原始 query 并做回归评测。

### 4.3 Reranking 与 Context Selection

典型漏斗：

```text
廉价召回 Top 100
 → cross-encoder/LLM rerank Top 20
 → 去重与证据覆盖
 → token 预算选择 Top N
```

排序靠前不等于应全部塞进上下文。需要兼顾多样性、相互支持/冲突、位置偏置和证据覆盖。

## 五、生成、引用与拒答

生成器应：

1. 明确区分指令与检索数据，防止知识库 Prompt injection
2. 只依据可用证据回答需要 grounding 的事实
3. 将引用绑定到具体 claim 和稳定来源锚点
4. 证据不足、过期或冲突时澄清或拒答
5. 对数字、日期、实体和跨文档推断进行额外校验

“要求引用”不能保证忠实性。模型可能生成不存在或不支持 claim 的引用，因此必须验证：引用存在、内容支持断言、版本和权限正确。

## 六、RAG 变体的正确定位

- **Graph RAG**：适合关系密集、全局主题或多跳实体问题；建图和更新成本高，并非普通 RAG 的必然升级
- **Multi-hop RAG**：根据中间证据迭代查询，需要管理错误传播和停止条件
- **Self-RAG/Corrective RAG**：研究与实现家族，通过模型判断检索必要性/质量；不能把模型自评当独立真值
- **Agentic RAG**：由 Agent Loop 动态选择知识源、拆解查询、重试和调用非检索工具；增加灵活性，也增加延迟、成本和不可预测性
- **Long-context generation**：不是 RAG 的替代品；长上下文仍有权限、更新、成本、位置偏置和证据发现问题，可与 RAG 组合

## 七、评测：先拆阶段

### 7.1 检索评测

需要 query、corpus 和 relevance judgement：

- Recall@K：相关证据是否进入候选集
- Precision@K：Top-K 噪声比例
- MRR：首个相关结果的位置
- nDCG@K：有分级相关性时的排序质量
- Hit rate、覆盖率、延迟与成本

### 7.2 Context 与生成评测

| 维度 | 问题 |
|------|------|
| Context precision | 注入的上下文有多少真正有用 |
| Context recall | 回答所需证据是否都在上下文中 |
| Faithfulness/groundedness | 答案断言是否被上下文支持 |
| Answer correctness | 最终答案是否事实正确 |
| Relevance | 是否回答了用户问题 |
| Citation correctness | 引用是否存在且支持对应 claim |
| Abstention quality | 证据不足时是否正确拒答 |

原文把 Context Recall/Precision 归为“端到端”并不严谨：它们主要评价检索上下文；端到端还应包含答案正确性、任务成功、延迟、成本和安全。

LLM-as-a-judge 可降低标注成本，但应抽样人工校准、固定 judge 版本，并防止 judge 与 generator 的同源偏差。

## 八、安全与数据治理

- 文档级/段落级 ACL 必须在检索阶段生效
- 多租户索引要防过滤旁路和 embedding 泄漏
- 外部文档视为不可信数据，净化脚本、隐藏文本和指令注入
- 记录来源、版本、解析器和 embedding 模型，支持重建索引
- 删除请求必须传播到原文、派生 chunk、缓存、向量和摘要
- 对个人信息、秘密和版权数据设置保留与输出策略

## 九、上线方法

1. 建立简单 BM25 和 dense baseline
2. 先评解析与 Recall@K，再优化生成 Prompt
3. 用失败查询决定是否加入 hybrid、reranker 或 query decomposition
4. 建立 golden set、困难负例、时间敏感和权限测试
5. 观测 zero-hit、低分、冲突证据、拒答、引用错误、延迟和成本
6. 对索引、模型、Prompt、解析器的每次变更做离线回归和小流量发布

> RAG 的质量上限由“可检索证据的质量”决定，而不是由最后一次 LLM 调用决定。
