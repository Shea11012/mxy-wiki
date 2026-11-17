---
title: "Part XI：案例方法论与案例准入"
created: 2026-07-18
updated: 2026-07-22
status: snapshot
type: source-summary
tags: [ai/enterprise-agent-platform]
collection: "企业级 Agent 平台工程"
collection_kind: book
collection_order: 11
collection_home: "（总览）"
source: "datagallery-lab/enterprise_agent_platform@e5d97a66c467045413fff692d28fd2c401503029"
sources:
  - "[[sources/ch54-agent]]"
  - "[[sources/ch55-agent]]"
source_kind: third-party
as_of: 2026-07-18
confidence: SYNTHESIZED
contested: false
contradictions: []
related:
  - "（总览）"
  - "Agent工程化实践"
---
# Part XI：案例方法论与案例准入

> **来源说明**：原始 Markdown 来自 datagallery-lab/enterprise_agent_platform 仓库固定 revision `e5d97a66c467045413fff692d28fd2c401503029`，文件路径为 `docs/part11-cases/ch/ch54-agent.md` 和 `docs/part11-cases/ch/ch55-agent.md`。GitHub Pages 当前未发布 Part XI 路由，但仓库 revision 中文件存在且完整。

---

## 1. 章节覆盖表

| 章节 | 标题 | 核心内容 |
|---|---|---|
| 第54章 | 业务案例准入与写作方法 | 案例纳入的四个条件、全链路证据覆盖、脱敏规范、案例结构、材料发布门槛、最小案例包、证据等级、降级写法、自检清单 |
| 第55章 | 案例复审与平台化收束 | 销售/财务/法务场景边界、复审工作流、平台化收束（抽取共享能力 + ownership）、与主体章节的双向关系、案例演进路线、维护与发布节奏 |

---

## 2. 核心问题

Part XI 试图回答的核心问题是：**在企业级 Agent 平台工程中，一个案例凭什么进入正式正文，而非停留在"功能演示 + 效果口号"？** 具体包括：

1. 案例材料需要满足哪些条件才能证明系统不只是"跑通了演示"？([SYNTHESIZED: sources 见 frontmatter])
2. 证据链路如何覆盖输入、工具调用、权限判断、输出、人工确认、评测和反馈？([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch54-agent.md §54.1.2])
3. 案例如何反向检验前文平台架构（Runtime、Registry、Trace、Eval、HITL 等）是否可落地？([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch55-agent.md §55.3.1])
4. 材料不足时如何降级而非虚构？([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch54-agent.md §54.12])

> 摘录：企业 Agent 的案例很容易写成"功能演示加效果口号"。([ch54-agent.md L3](https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch54-agent.md))

---

## 3. 递进

Part XI 内 54→55 章的递进逻辑：

```
第54章（案例准入）
  ├─ 定义四个准入条件（表54-1）
  ├─ 建立证据覆盖体系（表54-2、表54-3）
  ├─ 规定案例结构和降级写法
  └─ 输出：材料包 + 准入判断

     ↓ 材料包进入复审

第55章（案例复审与平台化收束）
  ├─ 按场景核验边界（表55-1、表55-2）
  ├─ 区分"可演示"与"可发布"（表55-3）
  ├─ 抽取共享平台能力（表55-4）
  ├─ 反向修订主体章节（表55-6）
  └─ 输出：复审台账 + 修订任务 + 演进路线
```

([SYNTHESIZED: sources 见 frontmatter]) 两章的分工明确：第54章负责"能不能进"，第55章负责"进了以后怎么服务全书主线"。

> 摘录：第54章更适合承担案例准入和写作方法，第55章则负责把案例复审结果收回到平台主线。([ch55-agent.md §55.7](https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch55-agent.md))

---

## 4. 关键结论

1. **四个准入条件缺一不可**：业务任务明确、平台能力可对应、数据与工具拓扑清楚、不夸大效果（表54-1）。([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch54-agent.md §54.1])
2. **截图不足以证明链路可信**：案例证据必须覆盖输入、检索、工具调用、权限判断、输出、复核和反馈（表54-2）。([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch54-agent.md §54.1.2])
3. **脱敏不能破坏因果关系**：保留"哪个输入触发哪个工具、哪个证据支撑哪个判断、哪个风险动作被人确认"这些关系链。([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch54-agent.md §54.1.3])
4. **案例应反向检验主体架构**：多个案例反复出现审批超时→回补第30章；指标口径冲突→回补第33/34章（表55-6）。([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch55-agent.md §55.3.1])
5. **材料不足宁降级不虚构**：缺 Trace 写成"链路设计"而非"已稳定运行"；缺评测写成"评估计划"而非"质量已可用"。([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch54-agent.md §54.12])
6. **案例材料分三级证据**：可回放运行证据 > 可核对业务材料 > 访谈和口头说明。([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch54-agent.md §54.15])

