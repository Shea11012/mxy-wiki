---
title: "Part II：模型与推理层"
created: 2026-07-18
updated: 2026-07-22
status: snapshot
type: source-summary
tags: [ai/enterprise-agent-platform]
collection: "企业级 Agent 平台工程"
collection_kind: book
collection_order: 2
collection_home: "（总览）"
source: "datagallery-lab/enterprise_agent_platform@e5d97a66c467045413fff692d28fd2c401503029"
sources:
  - "[[sources/ch05]]"
  - "[[sources/ch06]]"
  - "[[sources/ch07]]"
  - "[[sources/ch08]]"
  - "[[sources/ch09]]"
source_kind: third-party
as_of: 2026-07-18
confidence: SYNTHESIZED
contested: false
contradictions: []
related:
  - "（总览）"
  - "模型路由与工程化"
  - "从零实现LLM"
---
# Part II：模型与推理层

> 《企业级 Agent 平台工程》Part II 完整工程理解
> 覆盖：第5章—第9章
> 标注规范：[EXTRACTED: 章URL] 为原文事实；[INFERRED] 为跨章推导

---

## 1) 章节覆盖表

| 章 | URL | 核心命题 | 关键证据（原文 ≤50 字子串） |
|---|---|---|---|
| **第5章 大模型选型** | [ch05](https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch05.md) | 模型选型应从业务任务出发，构建可治理的模型矩阵，而非选"最强模型" | "选型的目标是为不同任务匹配合适模型，并让这套匹配能随业务和模型迭代调整" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch05.md]; "一个常见失败案例是把复杂 DataAgent 和简单分类任务都路由到同一个强模型" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch05.md]; "把模型视为可治理资源，把能力、成本、合规和生命周期写入平台配置" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch05.md] |
| **第6章 本地推理引擎** | [ch06](https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch06.md) | 本地推理服务的核心矛盾是吞吐与延迟的取舍，先建部署边界再选引擎 | "企业自建推理服务的核心矛盾，是在有限 GPU 上同时满足并发吞吐和单请求延迟，而这两者往往相互拉扯" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch06.md]; "本地推理的早期设计最好把'可运营'放在'极致性能'之前" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch06.md]; "业务应用不直接认识推理引擎，引擎不直接承担租户和合规责任" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch06.md] |
| **第7章 推理优化技术** | [ch07](https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch07.md) | 推理优化必须先定位瓶颈再选机制；KV Cache 是显存和长上下文的主约束 | "优化事故往往来自机制用在了错误瓶颈上" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch07.md]; "推理优化不能按'能开就开'的思路做" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch07.md]; "若只压低延迟和成本，系统可能更快地给出更难解释的结果" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch07.md] |
| **第8章 结构化输出与提示工程** | [ch08](https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch08.md) | 结构化输出是接口治理问题，需 Prompt/schema/工具调用三层契约绑定发布 | "结构化输出属于接口治理问题，而非单纯的格式问题" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch08.md]; "真正可靠的结构化任务，需要把这些资产当作同一个版本发布" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch08.md]; "企业里常见的事故并非 JSON 语法错误，而是语义正确性和执行边界混在一起" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch08.md] |
| **第9章 模型能力定制与知识增强** | [ch09](https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch09.md) | 能力定制先分诊再选路线：Prompt→RAG→微调→对齐，按可回滚程度递进 | "定制路线要服务失败原因，而非服务团队偏好的技术路线" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch09.md]; "先做可解释、可回滚、链路局部的改动" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch09.md]; "把动态知识写进模型参数，会导致更新慢、不可追溯、难删除" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch09.md] |

---

## 2) Part 核心问题

**Part II 回答的根本问题是：企业如何在多业务线、多模型、多部署形态的条件下，构建一个可治理、可观测、可回滚的模型推理层？**

展开为四个子问题：

1. **选什么模型？**（Ch05）——不选"最强"，选"合适"。以任务画像驱动模型矩阵决策，覆盖能力、成本、合规、生命周期。
2. **模型跑在哪、怎么管？**（Ch06）——本地推理不是把模型下载到 GPU 上跑起来就完事，需要五层部署形态递进，以及统一网关/配额/审计的治理边界。
3. **怎么跑更快而不破坏质量？**（Ch07）——优化必须从瓶颈定位开始：KV Cache 管显存，Prefix Cache 降 TTFT，Speculative Decoding 提 Decode，量化省显存——每项都要做质量回归。
4. **输出怎么变成可消费接口？模型能力怎么定制？**（Ch08+Ch09）——结构化输出是接口治理（三层契约绑定发布），模型定制必须按"分诊→Prompt→RAG→微调→对齐"的递进顺序。

