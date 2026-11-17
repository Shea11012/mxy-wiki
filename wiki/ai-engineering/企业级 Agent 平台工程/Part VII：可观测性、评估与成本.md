---
title: "Part VII：可观测性、评估与成本"
created: 2026-07-18
updated: 2026-07-22
status: snapshot
type: source-summary
tags: [ai/enterprise-agent-platform]
collection: "企业级 Agent 平台工程"
collection_kind: book
collection_order: 7
collection_home: "（总览）"
source: "datagallery-lab/enterprise_agent_platform@e5d97a66c467045413fff692d28fd2c401503029"
sources:
  - "[[sources/ch38-trace]]"
  - "[[sources/ch39-dataagent-eval-benchmark]]"
  - "[[sources/ch40-llm-as-judge]]"
  - "[[sources/ch41-cost-governance-cache]]"
  - "[[sources/ch42-slo]]"
source_kind: third-party
as_of: 2026-07-18
confidence: SYNTHESIZED
contested: false
contradictions: []
related:
  - "（总览）"
  - "Agent工程化实践"
  - "模型路由与工程化"
---
# Part VII：可观测性、评估与成本

> 来源摘要：Part VII 覆盖 Agent 平台从「能跑」到「可生产化」的五个关键维度——Trace 可观测性（ch38）、离线 Benchmark 评测（ch39）、在线评测与 LLM-as-Judge（ch40）、成本治理与缓存优化（ch41）、SLO 管理与系统韧性（ch42）。核心主张：Agent 能给出答案 ≠ 具备生产可用性；证据链、评测基线和成本归因缺一不可。

---

## 一、章节覆盖表

| 章节 | 标题 | 核心主题 | 关键输出物 |
|------|------|----------|-----------|
| ch38 | Agent 可观测性与运行诊断 | Session/Run/Trace/Artifact 边界定义；Run 采集点设计；前后端时间线映射；多轮回放；AgentOps 持续改进 | Trace 采集框架、诊断路径表(8类失败)、Context Package 结构、证据链三层次 |
| ch39 | DataAgent 评测体系与 Benchmark | Benchmark 演进史(ML→LLM→Agent→DataAgent)；评测三要素(结果/语义/轨迹)；任务空间矩阵；综合评分公式 | Score_agent 公式、任务空间矩阵(意图×复杂度×企业约束)、能力边界与失败标签映射 |
| ch40 | 在线评测、模型裁判与持续优化 | 离线/在线/人工复核闭环；反馈信号→证据转换；LLM-as-Judge 适用边界；裁判校准与防止 Goodhart | FeedbackSignal 公式、裁判多维评分 rubric、分层评测策略(规则/轻量Judge/强Judge/人工) |
| ch41 | 成本治理与缓存优化 | 成本归因到 Step；模型路由策略；五类缓存边界(提示词/前缀KV/结果/语义)；自建vs外部API token 策略 | Cost_run 公式、Utility(model,task) 路由函数、缓存键权限隔离模型、供应商适配层 |
| ch42 | SLO 管理、限流与系统韧性 | SLO 从接口→Run 的迁移；五场景等待语义；错误预算驱动发布；降级-恢复-清理闭环 | SLO objectives+guardrails 结构、目标/护栏分离、错误预算消耗公式、降级用户语言 |

---

## 二、Part 核心问题

Part VII 回答一个根本问题：**Agent 平台从演示到生产，缺少的不是模型能力，而是「凭什么相信它」的证据系统。**

具体分解为五个子问题：

1. **可观测性 (ch38)**：Agent 失败藏在上下文打包、Planner 决策、工具参数、权限策略中，不在最终回答里。没有 Trace，排障只能靠猜。
   > 摘录：「Agent 的失败常常藏在最终回答之前。用户看到一个结论，背后可能经过上下文检索、Planner 决策、工具调用、权限过滤、SQL 执行、图表生成和报告整理。任何一步出错，最终答案都可能看起来像模型问题。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch38-trace.md]

2. **离线评测 (ch39)**：SQL 可能蒙对口径，解释可能错了归因，数字却恰好正确。评测必须同时看结果、语义上下文和轨迹。Benchmark 是发布门禁，不是演示材料。
   > 摘录：「企业级 DataAgent 评测要同时看答案、上下文和执行轨迹。SQL 可能蒙对了口径，解释可能错了归因，过程可能绕了远路，数字却恰好正确。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch39-dataagent-eval-benchmark.md]

3. **在线评测 (ch40)**：离线通过≠线上可靠。真实流量带噪声，人工全量复核不现实。LLM-as-Judge 扩大覆盖面但不能替代专家判断。
   > 摘录：「LLM-as-Judge 的价值在于扩大初筛覆盖面，不在于替代专家判断。裁判本身也要被校准。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch40-llm-as-judge.md]

4. **成本治理 (ch41)**：账单异常时只盯模型单价往往找错方向。成本失控根因常是长上下文、失控重试、缓存未命中，而非模型太贵。
   > 摘录：「账单异常时，只盯模型单价往往会找错方向：成本可能来自过长上下文、失控重试，也可能来自本可命中的缓存被绕过。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch41-cost-governance-cache.md]

