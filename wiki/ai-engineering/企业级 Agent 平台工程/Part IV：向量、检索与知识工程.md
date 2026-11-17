---
title: "Part IV：向量、检索与知识工程"
created: 2026-07-18
updated: 2026-07-22
status: snapshot
type: source-summary
tags: [ai/enterprise-agent-platform]
collection: "企业级 Agent 平台工程"
collection_kind: book
collection_order: 4
collection_home: "（总览）"
source: "datagallery-lab/enterprise_agent_platform@e5d97a66c467045413fff692d28fd2c401503029"
sources:
  - "[[sources/ch16]]"
  - "[[sources/ch17]]"
  - "[[sources/ch18]]"
  - "[[sources/ch19-ocr]]"
  - "[[sources/ch20-rag]]"
  - "[[sources/ch21]]"
source_kind: third-party
as_of: 2026-07-18
confidence: SYNTHESIZED
contested: false
contradictions: []
related:
  - "（总览）"
  - "RAG检索增强生成全链路"
  - "RAG检索增强生成"
---
# Part IV：向量、检索与知识工程

> 来源：企业级 Agent 平台工程，Part IV（第16-21章）
> 置信度：EXTRACTED（带 [EXTRACTED: URL] 标记）与 INFERRED（标记为 [INFERRED]）

---

## 一、章节覆盖表

| 章节 | 标题 | 核心主题 | 覆盖状态 |
|------|------|----------|----------|
| 16 | 嵌入模型 | 向量表示、相似度度量、模型选型、多模态嵌入、评估框架、质量运营、资产复审、退化监控、跨链路验收、向量空间迁移、上线判断 | ✅ **完整（v3 补全 §16.4-16.14）** |
| 17 | 嵌入微调与重排 | 领域语义适配、对比学习、hard negative、微调路线、reranker、灰度发布、样本合规、停止条件、跨团队复盘、检索复测 | ✅ **完整（v2 补全 §17.6-17.20）** |
| 18 | 向量数据库与索引算法 | 向量库平台定位、ANN 索引谱系(HNSW/IVF/PQ/DiskANN)、选型、metadata 过滤、多租户、索引生命周期、运行回放、容量成本、删除一致性、归因分层、索引退役 | ✅ **完整（v2 补全 §18.6-18.20）** |
| 19 | 文档解析与多模态 OCR | 解析失败模式、版面结构、工具链选型、OCR/VLM 边界、质量门禁、回放边界、分层标准、人工复核、事故复盘、分离验收 | ✅ **完整（v2 补全 §19.6-19.18）** |
| 20 | RAG 工程与高级检索 | RAG 六层工程体系、chunk 策略、混合检索、查询理解、多跳检索、可信回答、证据链、诊断路径、运营指标、发布门禁、弱证据降级、冻结窗口 | ✅ **完整（v2 补全 §20.6-20.17）** |
| 21 | 知识工程 | 知识工程定位、本体建模、信息抽取、实体链接、GraphRAG、知识资产治理、发布/回滚、联合验收 | ✅ 完整 |

---

## 二、Part 核心问题

Part IV 回答一个总问题：**企业知识如何被准确召回、验证和组织，使 Agent 的回答有据可查、错了能定位？**

六个子问题构成递进：
1. **向量表示**：什么内容需要变成向量？向量能否解决语义鸿沟？[EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md] "embedding 的价值也容易被夸大。它负责把文本、代码、图片和结构化语义表示为可检索的向量，先找"
2. **语义适配**：通用 embedding 在企业术语、字段别名、长尾表达上失效怎么办？[EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch17.md] "公开模型能召回大部分高频问题，却在内部术语、字段别名、合同条款、工单根因和长尾表达上出错。"
3. **向量存储**：向量库在平台中的职责边界是什么？[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch18.md] — 概括：向量库事故通常不会以'数据库坏了'的形式出现。更常见的是业务用户看到一段无权文档的影子，DataAgent 召回了过期字段。
4. **文档解析**：为什么很多 RAG 失败发生在解析阶段而非模型阶段？[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch19-ocr.md] — 概括：向量库和 reranker 无法修复这些底层噪声。它们只能在已经被解析出来的内容里排序，不能找回被漏掉的表格单元格。
5. **RAG 工程**：如何把 RAG 从 demo 变成证据工程？[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch20-rag.md] — 概括：RAG 容易做出演示，却很难做到企业要求的'答案有据、引用可查、错了能定位'。
6. **知识工程**：何时需要从向量检索升级到知识图谱？[EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch21.md] "RAG 擅长从文档里找证据，但企业里很多问题还需要沿着对象关系继续追问。"

---

## 三、章节递进

```
ch16 嵌入模型          → 语义表示层：把内容变成可检索向量
  ↓
ch17 嵌入微调与重排    → 适配层：让向量空间理解企业术语
  ↓
ch18 向量数据库        → 存储层：带 metadata/权限/版本的索引基础设施
  ↓
ch19 文档解析          → 输入层：将非结构化文档转为可检索对象（前置）
  ↓
ch20 RAG 工程          → 编排层：chunk→召回→重排→引用→答案的证据链路
  ↓
ch21 知识工程          → 推理层：本体+抽取+图谱，处理关系型追问
```

