---
title: "Part V：Agent 能力百科"
created: 2026-07-18
updated: 2026-07-22
status: snapshot
type: source-summary
tags: [ai/enterprise-agent-platform]
collection: "企业级 Agent 平台工程"
collection_kind: book
collection_order: 5
collection_home: "（总览）"
source: "datagallery-lab/enterprise_agent_platform@e5d97a66c467045413fff692d28fd2c401503029"
sources:
  - "[[sources/ch22-agent-runtime]]"
  - "[[sources/ch23-tool-registry-function-calling]]"
  - "[[sources/ch24-mcp]]"
  - "[[sources/ch25-planner]]"
  - "[[sources/ch26-agentic-workflow]]"
  - "[[sources/ch27-memory]]"
  - "[[sources/ch28-agent]]"
  - "[[sources/ch29-agent]]"
  - "[[sources/ch30-human-in-the-loop]]"
  - "[[sources/ch31]]"
source_kind: third-party
as_of: 2026-07-18
confidence: SYNTHESIZED
contested: false
contradictions: []
related:
  - "（总览）"
  - "Agent工程化实践"
  - "多智能体架构与框架"
  - "Agent记忆系统设计"
  - "MCP协议与自定义服务器"
  - "业务 Agent 的受控执行与上下文治理"
---
# Part V：Agent 能力百科

> 来源：企业级 Agent 平台工程，第22-31章  
> 总结日期：2026-07-18  
> confidence: SYNTHESIZED

---

## 一、章节覆盖表

| 章节 | 标题 | URL | 核心主题 |
|------|------|-----|----------|
| ch22 | Agent Runtime | [链接](https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch22-agent-runtime.md) | Run/Step/Tool Call对象模型、Run六态、SSE事件流、检查点与恢复、失败分类 |
| ch23 | Tool Registry & Function Calling | [链接](https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch23-tool-registry-function-calling.md) | ToolSpec注册模型、版本治理、参数校验、调用链路 |
| ch24 | MCP与企业工具生态 | [链接](https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch24-mcp.md) | Host/Client/Server架构、三类能力、MCP→Registry集成 |
| ch25 | Planner与编排模式 | [链接](https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch25-planner.md) | PlannerDecision、ReAct、Plan-and-Execute、状态图映射 |
| ch26 | Agentic Workflow | [链接](https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch26-agentic-workflow.md) | Reflexion/Self-Refine/ToT增强机制、计量口径、开关设计 |
| ch27 | Memory系统 | [链接](https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch27-memory.md) | Working/Episodic/Profile/Org Context四层记忆、与RAG分工 |
| ch28 | 多Agent协作 | [链接](https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch28-agent.md) | 角色分工、Handoff契约、路由与能力发现、冲突仲裁 |
| ch29 | Agent协议与标准 | [链接](https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch29-agent.md) | MCP/A2A/Agent Card/ACP版图、协议治理、安全审查 |
| ch30 | Human-in-the-Loop与长任务 | [链接](https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch30-human-in-the-loop.md) | 审批模式、waiting_human状态、双检查点、异步队列 |
| ch31 | 框架横向对标 | [链接](https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch31.md) | LangGraph/AutoGen/CrewAI/Dify/Coze能力矩阵、自研vs采购 |

---

## 二、Part 核心问题

Part V 回答的核心问题是：**企业Agent平台如何在生产环境中安全、可审计、可恢复地运行Agent**。

> [EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch22-agent-runtime.md] "Agent Runtime 是很多企业试点走向生产时才真正意识到的缺口。"

具体分解为五个子问题：

1. **执行契约**：一次Agent任务如何跨多轮规划、工具调用、审批等待和进程重启后继续？(ch22 Runtime)
2. **工具治理**：模型提出工具调用意图后，平台如何确保工具存在、参数合法、版本正确、权限充分？(ch23 Registry)
3. **能力接入**：外部工具和Agent如何通过MCP/A2A等协议接入，而不破坏内部治理？(ch24, ch29)
4. **决策编排**：Planner如何组织推理与行动、何时反思或自修，同时不绕过Runtime的安全边界？(ch25, ch26)
5. **人与状态**：人工审批如何嵌入Agent链路？长任务如何保证可恢复？(ch30)