5. **SLO 与韧性 (ch42)**：接口200≠业务可用。Agent SLO 必须从 HTTP 成功率迁移到 Run 级别，覆盖首响应、任务完成、权限不让步、成本不失控。
   > 摘录：「Agent 平台的 SLO 不能只写 HTTP 成功率。它要回答用户到底在等什么：首响应、任务完成、产物可恢复、权限不让步，还是成本不失控。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch42-slo.md]

---

## 三、章节递进逻辑

```
ch38 可观测性 ──→ ch39 离线评测 ──→ ch40 在线评测+裁判 ──→ ch41 成本治理 ──→ ch42 SLO与韧性
    │                  │                    │                    │                  │
    ▼                  ▼                    ▼                    ▼                  ▼
 先有 Trace         用 Trace 构建        线上问题沉淀回       成本归因依赖        SLO 绑定 run_id
 才能诊断           Benchmark           离线回归集           Trace 字段         依赖 Trace 计算
```

**递进关系**：

1. **ch38 → ch39**：Trace 采集的 run_id / trace_id / Context Package 是 benchmark 构建的底座。没有 Trace，评测只能看最终答案，无法区分「SQL 碰巧正确」和「路径正确」。[INFERRED]

2. **ch39 → ch40**：离线 benchmark 固定基线 → 在线评测覆盖真实分布 → 人工复核提供高可信标尺。三者形成质量校准闭环。线上问题沉淀回离线样本，离线修复通过后灰度验证。
   > 摘录：「离线 benchmark 像回归测试，适合防止已知问题反复出现；在线评测像生产观测，适合发现未知问题。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch40-llm-as-judge.md]

3. **ch40 → ch41**：在线评测中裁判的成本需要纳入治理。分层评测策略（规则校验全覆盖 / 轻量 Judge 大比例 / 强 Judge 高风险）本身就是成本治理实践。[INFERRED]

4. **ch41 → ch42**：成本治理的归因字段（run_id / step_id / trace_id）复用了 ch38 的 Trace 体系。SLO 中的「成本可控」目标直接对应 ch41 的单位任务成本指标。[INFERRED]

5. **全 Part 闭环**：Trace(ch38) 采集证据 → Benchmark(ch39) 建基线 → 在线评测+裁判(ch40) 发现长尾 → 成本归因(ch41) 定位费用 → SLO(ch42) 约束稳定性 → 错误预算耗尽时暂停发布，回到 ch39 修复回归。[INFERRED]

---

## 四、关键结论

### 4.1 可观测性 (ch38)

1. **Session ≠ Run ≠ Trace**：Session 是用户连续会话（多 Turn），Run 是一次任务执行，Trace 描述 Run 内步骤。混成"大日志表"后期会在权限、生命周期和排障上付出代价。
   > 摘录：「把它们混成一张'大日志表'，早期实现简单，后期会在权限、生命周期和排障上付出代价。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch38-trace.md]

2. **Context Package 是关键但常被忽略的对象**：模型没有"看到整个会话"，它看到的是 Runtime 组装的一包上下文（最新原文 + 最近几轮 + 早期摘要 + Memory + Schema + Tool Spec + Policy）。排查多轮错误必须看 Context Package。
   > 摘录：「模型没有'看到整个会话'，它看到的是 Runtime 当时组装出来的一包上下文。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch38-trace.md]

3. **Trace 采集原则：默认保存摘要/hash/版本/引用，按权限查看原文**。Prompt、工具返回、数据库结果可能含敏感信息，Trace 不应成为敏感数据的第二份副本。
   > 摘录：「采集时要遵守一个原则：默认保存摘要、hash、版本和引用，按权限查看原文。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch38-trace.md]

4. **前端时间线 ≠ 后台 Trace**：前端服务 UX（少而稳定），后台服务诊断（细粒度 Span）。前端每个状态都应有后台依据，但不能暴露内部实现细节。[INFERRED]

5. **诊断不要用"模型幻觉"盖住所有问题**：很多看似模型胡说的结果，根因是语义层缺字段、工具描述不清、上下文摘要漏约束或权限反馈含糊。
   > 摘录：「诊断时要避免一句'模型幻觉'盖住所有问题。很多看似模型胡说的结果，根因其实是语义层缺字段。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch38-trace.md]

6. **AgentOps 持续改进**：失败 Run / 超时 Run / 高成本 Run / 用户点踩 / 人工接管 / 审批驳回 → 样本池 → 聚类 → benchmark → 修复 → 回归 → 灰度 → 观察。Trace 从"日志"变成"工程资产"。
   > 摘录：「AgentOps 可以按一条简单链路运行：采集 Run 轨迹，筛选高价值样本，按失败类别聚类，沉淀成 benchmark。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch38-trace.md]

### 4.2 离线评测 (ch39)