[INFERRED] 递进逻辑：从"如何表示"（ch16）到"如何适配"（ch17）到"如何存储"（ch18）到"如何解析输入"（ch19）到"如何编排检索"（ch20）到"如何推理关系"（ch21）。每层既独立又可组合，前面每层的质量决定后面所有层的上限。

---

## 四、关键结论

1. **embedding 是候选而非决策** [EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md] "一个相似条款不能直接成为法务结论，一个相似工单不能直接成为故障根因，一个相似字段也不能直接授权 SQ"
2. **先归因，后训练** [EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch17.md] "只有当错误稳定、样本可标注、baseline 可复现时才微调；否则先补术语、字段说明和 query "
3. **向量库是知识基础设施而非缓存** [EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch18.md] "缓存错了可以清掉，向量索引错了会影响召回、引用、评测和审计。"
4. **解析是 RAG 的第一道质量门** [EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch19-ocr.md] "解析链路做不好，后续的 embedding、RAG、GraphRAG 和 EvidenceRef 都"
5. **RAG 是证据工程链路** [EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch20-rag.md] "模型没有凭空编造，问题发生在检索链路：chunk 切分、表格解析、召回、重排和引用校验都没有把关键证"
6. **图谱错时 GraphRAG 会更有说服力地错** [EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch21.md] "图谱一旦错连，GraphRAG 会沿着错误关系生成结构完整的错误解释，比普通 RAG 的一段错误引用"
7. **最小上线门槛** [INFERRED] 贯穿全 Part：权限过滤 + 模型版本 + 索引版本 + 失败样例 + 人工复核边界 + 可回放检索日志。
8. **🆕 多模态 embedding 是候选生成器而非解析替代** [EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #164] "更稳的架构是：OCR 和版面解析提供可引用、可校验的结构化内容；多模态 embedding 负责视觉"
9. **🆕 embedding 评估需要五类对象和失败样例** [EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #165] "对企业来说，失败样例往往比平均分更有价值：它能告诉团队问题在模型、chunk、OCR、权限过滤、字段"
10. **🆕 embedding 质量运营依赖样本回流而非一次选型** [EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #166-embedding] "用户点踩、人工纠错、答案引用错误、SQL 字段链接错误，都应转成 query、golden docs"
11. **🆕 向量空间迁移风险常表现为引用偏移而非系统不可用** [EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #1611] "迁移风险常常不会表现为系统不可用，而是表现为引用顺序变了、相近制度被混淆、字段链接更容易选错、权限过"
12. **🆕 检索需区分"先过滤后召回"与"召回后过滤"** [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #166-embedding] — 概括：企业检索链路必须区分'先过滤后召回'和'召回后过滤'的差异。前者更安全但可能牺牲召回，后者更容易暴露敏感候选。
13. **🆕 embedding 资产复审要五维联动** [EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #167-embedding] "Embedding 资产复审要把模型、索引、chunk 策略、metadata、过滤规则和评测集放在"
14. **🆕 跨链路验收覆盖三类任务链路** [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #1610-embedding] — 概括：验收样本应覆盖三类链路。第一类是知识问答链路…第二类是 DataAgent 链路…第三类是权限链路。
15. **🆕 微调需明确停止条件** [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch17.md §17.15] — 概括：停止条件可以从三个方面设定。第一是质量边界：高风险样本没有退化…第二是成本边界：索引重建、重排延迟、训练样本维护和灰度回滚的成本能被业务价值覆盖。第三是治理边界：模型版本、索引版本、样本版本和权限策略都能被记录和回放。
16. **🆕 检索事故要分层归因** [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch18.md §18.12] — 概括：第一层看知识是否存在…第二层看可见性…第三层看召回与排序…第四层看上下文组装…第五层才看生成：模型是否忠实使用了证据。
17. **🆕 弱证据是正式状态而非异常** [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch20-rag.md §20.15] — 概括：企业 RAG 的生产边界应把弱证据视为正式状态，而不是异常日志。降级动作要按证据缺口选择：缺少文档时可请求上传材料或转人工补证；版本过旧时可提示使用旧版结论或等待知识刷新。

---

## 五、概念边界

### 5.1 embedding ≠ 答案判断
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md] — 概括：Embedding 负责候选召回，再用 metadata filter、reranker、语义层、规则校验和人工复核收口。

### 5.2 微调 vs 重排的职责边界
[EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch17.md] "微调影响第一阶段向量召回，重排影响候选排序，权限和引用校验则属于平台控制面。"

### 5.3 向量库职责边界
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch18.md] 表18-1 明确："向量索引不负责生成 embedding；metadata 过滤不替代统一权限系统；近似检索不保证最终答案正确；生命周期治理不替代数据血缘与审计。"

### 5.4 OCR vs 版面解析 vs VLM 的边界
[EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch19-ocr.md] "OCR 把图像中的文字转出来，版面模型识别区域和阅读顺序，VLM 进一步理解页面里的图表、截图和视觉"

### 5.5 RAG vs 语义层 vs 知识图谱 vs GraphRAG 的边界
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch21.md] 表21-1：
- 文档 RAG："制度怎么说""合同条款在哪里" → 不擅长复杂关系、全局聚合、实体消歧
- DataAgent 语义层："这个指标怎么计算""查哪些表" → 不擅长非结构化关系和开放文本证据
- 知识图谱："哪些对象相互影响""关系链是什么" → 不擅长无证据文本的开放生成
- GraphRAG：跨文档、跨实体的综合回答 → 本体和抽取质量差时会放大错误

### 5.6 🆕 RAG vs 工具调用的边界
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch20-rag.md §20.8] — 概括：边界判断可以看三个信号。问题需要实时数据时，优先工具调用；问题需要改变外部状态时，必须工具调用；问题需要引用稳定文档时，优先 RAG。混合场景也很常见…平台要把两类证据分开：文档证据说明规则，工具证据说明事实。

### 5.7 🆕 解析结果四类产物边界
[EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch19-ocr.md §19.8] "正文片段用于检索，结构化字段用于过滤和校验，证据坐标用于回答溯源，风险标签用于权限和脱敏。"

### 5.8 🆕 多模态 embedding ≠ 视觉等价（visual equivalence）
[EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #164] "两张缺陷图相似，不代表根因相同；两页合同版式相似，不代表条款风险相同。多模态 embedding 在"

### 5.9 🆕 "先过滤后召回" vs "召回后过滤"
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #166-embedding] — 概括：企业检索链路必须区分'先过滤后召回'和'召回后过滤'的差异。前者更安全但可能牺牲召回，后者更容易暴露敏感候选。高风险场景应优先保证无权内容不进入模型上下文。

### 5.10 🆕 embedding 评测指标 vs 业务解释的边界
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #1612-embedding] — 概括：召回率、MRR、nDCG、向量维度和相似度分布，对业务用户来说很难直接判断价值。企业上线嵌入模型时，需要把评测结果翻译成业务解释。

---

## 六、平台架构/工作流

### 6.1 整体检索链路
```
原始文档 → 文档解析(ch19) → chunk + metadata → embedding(ch16) → 向量库(ch18)
                                                        ↓
用户问题 → 查询理解(ch20) → 候选召回 → 重排(ch17) → 上下文组装(ch20) → 答案生成 + 引用校验(ch20)
                                                        ↓
知识图谱(ch21) ← 信息抽取 ← 实体链接 ← GraphRAG ← 关系推理 + 文本证据
```

### 6.2 embedding 能力链路
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md] 图16-1：文档/查询输入 → 嵌入模型编码 → 向量入库 → 相似度检索 → 结果返回。

### 6.3 RAG 六层工程体系
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch20-rag.md] 图20-1：离线侧（解析→分块→嵌入→入库），在线侧（查询改写→混合检索→重排→引用校验→生成），通过向量库衔接。

### 6.4 知识工程技术栈
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch21.md] 图21-1：数据源 → 信息抽取 → 本体/实体链接 → 知识图谱存储 → GraphRAG 检索 → 上层应用。

### 6.5 微调发布流程
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch17.md §17.3] 图17-4：新模型离线编码 → 离线评测 → shadow query → 小流量灰度 → 保留旧索引回滚窗口。