> [EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch22-agent-runtime.md] "没有 Runtime，一次看似普通的经营分析就会变成一串无法复盘的异步动作。"

---

## 三、章节递进逻辑

```
ch22 Runtime (执行内核)
  ├─→ ch23 Tool Registry (工具治理层)
  │     └─→ ch24 MCP (外部工具协议接入)
  ├─→ ch25 Planner (决策编排)
  │     └─→ ch26 Agentic Workflow (增强机制)
  ├─→ ch27 Memory (上下文管理)
  ├─→ ch28 多Agent (角色协作)
  │     └─→ ch29 协议标准 (跨组织互操作)
  ├─→ ch30 HITL (人工介入与长任务)
  └─→ ch31 框架对标 (能力矩阵与选型)
```

递进链：**先建立Runtime执行内核(ch22)，再逐层叠加工具治理(ch23)、外部协议接入(ch24)、决策编排(ch25/26)、记忆(ch27)、多角色(ch28/29)、人工介入(ch30)，最后用框架对标(ch31)把前9章能力映射到业界工具上。**

> [INFERRED] Part V从底层Runtime到顶层框架选型的递进结构，体现了"先建内核、再扩能力、最后对标"的平台工程方法论。

---

## 四、关键结论

### 4.1 Runtime是平台的最小内核

> [EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch22-agent-runtime.md] "企业里最危险的做法是让每个 Agent 应用自己写执行循环。"

Run/Step/Tool Call三层对象分离是Runtime设计的基石：Run管任务生命周期，Step管推理轮次，Tool Call管副作用记录。Run六态(pending→planning→executing→waiting_human→succeeded/failed)对外稳定，内部编排可复杂。SSE事件流(action/result必须成对)、检查点(引擎+业务双检查点)、恢复策略共同构成生产级Runtime。

> [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch22-agent-runtime.md] — 概括：终态只能由 Runtime 触发。模型可以在文本中说'任务完成'，但 Runtime 还要确认工具队列为空。

### 4.2 提议与执行必须分离

> [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch25-planner.md] — 概括：Planner 只提议，不执行。

这是Part V最核心的设计原则。Planner产出PlannerDecision，Runtime负责执行、校验、记录和恢复。Planner的finish=true只是建议，Runtime确认无未完成Tool Call后才进入succeeded。

> [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch23-tool-registry-function-calling.md] — 概括：模型输出不能替代平台校验。

Function Calling只产生调用意图，Tool Registry在handler执行前做强制参数校验、版本解析和权限判断。

### 4.3 外部协议必须经过内部Registry

> [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch24-mcp.md] — 概括：MCP 不替代 Tool Registry。正确路径是 MCP 工具先注册进 Registry，再由 Runtime 通过 Registry 调用。

MCP解决跨进程工具发现和调用，Registry管企业内部命名、版本、权限、风险等级和审计。A2A同样：外部Agent Task映射为内部Tool Call或Handoff。

> [EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch29-agent.md] "协议解决跨边界互操作，平台内核解决运行时状态、权限、检查点、审计和恢复。"

### 4.4 增强机制要可关闭、可计量

> [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch26-agentic-workflow.md] — 概括：全部增强默认关闭，按 Agent 配置按需启用，且所有额外 LLM 轮次仍计入 Step，受 max_steps 与 Gateway 配额约束。

Reflexion/Self-Refine/ToT等增强能力叠加在Planner之上，不改变Run六态。三个计数器(step_index、llm_call_count、tool_call_count)区分编排级和增强级的LLM消耗。

### 4.5 Memory必须分层治理

> [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch27-memory.md] — 概括：Memory 的难点不在'记住'，而在'该记什么、记多久、谁能看、什么时候失效'。