[INFERRED] Part II 的核心工程立场是：**推理层是模型能力供给层，不是模型本身**。业务只关心能力/SLO，不关心底层跑的是 vLLM 还是 SGLang。

---

## 3) 章节递进关系

```
Ch05（选型）
  │  定义"用什么模型"——任务画像→模型矩阵→路由策略
  │  产出：Model Catalog + Capability Profiler + Policy Engine
  ▼
Ch06（部署引擎）
  │  定义"模型跑在哪"——五种部署形态，吞吐vs延迟的工程取舍
  │  产出：可运营的本地推理服务池（通过Gateway接入）
  ▼
Ch07（推理优化）
  │  定义"怎么跑更快"——瓶颈定位优先，四类优化机制的作用位置
  │  产出：KV Cache管理 / Prefix Cache / Speculative Decoding / 量化
  ▼
Ch08（结构化输出）
  │  定义"输出变成接口"——Prompt/schema/工具调用三层契约绑定
  │  产出：校验闭环（parse→schema→business→recover）
  ▼
Ch09（能力定制）
     定义"模型不够用怎么办"——分诊→RAG→微调→对齐的决策框架
     产出：持续闭环（失败样本→训练/知识更新→评测→灰度→回滚）
```

递进逻辑：[INFERRED] Ch05 提供决策框架（"用什么"）→ Ch06 提供基础设施（"跑在哪"）→ Ch07 提供性能调优（"怎么快"）→ Ch08 提供接口治理（"输出怎么消费"）→ Ch09 提供能力迭代（"怎么变好"）。五章构成从"选择→部署→优化→输出→定制"的完整推理层工程链路。

---

## 4) 关键结论

1. **模型选型是运行时决策，不是一次性采购。** "一家多业务线企业需要的是模型矩阵，不是'一个最佳模型'" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch05.md]。选型需要任务画像、模型注册表、策略引擎、评测系统和可观测性共同构成闭环。

2. **推理服务应先建治理边界再追求性能。** "一个稍慢但有稳定指标、灰度发布和回滚路径的服务，比一个压测数据漂亮却无人能解释排队和显存波动的服务更适合生产" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch06.md]。

3. **优化必须定位瓶颈→选择机制→质量回归三步走。** KV Cache 优化作用在显存与长上下文并发，Prefix Cache 作用在共享前缀的重复 Prefill，Speculative Decoding 作用在 Decode 阶段，量化作用在权重/KV Cache 的存储。组合必须有评测和回滚依据 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch07.md]。

4. **结构化输出不是"输出 JSON"，而是接口治理。** 需要 Prompt 模板（输入契约）+ JSON Schema（输出契约）+ 工具执行（执行契约）+ 版本治理（治理契约）四层绑定发布 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch08.md]。"additionalProperties: false 限制模型输出未定义字段，unknown 给模型一个合法的失败出口" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch08.md]。

5. **模型定制按可回滚程度递进：Prompt → RAG → 微调 → 对齐。** "先做可解释、可回滚、链路局部的改动" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch09.md]。微调适合学习任务模式，RAG 适合接入可更新知识，对齐适合调整拒答偏好，三者不能互相替代。

---

## 5) 核心概念与边界

| 概念 | 定义边界 | 出处 |
|---|---|---|
| **模型矩阵（Model Matrix）** | 多个模型按能力/成本/合规/生命周期组成可路由资源池，由网关统一调度，而非"一个统一模型" | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch05.md] |
| **任务画像（Task Profile）** | 将业务需求翻译为可评测约束：输入类型、输出消费者、失败成本、数据敏感度、延迟目标、部署边界 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch05.md §5.1.3] |
| **三种路由模式** | 显式模型（可复现）、策略路由（规模化，生产默认）、多模型仲裁（高风险复判） | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch05.md §5.2.1] |
| **连续批处理（Continuous Batching）** | 每个 Decode step 动态加入新请求、移除完成请求，而非固定 batch 等齐 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch06.md §6.3] |
| **KV Cache** | 自回归生成时缓存历史 Token 的 Key/Value，规模≈batch×seq_len×layers×2×kv_heads×head_dim×bytes | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch07.md §7.2] |
| **Prefix Caching** | 多请求共享相同前缀时复用已生成的 KV Cache，减少重复 Prefill；收益取决于 prompt 规范化（稳定内容在前缀，动态内容在后缀） | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch07.md §7.3] |
| **Speculative Decoding** | 小模型（draft model）先草拟多个候选 Token，大模型一次性验证接受/拒绝 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch07.md §7.4] |
| **结构化输出四层契约** | 语义契约（Prompt）+ 结构契约（Schema）+ 执行契约（Tool）+ 治理契约（版本/灰度/回滚） | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch08.md §8.1.2] |
| **能力定制闭环** | 线上失败→样本分诊→数据治理→训练/知识更新→评测→灰度→回滚→监控 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch09.md §9.2.1] |