### 6.6 🆕 微调生产验证四阶段
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch17.md §17.11] — 概括：第一步是离线回归，同一批 query 在旧/新 embedding、旧/新索引、旧/新 reranker 的组合上分别运行…第二步是 shadow query，把线上真实请求复制到新链路…第三步是小流量灰度，选择低风险知识库或内部用户…第四步才是按知识库、租户或任务类型逐步放量。

### 6.7 🆕 解析产物状态机
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch19-ocr.md §19.10] — 概括：一个文档可能处于 parsed、review_required、reviewed、approved_for_search、approved_for_high_risk_answer、rejected 等状态。普通检索可以使用 approved_for_search 的材料，高风险回答需要更高门槛。

### 6.8 🆕 RAG 弱证据降级状态
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch20-rag.md §20.15] — 概括：早期 RAG 可以把弱证据状态写进回答协议：insufficient_evidence、conflicting_evidence、stale_evidence、permission_blocked、partial_coverage。

### 6.9 🆕 多模态检索双路径架构
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #164] 图16-5："OCR/解析负责生成可引用文本和结构化字段，多模态 embedding 负责生成视觉相似候选。两条路径在证据校验和人工复核处汇合。"

### 6.10 🆕 embedding benchmark 评测数据流
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #165] 图16-7："同一批 query、golden docs、hard negatives 和 metadata filter 同时进入多个模型/索引组合，避免模型 A 用一批题、模型 B 用另一批题的不可比问题。"

### 6.11 🆕 embedding 变更跨链路验收框架
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #1610-embedding] — 概括：第一类是知识问答链路，检查 query、候选文档、引用段落、回答和用户可见证据。第二类是 DataAgent 链路，检查自然语言问题、字段候选、指标候选、语义层版本和最终 SQL。第三类是权限链路。

### 6.12 🆕 向量空间迁移发布流程
[EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #1611] "发布前跑固定样本，发布中保留旧索引，发布后观察失败反馈，稳定后再清理旧版本。"

---

## 七、方案权衡

### 7.1 文本 embedding 模型路线
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md] 表16-2：
- 商业 API（OpenAI/Cohere/Voyage）：接入快、稳定性好，适合 baseline → 代价：数据出域、供应商锁定
- 开源通用（BGE-M3/E5/GTE）：可私有化、可控性强 → 代价：需推理服务、资源调度、版本治理
- 国产生态（Qwen3 Embedding）：国产模型链路协同 → 代价：版本更新节奏、生态成熟度
- 行业专用（金融/医疗/法务）：垂域术语提升 → 代价：迁移成本、透明度、泛化能力需验证

