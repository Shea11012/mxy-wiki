---
title: "业务 Agent 的受控执行与上下文治理"
created: 2026-07-18
updated: 2026-07-27
status: stable
type: concept
tags: [ai/business-agent, ai/context, ai/memory, ai/governance]
confidence: INFERRED
contested: false
contradictions: []
related:
  - "Agent 在企业固定流程中的嵌入边界"
  - "Agent工程化实践"
  - "Agent记忆系统"
  - "Part V：Agent 能力百科"
  - "Part X：安全、合规与组织"
---

# 业务 Agent 的受控执行与上下文治理

## 一、业务 Agent 不是缩小版通用 Agent

两者可以共享“构造上下文 → 调用模型 → 工具反馈”的循环，但控制权不同：

| 维度 | 通用 Agent | 业务 Agent |
|---|---|---|
| 目标 | 开放式用户任务 | 明确业务目标 |
| 主流程 | 模型动态规划 | 确定性工作流编排 |
| 工具 | 通用、动态扩展 | 领域化、白名单化 |
| 状态中心 | 会话或任务 | 订单、工单、Case 等业务实体 |
| 权限 | 环境和能力权限 | 身份、租户、对象、字段、额度和审批 |
| 副作用 | 可由用户观察和纠正 | 必须经过策略、幂等、审批和审计 |

> 确定性工作流控制主流程，LLM 负责理解、分类、抽取、草拟和有限决策。

因此，LLM 适合承担理解、分类、抽取、草拟和候选建议；身份、金额、资格、状态转换以及不可逆副作用仍由确定性系统裁决。

## 二、模型只提出候选，系统负责执行

模型输出应分为三档：

1. **非约束性文本**：总结、解释、草稿；
2. **结构化建议**：分类、评分、动作建议；
3. **业务命令候选**：必须经过 Schema、权限、业务规则、幂等与审批校验。

> 模型输出是待验证的提案，不是业务事实或授权指令。

推荐链路是：

```text
LLM 生成提案
→ Schema 校验
→ 身份、租户与对象权限校验
→ 业务规则和金额重算
→ 幂等检查
→ 必要时人工审批
→ 业务服务执行
→ 查询最终状态
→ 写入结构化审计事件
```

对话中的“已批准”或“退款成功”都只是陈述；最终状态必须回到审批系统、支付系统等权威事实源确认。

## 三、业务上下文的八层模型

| 层 | 内容 | 推荐真相源 |
|---|---|---|
| L0 | 运行与安全策略 | 版本化策略、系统配置 |
| L1 | 身份、租户与授权 | 受信任的 Execution Context |
| L2 | 工作流与任务状态 | Workflow / Task Store |
| L3 | 结构化业务事实 | 业务数据库和领域服务 |
| L4 | 按需检索的业务知识 | SOP、政策库、知识库 |
| L5 | 当前会话与用户输入 | Conversation Store |
| L6 | 临时推理与工具结果 | 当前上下文、Working State |
| L7 | 长期偏好与历史经验 | 受治理的 Long-term Memory |

优先级不是简单的 Prompt 排序，而是信任边界：L0–L3 不能被用户输入、检索文本、Memory 或模型推断覆盖。Prompt 中的标签有助于模型识别来源，但真正的权限边界仍在 Policy Engine、工具网关和业务服务中。

## 四、Context Builder：最小视图与渐进式披露

Context Window 是一次调用的可见工作区，不是状态存储。Context Builder 应先识别任务和当前状态，再并行加载所需层级，完成权限过滤、字段裁剪、脱敏、版本冲突处理和 Token Budget，最后保留来源边界进行序列化。

常驻内容只包括精简策略、身份摘要、当前工作流状态、核心对象摘要和近期有效对话；订单明细、历史工单、合同、附件和相似案例按需加载。

```text
对象摘要
→ 相关记录列表
→ 单条记录明细
→ 原始附件或证据
```

大型工具结果应压缩为“确定性结论 + 证据引用”写入 Working State，原始数据继续留在业务系统或对象存储中。这样减少重复 token，也避免旧状态、无关内容和间接 Prompt Injection 淹没关键事实。

## 五、大 Context Window 不能替代状态治理

即使模型支持 1M tokens，也仍需预留 System Prompt、Tool Schema、后续 Tool Result、输出和安全余量。工程问题不只是“能否放下”，还包括延迟、成本、吞吐、缓存、冲突状态、注意力稀释和注入面扩大。

> Context Window 是一次模型调用的容量，不是业务状态存储。

Context 不能解决跨日恢复、服务重启、人工转交、供应商切换或审计举证；这些分别依赖 Workflow Store、Business Database、Conversation Store、Audit Store 和 Object Storage。

## 六、不要用一个 Memory 包打天下

业务 Agent 至少要区分三类状态：

| 状态 | 存什么 | 生命周期与约束 |
|---|---|---|
| Authoritative State | 订单、合同、支付、审批等业务事实 | 强 Schema、事务一致性，模型不能任意写 |
| Working State | 步骤、工具摘要、假设、证据、恢复点 | 绑定 Task / Case，任务结束后归档或清理 |
| Long-term Memory | 已确认偏好、弱结构化背景、审核经验 | 按需召回，不能覆盖权威事实 |

其优先级是：

```text
Authoritative State
> Working State
> Long-term Memory
> 模型当前推断
```

长期 Memory 只有在信息跨任务仍有价值、尚无权威系统保存、不属于当前 Case 状态，并且具备来源、作用域、权限、有效期和修正机制时才值得写入。

## 七、Memory 的写入与召回门控

写入不能由“模型觉得有用”直接触发：

```text
候选 Memory
→ Schema 与敏感数据检查
→ tenant / subject / scope 检查
→ 证据与冲突检查
→ 自动接受、人工审核或拒绝
```

召回时先按 tenant、ACL、subject、状态和时间做确定性过滤，再执行关键词与语义检索、重排、去重和冲突处理。Memory 必须支持 `supersede`、`revoke`、`expire`、冲突降权、用户纠正及租户删除，不能只增不减。

## 八、落地顺序

1. **先不建设长期 Memory**：建立权威事实源、Workflow State、Conversation Store、证据审计和按需知识检索；
2. **增加受控偏好**：只保存用户明确确认、低风险且可修改的偏好；
3. **增加 Case / Customer Summary**：必须带证据、版本、权限和人工修正；
4. **增加经验型 Memory**：只把审核后的经验晋升为组织知识，并评测过期、冲突和实际效果。

## 核心心智模型

```text
Context Window   = 本轮模型可见的工作区
Working State    = 可暂停和恢复的任务进度
Business State   = 权威业务事实
Long-term Memory = 受治理的跨任务辅助信息
```

业务 Agent 的目标不是让模型拥有更多自由，而是把非确定性判断嵌入可验证、可恢复、可审批和可审计的业务闭环。
