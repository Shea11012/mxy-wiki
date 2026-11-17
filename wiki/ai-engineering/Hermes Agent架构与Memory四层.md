---
title: "Hermes Agent架构与Memory四层"
created: 2026-07-15
updated: 2026-07-16
status: snapshot
type: case-study
tags: [ai/hermes, ai/agent-loop, ai/context, ai/memory, ai/gateway]
source: "workspace/agents/AI-Engineer-from-scrach/08-hermes-agent/（视频/课程快照）"
source_kind: course-material
as_of: 2026-07-16
confidence: SYNTHESIZED
related:
  - "AI工程师从零开始（总览）"
  - "Agent工程化实践"
  - "Agent记忆系统"
  - "Agent范式详解"
---

# Hermes Agent 架构与 Memory 四层

> 这是基于课程视频和实机笔记整理的**产品快照**，不是对 Hermes 当前版本的永久规格。路径、默认阈值、provider、渠道和更新时机都可能变化；部署前应以当前官方文档、配置和源码为准。

## 一、可复用的架构主线

```text
CLI / Desktop / Gateway / API
              ↓
        Session Runtime
              ↓
 Context Assembly → LLM ↔ Tools → Final Response
      ↑                         ↓
 transcript / compaction / curated memory / external memory
```

Hermes 案例的重要价值不在某个默认数字，而在以下分离：

- **入口与核心分离**：多渠道统一进入 Agent Runtime
- **Session 与 Turn 分离**：Session 管历史容器，Turn 执行一次模型—工具闭环
- **稳定前缀与可变消息分离**：便于 provider prompt cache
- **原始历史、压缩摘要和长期记忆分离**
- **内置记忆与外部 provider/知识库扩展分离**

## 二、Agent Loop

```text
用户消息
 → 构建上下文
 → 调用 LLM
 → toolCall？
    ├─ 是：执行工具 → toolResult 回灌 → 继续
    └─ 否：最终回复 → 持久化/可选记忆更新
```