1. **Benchmark 的四要素**：明确任务定义、可复现数据集、统一评测流程、可解释指标。只有问题+答案无任务边界/版本/脚本，只是练习题。
   > 摘录：「如果只有一批问题和答案，却没有任务边界、数据版本、评测脚本和指标口径，那更像练习题，不能用于工程回归和上线决策。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch39-dataagent-eval-benchmark.md]

2. **Benchmark 演进时间线**：ML 单点评测(BLEU/EM/F1) → LLM 通用(MMLU/BIG-Bench/HELM) → Agent 多步(SWE-bench/AgentBench/WebArena) → Text-to-SQL(Spider/BIRD) → 企业 DataAgent(BEAVER/Workspace-Bench)。本书定位为吸收前四阶段方法的生产质量系统。
   > 摘录：「企业 DataAgent benchmark 可以吸收公开 benchmark 的方法，但上线判断要回到自己的生产问题。它更像一张能力地图，不能只显示总分。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch39-dataagent-eval-benchmark.md]

3. **评测三要素**：结果评测(SQL执行/数值一致)、语义上下文评测(指标口径/schema/source/Memory/policy是否被正确提供)、轨迹评测(是否沿可审计链路执行)。三层都用确定性程序先判，不确定的再给模型裁判或专家。
   > 摘录：「企业评测要同时关注'结果''语义上下文'和'轨迹路径'。只看结果，会放过碰巧答对、越权取数、过度依赖旧记忆的风险。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch39-dataagent-eval-benchmark.md]

4. **综合评分公式**：Score_agent = Score_quality - w_cost·CostPenalty，其中 Score_quality = w_result·Score_result + w_semantic·Score_semantic + w_trajectory·Score_trajectory + w_safety·Score_safety。权重由任务类型决定。
   > 摘录：「这里的权重不是固定常数，应由任务类型决定。财务报表生成更重视安全、口径和可复核。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch39-dataagent-eval-benchmark.md]

5. **任务空间三维度**：任务意图(查询/对比/归因/预测/解释/报告) × 执行复杂度(单表/多表Join/多事实表/跨域/历史快照) × 企业约束(语义明确度/权限敏感度/数据冲突/审计要求/追问可能性)。500条高质量样本 > 5万条无约束自动生成题。[INFERRED]

6. **失败标签必须可修复**：metric_definition_missing → 补指标字典/改Context Builder；tool_affordance_missing → 补工具描述/参数示例；llm_ignored_observation → 改Planner prompt/加状态机约束。评测结果能反哺工程才有价值。
   > 摘录：「Benchmark 不应只告诉团队'这个 case 错了'，还要指出下次应该怎么做才对。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch39-dataagent-eval-benchmark.md]

### 4.3 在线评测与 LLM-as-Judge (ch40)

1. **在线评测 ≠ 收集用户反馈**：要回答四个问题——任务是否完成、完成路径是否可接受、线上信号是否异常、异常能否转化为可修复工程动作。
   > 摘录：「在线评测要回答四个问题：任务是否完成，完成路径是否可接受，线上信号是否异常，异常能否转化为可修复的工程动作。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch40-llm-as-judge.md]

2. **三类反馈信号**：显式反馈(点赞/点踩/评论)→覆盖率低但直观；隐式行为(复制/追问/重新生成/放弃)→覆盖广但解释困难；业务结果(报告进入周会/SQL保存为看板)→真实价值但滞后噪声大。
   > 摘录：「反馈记录不能停在一个 thumb_down 字段。它要把用户信号和执行证据绑在一起。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch40-llm-as-judge.md]

3. **LLM-as-Judge 适用边界**：确定性脚本判 SQL执行/数值一致/字段越权/格式schema；模型裁判判语义质量/论证质量/表达质量。不适合单独判指标口径和权限边界。
   > 摘录：「模型裁判要进入证据系统，而非代替证据系统。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch40-llm-as-judge.md]

4. **裁判三种模式**：单答案打分(批量筛查)、成对比较(A vs B 更稳定)、多维评分(企业DataAgent推荐——拆成正确性/证据支撑/完整性/可执行性/表达适配/安全合规)。
   > 摘录：「多维评分适合企业 DataAgent：它把'报告质量'拆成正确性、证据支撑、完整性、可执行性、表达适配和安全合规。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch40-llm-as-judge.md]

5. **防止"为评测器优化"(Goodhart's Law)**：Judge 分数上升但用户追问率也上升、延迟翻倍、证据引用隐藏更深 → 说明在迎合裁判而非改善体验。Judge 分数必须和真实行为一起看。
   > 摘录：「平台应把 Judge 分数和真实行为放在一起看：用户是否少追问，报告是否更常被采纳，人工退回率是否下降。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch40-llm-as-judge.md]

6. **四层制衡**：规则做硬校验 → 模型裁判做初筛 → 人工专家做校准 → 线上行为做现实反馈。四者互相制衡，不把质量判断交给另一个黑箱。[INFERRED]

### 4.4 成本治理与缓存 (ch41)