**关键边界区分：**

- **模型选型** ≠ 选购模型品牌。选型是按任务画像匹配能力资源 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch05.md]
- **本地推理** ≠ 下载模型+启动脚本。推理是"可治理资源池"，通过网关、配额、审计接入 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch06.md]
- **推理优化** ≠ 能开就开。每种优化有明确作用位置和质量风险 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch07.md]
- **结构化输出** ≠ 输出 JSON。需要 schema 校验 + 业务规则 + 恢复策略 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch08.md]
- **微调** ≠ 解决所有效果问题。知识过期走 RAG，格式不稳修 Prompt/schema [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch09.md]

---

## 6) 平台架构/工作流

### Part II 的整体架构位置

[INFERRED] Part II 处于平台的**模型能力供给层**，向上对接 Part V（Agent Runtime），向下占用 Part VIII（部署与基础设施）的 GPU/网络资源，左右通过 LLM 网关（LLM Gateway） 连接外部模型 API 和内部业务应用。

```
业务应用 (Part III 等)
        │
   ┌────▼─────────────────────────────────────┐
   │          LLM 网关（LLM Gateway）                      │
   │  (认证/配额/路由/审计/限流/灰度)           │
   └────┬──────────────┬──────────────┬────────┘
        │              │              │
   ┌────▼────┐   ┌────▼────┐   ┌────▼────┐
   │闭源API   │   │云模型平台│   │本地推理  │
   │(Ch05选型)│   │(Bedrock等)│  │(Ch06引擎)│
   └─────────┘   └─────────┘   └────┬────┘
                                    │
                    ┌───────────────┼───────────────┐
                    │               │               │
               ┌────▼────┐   ┌────▼────┐   ┌─────▼─────┐
               │优化层    │   │结构化输出│   │能力定制    │
               │(Ch07)   │   │(Ch08)   │   │(Ch09)     │
               │KV/Prefix│   │契约+校验 │   │RAG/微调/对齐│
               │量化/投机 │   │恢复闭环  │   │灰度+回滚   │
               └─────────┘   └─────────┘   └───────────┘
```

### 关键工作流

**选型到生产的工作流**（Ch05）：
```
任务画像 → 候选模型池 → 企业评测集验证 → 策略引擎(数据边界/租户/风险)
→ 模型选择器(质量/成本/延迟) → LLM 网关（LLM Gateway）(鉴权/限流/路由)
→ 实际调用 → 可观测性(质量/成本/延迟/错误) → 反馈校准策略
```

**推理优化的标准流程**（Ch07）：
```
定位瓶颈(指标拆解：网关排队/Prefill/Decode/TTFT/显存)
→ 选择对应机制(Prefill重→Prefix;Decode慢→Spec;显存紧→KV管理+量化)
→ 灰度验证(性能+质量+结构化输出+拒答+安全拦截)
→ 配置记录+回滚条件 → 全量发布
```

**结构化输出的校验闭环**（Ch08）：
```
RenderPrompt → CallModel → ParseOutput → ValidateSchema
→ ValidateBusiness → Accepted/HumanReview
失败分支: Repair(一次) → Retry(预算内) → HumanReview/Failed
```

**能力定制的持续闭环**（Ch09）：
```
线上失败样本 → 分诊(知识/格式/能力/偏好) → 选择路线
→ 数据治理(清洗/脱敏/去重/版本化)
→ 执行(RAG更新 / SFT/LoRA / DPO/KTO)
→ 评测(任务/安全/回归 gates)
→ 灰度(canary tenants) → 全量发布/回滚 → 监控 → (回到分诊)
```

---

