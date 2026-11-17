---
title: "Agent范式详解"
created: 2026-07-15
updated: 2026-07-16
status: stable
type: concept
tags: [ai/agent, ai/react, ai/plan-execute, ai/reflection, ai/lats, ai/multi-agent]
source: "workspace/agents/AI-Engineer-from-scrach/01-Agent/；workspace/agents/pi（当前源码）"
confidence: SYNTHESIZED
related:
  - "AI工程师从零开始（总览）"
  - "多智能体架构与框架"
  - "Agent记忆系统"
---

# Agent范式详解

> 本文以 5 种经典范式建立心智模型，并结合现代原生工具调用模型与 pi 的实现，区分“模型习得的行为策略”和“Agent Harness 提供的控制流”。这些范式不是互斥类别，也不是严格的线性演进阶段。

## 五种范式对比

| 范式 | 核心思路 | 适用场景 | 关键组件 |
|------|----------|----------|----------|
| **ReAct** | 推理/决策 → 行动 → 环境反馈循环 | 路径不确定、需频繁与环境交互 | LLM + 工具协议 + Agent Loop |
| **Plan-and-Execute** | 显式计划、分步执行、必要时重规划 | 长任务、强依赖、需跟踪与审计 | Planner + 计划状态 + Executor/Replanner |
| **Reflection** | 产出 → 评估 → 提炼改进 → 重试 | 有质量标准、可验证或可比较的任务 | Generator + Evaluator + 反思/经验状态 |
| **LATS** | 对候选行动/轨迹进行树搜索 | 单一路径易陷入局部最优、值得投入额外推理预算 | 候选生成 + 选择策略 + 价值评估 + 回传 |
| **Multi-Agent** | 多个自治上下文通过协议协作 | 可并行、需角色隔离或专业分工的任务 | 多 Agent + 编排/通信 + 状态与冲突管理 |

## 1. ReAct（基础循环）

**抽象核心循环**：

```text
Reason/Decide → Action → Observation → Reason/Decide → ... → Final Answer
```

- LLM 根据当前上下文决定回答还是行动
- Harness 执行工具，将结果作为环境反馈回灌给 LLM
- 循环持续到模型不再请求工具，或运行时满足终止条件

经典论文和早期教学实现常把协议写成显式文本：

```text
Thought: ...
Action: search(...)
Observation: ...
```

这不是 ReAct 的必要条件。现代系统通常使用原生 function/tool calling：模型输出结构化 `toolCall`，Harness 返回结构化 `toolResult`。推理可能位于隐藏状态、独立 thinking block 或普通响应中，无须解析 `Thought:` 文本。因此更稳妥的定义是**推理/决策与环境交互交错进行**，而不是特定输出格式。

还需区分两层：

- **模型策略**：何时调用什么工具、看到结果后如何调整，现代工具模型已部分习得这些行为
- **运行时机制**：真正执行工具、回灌 observation、再次调用模型、处理中止和终止，仍由 Agent Harness 提供

## 2. Plan-and-Execute（显式规划层）

> 教学实现可以采用“先由 LLM 分解为若干步骤，再逐步执行”，但步骤数量和 Executor 类型都不是范式定义的一部分。

**架构**：`plan_and_execute(task)` → `steps = plan(task)` → 逐步 `execute_step(step, state)`

- Planner 产生显式、可追踪的计划及完成条件
- Executor 可以是 ReAct Agent、确定性函数、工作流、人工操作或另一个模型，并不必然是 ReAct
- 执行结果更新计划状态；环境变化或步骤失败时可局部或整体重规划
- 价值主要来自可观察性、依赖管理、恢复能力和质量门禁，而不只是“先想再做”

**与 ReAct 对比**：

| 维度 | ReAct | Plan-and-Execute |
|------|-------|------------------|
| 规划状态 | 通常隐含在对话中 | 显式表示并持续更新 |
| 控制粒度 | 每轮根据最新观察决定 | 按步骤、依赖和完成条件调度 |
| 灵活性 | 天然适应即时反馈 | 取决于重规划策略，不必然较低 |
| 主要风险 | 局部短视、遗漏长期约束 | 初始计划僵化、状态与现实漂移 |