1. **先归因再降价**：费用归因到租户/Agent/任务类型/Step。30%账单花在失败重试 → 修复工具错误/超时策略；而不是先换便宜模型。归因依赖 ch38 的 session_id/run_id/step_id/trace_id。
   > 摘录：「成本治理要从归因开始。只有把费用贴回任务、用户、模型、工具和策略，团队才能决定该改路由、加缓存、限制预算还是重构流程。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch41-cost-governance-cache.md]

2. **Cost_run 公式**：Cost_run = ΣCost_model_i + ΣCost_tool_j + Cost_retrieval + Cost_storage + Cost_eval + Cost_retry。模型调用又可拆为 Cost_model = P_in·T_in + P_out·T_out - Discount_cache。
   > 摘录：「重试成本常常被低估：失败 SQL 的自动修复、模型超时后的重跑、用户短时间内重复点击，都会把一次任务的真实成本放大。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch41-cost-governance-cache.md]

3. **单位任务成本比总金额更有用**：每个成功Run平均成本、失败成本占比、缓存节省金额、P95成本。总金额适合财务结算，不适合工程决策。
   > 摘录：「成本看板不宜只展示总金额。总金额适合财务结算，却不适合工程决策。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch41-cost-governance-cache.md]

4. **模型路由函数**：Utility(model,task) = α·Quality - β·Cost - γ·Latency - δ·Risk。路由规则变更必须跑 ch39 Regression Set + Safety Set，再经 ch40 在线灰度观察。
   > 摘录：「路由策略如果没有失败标签和回归保护，就会把省钱变成质量债。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch41-cost-governance-cache.md]

5. **五类缓存，边界不同**：提示词缓存(稳定system prompt/tool spec/schema)、前缀KV缓存(推理服务层共享前缀计算)、结果缓存(确定查询)、语义缓存(相似问题复用)。不能混成一个开关。
   > 摘录：「Agent 平台里至少有五类缓存，它们复用的对象、失效条件和风险都不同。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch41-cost-governance-cache.md]

6. **缓存键第一层是权限边界**：tenant_id + prompt_version + tool_version + semantic_layer_version + policy_version + role + data_snapshot。缺权限隔离 → 性能优化变安全事故。
   > 摘录：「如果缓存没有权限隔离，就可能把高权限用户的上下文或结果复用给低权限用户。这会从性能优化变成安全事故。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch41-cost-governance-cache.md]

7. **自建 vs 外部API token策略**：自建模型关注算力/显存/吞吐/量化；外部API关注稳定前缀结构/供应商缓存参数。企业应有"供应商适配层"隔离差异。
   > 摘录：「业务团队不需要记住某个供应商当前的缓存规则，差异应隔离在LLM 网关里。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch41-cost-governance-cache.md]

### 4.5 SLO 与系统韧性 (ch42)

1. **入口200 ≠ 业务可用**：周一早高峰案例——API返回200但用户拿不到报告，反复点击制造更多重试，形成正反馈放大。Agent SLO 必须从接口级迁移到 Run 级。
   > 摘录：「平台没有完全宕机，接口甚至可能还在返回200，但用户拿不到报告，业务上已经不可用。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch42-slo.md]

2. **五维 SLO 目标**：首响应及时(95%交互式任务5秒内返回进度/澄清)、最终可交付(财务分析成功率≥98%)、质量可接受(回归集通过率不低于阈值)、安全不让步(高风险越权率=0)、成本可控(单成功任务平均成本不超过预算)。
   > 摘录：「企业 Agent 的稳定性，应该被定义为：用户在可接受时间和成本内，安全地完成有质量的任务。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch42-slo.md]

3. **目标(objectives)与护栏(guardrails)分离**：objectives 是希望优化的（延迟/成功率/成本），guardrails 是不能被牺牲的边界（安全违规率=0/回归集通过率≥97%/Judge分≥0.80）。降低延迟换小模型→质量抽样下降不算改善。
   > 摘录：「目标表达平台希望优化什么，护栏表达哪些边界不能被优化动作破坏。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch42-slo.md]

4. **五场景等待语义**：交互式分析→首响应；长任务报告→进度可查/检查点恢复；后台批处理→截止时间+成本；高风险审批→权限/审计/人工确认优先于延迟。
   > 摘录：「SLO 设计顺序：先判断用户如何等待，再决定系统如何承诺。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch42-slo.md]

5. **错误预算驱动发布节奏**：周10万次任务×98% SLO=2000次失败预算。周三已消耗1800次→暂停高风险发布。不同失败消耗不同预算（普通超时1次/越权泄漏按安全门禁单独处理）。
   > 摘录：「错误预算被耗尽后，团队要暂停高风险发布、收紧流量或增加人工复核，而非继续叠加新功能。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch42-slo.md]

6. **降级要有用户语言**："完整报告预计12分钟后完成，先返回关键指标摘要" > "服务繁忙"。用户理解系统做了什么取舍、任务没丢、后续在哪查看，才不会用刷新/重试放大压力。
   > 摘录：「降级也要有用户语言。系统可以说'完整报告预计 12 分钟后完成，先返回关键指标摘要'。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch42-slo.md]