四层记忆(Working/Episodic/Profile/Org Context)生命周期和权限完全不同，不应混为"一个记忆向量库"。检查点必须包含Working Memory快照。Memory与RAG分工明确：Memory面向用户和任务上下文，RAG面向文档和知识库。

### 4.6 多Agent的价值在责任拆分

> [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch28-agent.md] — 概括：多 Agent 不应成为默认架构。拆分只有在换来清晰的权限隔离、组织分工或并行专长时才值得。

Handoff必须是结构化的控制权转移(特殊Tool Call)，有handoff_id、payload、return_policy。路由采用混合策略：规则→分类模型→Agent Catalog→澄清。

### 4.7 HITL是组织授权机制，不是模型兜底

> [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch30-human-in-the-loop.md] — 概括：人工应按用来承接组织责任：谁批准导出客户名单，谁确认折扣写入主数据，谁决定报告可以对外发送。

HITL的waiting_human是Run六态的正式状态，审批通过后同一run_id恢复。双检查点(引擎+业务)保证崩溃恢复和业务进度展示对齐。

---

## 五、概念边界

### 5.1 Run vs 会话 vs Step vs Tool Call

> [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch22-agent-runtime.md] — 概括：会话面向 UI；Run 面向任务；Step 面向推理轮次；Tool Call 面向副作用。

| 概念 | 面向 | 关键ID | 边界 |
|------|------|--------|------|
| 会话 | UI | conversation_id | 聊天窗口 |
| Run | 任务 | run_id | 一次可审计任务的完整生命周期 |
| Step | 推理 | step_index | Planner的一轮决策 |
| Tool Call | 副作用 | tool_call_id | 一次实际工具执行 |

### 5.2 Planner vs Runtime

> [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch25-planner.md] — 概括：Planner 决定下一步调用哪个工具、是否请求澄清、是否结束任务，但它不应该绕过 Runtime 直接执行动作。

- **Planner**：读取上下文→生成PlannerDecision(finish/tool/args)
- **Runtime**：推进Run六态、执行工具、推送SSE、写检查点、失败恢复
- **Registry**：工具注册、版本解析、参数校验、handler调用

### 5.3 MCP vs Tool Registry vs Runtime

> [EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch24-mcp.md] "MCP 解决跨进程工具发现和调用；企业平台还要解决谁能调用、调用什么版本、结果进入哪些日志、失败如何"

- **MCP**：L3协议层，暴露Tools/Resources/Prompts
- **Tool Registry**：L1/L2，管命名、版本、权限、风险等级
- **Runtime**：L2，管Run状态、Tool Call记录、SSE事件

### 5.4 Memory vs RAG

> [EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch27-memory.md] "RAG 面向文档、知识库、表结构和政策，强调引用来源和可追溯事实。Memory 面向用户、任务和运行"

### 5.5 编排模式 vs Workflow增强

> [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch26-agentic-workflow.md] — 概括：编排模式决定'先规划再执行还是边想边做'；Workflow 增强决定'单步决策是否允许内部多轮 LLM 或并行候选'。

- **编排模式**(ch25)：ReAct、Plan-and-Execute、状态图
- **Workflow增强**(ch26)：Reflexion、Self-Refine、ToT——默认关闭

### 5.6 内部Handoff vs 外部A2A委托

> [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch28-agent.md /29] — 概括：内部 Handoff 只需要根据 agent_id 找到平台配置；外部委托要经过 A2A、Agent Card、TLS、mTLS 和出站 Policy。

### 5.7 引擎检查点 vs 业务检查点

> [EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch30-human-in-the-loop.md] "引擎检查点服务 Runtime 恢复，业务检查点服务 Console 展示和合规回放。"

- **引擎检查点**：state、step_index、memory_refs、tool_calls、handoff_stack
- **业务检查点**：label、artifact_ref、created_at、approval_status、metadata

---

## 六、平台架构/工作流

### 6.1 完整执行链路