### 7.2 embedding 微调路线
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch17.md] 表17-5：
- 不微调，只补 query 改写/术语表/metadata：风险低、上线快 → 默认第一阶段
- 用 sentence-transformers 做小规模对比学习：工程生态成熟 → 作为实验路线
- 用 LoRA/adapter：训练成本可控 → 部署复杂度高，作为后续扩展
- 不改 embedding，引入 reranker：回滚简单 → mini-platform 优先实现 reranker 插槽

### 7.3 ANN 索引算法
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch18.md] 表18-2：
- HNSW：召回和延迟稳定，内存占用较高
- IVF：大规模可控，参数不当会漏召回
- PQ/SQ 量化：节省内存，分数精度下降需重排补偿
- DiskANN：降低内存压力，延迟抖动和冷热数据敏感

### 7.4 向量库选型
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch18.md] 表18-3：
- pgvector：PostgreSQL 结合紧密 → 中小规模默认 baseline
- Milvus：大规模分布式 → 多业务共享平台候选
- Qdrant：payload filtering 友好 → 多租户 RAG 候选
- Weaviate：schema/模块化完整 → 产品化候选
- Vespa：搜索推荐平台 → 高级搜索候选

### 7.5 chunk 策略
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch20-rag.md] 表20-2：
- 固定长度：简单 baseline，易切断语义 → 仅作 baseline
- 结构化 chunk：保留标题/段落/表格/页码 → 默认策略
- Parent-child：小 chunk 召回，大 parent 上下文 → 高价值知识库
- Small-to-big：精确引用 + 扩展邻近上下文 → 作为高级策略

### 7.6 本体建模路线
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch21.md] 表21-4：
- 文档驱动抽取：启动快 → 早期探索，冷启动输入
- 业务对象驱动本体：结构稳定，便于治理 → 默认路线
- 关系优先建模：适合依赖分析和影响分析 → 场景扩展
- OWL/RDF 标准建模：语义严谨 → 先调研，不作为默认实现

### 7.7 🆕 检索路线对比
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch20-rag.md] 表20-3：
- 纯向量检索：语义召回强 → 专有名词、编号、字段名可能漏召回
- 纯关键词检索：精确词、编号表现好 → 同义表达和口语化问题召回弱
- 混合检索 + RRF：兼顾语义和关键词，工程解释性较好 → 参数、去重、融合策略需评估
- 混合检索 + reranker：前排证据质量更好 → 延迟和成本增加

### 7.8 🆕 多模态 embedding 场景与发布控制
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #164] 表16-5：
- 质检巡检：图片找相似缺陷 → 图片权限、拍摄规范、误判复核
- 合同票据：页面图像找相似条款 → 页码引用、金额校验、人工复核
- 数据看板截图：截图与指标说明对齐 → 截图脱敏、版本识别、字段映射
- 商品检索：图文联合召回 → 类目过滤、库存价格约束
- 设备运维：找相似设备状态 → 设备权限、时间地点、低质量图片处理

### 7.9 🆕 embedding 评估指标分层
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #165] 表16-7：
- recall@k：正确证据是否进入前 k 个候选 → 适合检索工程师
- MRR/nDCG：排序质量 → 适合评测负责人
- answer citation hit rate：回答是否引用正确证据 → 适合 RAG 负责人、业务 Owner
- p50/p95 latency：查询延迟 → 适合平台负责人
- cost/query：成本 → 适合 CTO
- permission violation rate：是否召回无权内容 → 适合安全/合规负责人

### 7.10 🆕 先过滤 vs 后过滤
[EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #166-embedding] "前者更安全但可能牺牲召回，后者更容易暴露敏感候选。高风险场景应优先保证无权内容不进入模型上下文，再通"

---

## 八、失败模式

### 8.1 embedding 典型事故（非模型分数低造成）
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md]
- 制度问答召回过期文档 → 索引里没有生效时间过滤
- DataAgent 选错字段 → 字段注释缺少业务别名
- 合同检索漏掉关键条款 → PDF 解析把页眉页脚和正文混在一个 chunk
- 客服相似工单越权 → 向量索引没有租户和部门 metadata

### 8.2 微调常见错误
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch17.md]
- 错把文档解析/chunk/权限问题归因给 embedding → "继续训练 embedding 只会掩盖问题"
- 没有 hard negative 就训练 → 模型只学会粗粒度主题
- 新模型直接写入旧索引 → 新旧向量空间混用导致漂移
- 微调提升平均分但高风险场景退化 → 需要按场景退出条件

### 8.3 向量库事故
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch18.md]
- metadata 缺少 department_id 写入共享 collection，post-filter 后无权候选已进入 trace
- 索引未绑定 embedding 模型版本，新旧向量混用，召回偏向历史样例
- 租户删除数据后评测样本和重排缓存未同步清理