7. **SLI/SLO/SLA 边界**：SLI=实际测量值（"过去7天成功率98.6%"），SLO=内部承诺（"不低于98%"），SLA=对外正式承诺（可能涉及补偿/合同）。先采准 SLI，用 SLO 驱动治理，稳定后再包装 SLA。[INFERRED]

8. **复盘不能只写"容量不足"**：背后可能是模型路由错误、重试放大、任务无法去重、缓存命中率低、SQL无限制流、GPU扩容慢、用户看不到进度反复提交。复盘沿 run_id 拆真实路径。
   > 摘录：「复盘要沿 run_id 拆解真实路径，把每一步的等待、失败和重试数量列出来。」[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part07-observability-eval/ch/ch42-slo.md]

---

## 五、概念边界

| 概念对 | 区分 | 混淆后果 |
|--------|------|---------|
| Session vs Run | Session=连续会话(多Turn)；Run=一次任务执行 | 把多轮会话塞进一个Run → 执行状态管理混乱 |
| Trace vs Log | Trace=结构化时间线(Step/Span/Event)；Log=文本流 | 把Trace当普通日志→无法下钻和分步诊断 |
| Context Package vs Session History | Context Package=模型当时实际看见的输入包；History=完整对话记录 | 排查多轮错误时看错对象→找不到根因 |
| Artifact vs Trace Payload | Artifact=业务产物引用(hash/ID/存储位置)；Payload=大文件正文 | 把大文件正文塞进Trace→观测系统变大对象存储 |
| Checkpoint vs Memory | Checkpoint=中断恢复点(临时)；Memory=偏好/事实(长期) | Checkpoint当长期记忆→清理困难/权限混乱 |
| 离线评测 vs 在线评测 | 离线=固定样本/固定版本/固定脚本(回归测试)；在线=真实流量观测(发现未知) | 只用离线→漏掉新问法/新数据；只用在线→无法防止已知问题回归 |
| LLM-as-Judge vs 确定性脚本 | Judge=语义/论证/表达质量；脚本=执行/数值/格式/权限 | Judge判SQL执行结果→不可靠；脚本判报告质量→规则膨胀 |
| 成本归因 vs 成本降价 | 归因=把费用贴回run_id/step_id/租户；降价=换模型/砍长度 | 先降价→隐性风险转移给用户和人工复核；先归因→知道该修哪里 |
| 语义缓存 vs 结果缓存 | 语义缓存=相似度匹配复用；结果缓存=确定查询复用 | 语义缓存命中错误→旧答案发给新问题；结果缓存跨权限→安全事故 |
| SLO vs SLA | SLO=内部承诺(可度量/可告警/可复盘)；SLA=对外正式承诺(补偿/合同) | 过早承诺SLA→事故变组织问题 |
| 目标(objectives) vs 护栏(guardrails) | 目标=要优化的；护栏=不能被牺牲的边界 | 降延迟换小模型但安全违规→不算改善 |

---

## 六、平台架构/工作流

### 6.1 AgentOps 持续改进闭环 (ch38→ch39→ch40 贯通)

```
生产流量 → Trace采集(run_id/trace_id/step_id/context_package) 
    │
    ▼
样本池：失败Run + 超时Run + 高成本Run + 用户点踩 + 人工接管 + 审批驳回
    │
    ▼
聚类(失败类型标签) → 沉淀为 benchmark 样本(三层脱敏)
    │
    ▼
离线评测(ch39)：Regression Set + Safety Set → 修复(Prompt/Tool/Semantic Layer/Policy)
    │
    ▼
灰度发布(ch40)：在线灰度观察(领先指标：点踩率/重新生成率/延迟/工具错误率)
    │
    ▼
全量上线 → 继续观察线上指标 → 回到样本池
```

### 6.2 质量校准三角 (ch40)

```
         离线 Benchmark (稳定尺子)
              │
    ┌─────────┼─────────┐
    ▼         ▼         ▼
在线评测    人工复核   模型裁判
(真实分布)  (高可信标尺) (规模覆盖)
    │         │         │
    └─────────┼─────────┘
              ▼
        质量资产回流：
        · 错误回答 → query+context+期望答案+错误类型+修复说明
        · 评估误判 → 裁判失败样本
        · 人工争议 → rubric更新
```

### 6.3 成本治理工作流 (ch41)

```
Trace 归因(run_id/step_id) 
    → 按租户/Agent/任务类型/Step 分摊成本
    → 诊断：哪个环节最贵？失败成本占比？缓存命中率？
    → 决策：
        ├─ 模型路由(弱→强模型按 Utility 函数)
        ├─ 缓存优化(五类缓存，权限隔离键)
        ├─ 预算控制(按租户/Agent设上限)
        └─ 重试策略(限制步数/退避)
    → 回归评测(ch39) + 灰度验证(ch40)
```

### 6.4 SLO 治理闭环 (ch42)

