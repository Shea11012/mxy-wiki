---
title: "Agent工程化实践"
created: 2026-07-15
updated: 2026-07-18
status: stable
type: synthesis
tags: [ai/engineering, ai/harness, ai/agent-ops, ai/security, ai/eval]
source: "workspace/agents/AI-Engineer-from-scrach/06-harnes/；workspace/agents/pi（对照）"
confidence: SYNTHESIZED
related:
  - "AI工程师从零开始（总览）"
  - "Agent范式详解"
  - "模型路由与工程化"
  - "Hermes Agent架构与Memory四层"
  - "业务 Agent 的受控执行与上下文治理"
---

# Agent工程化实践

> Agent 产品不是“多写几段 Prompt”，而是模型能力与 Harness 约束共同形成的系统。模型负责判断，Harness 负责提供可观察、可执行、可约束、可恢复的环境。

## 先纠正一个来源误区

课程中的 `learn-claude-code` 是一个 **Claude Code 风格 Harness 的教学重建项目**，不是 Anthropic Claude Code 官方源码，也不能据此断言 Claude Code 的真实内部实现。它适合学习设计模式；涉及具体产品行为时，应以对应产品当期官方文档和实测为准。

## 一、最小内核

```text
User message
  → 构建模型上下文
  → LLM response
  → 是否包含 toolCall？
      ├─ 是：鉴权/校验 → 执行 → toolResult 回灌 → 再调用 LLM
      └─ 否：结束
```

最小循环只解决“模型能行动”。生产系统还必须解决：状态、权限、故障、观测、评测和副作用。

## 二、教学项目的 20 个机制

| 阶段 | 课程模块 | 应学习的通用机制 |
|------|----------|------------------|
| 基础行动 | Agent Loop / Tool Use | 消息协议、工具注册、参数校验、结果回灌、终止条件 |
| 边界扩展 | Permission / Hooks | 默认拒绝、审批、调用前后拦截、扩展点与核心解耦 |
| 复杂任务 | Todo / Subagent / Skill | 显式任务状态、上下文隔离、按需加载知识 |
| 上下文生命周期 | Compaction / Memory / System Prompt | token 预算、结构化摘要、持久记忆、稳定前缀 |
| 可靠性 | Error Recovery / Task System | 错误分类、幂等重试、持久任务图、恢复点 |
| 长任务 | Background / Cron | 异步执行、取消、通知、租约、超时和调度 |
| 多 Agent | Teams / Protocols / Autonomous / Worktree | 消息协议、所有权、隔离、冲突和回收 |
| 外部能力 | MCP / Comprehensive | 协议适配、能力发现、认证、统一工具池 |

这些模块是教学实现，不代表每个 Agent 都应全部内置。例如 pi 有意保持核心极简，不默认内置 plan mode、sub-agent、todo 或 permission popup，而把它们留给扩展和部署环境。

## 三、生产级 Harness 分层

| 层 | 核心职责 | 关键问题 |
|----|----------|----------|
| 模型接入 | provider、模型能力、流式协议、重试 | 超时、限流、兼容性、成本、数据驻留 |
| Agent Loop | 消息、工具调用、反馈、终止 | 无限循环、重复调用、错误结果回灌 |
| Tool Runtime | schema、执行、取消、结果规范化 | 参数注入、副作用、输出过大、并发冲突 |
| Context | system prompt、历史、RAG、compaction | 信息污染、缓存失效、摘要失真、越权召回 |
| State/Memory | session、任务、工件、长期记忆 | 单一事实源、版本、过期、删除、租户隔离 |
| Orchestration | queue、workflow、sub-agent、HITL | 幂等、租约、恢复、背压、所有权 |
| Governance | 身份、策略、审批、审计 | 最小权限、密钥、供应链、不可抵赖性 |
| Observability/Eval | trace、metrics、dataset、grader | 仅有日志但不知道是否正确 |

## 四、工具设计

### 4.1 契约而非“函数列表”

每个工具至少需要：

- 唯一名称与清晰描述
- 尽可能收窄的参数 schema 和语义约束
- 结构化成功/失败结果
- `AbortSignal`、timeout 与资源上限
- 幂等性说明及 idempotency key（适用于写操作）
- 副作用、权限级别和审计字段
- 大输出分页、截断或工件化策略

Schema 校验只能证明“形状合法”，不能证明“操作被授权”或“业务语义安全”。