## 7) 方案权衡

| 权衡维度 | 选项A | 选项B | 权衡要点 | 出处 |
|---|---|---|---|---|
| **统一模型 vs 模型矩阵** | 全公司统一一个模型 | 按任务路由不同模型 | A简单但成本风险失衡：低风险任务浪费预算，高风险任务缺乏推理能力 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch05.md] |
| **闭源API vs 自托管** | 调用厂商API | 自建推理服务 | API省运维但数据出域不可控；自托管控制数据边界但需承担GPU/引擎/运维成本 | [INFERRED: ch05+ch06] |
| **吞吐优先 vs 延迟优先** | 增大batch提高吞吐 | 降低batch降低单请求延迟 | 后台摘要追求吞吐，前台助手追求TTFT。"混在同一个服务池里，任何一个场景的优化都可能伤到另一个场景" | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch06.md] |
| **Prompt约束 vs 约束解码 vs 后置校验** | 只靠Prompt控制格式 | 用guided decoding强制schema | Prompt灵活但不可靠；约束解码可靠但可能影响自然语言质量；后置校验兜底但增加延迟和重试成本 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch08.md §8.3.1] |
| **微调 vs RAG vs 对齐** | 训练模型参数 | 注入外部上下文 | 微调学习任务模式(周~月级更新)；RAG接入可更新知识(小时~天级)；对齐调偏好边界 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch09.md §9.1.1] |
| **全量微调 vs LoRA/QLoRA** | 全参数训练 | adapter式微调 | LoRA更小、更快发布、更易回滚、支持多业务域共享基座；但能力上限受基座模型约束 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch09.md §9.3.2] |
| **Prefix Cache命中 vs 上下文膨胀** | 保留长前缀提高命中率 | 压缩上下文控制显存 | 多轮对话历史越长可复用越多，但KV Cache占用越大；需要会话摘要+历史裁剪平衡 | [INFERRED: ch07 §7.3] |

---

## 8) 工程风险/失败模式

| 风险/失败模式 | 触发条件 | 后果 | 出处 |
|---|---|---|---|
| **统一模型风险** | 全公司用同一强模型 | 低风险任务成本过高 + 版本升级导致全业务受影响 [EXTRACTED] | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch05.md] |
| **benchmark 选型陷阱** | 按公开榜单选模型 | "榜单高分模型如果在企业 schema 上频繁输出非法 JSON，也不能直接上线" | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch05.md §5.1.4] |
| **缺少容量控制** | 单机演示后直接上生产 | 客服批量任务占满GPU → DataAgent长上下文TTFT飙升 → 某个租户调大max_tokens导致全局排队 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch06.md 开篇] |
| **指标不可见** | 引擎未接入网关/监控/配额 | "模型团队说模型更好，平台团队看到P99 TTFT超标，业务团队抱怨流式断开" | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch06.md] |
| **优化机制错配** | 不定位瓶颈就开所有优化 | "团队会把所有开关都打开，最后得到一个更快但更难解释的系统" | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch07.md] |
| **量化质量退化** | 上线量化模型未做业务评测 | 数学/代码/长上下文检索/结构化输出/拒答可能退化，"尤其在企业场景需用自己的评测集验证" | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch07.md §7.2] |
| **Prompt不规范导致缓存失效** | 动态字段(trace_id/时间戳)放在前缀 | "同一段系统提示词里如果混入时间戳、随机 trace id 或不稳定字段，缓存命中率会明显下降" | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch07.md §7.3] |
| **JSON格式假象** | 演示阶段只看到JSON渲染成功 | 进入工单/合同/SQL执行器后，字段语义错误、枚举写成自然语言、证据引用不存在才会暴露 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch08.md] |
| **契约分离发布** | Prompt/schema/工具调用分开管理 | Prompt新增字段但schema未更新；schema枚举变了但评测样本仍用旧值 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch08.md] |
| **微调解决知识过期** | 用训练替代知识库更新 | "更新慢、不可追溯、难删除"；旧政策/错误话术/SQL反模式被模型学进去 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch09.md §9.1.4] |
| **评测污染** | 训练数据包含评测样本或近似改写 | 微调后分数虚高，线上泛化下降 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch09.md §9.2.3] |
| **发布不可回滚** | 模型/adapter/Prompt/索引版本未绑定 | 出问题时不知道回答用了哪个版本，无法快速回滚 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch09.md §9.2.3] |
| **过度拒答** | 对齐后正常问题也大量拒答 | "补充安全可答正例，拆分风险等级" | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch09.md §9.2.3] |