```
场景分类(交互式/长任务/批处理/审批)
    → SLO 定义(objectives + guardrails)
    → 指标展开(入口层/Run层/Step层/质量层/风险层/成本层)
    → 错误预算计算与消耗
    → 保护动作(限流/熔断/降级/排队/异步)
    → 恢复清理(任务状态/产物/用户通知一致性)
    → 复盘(沿run_id拆路径)
```

---

## 七、方案权衡

| 场景 | 方案A | 方案B | 权衡依据 |
|------|-------|-------|---------|
| 评测粒度 | 只看最终答案 | 同时看结果+语义上下文+轨迹 | A 简单但放过碰巧答对/越权取数/路径依赖；B 成本高但可分步修复 |
| Judge覆盖 | 每个Run都跑强Judge | 分层评测(规则全量/轻量Judge大比例/强Judge高风险) | A 成本高但全面；B 平衡覆盖与成本 |
| 成本控制 | 换便宜模型 | 先归因→改路由/加缓存/限重试 | A 可能让SQL错误增多/人工压力上升；B 根治但需要Trace基础 |
| 缓存策略 | 全开缓存 | 权限隔离+版本绑定缓存 | A 命中率高但有安全风险；B 安全但命中率可能降低 |
| SLO 粒度 | 接口级(HTTP成功率) | Run级(任务完成率+质量+安全+成本) | A 简单但漏掉大量真实失败；B 复杂但指导事故处理 |
| 高峰期处理 | 按入口请求限流 | 识别任务指纹(重复请求返回已有run_id) | A 可能制造更多重试；B 保护已开始任务 |
| 模型部署 | 全部用外部API | 自建小模型+外部强模型混合 | A 简单但token费用高；B 降低高频任务成本但增加维护复杂度 |
| benchmark规模 | 5万条自动生成题 | 500条高质量可回放可归因样本 | A 数量多但无业务约束；B 少量但真实覆盖风险场景 |

---

## 八、失败模式

| 失败模式 | 表现 | 根因 | 修复方向 | 来源章节 |
|----------|------|------|---------|---------|
| 上下文错误 | Summary 漏约束、Memory 注入错误、Artifact 参数缺失 | 上下文打包策略不当 | 调整上下文摘要策略、修复 Memory 写入逻辑 | ch38 |
| Schema Linking 失败 | 选错表/字段/Metric 版本 | 语义层字段说明不清晰 | 补指标字典、改进 Glossary 和 Linker | ch38 |
| 工具选择失败 | 调错工具或漏掉必要工具 | 工具描述不清、Planner 约束不够 | 改工具描述、增加 Planner 约束 | ch38 |
| "碰巧答对" | 错误口径但数字正确、绕开权限但结果正确 | 只看最终答案的评测盲区 | 增加语义上下文评测和轨迹评测 | ch39 |
| 路径依赖 | 每轮沿用上一轮SQL、默认使用老报表 | 缺少口径确认和反馈校验 | 增加指标确认步骤、引入 observation 回灌 | ch39 |
| Goodhart效应 | Judge分上升但用户追问率也上升 | 为评测器优化而非为质量优化 | Judge分数和真实行为一起看(追问率/采纳率/退回率) | ch40 |
| 缓存越权 | 低权限用户拿到高权限用户缓存结果 | 缓存键缺权限隔离 | 缓存键绑定 tenant_id + role + policy_version | ch41 |
| 重试放大 | SQL失败后反复自动重试、用户看不到进度反复点击 | 缺少重试上限和进度可见性 | 设置步数限制、前端展示进度、任务去重 | ch41/ch42 |
| 接口200但业务不可用 | API返回成功但用户拿不到结果 | SLO 对象错绑到接口而非 Run | SLO 绑定 run_id、展开任务链路指标 | ch42 |
| 降级绕过安全 | 高峰期为提速省略权限/脱敏/审批 | 降级策略无护栏 | 护栏作为门禁项(安全违规率=0)、目标/护栏分离 | ch42 |
| 容量不足（伪根因） | 复盘只写"容量不足" | 可能是路由错误/重试放大/缓存命中低/无去重 | 沿 run_id 拆真实路径，列出每步等待/失败/重试 | ch42 |

---

## 九、实践清单

### 可观测性 (ch38)
- [ ] 建立 Session/Run/Trace/Checkpoint/Artifact 五对象存储模型，生命周期分开管理
- [ ] 每次 Run 记录最小证据：run_id、Context Package、Model Call、Tool Call、State Event、Artifact 写入
- [ ] Trace 默认保存摘要/hash/版本/引用，原文按权限访问
- [ ] 前端时间线少而稳定，后台 Trace 细粒度，两者通过统一事件 ID 关联
- [ ] Context Package 记录每个 item 的来源、是否进入模型、进入方式、摘要版本、token 估算
- [ ] 不长期保存模型隐式推理（合规风险），只保存可审计决策摘要
- [ ] 多租户平台给每类对象绑定租户和权限（运维看结构/业务看产物/审计看脱敏证据链）