### 4.2 并行执行的约束

只并行执行相互独立、无共享写冲突的调用。文件写入、数据库迁移、部署和依赖同一临时状态的调用应顺序化或加锁。并行结果写回时应保持确定性顺序，避免模型上下文因完成时序漂移。

## 五、安全边界

“工具白名单 + Schema + 审计”只是起点。完整控制链应是：

```text
身份/租户
 → 能力发现过滤
 → 参数结构校验
 → 语义策略与资源作用域校验
 → 风险分级/人工审批
 → 沙箱或最小权限执行
 → 输出净化
 → 不可变审计与告警
```

重点风险：

- **Prompt injection**：外部文档和网页是非可信数据，不能提升为系统指令
- **Confused deputy**：模型借用户或工具身份访问不应访问的资源
- **数据外泄**：工具结果、trace、缓存和第三方 provider 都可能泄密
- **持久化投毒**：恶意内容进入 memory、skill、AGENTS.md 后跨会话生效
- **供应链**：扩展、MCP server、npm/pip 包都按可执行代码审查
- **高风险副作用**：付款、删除、发布、发信等应采用 scoped credential、审批和可回滚设计

## 六、可靠性与恢复

### 错误必须分类

| 类型 | 示例 | 策略 |
|------|------|------|
| 瞬时基础设施 | 429、连接中断、5xx | 有界退避、jitter、尊重 `Retry-After` |
| 确定性请求错误 | schema 错、权限拒绝 | 回灌明确错误，修改参数；不要盲重试 |
| Context overflow | 超窗口 | 精确识别、压缩/裁剪、最多一次恢复重试 |
| 工具副作用不确定 | 请求超时但可能已执行 | 先查状态；依靠幂等键，不能直接重放 |
| 模型行为失败 | 循环、偏题、遗漏 | 步数/成本预算、质量门禁、人工接管 |

持久任务需要 checkpoint、attempt、lease owner、heartbeat、取消状态和恢复语义。`retry` 不是可靠性的同义词；不安全重试会放大故障和副作用。

## 七、上下文与记忆

- System prompt 保持稳定、最小、分层组装，便于缓存和审计
- Skills/RAG 按需加载，不把所有知识常驻上下文
- Compaction 是有损压缩，应保存原始历史和摘要边界
- 长期记忆写入应经过提取、验证、去重、更新/冲突处理、过期和删除
- 会话历史、任务状态、长期事实和工件是不同数据类型，不应混成一个向量库

详见 [[Agent记忆系统]]。

## 八、可观测性不是评测

一次 run 应能关联：

```text
trace_id / session_id / run_id
  → model request/response metadata
  → tool call/result
  → queue/retry/approval
  → state change/artifact
  → final outcome/grader
```

应记录 token、缓存、成本、延迟、重试、工具错误、取消和状态变更；敏感内容需脱敏并设置保留期限。

Trace 回答“发生了什么”，Eval 回答“做得对不对”。生产发布必须有：

- 可复现任务集与隔离环境
- 确定性 grader 优先，LLM judge 作为补充
- 多次 rollout 的成功率、方差和成本
- prompt injection、工具故障、权限和动态环境测试
- 新模型、Prompt、工具或 Harness 的回归对比

## 九、部署清单

1. 明确任务边界、资产和威胁模型
2. 为工具定义最小权限、幂等和取消语义
3. 为循环设置 turn、tool、token、时间和费用预算
4. 状态机持久化，并支持重放、恢复和人工接管
5. 隔离不可信输入、执行环境、租户和并行工作区
6. 建立 trace、metrics、审计与敏感数据治理
7. 先建 eval 环境和发布门禁，再接真实副作用
8. 从单 Agent 开始，只有在并行、隔离或专业化有可测收益时才增加子 Agent

## 核心结论

> 模型能力决定“可能做到什么”，Harness 决定“在什么边界内，以多稳定、可审计和可恢复的方式做到”。

## 十、从通用 Harness 到业务 Agent

通用 Agent 可以让模型围绕开放目标动态选择工具；业务 Agent 则应把模型放入确定性工作流，只负责分类、抽取、草拟和候选建议。业务状态围绕订单、工单或 Case 持久化，而非从对话中推断；模型输出也只是待验证提案，不能成为授权、业务事实或副作用执行依据。业务场景中的上下文、状态与 Memory 边界详见 [[业务 Agent 的受控执行与上下文治理]]。