```
用户POST /agents/{id}/run
  → Runtime创建Run(pending)
  → Planner.next_step()  → PlannerDecision
  → Registry.get() + Policy.check()
  → Runtime执行工具  → SSE: action
  → Registry.invoke() → handler
  → Runtime记录结果  → SSE: result
  → [可选] Policy触发need_approval → waiting_human
  → 审批回调 → 恢复执行
  → Planner返回finish=true + 无未完成Tool Call
  → Runtime标记succeeded
```

### 6.2 L1/L2/L3分层

> [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch23-tool-registry-function-calling.md] — 概括：Registry 横跨 L1 与 L2：运维与 Console 通过 L1 注册工具；Runtime 在 L2 按 (name, version) 解析并 invoke。

- **L1 管控面**：注册工具、Agent配置、协议能力发现、Agent Card导入
- **L2 运行时**：RunLoop、Tool invoke、Policy拦截、SSE推送、检查点
- **L3 协议互通**：MCP Client/Server、A2A适配、Agent Card解析、ACP映射

### 6.3 多Agent Handoff工作流

```
Workflow Agent 接收任务
  → Router选择目标Agent (规则→分类→Catalog→澄清)
  → Handoff(from_agent_id, to_agent_id, handoff_id, payload, reason)
  → 目标Agent接收控制权(同一run_id、同一检查点)
  → 完成后Handoff回Workflow或下一Agent
  → 冲突检测(事实/口径/叙事/资源冲突)
  → Workflow汇总 → succeeded
```

> [EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch28-agent.md] "在 Run 六态上，多 Agent 切换不应改变状态模型。"

---

##七、方案权衡

### 7.1 ReAct vs Plan-and-Execute

| 维度 | ReAct | Plan-and-Execute |
|------|-------|-----------------|
| 适用 | 探索性问数、多跳未知 | 路径清楚、强合规、需预审计 |
| 计划可见性 | 分散在每轮Step | 先生成完整计划 |
| 人工审批 | 常按工具或artifact审批 | 可对计划先审批 |
| 风险 | 可能绕路或循环 | 计划错误影响后续步骤 |

> [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch25-planner.md] — 概括：Plan-and-Execute 不是默认更好的选项。对探索性任务，它可能先生成一个看似完整但很快过时的计划。

### 7.2 Reflexion vs Self-Refine vs ToT

| 机制 | 触发 | 目标 | 风险 |
|------|------|------|------|
| Reflexion | Tool failed/empty | 改进参数和行动 | 重复错误工具调用 |
| Self-Refine | finish前 | 改进文本/结构化答案 | 过度润色导致事实漂移 |
| ToT | 复杂多步推理 | 分支搜索最优路径 | 指数级token成本 |

> [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch26-agentic-workflow.md] — 概括：Reflexion 最适合处理'模型可以自己修正'的错误；不适合权限不足、下游系统故障和业务规则冲突。

### 7.3 单Agent vs 多Agent

> [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch28-agent.md] — 概括：单 Agent 足够的典型场景，是一次只读查询加简短解释。拆成 Question Agent、SQL Agent、Report Agent 只会增加路径长度。

拆分判断维度：工具权限域是否不同、Prompt角色是否过载、审计责任是否需分离、是否有并行需求。

### 7.4 自研 vs 采购 vs 混合

> [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch31.md] — 概括：自研适合平台内核：统一 /run、Run 六态、Tool Registry、MCP/A2A 适配、Policy、HITL、Trace。

- **自研**：Runtime内核、Registry、Policy、Trace、Export Bundle
- **采购**：低代码Console、Bot渠道、知识库UI、模型Gateway
- **混合**：框架做Planner原型，平台做运行治理

### 7.5 前置审批 vs 后置审批

| 模式 | 行为 | 适用 |
|------|------|------|
| 前置审批 | 展示拟执行工具和参数，批准后才执行 | 高风险写操作、删除、转账 |
| 后置审批 | 先生成草稿，批准后才发布或外发 | 报告、邮件、公告 |
| 分级审批 | 多个waiting_human节点串联 | 合同、采购、对外声明 |