现代推理模型具备较强的隐式任务分解能力，因此简单任务未必需要独立 Planner；但模型“心里有计划”不等于系统具备可持久化、可审计、可恢复的 Plan-and-Execute 控制流。

## 3. Reflection（评估与改进层）

> Reflection 是一类“根据对产出或轨迹的评价生成改进信息，并影响后续尝试”的闭环；是否跨任务持久化取决于具体实现。

**核心循环**：

```python
reflections = []
repeat until success or max_trials:
    output = Action(task, reflections, tools)
    score, feedback = Evaluator(task, output)
    if success: return output
    reflections.append(Reflector(feedback))
```

**关键设计**：

- **Evaluator**：优先使用测试、编译器、规则、环境结果等可验证信号；主观任务可使用 LLM-as-a-judge，但要防范同源偏差和评价漂移
- **Reflector**：把失败证据提炼成下一次可执行的改进，而不是泛泛地“再想一遍”
- **状态范围**：反思可以只用于当前重试、保留在当前会话，也可以写入跨任务经验库；持久化不是 Reflection 的必要条件
- **停止条件**：除成功外，还应限制尝试次数、成本和退化，避免无收益的自我循环
- **与普通自检的区别**：边界不是“是否持久”，而是是否存在显式评估信号、改进状态和受控重试闭环

工具报错后模型自行换参数属于基于 observation 的 ReAct 修正；只有当系统显式评价产出或轨迹、形成改进信息并驱动重试时，才更适合称为 Reflection。

## 4. LATS（树搜索）

**核心思路**：在单次贪心轨迹之外，展开并评估多条候选推理/行动路径。LATS 的代表性实现借鉴 MCTS，但实际 Agent 搜索系统不一定完整复刻经典 MCTS。

- 节点可表示对话、环境或部分解答状态，边表示候选思考或动作
- 通过候选生成、选择策略（如 UCB/PUCT）、价值评估和回传更新搜索偏好
- 环境反馈和反思都可以参与节点评价，但“反思节点”不是所有实现的必备结构
- 更高覆盖率的代价是 token、工具调用、延迟以及状态复制成本，通常不适合作为所有任务的默认循环

## 5. Multi-Agent（多智能体协作）

- Agent 通常拥有相对独立的上下文、角色、工具或权限边界
- 协作拓扑可以是主管—执行者、流水线、并行投票、辩论、市场或去中心化网络
- 通信不必依赖共享 thread；可以通过结构化消息、任务队列、工件存储或共享黑板传递结果
- 共享完整历史实现简单，但会造成上下文污染、权限泄漏和 token 膨胀；生产系统通常只传递必要状态和工件
- 多 Agent 的收益来自并行性、上下文隔离和能力边界，而不是角色提示词数量；若任务不可并行或无法验证，增加 Agent 可能只会增加成本与协调失败

## 6. pi：事件驱动的结构化 Tool-Calling Loop

pi 当前核心最准确的描述是：

> **事件驱动、模型主导、支持 steering/follow-up 的结构化工具反馈循环（Structured ReAct / Tool-Calling Feedback Loop）。**

主调用链：

```text
AgentSession.prompt
  → Agent.prompt
    → runAgentLoop
      → streamAssistantResponse
      → 提取结构化 toolCall
      → 校验并执行工具（默认可并行）
      → 将 toolResult 追加到上下文
      → 再次调用模型
      → 无工具调用且无排队消息时结束
```

源码依据：

- `packages/agent/src/agent-loop.ts`：两层循环、模型调用、工具执行和结果回灌
- `packages/agent/src/agent.ts`：Agent 状态、事件归约、steering/follow-up 队列
- `packages/coding-agent/src/core/agent-session.ts`：会话持久化、自动重试、上下文压缩与溢出恢复
- `packages/coding-agent/src/core/sdk.ts`：coding-agent 对核心 `Agent` 的实际装配
- `packages/coding-agent/src/core/system-prompt.ts`：默认提示词没有强制 Plan 或 Reflection 阶段