### 8.4 文档解析失败
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch19-ocr.md] 表19-1：
- 文本顺序错误：双栏 PDF/页眉页脚/脚注进入正文 → chunk 语义混乱
- 表格结构丢失：跨页表格被打散 → DataAgent 找不到指标口径
- OCR 漏识别：印章/手写/低清扫描件 → 证据缺失
- 权限和来源缺失 → 检索越权，审计无法追踪

### 8.5 RAG 失败
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch20-rag.md]
- 制度问答中漏掉附件表里的金额阈值 → chunk 切分、表格解析未把关键证据带回
- 上下文组装丢掉表头 → LLM 在缺证据的地方补全
- 引用覆盖不足 → 答案看似有出处，实际缺少决定性证据

### 8.6 知识图谱失败
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch21.md]
- 图谱错连 → GraphRAG 沿错误关系生成结构化错误解释
- 实体链接误合并 → 不同客户/供应商被错误关联
- 事实缺少证据/置信度 → 图推理变成另一种幻觉
- 权限未从文档传递到节点和边 → 敏感关系暴露

### 8.7 🆕 向量库变更事故
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch18.md §18.15] — 概括：索引算法、embedding 版本、分片策略、metadata 字段、过滤顺序和混合检索权重都会改变候选集合。变更看起来只是基础设施调整，用户感受到的却是'资料找不准'或'引用变了'。

### 8.8 🆕 解析事故的下游放大
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch19-ocr.md §19.11] — 概括：用户看到引用页不对、表格金额错列、合同条款缺失、图表解释与原图不一致时，问题可能已经经过了解析、chunk、索引、检索和生成多个环节。复盘时不能只问'RAG 为什么答错'，要先把证据链拆开。

### 8.9 🆕 RAG 三类失败
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch20-rag.md §20.12] — 概括：检索失败表现为相关材料没有进入候选集；证据失败表现为候选材料存在，但引用无法支撑结论；生成失败表现为证据正确但回答夸大、遗漏限制条件或把多个来源混合成错误结论。三类失败的修复方向不同。

### 8.10 🆕 解析"看起来成功"的隐性错误
[EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch19-ocr.md §19.1] "合同中一个跨页表格如果被拆成两段，付款条件和违约责任可能分别落进不同 chunk；制度文档的页眉如果"

### 8.11 🆕 多模态 embedding 典型事故
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #164] — 概括：第一个是把截图、票据和巡检照片直接进入共享索引，结果把客户姓名、地址、金额、设备编号一起扩散到检索系统。第二个是把视觉相似当成业务等价：两张缺陷图相似，不代表根因相同。
[INFERRED] 没有采集规范（设备编号、拍摄角度、时间地点）的视觉内容进入向量空间后，会成为无法追溯的噪声源。

### 8.12 🆕 资产复审遗漏导致的检索退化
[EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #167-embedding] "很多检索事故并非模型突然变差，而是文档重建时丢了字段、chunk 规则改变了段落边界、权限过滤把关键"

### 8.13 🆕 embedding 退化被平均指标掩盖
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #169-embedding] — 概括：退化往往先出现在业务样本里。新产品、新字段、新政策、新缩写和新文档格式进入系统后，原有向量空间可能仍能召回相似文本，却不能区分新的业务边界。这类退化很难通过平均 recall 发现。

### 8.14 🆕 跨链路变更的隐性影响
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #1610-embedding] — 概括：一次模型替换可能让制度问答更好，却让 DataAgent 字段链接更差；一次 chunk 策略调整可能让长文档引用更完整，却让短问题召回变慢。

### 8.15 🆕 向量空间迁移的静默漂移
[EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #1611] "这样的风险比接口报错更难发现，因为回答仍然流畅，用户需要事后核对才会发现证据偏移。"

---

## 九、实践清单

### 9.1 embedding 上线最低门槛
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md] — 概括：最小上线门槛：权限过滤、模型版本、索引版本、召回评测、失败样例和人工复核边界必须齐全。
- 每条向量记录含 model_name、model_version、dimension、normalized、metric、index_version
- metadata 含 tenant_id、acl、source_version、effective_at
- 检索日志记录 query、模型版本、索引版本、过滤条件、候选列表、reranker 结果、最终引用

### 9.2 微调决策检查
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch17.md] 表17-2：
- 先归因错误类型 → 检索错误分四类（召回不到/排序靠后/相似但不可答/权限过滤不足）
- 正确证据已进 top-k 但排序靠后 → 优先 reranker
- 错误稳定、样本可标注、baseline 可复现 → 才微调
- 必须有 hard negative、失败样例、索引回滚、模型版本和候选日志

### 9.3 向量库选型检查
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch18.md]
- 数据规模、强过滤需求、多租户、事务一致性、索引重建窗口、混合检索兼容
- 每个索引记录 embedding 模型版本、chunk 策略、metadata schema、构建时间、评测结果、回滚窗口
- metadata 最小字段：tenant_id、acl、source_type、source_version、effective_at、index_version

### 9.4 文档解析质量门禁
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch19-ocr.md]
- chunk 能回到页码和 bbox，表格有结构，低置信区域可进入复核队列，解析版本可追踪
- 文档分流：制度手册自动解析；合同票据设质量门禁和人工复核
- VLM 输出只作为候选解释，金额/日期/条款/指标必须回到 OCR+规则+结构化校验