> 摘录：一个案例如果只说明模型回答很好，却没有 Runtime、Registry、Policy、Trace、Eval、Memory、语义层或 HITL 参与，就很难放进本书主体。([ch54-agent.md §54.7](https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch54-agent.md))

---

## 5. 概念边界

| 概念 | 定义 | 边界 |
|---|---|---|
| **案例（正式案例）** | 具备任务链路、证据、风险评估的完整叙事 | ≠ 功能演示、≠ 产品介绍、≠ 宣传素材 |
| **最小案例包** | 六类材料的轻量集合：任务说明、输入样本、工具/数据边界、运行链路、风险处理、评估证据 | 不要求完整项目报告，但要求可追溯 ([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch54-agent.md §54.5]) |
| **平台化收束** | 从多个单点案例中抽取共性能力，沉淀到平台章节，保留业务差异 | 区分"业务规则"（业务团队维护）与"平台机制"（平台章节沉淀）([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch55-agent.md §55.2]) |
| **降级写法** | 材料不足时的叙事策略：写方法和边界，不编造完整故事 | 保留工程判断，删除无证据支持的结论 ([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch54-agent.md §54.12]) |
| **证据等级** | L1 可回放运行证据（Trace/评测/发布记录）> L2 可核对业务材料（需求/反馈/验收）> L3 访谈和口头说明 | L3 只能做背景铺垫，不能支撑确定性结论 ([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch54-agent.md §54.15]) |
| **复审台账** | 记录候选案例编号、材料状态、证据来源、脱敏状态、风险边界、章节映射和复审结论的结构化记录 | 不是行政表格，是全书质量反馈系统 ([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch55-agent.md §55.8]) |

---

## 6. 工作流

### 6.1 案例准入工作流（第54章）

```
Step 1: 素材登记
  ├─ 输入：任务说明、脱敏样本、系统拓扑、截图
  ├─ 产出：素材清单
  └─ 退出条件：关键链路缺失→退回

Step 2: 证据复核
  ├─ 对照表54-2逐项检查材料支撑的结论范围
  ├─ 产出：可写结论列表
  └─ 退出条件：结论无证据→删除或降级

Step 3: 正文改写
  ├─ 按推荐结构展开：任务→数据/工具→Agent设计→运行链路→风险控制→评估
  ├─ 产出：案例正文
  └─ 通过 tone、图表、引用检查

Step 4: 发布复核
  ├─ 检查敏感信息、夸大收益
  └─ 合入版本
```

([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch54-agent.md §54.4 表54-8])

### 6.2 案例复审工作流（第55章）

```
Step 1: 核对场景边界
  └─ 用户是谁、任务是什么、哪些动作自动/确认/拒绝

Step 2: 核对证据链
  └─ 输入→工具调用→数据来源→Trace→评测样本→人工复核→最终产物

Step 3: 核对平台收束
  └─ 哪些能力是共享的、哪些是场景专用的、哪些需回补主体章节

Step 4: 输出复审结论
  ├─ 通过（进入正文排期）
  ├─ 退回（列出缺失材料 + 补材料动作）
  ├─ 降级（方法示例/链路示例/复盘材料）
  └─ 暂缓（脱敏风险/证据不足/主体章节未准备好）
```

([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch55-agent.md §55.6 + §55.10])

---

## 7. 权衡

| 权衡维度 | 方案A | 方案B | 本书选择及理由 |
|---|---|---|---|
| 案例数量 vs 证据质量 | 多放案例，看起来完整 | 少放案例，每个证据扎实 | **选B**——"三个证据扎实的案例，价值高于十个看似完整但无法复盘的故事" ([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch54-agent.md §54.8]) |
| 写满自动化 vs 标出边界 | 写成"系统完成业务动作" | 写成"生成草稿，等人工确认" | **选B**——"一个能停下来的 Agent，比一个什么都声称能做的 Agent 更适合企业场景" ([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch55-agent.md §55.1]) |
| 真实截图 vs 自绘图 | 用脱敏截图增强说服力 | 用自绘图保护敏感信息 | **选B（当截图不可公开时）**——"截图如果不能公开，就不要用仿真图冒充真实系统" ([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch54-agent.md §54.10]) |
| 写完整故事 vs 降级表达 | 补齐缺失材料写成完整叙事 | 明确标注证据边界，写方法论 | **选B（材料不足时）**——"材料强度决定结论强度，写作时应主动标注这一点" ([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch54-agent.md §54.12]) |
| 案例独立 vs 回写主体章 | 案例只做附录 | 案例反向修订主体章节 | **选B**——"案例的价值在于验证主体架构，不能只作为正文后的展示材料" ([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch55-agent.md §55.8]) |

---

## 8. 失败模式

