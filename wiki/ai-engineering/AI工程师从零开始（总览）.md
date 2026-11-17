---
title: "AI工程师从零开始（总览）"
created: 2026-07-15
updated: 2026-07-16
status: stable
type: synthesis
tags: [ai/agent, ai/llm, ai/rag, ai/memory, ai/engineering]
source: "workspace/agents/AI-Engineer-from-scrach/；workspace/agents/pi；2026-07框架官方资料"
confidence: SYNTHESIZED
related:
  - "Agent范式详解"
  - "RAG检索增强生成全链路"
  - "Agent记忆系统"
  - "多智能体架构与框架"
  - "模型路由与工程化"
  - "Agent工程化实践"
  - "从零实现LLM"
  - "Hermes Agent架构与Memory四层"
  - "个人思考与扩展"
---

# AI工程师从零开始（总览）

> 这组笔记从课程代码出发，已按当前工程实践重新整理。目标不是背框架 API，而是理解模型、检索、状态、工具、控制流、评测和治理如何组成可运行系统。

## 一、完整知识地图

```text
模型基础与后训练
        ↓
Prompt / Structured Output / Tool Calling
        ↓
Agent Loop ─────── Context / RAG / Memory
        ↓                    ↓
Task/Workflow / Multi-Agent / Routing
        ↓
Security / Reliability / Observability / Eval
        ↓
Deployment / Feedback / Data Flywheel
```

原课程的“基础概念 → RAG → 记忆 → 多智能体 → 路由 → 工程化 → LLM”适合教学顺序，但生产系统不是单向流水线：评测、安全、数据治理和可观测性应从第一天介入。

## 二、页面索引

| 页面 | 关键问题 |
|------|----------|
| [[Agent范式详解]] | Tool-Calling/ReAct、显式计划、Reflection、搜索和多 Agent 如何组合；pi 属于什么范式 |
| [[RAG检索增强生成全链路]] | 从解析、索引、混合检索、rerank 到引用验证、ACL 和评测 |
| [[Agent记忆系统]] | 工作/情节/语义/程序记忆；写入、召回、冲突、删除和安全 |
| [[多智能体架构与框架]] | 何时需要多个 Agent；拓扑、协议、状态、框架和评测 |
| [[模型路由与工程化]] | 能力/质量/成本/合规路由；缓存、熔断、发布与治理 |
| [[Agent工程化实践]] | Harness 分层、工具契约、权限、恢复、观测和上线门禁 |
| [[从零实现LLM]] | Transformer、数据、训练系统、MoE、SFT、偏好优化和 RL |
| [[Hermes Agent架构与Memory四层]] | Hermes 课程产品快照及可迁移的 Session/Context/Memory 设计 |
| [[个人思考与扩展]] | AI Coding、Memory 和 Agent Eval 的统一工程方法 |

## 三、核心概念边界

### Model、Agent 与 Harness

```text
Model：学习到的推理/生成/行动策略
Agent Loop：模型决策 → 工具行动 → 环境反馈 → 再决策
Harness：工具、上下文、权限、状态、运行时、观测和评测
Agent Product：Model + Harness + Domain Environment + Governance
```

现代模型已经部分习得规划、工具选择和错误修正，但工具真实执行、权限、状态持久化、验证和恢复仍由系统提供。

### RAG、Memory 与 Context

- **Context**：本次模型调用可见的信息
- **RAG**：从外部知识源检索证据并注入 context
- **Memory**：跨时间保存和召回用户/任务/经验状态的生命周期系统
- **Compaction**：有损压缩历史以控制 context，不等于长期记忆

### Workflow 与 Agent

- Workflow 用显式状态和分支保证过程
- Agent 让模型根据环境动态决策
- 生产系统通常是混合体：确定性外壳 + 局部 agentic loop

## 四、Agent 范式速览

| 范式 | 运行时特征 | 适用 |
|------|------------|------|
| Tool-Calling/ReAct | 模型行动与 observation 交错 | 路径未知、工具反馈频繁 |
| Plan-and-Execute | 显式计划、步骤状态和重规划 | 长任务、依赖和审计 |
| Reflection | 评价、改进状态和受控重试 | 有 grader、值得重复尝试 |
| Tree Search/LATS | 展开并比较多条候选轨迹 | 单轨迹易局部最优、预算充足 |
| Multi-Agent | 多个隔离执行上下文通过协议协作 | 并行、权限/上下文隔离、专业化 |

它们不是成熟度阶梯，可以组合。详见 [[Agent范式详解]]。

## 五、LLM 生命周期

```text
数据治理与 Tokenizer
 → Pretraining / Continued Pretraining
 → SFT
 → Preference Optimization（DPO 等）和/或 RL（PPO、GRPO/RLVR 等）
 → Capability/Safety Eval
 → Serving Optimization
 → 线上反馈与下一轮数据
```

Reward Model + PPO 只是经典 RLHF 路径，不是唯一后训练方案。GRPO 也不是已经取代 PPO 的终局算法；group-based RL 正持续演进。

## 六、生产级横切能力

| 能力 | 最低要求 |
|------|----------|
| 安全 | 身份、最小权限、沙箱、审批、secret、注入和数据治理 |
| 可靠性 | deadline、幂等、取消、错误分类、重试/熔断、checkpoint、人工接管 |
| 可观测性 | session/run/trace、模型/工具/状态 span、token/成本/延迟、脱敏 |
| 评测 | 隔离环境、任务集、确定性 grader 优先、多次 rollout、回归和安全测试 |
| 发布 | shadow/canary、切片指标、自动回滚、版本与 lineage |
| 经济性 | 每成功任务成本，而不是单次模型价格 |

## 七、推荐学习路径

1. 实现最小 tool-calling loop，理解消息和终止条件
2. 学习结构化输出、工具 schema、权限和错误处理
3. 建一个有 retrieval eval 的基础 RAG，而不只做向量 Demo
4. 区分 transcript、task state、compaction 和长期 memory
5. 为单 Agent 建环境、grader、trace 和成本基线
6. 再引入显式计划、router、sub-agent 或 multi-agent，并证明净收益
7. 深入 Transformer、训练数据、分布式训练和后训练目标
8. 将安全、可靠性和 eval 作为持续发布系统，而不是上线前补丁

## 八、最重要的工程判断

- 不因模型上下文更长就放弃检索和上下文治理
- 不因模型会规划就假设系统拥有可恢复计划
- 不因模型会自评就把 Reflection 当成独立验证
- 不因加了引用就认为答案已 grounded
- 不因多个 Agent 一致就认为结论正确
- 不因 trace 完整就认为系统通过 eval
- 不因 schema 合法就认为工具调用被授权
- 不因 benchmark 上升就认为生产可靠性上升

> AI Engineering 的核心是把概率模型放入一个可验证、可约束、可观察、可恢复的系统，而不是堆叠更多抽象名词。