### 离线评测 (ch39)
- [ ] Benchmark 包含四要素：任务定义、可复现数据集、统一评测流程、可解释指标
- [ ] 评测覆盖三层：结果(确定性脚本) + 语义上下文(指标口径/schema校验) + 轨迹(动作序列)
- [ ] 任务空间按三维设计：任务意图 × 执行复杂度 × 企业约束
- [ ] 权重由任务类型决定：财务报告→安全+口径权重高；探索分析→可容忍格式不完美
- [ ] 失败标签可修复：每个失败 case 指向具体工程动作（补字典/改Prompt/加状态机约束）
- [ ] Benchmark 像数据产品一样运营：版本/负责人/覆盖率/失效处理流程

### 在线评测 (ch40)
- [ ] 建立离线→在线→人工复核闭环，三者通过 run_id/trace_id 关联
- [ ] 反馈信号分三层采集：显式(点赞/点踩/评论) + 隐式(复制/追问/重试/放弃) + 业务(采纳率/工单降级)
- [ ] LLM-as-Judge 用于语义/论证/表达质量，不用于指标口径和权限边界
- [ ] 裁判输入包含结构化证据（SQL结果/ArtifactRef/EvidenceRef），按 rubric 输出多维评分
- [ ] 混合采样：固定比例随机 + 全部安全命中 + 低分 + 灰度 + 用户显式反馈
- [ ] 分层评测策略：规则全量 → 轻量Judge大比例 → 强Judge高风险 → 人工抽检校准
- [ ] Judge 版本化管理（rubric版本/Judge模型版本/样本来源/上线策略），偏差通过黄金样本和双裁判监控

### 成本治理 (ch41)
- [ ] 成本事件贴到 run_id/step_id/trace_id，关联 prompt_version/tool_version/semantic_layer_version
- [ ] 优先看单位任务成本（成功Run平均/失败成本占比/P95成本），不只看总金额
- [ ] 模型路由按 Utility(task) = α·Quality - β·Cost - γ·Latency - δ·Risk
- [ ] 缓存键绑定权限：tenant_id + prompt_version + tool_version + semantic_layer_version + policy_version + role + data_snapshot
- [ ] 供应商适配层隔离缓存策略差异：Runtime 生成逻辑上下文包 → 网关按供应商转换
- [ ] 自建模型关注算力/显存/吞吐/量化；外部API关注稳定前缀/缓存命中监控

### SLO 与韧性 (ch42)
- [ ] SLO 从接口级迁移到 Run 级：首响应 + 最终可交付 + 质量可接受 + 安全不让步 + 成本可控
- [ ] 目标(objectives)和护栏(guardrails)分开写，护栏是门禁项
- [ ] 延迟拆到链路阶段：入口受理→首响应 / 首响应→最终答案 / 下游各调用内部耗时 / 排队等待
- [ ] 成功率分层：入口层 / Run层 / Step层 / 质量层 / 风险层 / 成本层
- [ ] 错误预算绑定任务风险：低风险可更多降级，高风险宁愿延迟不牺牲权限
- [ ] 限流优先保护已开始任务（识别任务指纹，重复请求返回已有 run_id）
- [ ] 降级用用户语言（说明任务没丢/做了什么取舍/后续在哪查看）
- [ ] 复盘沿 run_id 拆真实路径，列出每步等待/失败/重试数量
- [ ] 恢复后清理：任务状态/产物状态/用户通知一致性

---

## 十、跨 Part 关系

| 关系 | 说明 | 置信度 |
|------|------|--------|
| Part VII → Part I（总论）、Part II（模型）、Part III（数据）、Part IV（知识） | ch38 的 Session/Run/Trace 对象模型是所有 Part 中 Agent 运行的基础抽象；Context Package 是 Part III(上下文工程) 的消费端 | INFERRED |
| Part VII → Part X（安全与权限） | ch42 的安全不让步 SLO、ch38 的权限过滤 Trace、ch41 的缓存权限隔离，都依赖 Part V 的权限策略和脱敏体系 | INFERRED |
| Part VII → Part V ch28（多 Agent 协作） | 多 Agent 编排场景下，Trace 需跨 Agent 关联(run_id→parent_run_id)，SLO 需按编排链路展开；ch38-42 的框架可直接扩展 | INFERRED |
| ch38 ↔ ch39 | Trace 是 benchmark 构建的输入（回放验证轨迹），benchmark 是 Trace 的消费端（评测 Trace 质量） | EXTRACTED |
| ch39 ↔ ch40 | 离线 benchmark 提供稳定基线，在线评测发现长尾，人工复核校准裁判——三者互相回流形成质量校准闭环 | EXTRACTED |
| ch40 ↔ ch41 | 在线评测的模型裁判成本是成本归因的一部分；分层评测策略(规则/轻量Judge/强Judge)本身是成本优化实践 | INFERRED |
| ch41 ↔ ch42 | 成本归因依赖 ch38 Trace 字段；SLO 中"成本可控"目标直接对应 ch41 单位任务成本指标 | EXTRACTED |
| ch38 ↔ ch42 | SLO 绑定 run_id 依赖 Trace 体系；Trace 的 Step/Span 时间是 SLO 延迟拆分的底座 | EXTRACTED |