### 为什么归入 ReAct 家族

pi 具备 `模型决策 → 工具行动 → 环境结果 → 下一轮模型决策` 闭环，但不解析 `Thought/Action/Observation` 文本，而使用 provider 原生的 `toolCall/toolResult`。一条 assistant message 可以包含多个工具调用，默认在顺序预检后并行执行，因此比经典单动作教学版更接近 **Batched/Parallel Structured ReAct**。

### 为什么不是默认 Plan-and-Execute

核心没有 Planner、标准计划对象、步骤状态机或 Replanner。模型可以自行隐式规划，项目指令和扩展也可以要求显式计划，但那是模型行为或上层扩展，不是默认 Agent Loop 的控制结构。pi README 也明确说明核心不内置 plan mode。

### 为什么不是默认 Reflection

核心没有固定的 Evaluator/Critic/Reflector 阶段。工具错误回灌后模型修正属于 ReAct 的 observation feedback；自动重试是可靠性机制；compaction 是上下文与记忆管理。三者都不应直接等同于 Reflection。

### pi 的实际分层

| 层次 | 职责 |
|------|------|
| 模型策略层 | 隐式推理、任务分解、工具选择、根据反馈修正 |
| Agent Loop | 模型调用、工具执行、结果回灌、终止判断 |
| Session 编排层 | 持久化、事件、错误重试、compaction、队列续跑 |
| 扩展层 | 可加入 plan mode、critic、sub-agent、权限门禁或自定义工作流 |

这体现了 pi 的设计取向：核心只提供最小但完整的交互闭环，把更高层认知工作流留给模型、prompt、skill 和 extension。

## 组合关系（不是线性演进）

```text
单轨迹执行：ReAct / Tool-Calling Loop
  ├── 外接显式计划状态 → Plan-and-Execute
  ├── 外接评估与改进闭环 → Reflection
  └── 扩展为多候选轨迹搜索 → LATS / Tree Search

横向编排：Multi-Agent
  └── 每个子 Agent 内部仍可采用 ReAct、Plan-and-Execute、Reflection 或搜索
```

这些维度可以组合。例如 Planner 制定计划，多个 Executor 并行执行，每个 Executor 内部使用 ReAct，Reviewer 用外部测试做 Reflection，关键步骤再用树搜索比较候选方案。

## 现代模型“内化”了什么

现代推理与工具模型通过预训练、指令微调和强化学习，已经部分习得：

- 隐式任务分解与动态规划
- 选择工具并生成结构化参数
- 根据工具结果和错误调整策略
- 在提示或证据驱动下进行自我检查

但应避免说“ReAct 或 Plan-and-Execute 已完全被模型内化”。更准确的边界是：

| 已部分内化到模型 | 仍需 Harness/工作流提供 |
|------------------|------------------------|
| 规划和工具选择策略 | 工具的真实执行与权限控制 |
| 根据 observation 修正 | observation 回灌与再次调度模型 |
| 局部自检和错误恢复 | 可审计计划、步骤状态和恢复点 |
| 生成评价或 critique | 独立验证信号、质量门禁和重试预算 |

模型能力降低了每个任务都搭建重型 Planner/Critic 的必要性，但高风险、长周期、可审计任务仍需要把计划、验证和状态显式外化。

## 学习路径

1. **Tool-Calling Loop / ReAct** — 理解模型决策与环境反馈的基本闭环
2. **状态与终止条件** — 理解消息、工具结果、错误、预算和停止机制
3. **Plan-and-Execute** — 对比模型隐式规划与系统显式计划状态
4. **Reflection** — 用外部验证和受控重试解决“做错了怎么办”
5. **Tree Search** — 用额外计算预算解决“单一路径局部最优”
6. **Multi-Agent** — 在确有并行性或隔离收益时引入协作与编排