---

## 八、失败模式

### 8.1 Runtime层失败

| 失败模式 | 表现 | 根因 | 对策 |
|----------|------|------|------|
| 会话/Run混淆 | 审计无法追踪工具调用 | 只用会话ID管理任务 | Run/Step/Tool Call分层 |
| 模型口头完成 | 模型说"完成"但工具队列非空 | 未校验finish条件 | Runtime确认无未完成Tool Call |
| 检查点不完整 | 恢复后重复查数或结果不一致 | 未保存Working Memory快照 | 检查点包含working_snapshot |
| 审批后重开Run | 发布工具执行两次 | 审批通过后新POST /run | 同一run_id恢复 |

> [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch22-agent-runtime.md] — 概括：只保存 state=executing 不够。

### 8.2 Tool Registry层失败

| 失败模式 | 表现 | 根因 | 对策 |
|----------|------|------|------|
| 版本漂移 | 新旧Agent调用不同handler | 工具升级未同步 | (name,version)主键+pin版本 |
| 绕过Registry | 无法统一审计 | 直连数据库/HTTP | 强制invoke入口 |
| schema同源断裂 | Planner看到A、校验用B | 两套schema | to_openai_tool同源导出 |
| 参数校验缺失 | 漏tenant_id导致数据泄露 | 信赖模型输出 | Registry强制校验 |

> [EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch23-tool-registry-function-calling.md] "JSON Schema 负责形状，Policy 负责权限，handler 负责业务执行，三者各有位置"

### 8.3 MCP接入失败

| 失败模式 | 表现 | 根因 | 对策 |
|----------|------|------|------|
| MCP替代Registry | 工具治理失控 | 绕过内部权限与版本 | MCP→Registry→Runtime |
| Run直连MCP Server | 责任归属不清 | Runtime持有MCP Client | Registry handler封装 |
| Resources替代RAG | 检索质量下降 | 混淆"读指定版本"和"搜索" | RAG+Resource并存 |
| Server审计缺失 | 外部调用无记录 | MCP无IAM | Server侧记录run_id+tenant_id |

> [EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch24-mcp.md] "协议统一了，风险却从分散脚本变成分散 Server。"

### 8.4 Planner与增强失败

| 失败模式 | 表现 | 根因 | 对策 |
|----------|------|------|------|
| 无限反思 | SQL错→反思→再错→再反思 | 无max_reflection_rounds | 绑定预算+max_steps |
| 事实漂移 | Self-Refine改写数字 | 润色改变事实槽位 | 事实锚定+只改表达槽位 |
| ToT成本爆炸 | 每次Tool Call跑满深度搜索 | branching×depth | 限高风险写操作前 |
| Planner直调工具 | 无action事件 | 绕过Registry | 提议与执行分离 |

> [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch26-agentic-workflow.md] — 概括：生产环境最怕没有上限的自我修正。

### 8.5 多Agent失败

| 失败模式 | 表现 | 根因 | 对策 |
|----------|------|------|------|
| Handoff循环 | A→B→A反复交接 | 无栈深度限制 | handoff_stack max_depth |
| 低置信度路由 | 错误Agent访问数据 | 勉强选择近似Agent | 进入澄清而非勉强路由 |
| Catalog漂移 | 路由命中废弃Agent | 元数据未更新 | 定期刷新+标记stale |

> [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part05-agent-capabilities/ch/ch28-agent.md] — 概括：低置信度仍然直连数据工具，是多Agent系统中最容易造成权限事故的路径之一。

---

## 九、实践清单

### Runtime落地检查项

- [ ] Run/Step/Tool Call三层对象是否分离？
- [ ] Run六态状态机是否对外暴露稳定？
- [ ] SSE的action/result是否成对出现？
- [ ] 检查点是否包含状态、上下文、工具结果、Memory快照、事件位置？
- [ ] 是否区分run_id(业务)和trace-id(可观测性)？
- [ ] 失败分类是否区分可重试/需审批/终止？