---

## 9) 可执行实践

### 选型评审清单（Ch05）
1. 为每个业务场景写清任务画像：输入、输出消费者、失败成本、数据敏感度、延迟目标、部署边界 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch05.md §5.1.3]
2. 从八个维度评估候选模型：任务能力、输出可控性、事实可靠性、延迟与吞吐、成本、数据边界、可运维性、生态兼容 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch05.md §5.1.2]
3. 模型目录至少记录：model_id、provider、deployment、data_boundary、capabilities、limits、slo、eval、release [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch05.md §5.2.2]
4. 网关支持三种路由模式：显式模型（评测/复现）、策略路由（生产默认）、多模型仲裁（高风险复判） [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch05.md §5.2.1]

### 推理服务上线检查（Ch06）
1. 定义五组接口条件：模型边界、请求边界、租户边界、性能边界、数据边界 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch06.md §6.1]
2. TTFT 和 TPOT 拆成网关排队、调度等待、Prefill、首Token、Decode、工具等待、流式传输分别监控 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch07.md §7.1]
3. 应用侧不直接绑定引擎私有参数，"只关心我要哪个能力、最大延迟是多少" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch06.md §6.1]

### 优化上线流程（Ch07）
1. 先定位瓶颈（TTFT高→看排队/Prefill；Decode慢→看模型大小/Spec；显存紧→看KV/量化）
2. 每项优化记录：改善了哪个指标、影响了哪些业务样本、失败时如何回滚 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch07.md]
3. KV Cache 优化先减少无效上下文（检索质量/去重/片段压缩），再讨论压缩技术 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch07.md §7.2]
4. Prefix Cache：稳定内容放前缀（系统角色、安全规则、工具schema、企业术语），动态内容放后缀（用户问题、trace_id、时间） [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch07.md §7.3]

### 结构化输出生产要求（Ch08）
1. schema设计：`additionalProperties: false` + 提供 `unknown` 作为合法失败出口 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch08.md §8.2.2]
2. 失败分层处理：JSON解析失败→repair一次；schema失败→带错误重试；证据缺失→补检索或人工；参数越权→直接拒绝+安全事件 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch08.md §8.2.3]
3. 所有结构化任务资产（Prompt模板/schema/工具契约/模型版本/评测样本/回滚策略）作为同一版本发布 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch08.md §8.1.2]

### 能力定制决策框架（Ch09）
1. 失败样本先分诊（知识/格式/能力/偏好四分类），不直接进训练集 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch09.md §9.1.2]
2. 工程顺序：先修 Prompt + schema + RAG + 评测，再决定是否微调 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch09.md §9.3.1]
3. 微调前回答三问："是否有更低成本修复方式""修复后如何评测""失败时如何回滚" [INFERRED: ch09 §9.3.1]
4. 训练数据与评测数据必须隔离，"数据指纹、去重、评测集隔离" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part02-model-inference/ch/ch09.md §9.2.3]

---

## 10) 与其他 Part 的关系

[INFERRED] 所有以下关系均为跨 Part 推导：

| 与Part X的关系 | 说明 |
|---|---|
| **Part VIII（部署与基础设施）** | Part I 提供 GPU 集群、网络、存储等硬件底座；Part II 在此之上构建推理服务池。Ch06 的五种部署形态直接依赖 Part I 的资源调度能力（Kubernetes/Ray等）。 |
| **Part V（Agent Runtime）** | Part II 是 Part III 的模型能力供给层。Agent Runtime 通过 LLM 网关（LLM Gateway） 调用模型，不直接绑定引擎。Ch05 的任务画像（客服/DataAgent/合同/代码）正是 Part III 各 Agent 场景的模型需求来源。 |
| **Part IV（数据与知识）** | Ch09 的 RAG 管线和知识库快照直接对接 Part IV 的向量库/文档索引。能力定制的知识更新链路依赖 Part IV 的数据治理基础设施。 |
| **Part X（安全、合规与组织）** | Ch05 的数据边界/数据出域、Ch08 的工具权限校验、Ch09 的训练数据脱敏均依赖 Part V 的安全策略。LLM 网关（LLM Gateway） 的审计日志是 Part V 合规审计的输入。 |
| **全书贯穿的 LLM 网关（LLM Gateway）** | LLM 网关（LLM Gateway） 是 Part II 的核心组件，同时贯穿全书：Ch05 用它路由，Ch06 用它隔离引擎，Ch08 用它记录 trace，Ch09 用它灰度发布。 |

