---
title: "RAG检索增强生成"
created: 2026-06-09
updated: 2026-06-11
status: stable
type: concept
tags: [agent/rag, agent/hello-agents]
source: "raw/2026/06/09/第八章 记忆与检索.md"
confidence: EXTRACTED
contested: false
contradictions: []
related:
  - "Agent记忆系统设计"
  - "Agent记忆与检索（总览）"
  - "稠密检索、稀疏检索与混合检索"
---

# RAG检索增强生成

> 摘录：检索增强生成（Retrieval-Augmented Generation，RAG）是一种结合了信息检索和文本生成的技术。

> 摘录：**检索** 是指从知识库中查询相关内容； **增强** 是将检索结果融入提示词，辅助模型生成； **生成** 则输出兼具准确性与透明度的答案。

## RAG 发展三阶段

| 阶段 | 时期 | 检索方式 | 生成优化 |
|------|------|---------|---------|
| 朴素 RAG | 2020-2021 | TF-IDF / BM25 关键词匹配 | 直接拼接文档 |
| 高级 RAG | 2022-2023 | 稠密嵌入（Dense Embedding）语义检索 | 查询重写、文档分块、重排序 |
| 模块化 RAG | 2023-至今 | 混合检索、MQE、HyDE | CoT 推理、自我反思与修正 |

稠密检索、BM25 稀疏检索、结果融合与神经重排序的区别详见 [[稠密检索、稀疏检索与混合检索]]。

## 五层七步架构

> 摘录：我们的RAG系统的核心架构可以概括为"五层七步"的设计模式

```
用户层：RAGTool统一接口
  ↓
应用层：智能问答、搜索、管理
  ↓  
处理层：文档解析 → 分块 → 向量化
  ↓
存储层：向量数据库、文档存储
  ↓
基础层：嵌入模型、LLM、数据库
```

## 文档处理管道

> 摘录：任意格式文档 → MarkItDown转换 → Markdown文本 → 智能分块 → 向量化 → 存储检索

### MarkItDown 统一转换

> 摘录：使用MarkItDown作为统一的文档转换引擎，支持几乎所有常见的文档格式

微软开源工具，将 PDF/Word/Excel/PPT/图片/音频等统一转为标准 Markdown。PDF 文件额外调用增强处理（`_enhanced_pdf_processing`）。

### 智能分块（Markdown 结构感知）

> 摘录：根据标题层次分割段落，保持语义完整性

分块流程：
1. 解析 Markdown 标题层次（`#`/`##`/`###`），按标题边界分割段落
2. 按 Token 数量合并段落（`chunk_tokens` + `overlap_tokens`）
3. 中英文混合 Token 估算（CJK 字符计 1 token，非 CJK 按空白分词计）

## 高级检索策略

### 多查询扩展 MQE

> 摘录：通过生成语义等价的多样化查询来提高检索召回率的技术

> 摘录：对每个扩展查询并行执行向量检索，获取候选文档池

用 LLM 将用户查询扩展为多个语义等价表述（如 "如何学习Python" → "Python入门教程" / "Python学习方法" / "Python编程指南"），并行检索后合并去重，排序返回 top-k。

### 假设文档嵌入 HyDE

> 摘录：让LLM先生成一个假设性的答案段落，然后用这个答案段落去检索真实文档

"用答案找答案"——先生成假设答案段落，再用该段落的向量去检索真实文档，缩小 query 与 document 的语义鸿沟。适用于专业领域查询。

### 扩展检索框架

> 摘录：扩展检索的核心机制是"扩展-检索-合并"三步流程。

```
原始查询 → [MQE扩展×N] + [HyDE假设文档] → 并行向量检索 → 去重合并 → 按分排序 → top_k
```

关键参数 `candidate_pool_multiplier=4`：候选池 = top_k × 4，确保召回多样性。

### 策略选择建议

>[!note] INFERRED — 基于源文使用场景推断
| 场景 | 策略 | 原因 |
|------|------|------|
| 一般查询 | 启用 MQE | 用词多样性 |
| 专业领域查询 | MQE + HyDE | 语义鸿沟 |
| 性能敏感 | 基础检索 | 减少 LLM 调用 |

## RAGTool 接口

```python
rag_tool.execute("add_text", text="...", document_id="...")     # 添加文本知识
rag_tool.execute("add_document", file_path="pdf", chunk_size=1000)  # 加载文档
rag_tool.execute("search", query="...", limit=5)                # 语义搜索
rag_tool.execute("ask", question="...", enable_advanced_search=True)  # 增强问答
rag_tool.execute("stats")                                       # 统计信息
```

> 摘录：RAGTool作为RAG系统的统一入口，提供了简洁的API接口

RAGTool 的输出由 [[上下文工程|ContextBuilder]] 的 Gather 阶段作为证据源接入 Select 评分，与记忆系统、对话历史共同构成智能体的完整上下文。
