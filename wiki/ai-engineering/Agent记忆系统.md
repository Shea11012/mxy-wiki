---
title: "Agent记忆系统"
created: 2026-07-15
updated: 2026-07-18
status: stable
type: concept
tags: [ai/memory, ai/agent, ai/retrieval, ai/context]
source: "workspace/agents/AI-Engineer-from-scrach/03-memory/"
confidence: SYNTHESIZED
related:
  - "AI工程师从零开始（总览）"
  - "Agent范式详解"
  - "Hermes Agent架构与Memory四层"
  - "RAG检索增强生成全链路"
  - "业务 Agent 的受控执行与上下文治理"
  - "Agent自我模型记忆"
---

# Agent记忆系统

> Agent Memory 不是一个单独数据库，而是一套围绕“写什么、存在哪里、何时召回、如何更新和删除、怎样验证”的状态管理系统。

## 一、先区分四类状态

| 类型 | 内容 | 生命周期 | 典型载体 |
|------|------|----------|----------|
| **工作记忆** | 当前任务所需的近期消息、工具结果和 scratch state | 当前请求/会话 | 模型上下文、短期 state |
| **情节记忆** | 发生过的会话、动作、时间和结果 | 跨会话 | append-only transcript、事件日志 |
| **语义记忆** | 用户偏好、事实、概念、约束和决策 | 中长期 | 结构化数据库、文档、KV、图或向量索引 |
| **程序记忆** | 如何完成任务：流程、规则、技能和操作经验 | 中长期 | skill、runbook、代码、工作流定义 |

“窗口记忆、摘要记忆、长期记忆”是实现策略；上表是信息类型。二者不要混为一谈。

### 另一条分类轴：自我模型记忆

工作、情节、语义和程序记忆按“保存什么”分类。[[Agent自我模型记忆]]则按“记忆服务什么”分类：它用于维持 Agent 对自身身份、关系、价值边界和行为校准方式的连续理解，内容上横跨语义、情节和程序记忆。因此它适合单独建模，但不应直接列为与上述四类平行的第五种底层记忆。

## 二、课程中的四种策略

| 策略 | 作用 | 局限 |
|------|------|------|
| 滑动窗口 | 保留最近 N 轮或 token 预算内消息 | 被移出的信息不是“永久丢失”，前提是原始 transcript 仍持久化；只是当前上下文不可见 |
| 长期存储+检索 | 跨会话保存并按需召回 | 不等于必须使用向量库；结构化查询、全文检索和图检索都可能更合适 |
| 摘要记忆 | 用摘要替代部分原始历史 | 有损、会累积漂移；必须保留原始记录和摘要边界 |
| 三因子打分 | 综合时效、重要性和相关性排序 | 线性加权只是教学起点，分数尺度、过滤和多样性同样重要 |

## 三、写入流水线

```text
候选事件/对话
 → 提取候选事实或经验
 → 来源与权限检查
 → 规范化实体/时间/作用域
 → 与已有记忆比较
 → add / update / supersede / delete / skip
 → 建索引并记录 provenance
```

每条长期记忆建议包含：

```text
id
content / structured value
memory_type
subject / tenant / project scope
source_event_ids
created_at / observed_at
confidence
valid_from / valid_until
supersedes / contradicted_by
sensitivity / retention policy
```

### 为什么不能只 append

用户偏好和环境事实会变化：

```text
旧：用户默认使用 Python 3.11
新：该项目已迁移到 Python 3.13
```

正确行为通常是保留历史与来源，同时将旧事实标记为 superseded，而不是存两条互相冲突的“真相”。

## 四、召回流水线

```text
当前任务
 → 判断是否需要记忆
 → 按租户/用户/项目/时间/权限过滤
 → 多路候选召回
 → 合并、去重、重排、多样化
 → token 预算压缩
 → 带来源注入上下文
```

### 多路检索

- **精确/结构化查询**：ID、配置、日期、明确属性
- **全文/BM25**：专名、代码符号、原句和稀有关键词
- **向量检索**：改写表达和语义近邻
- **图检索**：实体关系、多跳关联
- **时间查询**：最近一次、某个阶段、有效期

向量相似度只是相关性的近似，不应该绕过访问控制、有效期或事实冲突处理。

## 五、从“三因子”扩展为约束排序

教学公式：

```text
score = w_r × recency + w_i × importance + w_s × semantic_relevance
```

生产中更适合：

```text
先硬过滤：tenant / ACL / type / validity / sensitivity
再召回：lexical + semantic + structured
再排序：relevance + recency + importance + confidence + source quality
再约束：dedup + diversity + token budget
```

注意：

- 不同因子必须归一化，否则权重没有可比意义
- recency 不应覆盖稳定但久远的关键偏好
- importance 若完全由同一个 LLM 自评，容易漂移
- MMR 或聚类可避免 Top-K 全是同一事实的近似复述

## 六、摘要与 Compaction

摘要记忆的目标是维持任务连续性，不是创造新的事实源。

生产要求：

1. 保存原始 transcript，摘要可追溯
2. 在完整 turn 边界切分，避免 tool call/result 断裂
3. 区分事实、决策、未完成事项和推测
4. 记录摘要覆盖范围与版本
5. 重要标识符、错误、文件和验证结果尽量原样保留
6. 多次压缩时防止摘要套摘要造成信息漂移
7. 支持用户删除及重新生成受影响摘要

## 七、安全与隐私

记忆会跨会话影响模型，因此是高价值攻击面：

- 外部内容不能未经验证写成用户偏好或系统规则
- Prompt injection、不可见 Unicode、凭据和敏感个人信息需要扫描
- 每个召回查询先做租户/ACL 过滤，不能“先向量搜再过滤”导致旁路泄漏
- 支持同意、查看、更正、遗忘、导出和保留期限
- 记忆内容进入 Prompt、日志、缓存和第三方模型时都要重新评估数据边界
- Skill/程序记忆可能包含可执行指令，安全级别高于普通事实

## 八、记忆评测

| 阶段 | 指标 |
|------|------|
| 写入 | 提取 precision/recall、重复率、冲突更新正确率、敏感信息误写率 |
| 检索 | Recall@K、MRR/NDCG、权限过滤正确率、时间有效性、召回延迟 |
| 使用 | 答案提升、引用正确率、错误记忆依赖率、任务成功率变化 |
| 生命周期 | 过期清理、删除传播、摘要重建、存储与 token 成本 |

评测集应包含：偏好变化、同名实体、多租户隔离、恶意注入、长期未使用事实、相互矛盾来源和删除请求。

## 九、设计原则

1. 原始事件是证据，摘要和长期事实是派生视图
2. 写入与召回都要有明确策略，不能“所有消息都 embed”
3. 先做结构化过滤和权限，再做相关性排序
4. 记忆必须可追溯、可更新、可删除、可评测
5. 当前会话、任务状态、长期事实和技能分开存储
6. 只在可测地提升任务表现时增加复杂记忆层

> 好的记忆系统不是“永不忘记”，而是能够在正确权限和时间范围内，找回正确、仍然有效且足够有用的信息。

## 十、业务系统中的状态优先级

在业务 Agent 中，“Memory”还要进一步拆成三个治理域：Authoritative State 保存订单、支付、合同和审批等权威事实；Working State 保存步骤、工具摘要、证据与恢复点；Long-term Memory 只保存已确认偏好、弱结构化背景和审核经验。其优先级为 `Authoritative State > Working State > Long-term Memory > 模型当前推断`。大 Context Window 不能替代这三类状态，详见 [[业务 Agent 的受控执行与上下文治理]]。