### 9.5 RAG 最小接口
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch20-rag.md]
- `retrieve(query, filters)` → `rerank(query, candidates)` → `assemble_context(candidates, budget)` → `generate_answer(context)` → `verify_citations(answer, context)`
- 可信回答门禁：引用覆盖、引用一致、权限有效、时间有效、冲突证据

### 9.6 知识工程最小上线门槛
[EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch21.md] "每条事实有来源、证据、置信度、抽取器版本、复核状态和生命周期。"
- 本体最小对象：Entity Type、Relation Type、Attribute、Constraint、Evidence
- 事实 JSON 含 evidence、confidence、extractor、review_status
- 发布前影子图运行，事实回滚按 fact_batch 粒度

### 9.7 🆕 微调发布台账
[EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch17.md §17.12] "台账至少记录 embedding 模型、reranker、索引版本、chunk 策略、评测集版本、灰"

### 9.8 🆕 检索事故五层归因法
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch18.md §18.12] — 概括：第一层看知识是否存在：文档是否入库、解析是否成功、chunk 是否保留关键结构。第二层看可见性：租户、部门、密级、时间范围和文档状态是否把候选正确纳入或排除。第三层看召回与排序。第四层看上下文组装。第五层才看生成。

### 9.9 🆕 知识库容量分层
[EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch18.md §18.10] "平台应把内容分成热知识、温知识和冷证据：热知识进入在线索引，温知识可以降低副本或使用较慢索引，冷证据"

### 9.10 🆕 索引发布四类版本
[EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch18.md §18.11] "一个成熟的发布流程至少要记录四类版本：文档解析版本、chunk 策略版本、embedding 模型版"

### 9.11 🆕 解析产物四类分层
[EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch19-ocr.md §19.8] "正文片段用于检索，结构化字段用于过滤和校验，证据坐标用于回答溯源，风险标签用于权限和脱敏。"

### 9.12 🆕 RAG 知识冻结窗口
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch20-rag.md §20.14] — 概括：高风险知识域需要冻结窗口。冻结窗口内，平台限制自动发布新答案，或要求答案进入人工复核，直到文档版本、索引版本和权限状态重新稳定。解冻条件应包括解析完成、索引构建完成、关键样本通过、权限过滤记录完整、人工抽检通过。

### 9.13 🆕 微调样本合规与数据卡
[EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch17.md §17.18] "早期可以为检索微调样本建立数据卡。数据卡记录来源、范围、脱敏方式、owner、保留时间、删除入口和可"

### 9.14 🆕 微调项目停止条件三元组
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch17.md §17.15] — 概括：第一是质量边界：高风险样本没有退化，核心任务样本达到发布阈值，且不可回答样本没有被强行匹配。第二是成本边界：索引重建、重排延迟、训练样本维护和灰度回滚的成本能被业务价值覆盖。第三是治理边界：模型版本、索引版本、样本版本和权限策略都能被记录和回放。三类条件中任意一类不满足，都应暂停发布。

### 9.15 🆕 索引退役影响扫描
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch18.md §18.20] — 概括：一个旧索引可能已经不再服务主路径，却仍被某些低频 Agent、评测样本、历史报告或灰度租户引用。直接删除会让问题在很晚之后才暴露。平台应先扫描调用日志、配置、缓存、评测集和文档链接。

### 9.16 🆕 多模态嵌入上线最低门槛
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #164] — 概括：巡检照片要有设备编号、拍摄角度和时间地点；票据图片要有脱敏规则和金额校验；报表截图要能绑定报表版本和字段口径；合同扫描页要保留页码和原文件引用。没有这些元数据，多模态 embedding 只能找'长得像'的内容，无法支撑责任判断。
[INFERRED] 上线门槛与文本检索对齐：权限 + 引用 + 版本 + 人工复核 + 采集规范。

### 9.17 🆕 embedding 评测五类基本对象
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #165] 表16-6：Query（真实用户问题）、Golden docs（应被召回的 chunk）、Hard negatives（语义相近但不能回答的材料）、Metadata filter（部门/租户/权限/时间）、Judgment（相关/部分相关/不相关；是否支持最终答案）。

### 9.18 🆕 embedding 发布前六项工程检查
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #165] 表16-8：权限过滤、模型版本、索引重建、成本口径、可观测性、人工复核。