这属于 [[Agent范式详解#6. pi：事件驱动的结构化 Tool-Calling Loop|结构化 Tool-Calling/ReAct 家族]]。Memory update 是外围生命周期机制，不改变核心工具反馈循环的分类。

## 三、Session、Turn 与 Context

| 概念 | 作用 |
|------|------|
| Turn | 一条用户输入触发的一次完整 Agent run，可含多次工具调用 |
| Session | 多个 turns 共用的对话容器、历史和运行配置 |
| System Prompt | 身份、工具说明、项目上下文和部分长期记忆的组装结果 |
| Message History | 当前 session 的 user/assistant/tool transcript |
| Ephemeral Context | 当前 turn 的检索结果、环境状态或临时提示 |

课程描述了 **frozen snapshot**：session 启动时读取 `SOUL.md`、用户/长期记忆等并形成相对稳定的 system prompt；同 session 中写盘不必立即改写该前缀，活状态可通过 tool result 或消息侧传递。

这是合理的缓存优化模式，但不应普遍化为所有 Agent 产品都必须如此。它的取舍是：

- 优点：prompt cache 稳定、成本与延迟更低、行为易复现
- 代价：持久记忆的新鲜度降低，需要明确“写盘何时对模型可见”

## 四、Context Compression

课程快照中的实现包含：

```text
旧消息
 → 廉价裁剪大 tool output
 → 保留头部/近期尾部
 → 对中间段生成结构化摘要
 → 修复 tool call/result 配对
 → 组装新上下文
```

具体的 50% 触发阈值、20% 目标、保护条数和摘要长度属于版本配置，不是通用最佳实践。生产设计应根据模型窗口、输出预算、工具输出和缓存策略测试。

必须保留：

- 原始 transcript 和摘要覆盖范围
- 关键决定、未完成工作、文件、命令、错误和验证证据
- tool call/result 的原子性
- 摘要失败、取消和 context overflow 的恢复路径

Compaction 是有损上下文管理，不等同于长期 Memory，也不等同于 Reflection。

## 五、Gateway

Gateway 比本地 CLI 多处理：

- 渠道身份到 session key 的映射
- 从持久存储重建历史
- 同 session 并发消息的 queue/steer/interrupt
- 渠道格式、长度和附件适配
- 出站重试、幂等、速率限制与通知

课程中列出的 Telegram、Email、Slack、Discord、WhatsApp、SMS 等是当时能力快照，不能据此保证当前版本全部可用。安全上还需要验证 webhook、渠道身份、租户隔离和出站权限。

## 六、Memory 四层：产品视角与通用视角

| 课程层次 | 产品快照中的载体 | 更通用的类型 |
|----------|------------------|--------------|
| Layer 1 | `memory.md` / `user.md` | 策展后的短语义记忆/用户画像 |
| Layer 2 | `state.db` + session search/FTS | 情节记忆与全文会话检索 |
| Layer 3 | Honcho、Mem0、Supermemory、Hindsight 等 | 可选外部语义/图/用户建模服务 |
| Layer 4 | Obsidian Skill | 人工可编辑的项目知识库/程序化技能 |

“四层”是 Hermes 课程的产品组织方式，不是学术或行业统一标准。更通用的分类见 [[Agent记忆系统]]：工作、情节、语义和程序记忆。

### Layer 1：有界 Markdown

优点：透明、可手改、可审计、成本可控。硬上限能迫使策展，但具体字符数不应固化为通用事实。

写入需要：来源、去重、冲突更新、敏感信息和 Prompt injection 扫描。纯 substring replace 在并发和歧义下脆弱，生产实现还需版本或 CAS。

### Layer 2：Transcript Search

全文检索擅长原句、专名和历史细节；它不等于“无限容量”，数据库仍需保留、归档、删除、索引和隐私策略。自动 prune/vacuum 是否存在及具体行为是版本事实，应实机核对。

### Layer 3：External Provider

Provider 生态和功能变化很快，不应把“同时只能启用一个”“固定五件套生命周期”写成永久事实。选型应关注：

- 数据是否离开本机、驻留和删除能力
- 记忆抽取与召回的可解释性
- namespace/tenant 隔离
- 成本、延迟、锁定和迁移/导出
- 注入、错误事实和冲突处理

### Layer 4：Obsidian

Obsidian/vault 更适合作为人工可见、结构化的项目知识库。它不是自动记忆算法；文件系统权限、索引、链接一致性、并发编辑和恶意笔记注入仍需治理。

## 七、原笔记中需要降级的断言

以下内容只能保留为课程/实测快照，不能提升为稳定结论：

- 固定目录、字符 cap 和压缩百分比
- 某 provider 的内部算法和“社区热门”程度
- external memory 首轮/每轮的固定调用时机
- cron 使用 JSON 还是 SQLite
- 特定渠道和 Desktop/CLI 的路径共享方式
- “完全本地”“所有数据不落云”——实际还取决于模型 provider、外部 memory、telemetry、搜索和 Gateway
- “Semantic 用纯文本关键词而非 RAG”——只描述某层实现，不能概括整个产品的所有记忆能力

## 八、安全与评测补充

- 记忆、skill、项目 context 和网页都是持久/间接 Prompt injection 入口
- system prompt 缓存不能以牺牲撤权、删除和紧急策略更新为代价
- Gateway 渠道必须绑定身份和 scoped permissions
- cron/autonomous task 需要预算、租约、取消、审批和 dead-letter
- 应评测记忆正确率、召回率、错误更新、跨租户泄漏、压缩恢复和动态消息并发

## 九、可迁移的设计原则

1. 稳定前缀与 turn 级动态上下文分离
2. 原始 transcript、摘要、策展事实和知识库各自建模
3. 写入和“何时对当前模型可见”必须有明确语义
4. Gateway 的 session 重建与并发控制属于核心可靠性能力
5. 产品默认值是快照，架构原则才适合长期复用

> Hermes 是一个有价值的案例，但案例笔记必须把可迁移设计与版本相关产品事实分开。
