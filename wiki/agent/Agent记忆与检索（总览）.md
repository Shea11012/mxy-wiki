---
title: "Agent记忆与检索（总览）"
created: 2026-06-09
updated: 2026-06-11
status: stable
type: synthesis
tags: [agent/memory, agent/rag, agent/hello-agents]
source: "raw/2026/06/09/第八章 记忆与检索.md"
confidence: INFERRED
contested: false
contradictions: []
related:
  - "Agent记忆系统设计"
  - "RAG检索增强生成"
---

# Agent记忆与检索（总览）

> 源文摘要：Datawhale HelloAgents 框架第八章，借鉴人类认知科学（感觉记忆、工作记忆、长期记忆）为智能体实现记忆系统（MemoryTool + MemoryManager），同时引入 RAG 检索增强生成（五层七步管道 + MQE + HyDE），最后以 Gradio PDF 学习助手为完整案例展示两者协同。

## 拆分页面

- [[Agent记忆系统设计]] — 四种记忆类型（工作/情景/语义/感知）架构、评分公式、MemoryTool 操作
- [[RAG检索增强生成]] — RAG 三阶段发展、五层七步管道、MQE + HyDE 高级检索、MarkItDown 文档处理

## 关联页面

- [[上下文工程]] — ContextBuilder 将记忆与 RAG 输出整合为 GSSC 流水线，NoteTool 实现结构化笔记

## 四种记忆评分速查

| 记忆类型 | 相似度计算 | 时间因子 | 存储 |
|----------|-----------|---------|------|
| 工作记忆 | TF-IDF×0.7 + 关键词×0.3 | ✓ | 纯内存 |
| 情景记忆 | 向量×0.8 + 时间近因性×0.2 | ✓ | SQLite+Qdrant |
| 语义记忆 | 向量×0.7 + 图相似度×0.3 | ✗ | Neo4j+Qdrant |
| 感知记忆 | 向量×0.8 + 时间近因性×0.2 | ✓ | Qdrant(按模态分 collection) |

> **语义记忆无时间衰减**：知识性内容不随时间贬值，这与情景记忆（事件越近越重要）的设计哲学互补。