### 9.19 🆕 失败样本回流格式
[EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #167-embedding] "样本要记录 query、用户角色、可访问范围、期望文档、错误候选、错误原因、chunk id、文档版"

### 9.20 🆕 embedding 资产生命周期管理
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #167-embedding] — 概括：旧索引在双版本并行结束后要归档或删除，过期文档要从向量库清理，失效 metadata 要重建，长期不用的租户索引要进入成本复审。删除动作不能只由存储团队判断…需要知识库 owner、检索 owner、安全 owner 和业务 owner 共同确认。

### 9.21 🆕 embedding 退化监控六类信号
[EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #169-embedding] "正确证据排名下降，hard negative 排名上升，权限过滤后候选不足，用户频繁改引用，报告复核"

### 9.22 🆕 embedding 变更跨链路验收记录
[EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #1610-embedding] "每次 embedding 变更生成一份跨链路验收记录：模型版本、索引版本、chunk 策略、meta"

### 9.23 🆕 向量空间迁移局部回滚
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #1611] — 概括：平台应允许按知识库、租户或任务类型局部回滚，而不是要求全部回到旧模型。局部回滚需要发布台账记录模型版本、索引版本、chunk 版本和路由规则。

### 9.24 🆕 embedding 选型后业务解释一页纸
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #1613-embedding] — 概括：解释材料要连接评测样本。平台可以说明新模型在哪些查询类型上改善召回…也要说明哪些场景没有改善…若选型只报告整体 nDCG 或 Recall@K，业务 owner 很难判断它是否适合自己的知识域。

---

## 十、跨 Part 关系

| 关联 Part/章节 | 关系 |
|---------------|------|
| Part V ch26（Agentic Workflow） | RAG 和知识图谱检索结果需进入 Planner/SQL 生成器/工具选择等决策链路 |
| Part VI（DataAgent） | DataAgent 是 Part IV 的主要消费者：embedding 用于 schema linking，向量库用于字段检索，语义层 + 知识图谱用于口径和关系推理 |
| ch33 语义层 | [INFERRED] 语义层服务问数和指标口径，知识图谱服务实体关系和多跳证据，两者协同 DataAgent 才能同时回答"指标怎么算"和"异常关联哪些对象" |
| ch38 Trace | 检索链路 Trace 需记录本体版本、事实批次、召回链路、引用证据；🆕 ch19 §19.10 明确"Trace 中应能看到 chunk 来自哪个 parser version、哪次人工复核" |
| ch39 Eval | RAG 评估指标（context precision/recall/faithfulness）+ 知识图谱抽取准确率/链接准确率/冲突率 |
| Part IX（前端与多模态） | 多模态 embedding + VLM 解析与 Part V 的视觉理解能力衔接 |
| 🆕 ch50 安全治理 / ch52 合规审计 | ch17 §17.18 明确"第17章的检索优化会和第50章安全治理、第52章合规审计形成连接，避免把质量提升建立在不可控样本上" |
| 🆕 Ch16 → Ch41 成本治理 | ch16 §16.6 明确"在线查询成本、离线重建成本、reranker 成本、索引存储成本要分开统计。若只看单次 embedding API 价格，团队会低估全量重建、双版本并行和多租户隔离带来的开销。" |
| 🆕 Ch16 → Ch17/18/20 跨链路验证 | ch16 §16.10明确"embedding 变更会影响 RAG、知识库、DataAgent、报告引用和权限过滤，验收不能只停在召回分数。"；§16.14 明确"和第17章微调、第18章向量库和第20章 RAG 连起来" |

---

## 十一、开放问题

1. **embedding 评测集如何持续迭代？** [INFERRED] 公开 benchmark（MTEB）不能替代内部评测，但内部评测集的构建成本和更新节奏尚未有标准方案。
2. **reranker 和 embedding 微调的投入边界？** [INFERRED] 书中给出"正确证据进 top-k 就优先 reranker"，但 top-k 的具体阈值（50? 100?）需要按场景实验确定。
3. **混合检索（BM25 + 向量）的融合参数如何自动调优？** [INFERRED] RRF 的 k 参数和权重配比目前依赖人工经验。
4. **GraphRAG 的生产维护成本？** [INFERRED] 书中强调本体需业务 owner 维护，但实际组织中谁来长期投入是个组织问题。
5. **多跳检索的自动停止条件？** [EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch20-rag.md] "系统不能无限拆问题，也不能把所有中间结果塞进上下文"，但具体停止策略未详细展开。
6. **知识图谱的退役与物理删除边界？** [INFERRED] 书中提出 "retired/inactive" 标记 vs 物理删除的区分，但 GDPR 等法规下的具体实现未展开。
7. **中英文双语知识库的引用一致性？** [INFERRED] ch21 提到多语言版本复盘，但跨语言的 chunk 对齐和证据一致性尚无成熟方案。
8. **🆕 解析器升级后旧 chunk 与旧回答的兼容？** [EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch19-ocr.md §19.14] "若 chunk 边界变化导致原 citation span 失效，索引可以重建，但历史回答需要保留旧"——保留双版本的成本和策略未具体展开。
9. **🆕 弱证据场景中"补证材料"如何避免污染全局索引？** [EXTRACTED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch20-rag.md §20.15] "若补证材料未经过解析和权限检查，只能作为当前任务的人工证据，不能自动成为全局知识。"——但实际操作中如何隔离临时证据与正式知识，未给出工程方案。
10. **🆕 多模态嵌入与OCR/解析的同步更新机制？** [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #164] 图16-5 设定两条互补路径在证据校验处汇合，但 OCR 版本、解析版本、embedding 版本三者不同步时，视觉候选和文本候选如何对齐未详细展开。
11. **🆕 embedding 资产复审的触发频率？** [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #167-embedding] 资产需定期复审（文档、术语、权限变化），但复审周期多久一次、是事件驱动还是时间驱动，未给出具体建议。
12. **🆕 向量空间局部回滚的工程复杂度？** [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #1611] 按知识库/租户/任务类型局部回滚需要路由规则，但多版本索引并存的路由开销和一致性保证未展开。
13. **🆕 embedding 退化监控的自动化阈值？** [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md  #169-embedding] 书中建议"每月抽样一批真实 query"，但 correct evidence rank drop 到多少算退化、hard negative rank rise 多少触发告警，未量化。

---

## 十二、suggested_tags

```yaml
suggested_tags:
  - embedding/企业
  - embedding/选型
  - embedding/微调
  - embedding/停止条件
  - reranker/重排
  - reranker/灰度
  - vector-db/选型
  - vector-db/HNSW
  - vector-db/pgvector
  - vector-db/metadata-filter
  - vector-db/容量治理
  - vector-db/索引退役
  - document-parsing/OCR
  - document-parsing/版面分析
  - document-parsing/质量门禁
  - document-parsing/人工复核
  - RAG/工程体系
  - RAG/chunk策略
  - RAG/混合检索
  - RAG/可信回答
  - RAG/弱证据降级
  - RAG/冻结窗口
  - knowledge-engineering/本体建模
  - knowledge-engineering/信息抽取
  - knowledge-engineering/知识图谱
  - knowledge-engineering/GraphRAG
  - knowledge-engineering/实体链接
  - enterprise-agent-platform/part4
```

---

## 十三、source_summary

本总结基于「企业级 Agent 平台工程」在线书籍 Part IV 全部六章原文完整获取：

| 章节 | URL | 获取方式 | 状态 |
|------|-----|---------|------|
| ch16 嵌入模型 | https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch16.md | ax --md --all + **web_extract (完整 §16.1-16.14)** | ✅ **v3 完整** |
| ch17 嵌入微调与重排 | https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch17.md | ax (截断§17.4) + **web_extract (完整§17.1-17.20)** | ✅ **v2 完整** |
| ch18 向量数据库与索引算法 | https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch18.md | ax (截断§18.5) + **web_extract (完整§18.1-18.20)** | ✅ **v2 完整** |
| ch19 文档解析与多模态 OCR | https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch19-ocr.md | ax (截断§19.5) + **web_extract (完整§19.1-19.18)** | ✅ **v2 完整** |
| ch20 RAG 工程与高级检索 | https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch20-rag.md | ax (截断§20.5) + **web_extract (完整§20.1-20.17)** | ✅ **v2 完整** |
| ch21 知识工程 | https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part04-vector-knowledge/ch/ch21.md | web_extract (完整) | ✅ 完整获取 21.1-21.18 |

**v2 更新说明**：v1 总结中 ch17-20 各章仅覆盖前 4-5 节导言 + 核心架构内容，尾节（约占全书相应章节 60-75% 篇幅）截断。v2 通过 web_extract（char_limit=80000）全量获取所有章节内容，补全了以下尾节：

- **ch17**：补全 §17.6-17.20（Hard Negative 治理、重排上线边界、评测闭环、索引协同、成本边界、生产验证、发布台账、回放样本、向量库迁移复盘、停止条件、发布纪律、跨团队复盘、样本合规、检索复测）——新增 15 节
- **ch18**：补全 §18.6-18.20（benchmark 工程实践、运行回放、生产判断、质量回放、容量成本治理、发布删除一致性、归因分层、质量复盘、目录责任分工、灰度窗口、分片隔离校验、缓存版本证据、线上灰度、健康信号、索引退役扫描）——新增 15 节
- **ch19**：补全 §19.6-19.18（回放边界、OCR 分层标准、入库边界、人工复核样本库、复核运营、事故复盘、OCR 门禁抽检、业务验收样本、策略变更回放、版式回归、OCR 纠错、业务复核、分离验收）——新增 13 节
- **ch20**：补全 §20.6-20.17（证据链运行标准、诊断路径、工具调用边界、运营指标、发布门禁、证据质量台账、失败诊断与回收、用户可见性、冻结窗口、弱证据降级与补证、答案争议复核、用户呈现）——新增 12 节

总计新增 **55 节尾节**内容，从 v1 的约 19KB 扩展到 v2 的约 35KB，覆盖从 "能力怎么建" 到 "上线后怎么运营/复盘/退役" 的完整生命周期链路。

**v3 更新说明**：v2 总结中 ch16 仅覆盖 §16.1-16.3（约 3/14 节），尾节 §16.4-16.14 全部缺失。v3 通过 web_extract（char_limit=80000）全量获取 ch16 全章并补全了以下尾节：

- **ch16**：补全 §16.4-16.14（多模态嵌入与视觉检索、评估框架、质量运营、资产复审与样本回流、评测集业务覆盖、业务退化监控、跨链路验收、向量空间迁移风险复盘、评测结果业务解释、选型后业务解释、上线判断）——新增 **11 节**

总计 v3 新增 **11 节** ch16 尾节内容，新增 9 个实践清单项（9.16-9.24）、5 个失败模式（8.11-8.15）、3 个概念边界（5.8-5.10）、4 个架构工作流（6.9-6.12）、3 个方案权衡（7.8-7.10）、7 个关键结论（8-14）、2 个跨 Part 关系、4 个开放问题（10-13），文件从 v2 约 29KB 扩展到 v3 约 38KB。

**最终边界确认**：Part IV 全部六章（ch16-21）均已全量覆盖至章节末尾。ch16 尾节止于 §16.14 + 本章小结 + 参考文献；ch17 止于 §17.20；ch18 止于 §18.20；ch19 止于 §19.18；ch20 止于 §20.17；ch21 止于 §21.18。无剩余未覆盖内容。