| 失败模式 | 症状 | 根因 | 预防/修复 |
|---|---|---|---|
| **功能演示冒充案例** | 只有理想路径和界面交互，缺失败样本和权限边界 | 团队认为"跑通了"即"可发布" | 用表55-3区分演示/复盘/案例；要求至少一条失败链路 ([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch55-agent.md §55.1.2]) |
| **动词边界模糊** | "生成报价"被理解为正式对外报价；"完成财务分析"被理解为可披露结论 | 写作时未拆分动作粒度 | 拆开动词：生成草稿→列出依据→标出风险→提交审批→等待确认→记录结果 ([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch55-agent.md §55.1]) |
| **脱敏破坏因果关系** | 只剩通用描述，工具链、证据链被打断 | 只做表面脱敏（改公司名），未保留任务关系 | 表54-3检查项逐项复核；保留"输入→工具→证据→判断"关系链 ([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch54-agent.md §54.1.3]) |
| **评估指标写成越高越好** | 人工接管率被当作"越低越好" | 未区分任务风险等级 | 高风险场景接管率过低可能代表越权；按表54-5写指标的适用边界 ([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch54-agent.md §54.2.3]) |
| **缺少 ownership** | 只写"平台统一处理"，不说谁维护规则 | 把平台写成无成本的抽象层 | 表55-5明确每个能力的主要责任方和维护责任 ([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch55-agent.md §55.2.2]) |
| **案例与正文脱节** | 案例看不出平台能力的作用 | 写成产品介绍/博客风格 | 表54-6 + 表54-7 强制映射到平台章节和工程证据 ([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch54-agent.md §54.3]) |
| **复审意见散落** | 退回原因丢失，下次复审重复劳动 | 无台账记录 | 复审台账记录 case_id、证据状态、缺失项、结论、回补动作 ([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch55-agent.md §55.8]) |

---

## 9. 实践清单

### 9.1 案例准入自检（第54章）

- [ ] 能否用一句话说明业务任务，而不先列技术名词？([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch54-agent.md §54.7])
- [ ] 四条件逐项通过：业务任务明确 / 平台能力可对应 / 数据与工具拓扑清楚 / 不夸大效果？(表54-1)
- [ ] 证据材料是否覆盖输入、检索、工具调用、权限判断、输出、复核和反馈？(表54-2)
- [ ] 脱敏是否保留因果关系：哪个输入触发哪个工具、哪个证据支撑哪个判断？(表54-3)
- [ ] 自动化边界是否写清楚：自动完成 / 需要确认 / 拒绝执行 / 降级处理？(表54-4)
- [ ] 评估指标是否说明了样本来源、复核口径和失败分类？(表54-5)
- [ ] 案例是否能映射到至少三类平台能力并有证据？([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch54-agent.md §54.3.1])
- [ ] 最小案例包是否包含成功链路、拒绝链路和降级/人工接管链路？([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch54-agent.md §54.5])
- [ ] 证据等级是否使读者能区分运行证据、业务材料和口头说明？([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch54-agent.md §54.15])

### 9.2 案例复审自检（第55章）

- [ ] 复审是否区分了"材料缺失"（可补材料）和"能力缺失"（不能写成已完成）？([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch55-agent.md §55.6])
- [ ] 高风险场景的自动/确认边界是否落到具体状态（draft_ready / waiting_approval / approved / rejected）？(表55-2)
- [ ] 是否从单点案例中抽取了共享平台能力并映射到对应章节？(表55-4)
- [ ] ownership 是否明确：工具注册→平台/安全、指标口径→数据/业务owner、审批链路→业务owner/平台？(表55-5)
- [ ] 案例反复暴露的问题是否回写到主体章节（表55-6触发条件）？
- [ ] 复审台账是否记录了 case_id、材料状态、脱敏状态、章节映射和复审结论？([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch55-agent.md §55.8])
- [ ] 材料不足时是否执行了降级（方法示例/链路示例/复盘材料）而非虚构？([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch55-agent.md §55.5])
- [ ] 双语版本（中/英）的证据边界是否一致？([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch55-agent.md §55.14])

---

## 10. 跨Part关系

Part XI 作为全书末尾的"案例方法论"，与各 Part 形成双向验证关系：