### Tool Registry落地检查项

- [ ] 工具是否以(name,version)主键注册？
- [ ] 是否在invoke前做强制schema校验(非仅依赖模型输出)？
- [ ] handler部署方式是否与Registry invoke解耦？
- [ ] 是否区分注册面(L1管控)和调用面(L2运行时)？
- [ ] 工具描述是否写明适用场景、禁止场景和副作用？

### MCP接入检查项

- [ ] MCP工具是否先注册进Registry再由Runtime调用？
- [ ] Run主循环是否直连MCP Server(不应直连)？
- [ ] Resources是否替代了RAG(不应替代)？
- [ ] Server侧是否记录调用方身份、tenant_id、run_id？
- [ ] 是否保存tools/list快照用于历史回放？

### Planner检查项

- [ ] Planner是否只返回PlannerDecision而不直接执行工具？
- [ ] finish=true是否经Runtime确认无未完成Tool Call？
- [ ] 工具列表是否按状态裁剪(非全量暴露)？
- [ ] 增强机制是否默认关闭且有预算上限？

### Memory检查项

- [ ] 四类Memory(Working/Episodic/Profile/Org Context)是否分治？
- [ ] 检查点是否包含Working Memory快照？
- [ ] 长期记忆是否有晋升审批流程？
- [ ] Memory与RAG索引和权限模型是否分离？
- [ ] 是否支持删除和导出？

### 多Agent检查项

- [ ] Handoff是否为结构化Tool Call(handoff_id/payload/reason)？
- [ ] 路由是否采用混合策略(规则+分类+Catalog+澄清)？
- [ ] 是否有Conflict检测和权威源定义？
- [ ] 角色是否映射到组织责任？

### HITL检查项

- [ ] 审批是否进入waiting_human状态(非if/else)？
- [ ] 审批通过后是否同一run_id恢复(非新开Run)？
- [ ] 是否区分引擎检查点和业务检查点？
- [ ] 审批记录是否绑定artifact hash和策略版本？
- [ ] 审批超时是否有策略(提醒/升级/自动关闭)？

### 框架选型检查项

- [ ] 框架是否可包进企业/run？
- [ ] 工具调用是否强制走Registry？
- [ ] 状态和日志是否进入Trace和Export Bundle？
- [ ] 是否有退出路径？

---

## 十、跨Part关系

| 关联Part/章节 | 关系 | 说明 |
|--------------|------|------|
| Part I (ch1-3) 平台基础 | 概念引用 | L1/L2/L3分层贯穿Part V全部章节 |
| Part II (ch8) 提示工程 | 工具链引用 | Structured Outputs、ToT在提示层的定义 |
| Part IV (ch20) RAG | 分工 | Memory≠RAG(ch27)，Resources≠RAG(ch24) |
| Part VI (ch33) 语义层 | 治理依赖 | 语义层版本影响SQL工具和指标口径 |
| Part VII (ch38) 可观测性 | 工具链 | Trace span记录tool_call、llm_call |
| Part VII (ch39) 评测 | 质量闭环 | Planner评测、路由评测、HITL拒绝样本 |
| Part VII (ch41) 冒烟测试 | 发布门禁 | 协议适配器契约测试 |
| Part VIII (ch45) LLM 网关（LLM Gateway） | 调用链 | Gateway传tools给模型、Planner调Gateway |
| Part X (ch50) Policy | 权限依赖 | Registry假定调用已通过Policy；HITL触发依赖Policy |

> [INFERRED] Part V是全书的能力集成层，它把前四部分定义的平台基础(Part I)、知识管理(Part III)、语义层(Part IV)转化为可执行的Agent运行时，同时为后续的可观测性(Part VI)和治理(Part VI)提供统一的事件、检查和审计数据结构。

---

## 十一、开放问题

1. **Runtime的SLA分级**：不同风险等级的Run是否需要不同SLA？例如只读分析vs写入工单的恢复窗口差异。

