---
title: "Agent记忆系统设计"
created: 2026-06-09
updated: 2026-06-11
status: stable
type: concept
tags: [agent/memory, agent/hello-agents]
source: "raw/2026/06/09/第八章 记忆与检索.md"
confidence: EXTRACTED
contested: false
contradictions: []
related:
  - "RAG检索增强生成"
  - "Agent记忆与检索（总览）"
---

# Agent记忆系统设计

> 摘录：借鉴人类记忆系统的设计，我们可以理解为什么智能体也需要类似的记忆能力。

## LLM 的两大根本局限

> 摘录：当前的大语言模型虽然强大，但设计上是 **无状态的** 。

> 摘录：除了遗忘对话历史，LLM 的另一个核心局限在于其知识是 **静态的、有限的** 。

基于 LLM 的智能体面临两个根本问题：**无状态导致的对话遗忘**（每次 API 调用互不关联，上下文窗口有限），以及**知识时效性和领域深度不足**（训练数据截止、容易幻觉）。这两个问题分别由记忆系统和 RAG 解决，而 [[上下文工程|ContextBuilder]] 在运行时将两者整合为结构化上下文。

## 架构分层

记忆系统采用四层架构：

```
HelloAgents记忆系统
├── 基础设施层 → MemoryManager（统一调度）、MemoryItem、MemoryConfig、BaseMemory
├── 记忆类型层 → WorkingMemory / EpisodicMemory / SemanticMemory / PerceptualMemory
├── 存储后端层 → QdrantVectorStore / Neo4jGraphStore / SQLiteDocumentStore
└── 嵌入服务层 → DashScopeEmbedding / LocalTransformerEmbedding / TFIDFEmbedding
```

> [!note] MemoryManager 作为协调者，初始化时根据 `enable_*` 参数按需创建各记忆模块，并提供统一的 `add_memory`、`retrieve_memories` 等接口。

## 四种记忆类型

### 工作记忆（Working Memory）

容量 50 条，纯内存存储，TTL 自动清理。检索采用 TF-IDF + 关键词混合评分。

> 摘录：final_score = base_relevance * time_decay * importance_weight

评分公式：`(TF-IDF相似度×0.7 + 关键词匹配×0.3) × 时间衰减 × (0.8 + 重要性×0.4)`。向量检索失败时回退到关键词匹配。

### 情景记忆（Episodic Memory）

SQLite + Qdrant 混合存储，记录具体事件和时间序列。支持结构化预过滤 + 向量语义检索。

> 摘录：base_relevance = vec_score * 0.8 + recency_score * 0.2

> 摘录：importance_weight = 0.8 + (importance * 0.4)

评分公式：`(向量相似度×0.8 + 时间近因性×0.2) × (0.8 + 重要性×0.4)`。时间近因性权 0.2 体现情景记忆强调"最近发生的事"。

### 语义记忆（Semantic Memory）

Neo4j + Qdrant 双引擎，存储抽象知识、概念和规则。添加记忆时自动提取实体和关系构建知识图谱。

> 摘录：向量相似度 × 0.7 + 图相似度 × 0.3

> 摘录：**重要性权重范围\[0.8, 1.2\]** ：避免重要性过度影响相似度排序，保持检索的准确性

评分公式：`(向量相似度×0.7 + 图相似度×0.3) × (0.8 + 重要性×0.4)`。图检索权重 0.3 用于发现概念间的隐含关联，适合需要关系推理的查询。

### 感知记忆（Perceptual Memory）

多模态存储（文本/图像/音频），按模态分离向量空间（CLIP 编码图像、CLAP 编码音频），Qdrant 存储。

> 摘录：base_relevance = vector_score * 0.8 + recency_score * 0.2

评分公式同情景记忆：`(向量相似度×0.8 + 时间近因性×0.2) × (0.8 + 重要性×0.4)`。支持同模态和跨模态检索，跨模态时自动调整评分权重。

## 四种记忆评分公式对比

| 记忆类型 | 相似度计算 | 时间因子 | 存储 |
|----------|-----------|---------|------|
| 工作记忆 | TF-IDF×0.7 + 关键词×0.3 | ✓ | 纯内存 |
| 情景记忆 | 向量×0.8 + 时间近因性×0.2 | ✓ | SQLite+Qdrant |
| 语义记忆 | 向量×0.7 + 图相似度×0.3 | ✗ | Neo4j+Qdrant |
| 感知记忆 | 向量×0.8 + 时间近因性×0.2 | ✓ | Qdrant(按模态分 collection) |

> **语义记忆无时间衰减**：知识性内容不随时间贬值，这与情景记忆（事件越近越重要）的设计哲学互补。

## 记忆操作

MemoryTool 提供统一的 `execute(action, **kwargs)` 接口：

| 操作 | 功能 |
|------|------|
| `add` | 添加记忆，自动补全 session_id + timestamp |
| `search` | 多类型混合检索，支持 min_importance 过滤 |
| `forget` | 三种策略：重要性/时间/容量 |
| `consolidate` | 短期→长期整合（默认重要性>0.7） |
| `summary` / `stats` / `clear_all` | 管理操作 |

### 遗忘机制

> 摘录：删除不重要或过时的信息

三种策略：
- **重要性**：低于阈值 → 删除
- **时间**：超过 max_age_days → 删除
- **容量**：超限时删除最不重要的

### 时间近因性衰减

> 摘录：recency_score = math.exp(-decay_factor * age_hours / 24)

指数衰减模型：24 小时内保持高分，之后逐渐衰减，最低 0.1。模拟人类遗忘曲线。

## 设计要点

- **记忆与 RAG 封装为独立 Tool**：`memory_tool` 和 `rag_tool`，Agent 直接调用，遵循"框架扩展"原则（不新建 Agent 类）
- **数据隔离**：`user_id` 隔离 Memory，`rag_namespace` 隔离 RAG 知识库
- **会话管理**：`session_id` 追踪单次对话过程