---

## 11) 开放问题

1. [INFERRED] **模型版本生命周期治理的自动化程度**：Ch05 提到需要"模型版本 pinning、灰度和回归评测"来应对厂商变更，但未展开如何自动化检测上游模型弃用（deprecation notice）、自动触发回归评测、自动更新路由规则。
2. [INFERRED] **多模型仲裁的裁决逻辑**：Ch05 提到高风险场景用多模型投票/裁决，但未展开当多个模型结论冲突时的裁决机制（多数投票？置信度加权？引入人工？）。
3. [INFERRED] **KV Cache 跨请求共享的安全隔离**：Ch07 §7.3 提到 Prefix Cache 需要按租户/安全域隔离，但未展开具体实现方案（key 设计细节、不同安全等级之间的隔离强度）。
4. [INFERRED] **结构化输出的约束解码与通用能力的平衡**：Ch08 提到约束过强会影响自然语言质量，但未给出不同类型任务（分类 vs 生成式 UI vs DataAgent）的约束强度建议范围。
5. [INFERRED] **LoRA adapter 的组合与冲突**：Ch09 提到同一基座模型服务多个业务域的 adapter，但当多个 adapter 同时作用于一个请求时（如客服助手同时需要分类 adapter + 话术 adapter），如何管理优先级和冲突？
6. [INFERRED] **RAG 与微调的混合策略边界**：Ch09 区分了 RAG（知识）和微调（任务模式），但在长尾场景（如需要同时理解最新政策并按固定格式输出），两者如何协同？谁先谁后？

---

## 12) suggested_tags

```
企业Agent平台, 模型选型, 模型矩阵, 本地推理, 推理引擎,
KV-Cache, Prefix-Cache, Speculative-Decoding, 量化, 推理优化,
结构化输出, 提示工程, 接口契约, JSON-Schema, 工具调用,
模型微调, LoRA, RAG, 偏好对齐, 能力定制,
LLM-Gateway, 任务画像, 连续批处理, 灰度发布, 失败分诊
```

---

## 13) source_summary

本总结覆盖《企业级 Agent 平台工程》Part II（第5—9章）全部内容。五章构成从"选择→部署→优化→输出→定制"的完整模型推理层工程链路：

- **Ch05（大模型选型）**：论述企业模型选型应从一次性采购转变为由任务画像、数据边界、SLO、成本和生命周期共同约束的运行时决策。核心产出是模型矩阵 + 三种路由模式 + 七组件治理体系（Model Catalog / Capability Profiler / Policy Engine / Model Selector / Provider Adapter / Release Controller / Cost & Quality Monitor）。
- **Ch06（本地推理引擎）**：分析本地推理的五种部署形态（从 Ollama 单机到分布式集群），论证吞吐与延迟的核心矛盾，提出推理服务的治理边界应优先于极致性能。展开连续批处理、KV Cache、PagedAttention、Prefix Caching 等引擎底座概念。
- **Ch07（推理优化技术）**：建立"瓶颈定位→选择机制→质量回归"三步法。深度分析四类优化机制的作用位置和质量风险：KV Cache 管理（显存&长上下文）、Prefix Cache（降低重复 Prefill 的 TTFT）、Speculative Decoding（加速 Decode）、量化（节省显存与带宽）。
- **Ch08（结构化输出与提示工程）**：将结构化输出定义为接口治理问题，提出 Prompt（输入契约）、Schema（输出契约）、工具调用（执行契约）、版本治理（治理契约）四层绑定发布模型。给出完整的状态机（parse→schema→business→recover）和失败分层处理策略。
- **Ch09（模型能力定制）**：建立"分诊→Prompt→RAG→微调→对齐"的能力定制递进框架，强调"先做可解释、可回滚、链路局部的改动"。展开能力定制持续闭环（样本→训练/知识更新→评测→灰度→回滚），并给出 LoRA/QLoRA 在企业场景的治理优势。

Part II 的核心工程立场：模型推理层是"可治理的模型能力供给层"，业务应用通过 LLM 网关（LLM Gateway） 消费能力/SLO，不直接绑定具体引擎或模型——这是全书架构中承上（基础设施）启下（Agent 运行时）的关键抽象。