2. **工具跨Agent共享的租户隔离**：同一工具被多个Agent、多租户调用时，Registry如何在handler层实现租户过滤？是handler内自主过滤还是Registry中间件统一过滤？

3. **Memory与模型上下文窗口的博弈**：随着上下文窗口增长(128K→1M+)，全文塞入历史vs分层裁剪的权衡点是否会移动？

4. **MCP/A2A协议收敛**：MCP和A2A是否会长期共存还是合并？企业如何避免对特定协议的过度绑定？

5. **Agent-to-Agent的自动化协商**：当两个Agent对同一事实产生冲突时，是否允许自动协商还是必须人介入？

6. **长任务跨组织边界**：外部A2A Task超时后，如何补偿已发生的副作用？跨组织的事务补偿尚无标准。

7. **HITL的人因工程**：审批体验(证据展示、决策辅助、拒绝原因)直接影响HITL是否被业务接受，但本书未展开讨论UI/UX。

8. **框架与平台的演进节奏**：平台内核(Registry/Runtime/Policy)的迭代速度通常慢于框架(LangGraph/AutoGen)，如何管理版本兼容？

### 业务 Agent 补充：平台能力如何落到业务实体

Part V 主要从平台内核解释 Runtime、Registry、Policy、Memory 与 HITL；[[业务 Agent 的受控执行与上下文治理]] 补充了应用侧边界：业务状态应围绕订单、工单或 Case 持久化，领域工具从受信任的 Execution Context 获取租户与身份，模型只产生提案，业务服务负责权威事实和副作用。两者分别回答“平台如何运行 Agent”和“Agent 如何进入具体业务流程”。

> [INFERRED] 以上开放问题源于各章节在工程实践中暴露但未完全解决的边界条件。

---

## 十二、suggested_tags

```
enterprise-agent-platform, agent-runtime, tool-registry, function-calling,
mcp, model-context-protocol, planner, react, plan-and-execute,
agentic-workflow, reflexion, self-refine, tree-of-thoughts,
agent-memory, working-memory, episodic-memory, hitl, human-in-the-loop,
multi-agent, agent-handoff, a2a, agent-card, agent-protocol,
langgraph, autogen, crewai, dify, coze, framework-comparison,
checkpoint, sse, run-state-machine, long-running-tasks, approval-workflow
```

---

## 十三、source_summary

Part V「Agent能力百科」共10章(ch22-ch31)，覆盖了企业Agent平台的完整能力栈：从底层的Runtime执行内核(ch22)、工具治理(ch23)、外部协议接入(ch24)、决策编排(ch25/26)、记忆系统(ch27)、多Agent协作(ch28/29)、人工介入与长任务(ch30)，到顶层的框架横向对标(ch31)。全书始终以DataAgent经营分析场景为贯穿案例（华东区SKU下滑→SQL查询→报告生成→审批发布），通过mini-platform的`core/runtime/`、`core/registry/`、`core/planner/`、`core/memory/`、`core/mcp/`以及`projects/multi-agent-workflow/`等代码目录展示了各章设计的具体实现落点。

核心设计原则：**提议与执行分离**(Planner只提议，Runtime执行)、**协议经过Registry**(MCP/A2A不直连Runtime)、**增强默认关闭**(Reflexion等按需开启有预算)、**记忆分层治理**(Working/Episodic/Profile/Org四层不混用)、**多Agent按责任拆分**(非默认架构)、**HITL是Runtime状态**(非业务if/else)、**框架作为Planner实现而非替代平台内核**。

各章平均引用学术文献2-3篇(Yao et al. 2023的ReAct和ToT、Shinn et al. 2023的Reflexion、Madaan et al. 2023的Self-Refine、Amershi et al. 2014的交互式ML等)，以及业界标准(OpenAI Function Calling、MCP规范、JSON Schema、W3C traceparent等)，为设计决策提供理论基础。

> confidence: SYNTHESIZED — 综合ch22-ch31全文、10章完整覆盖，所有EXTRACTED摘录均可回溯到原文对应段落。