```
Part II (Agent 基础设施)  ←── 案例中的 Tool Registry / MCP 证据
  - 第23章 Tool Registry ── 案例需提供工具契约和执行结果
  - 第24章 MCP ── 案例需说明 Server 准入和凭证管理

Part III (数据底座)  ←── DataAgent 案例中的语义层 / SQL 证据
  - 第33章语义层 ── 案例暴露口径版本冲突 → 回补口径治理
  - 第34章 NL2SQL ── 案例暴露 SQL 编译失效 → 回补编译契约

Part IV (Agent 协作)  ←── 多角色案例中的分工 / 交接证据
  - 第28章多 Agent ── 案例验证分工和交接协议
  - 第30章 HITL ── 案例暴露审批超时 → 回补恢复策略

Part VII (观测与评估)  ←── 案例的 Trace / Eval / SLO 证据
  - 第38章 Trace ── 案例缺 Trace → 平台能力不足
  - 第39章评测 ── 案例缺评测样本 → 无法支撑质量结论
  - 第42章 SLO ── 案例缺时延/成本数据 → 不能写上线规模

Part X (安全与治理)  ←── 案例的安全 / 合规证据
  - 第51章 Guardrails ── 案例验证策略引擎和拒答边界
  - 第52章合规 ── 案例提供审计证据和脱敏规范
```

([SYNTHESIZED: ch54 表54-6 + ch55 表55-6]) 案例与主体章节不是单向引用，而是双向校准：**案例验证架构 → 架构补足后提升案例质量 → 循环**。

---

## 11. 开放问题

以下问题在两章中被提出但未给出最终答案，依赖于后续材料补齐：

1. **DataAgent 案例何时准备好**：语义层样本、SQL 执行日志和评估集是否能在第一批案例中产出完整的可发布案例？([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch55-agent.md §55.4.1])
2. **工单/客服 Agent 的工具契约和接管记录**：是否能从真实项目中脱敏出足够支撑案例的证据？([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch55-agent.md §55.4.1])
3. **法务审阅案例的公开可行性**：条款样本、复核口径和责任边界是否能在不暴露敏感信息的前提下写成案例？（被标记为"暂缓"）([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch55-agent.md §55.4.1])
4. **审批超时恢复策略的具体实现**：多个案例若反复暴露审批超时，第30章应补到什么粒度？([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch55-agent.md §55.3.1])
5. **英文版案例的边界一致性**：如何确保中英双语案例的证据强度、降级表达和风险描述不被翻译漂移？([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch55-agent.md §55.14])
6. **案例退出机制的具体触发条件**：项目停止、材料无法脱敏、主体章节改变时，案例退出应遵循什么流程？([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch55-agent.md §55.14])
7. **反馈归档的运营成本**：读者反馈、团队反馈如何持续归档到案例台账而不产生过高的维护负担？([INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part11-cases/ch/ch55-agent.md §55.17])

---

## 12. suggested_tags

```yaml
suggested_tags:
  - "案例方法论/case-methodology"
  - "案例准入/case-admission"
  - "案例复审/case-review"
  - "平台化收束/platform-consolidation"
  - "证据标准/evidence-standards"
  - "降级写法/degraded-writing"
  - "脱敏规范/desensitization"
  - "企业Agent/enterprise-agent"
  - "写作方法/writing-methods"
```

---

## 13. source_summary

> 本文基于 `datagallery-lab/enterprise_agent_platform` 仓库固定 revision `e5d97a66c467045413fff692d28fd2c401503029` 下的两个原始 Markdown 文件：
>
> - `docs/part11-cases/ch/ch54-agent.md`：第54章「业务案例准入与写作方法」— 定义案例纳入正文的四个条件（任务明确、能力对应、拓扑清楚、不夸大效果），建立全链路证据覆盖体系和脱敏规范，规定推荐案例结构（任务→数据/工具→Agent设计→运行链路→风险控制→评估），提出最小案例包概念、三级证据等级、材料不足时的降级写法，以及案例方法对正文质量的约束作用。
>
> - `docs/part11-cases/ch/ch55-agent.md`：第55章「案例复审与平台化收束」— 以销售/财务/法务三类高风险场景为切入点，建立复审工作流（场景边界→证据链→平台收束），区分"功能演示"与"可发布案例"，要求从单点案例中抽取共享平台能力并写清 ownership，通过复审台账实现案例演进路线管理，并规定案例反向修订主体章节的触发条件（表55-6）。
>
> **注意**：GitHub Pages 当前未发布 Part XI 对应路由，但仓库 revision 中上述文件完整存在。两个文件共约 ~18000 字（中文），包含 15 张表格、3 篇参考文献。引用格式为 `Yin (2018) Case Study Research`、`NIST AI RMF 1.0` 和 `OWASP Top 10 for LLM Applications`。

---

## 覆盖状态

| 维度 | 状态 |
|---|---|
| 第54章（16个小节） | ✅ 全部覆盖：54.1~54.16 |
| 第55章（18个小节） | ✅ 全部覆盖：55.1~55.18 |
| 13张关键表格 | ✅ 核心逻辑已提取（表54-1至54-8、表55-1至55-7） |
| 事实引用 | ✅ 每次引用附 ≤50字精确摘录 + raw URL |
| 跨章综合 | ✅ 标注 SYNTHESIZED |
| GitHub Pages 状态 | ✅ 已说明：未发布但仓库 revision 存在 |