---

## 十一、开放问题

1. **Trace 存储成本与保留策略的平衡点**：ch38 提出三层脱敏（公开任务结构/内部脱敏工具结果/需审批原始证据），但具体保留周期、存储分层（热/温/冷）未展开。[INFERRED]

2. **Benchmark 自动演化机制**：ch39 要求 benchmark 像数据产品运营，但线上新失败模式如何自动进入回归集（vs 人工筛选）未给出自动化方案。[INFERRED]

3. **Judge 跨领域泛化**：ch40 的模型裁判按 Rubric 打分，但不同业务域（财务 vs 客服 vs 供应链）的 Rubric 是否可复用，还是每个域单独维护一套 Judge，未讨论。[INFERRED]

4. **缓存命中率与质量的量化关系**：ch41 区分五类缓存但未给出"某类缓存命中率低于X%时需要检查"的阈值经验值。[INFERRED]

5. **多模型供应商的 SLO 聚合**：ch42 的 SLO 假设统一后端，当 Agent 同时依赖多个模型供应商（各有各的 SLA）时，Run 级 SLO 如何折算未展开。[INFERRED]

6. **Agent 下线标准量化**：ch42 提到"长期低使用率、高成本、高失败率的 Agent 应进入整改或下线"，但未给出具体阈值（如连续N周采纳率<X%）。[INFERRED]

7. **Trace 与 Privacy/Security 审计的集成**：ch38 提到审计人员可按流程查看脱敏证据链，但与 GDPR/数据本地化等合规框架的具体对接方式未展开。[INFERRED]

---

## 十二、suggested_tags

```
eap/part7, eap/observability, eap/trace, eap/evaluation, eap/benchmark, 
eap/llm-as-judge, eap/online-eval, eap/cost-governance, eap/cache, 
eap/model-routing, eap/slo, eap/resilience, eap/error-budget, 
eap/dataagent, eap/agentops, eap/context-package, eap/session-run-trace,
eap/quality-calibration, eap/feedback-signal, eap/degradation
```

---

## 十三、source_summary

Part VII 共 5 章，围绕「Agent 平台从演示到生产」构建完整的质量保障与治理框架：

- **ch38(可观测性)**：建立 Session/Run/Trace/Checkpoint/Artifact 五对象边界模型，定义 Run 采集点最小证据表（8个采集点×6列），提供 8 类失败诊断路径(上下文错误→意图理解→Schema Linking→工具选择→参数→下游→权限→成本)。Context Package 被定位为多轮排障的核心对象。AgentOps 将 Trace 从"日志"升级为"工程资产"。

- **ch39(离线评测)**：从 ML 单点评测→LLM 通用评测→Agent 多步评测→Text-to-SQL→企业 DataAgent 的完整演进。提出评测三要素（结果/语义上下文/轨迹），综合评分公式 Score_agent = Score_quality - w_cost·CostPenalty，任务空间三维矩阵(意图×复杂度×企业约束)，以及将失败标签映射为可修复工程动作的设计原则。

- **ch40(在线评测与 LLM-as-Judge)**：建立离线/在线/人工复核质量校准闭环。三类反馈信号（显式/隐式/业务结果）→ FeedbackSignal 公式。LLM-as-Judge 三种模式中多维评分为企业 DataAgent 推荐方案。裁判输入需包含结构化证据，输出需按 rubric 给维度分+失败证据。防止 Goodhart 效应的关键是 Judge 分与真实行为（追问率/采纳率/退回率）一起看。

- **ch41(成本治理)**：Cost_run = ΣCost_model + ΣCost_tool + Cost_retrieval + Cost_storage + Cost_eval + Cost_retry。先归因再降价，单位任务成本比总金额更有工程决策价值。模型路由用 Utility(task) 函数平衡质量/成本/延迟/风险。五类缓存（提示词/前缀KV/结果/语义）边界区分，缓存键第一层是权限隔离。自建 vs 外部API 策略差异通过供应商适配层隔离。

- **ch42(SLO 与韧性)**：Agent SLO 从 HTTP 接口级迁移到 Run 任务级，覆盖首响应及时/最终可交付/质量可接受/安全不让步/成本可控。五场景等待语义(交互式/长任务/批处理/审批/报表)。目标与护栏分离、SLI/SLO/SLA 边界区分。错误预算驱动发布节奏，降级用用户语言，复盘沿 run_id 拆真实路径。

**全 Part 贯穿**：run_id / trace_id / step_id / Context Package 作为治理主线，ch38 定义采集 → ch39 构建离线基线 → ch40 覆盖在线真实分布 → ch41 归因成本 → ch42 约束稳定性，五环相扣形成 Agent 平台生产化质量闭环。

---

*本总结基于 ax --md --all + web_extract(char_limit=80000) 双路径读取的原始章节内容。所有 [EXTRACTED] 标记的摘录均为原文精确子串（≤50字或截断保留前缀）。[INFERRED] 为从原文推理的结论。*
