---
title: "Part VIII：部署与基础设施"
created: 2026-07-18
updated: 2026-07-22
status: snapshot
type: source-summary
tags: [ai/enterprise-agent-platform]
collection: "企业级 Agent 平台工程"
collection_kind: book
collection_order: 8
collection_home: "（总览）"
source: "datagallery-lab/enterprise_agent_platform@e5d97a66c467045413fff692d28fd2c401503029"
sources:
  - "[[sources/ch43-gpu-kubernetes]]"
  - "[[sources/ch44]]"
  - "[[sources/ch45-llm]]"
  - "[[sources/ch46-gitops-iac]]"
source_kind: third-party
as_of: 2026-07-18
confidence: SYNTHESIZED
contested: false
contradictions: []
related:
  - "（总览）"
  - "模型路由与工程化"
---
# Part VIII：部署与基础设施

## 一、章节覆盖表

| 章节 | 标题 | 核心主题 | 读取方式 |
|------|------|----------|----------|
| ch43 | GPU 调度与 Kubernetes | GPU节点池、队列调度(Volcano/Kueue)、算力承诺、多租户资源隔离 | ax --budget 300000 (截断) + web_extract char_limit=60000 |
| ch44 | 模型部署 | KServe InferenceService、金丝雀发布/回滚、服务形态(在线/批/Embedding)、部署框架对比 | ax --budget 300000 (截断) + web_extract char_limit=60000 |
| ch45 | LLM 网关与多租户 | LiteLLM统一入口、多租户认证/授权/配额、路由决策链、fallback降级、缓存隔离 | ax --budget 300000 (截断) + web_extract char_limit=60000 |
| ch46 | GitOps、IaC 与边缘推理 | Terraform+Helm+ArgoCD声明式交付、环境Promotion、边缘OTA同步、漂移检测 | ax --budget 300000 (截断) + web_extract char_limit=60000 |

> **注**：所有4章 ax 命令均触发内部行限制（stderr显示 "N more result(s) hidden"），按 web-reading skill 指引补充 web_extract 获取完整内容。详见报告末尾。

## 二、Part 核心问题

Part VIII 回答企业 Agent 平台从"模型能跑"到"生产可用"的四个递进问题：

1. **算力如何承诺？**（ch43）—— GPU 昂贵且稀缺，多租户、多模型抢同一批卡时，如何保证在线推理不被训练Job挤占？[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch43-gpu-kubernetes.md] — 概括：GPU调度层要回答一个具体承诺：哪类业务、在什么优先级和配额下，可以在多长时间内拿到哪类GPU
2. **模型如何发布？**（ch44）—— "有卡能跑"与"生产可调、可发、可回滚"之间缺一层模型服务层 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch44.md] — 概括：模型服务层要把权重、镜像、GPU节点、就绪探针、灰度和回滚收敛成稳定API
3. **调用如何治理？**（ch45）—— 没有统一网关时，每个 Agent 各自接模型 API "成本不透明、限流无法统一、审计也不可查" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch45-llm.md]
4. **交付如何可重现？**（ch46）—— 手工部署 "谁改了什么、何时改的、现在与期望状态是否一致，往往说不清" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch46-gitops-iac.md]

## 三、章节递进逻辑

```
ch43 (算力底座) → ch44 (模型服务层) → ch45 (统一网关入口) → ch46 (声明式交付)
```

- **ch43→ch44**：ch43 分配 GPU 节点池和队列，ch44 的 InferenceService 消费 GPU 配额并暴露稳定 API。"第43章保证'1张A100在60秒内绑定到推理Pod'，第44章保证'该Pod就绪后对外只暴露 /v1/chat/completions 与明确的 /ready 语义'" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch44.md]
- **ch44→ch45**：ch44 交付可发布的模型服务，ch45 的网关在此之上叠加租户、配额与合规路由。"第45章保证'这条请求以哪个tenant身份、走哪条backend、失败时退到哪'" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch45-llm.md]
- **ch45→ch46**：前三章分别交付算力、模型服务和网关，ch46 回答 "这些组件如何作为整体交付到 dev、staging、prod，如何保持版本化和可晋升" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch46-gitops-iac.md]
- **闭合关系**：Part VIII前半段在ch44闭合——"第43章供卡，第44章供模型API，第45章供统一入口" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch44.md]

## 四、关键结论

### ch43 GPU 调度
1. **专用节点池优于统一池**：将延迟敏感推理、批处理、训练混在同一默认节点池，"短期利用率可能看起来更高，长期会把SLO、成本和事故责任混在一起" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch43-gpu-kubernetes.md]
2. **Gang Scheduling是AI作业刚需**：默认kube-scheduler不理解"一把糖必须同时发到N个小朋友手里"，需要Volcano的PodGroup.minMember [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch43-gpu-kubernetes.md]
3. **GPU共享有代价**：MIG适合同规格同SLA多租户推理；Time-Slicing适合dev/test；"在线推理默认独占整卡仍是多数规模化企业的书面策略" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch43-gpu-kubernetes.md]
4. **Slurm与K8s可共存**：Slurm管重训练与科学计算，K8s管推理与Agent Runtime，"统一台账比分仓采购更有价值" [INFERRED ch43]
5. **Autoscaler对GPU滞后**：GPU节点从Pending到可调度常需8-25分钟，可预测尖峰需提前预热 gpu-burst 缓冲池 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch43-gpu-kubernetes.md]

### ch44 模型部署
6. **模型部署≠单机跑通**："单机把权重跑起来只是能力验证"，生产还要处理镜像、权重、GPU节点、就绪探针、流式连接、灰度切流和多环境晋升 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch44.md]
7. **Agent不应知道Pod IP**：某次NL2SQL批量报错根因是Runtime直连vLLM Pod内网IP，HPA重建后IP变更 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch44.md]
8. **readiness探针要区分/health和/ready**：/health只说明进程监听端口，/ready才说明权重已加载、KV Cache已初始化 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch44.md]
9. **Embedding版本静默失败比503更危险**："一旦权重版本与索引不一致，检索会'静默失败'" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch44.md]
10. **模型版本=权重digest+Runtime镜像+关键参数**：Git里改一行YAML不等于新版本，直到InferenceService产生新Revision并通过Canary门禁 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch44.md]

### ch45 LLM 网关
11. **网关是治理执行点，不替代治理策略**："第41章讲Token成本核算策略，第42章讲SLO与熔断。网关是这些策略的执行点，不是替代者" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch45-llm.md]
12. **租户隔离需要四层叠加**：认证（谁在调）、授权（允许调什么）、配额（能调多少）、观测（调了什么），"少掉任何一层，问题不会立刻出现，但规模一上来就会暴露" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch45-llm.md]
13. **tenant应来自Key映射，不信任客户端Header**：某次合规扫描发现finance Agent请求出现在OpenAI账单，根因是网关信任客户端 X-Tenant-Id [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch45-llm.md]
14. **fallback链必须无环**：某次backend故障时网关QPS翻倍，因为fallback指向同一KServe Service [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch45-llm.md]
15. **缓存key必须含tenant_id**：某次零售客服Agent返回含金融术语的缓存片段，根因是cache key=hash(prompt)未含tenant [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch45-llm.md]

### ch46 GitOps/IaC
16. **"改生产"的唯一合法路径是merge PR**：某次运维手工kubectl edit LiteLLM ConfigMap排查延迟后未改回，"没有任何PR、没有ArgoCD Sync记录、没有Terraform state变更，复盘无法回答'是谁在何时改了什么'" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch46-gitops-iac.md]
17. **Terraform管云、Helm管应用、ArgoCD管同步**：四层职责不能互相替代，"用Terraform硬写Deployment模板可维护性不如Helm；用Helm创建VPC则状态与模块复用差" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch46-gitops-iac.md]
18. **prod禁止automated sync**：某次工程师merge到main，prod网关配置5分钟内变更，"业务高峰期间误关fallback" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch46-gitops-iac.md]
19. **边缘推理是GitOps的特殊overlay**：OTA下载到.staging，sha256校验通过后才atomic rename到active [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch46-gitops-iac.md]
20. **边缘OTA必须原子切换**："下载不完整的模型文件如果被直接加载，可能表现为回答质量异常或进程随机崩溃" [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch46-gitops-iac.md]

## 五、概念边界

| 概念 | 定义 | 边界（不属于什么） | 来源 |
|------|------|-------------------|------|
| GPU调度 | 在集群内按优先级、配额和拓扑把GPU分配给Pod/Job | 不同于ch44的模型版本管理与Canary流量 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch43-gpu-kubernetes.md] |
| Kubernetes | 容器编排平台，提供Pod、节点、命名空间等抽象 | 不同于Slurm的HPC批作业语义与sbatch文化 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch43-gpu-kubernetes.md] |
| 队列调度器 | 在K8s之上管理AI作业排队、优先级与Gang Scheduling | 不同于默认kube-scheduler的"来一个Pod调度一个" | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch43-gpu-kubernetes.md] |
| GPU共享 | 多工作负载复用同一张物理卡（MIG、时间片等） | 不同于HPA增加Pod副本（那是水平扩展） | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch43-gpu-kubernetes.md] |
| 模型服务 | 对外暴露稳定API的推理部署单元 | 不等同于裸vLLM进程 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch44.md] |
| Serving Runtime | 执行推理的容器镜像和启动参数 | 不等同于训练Job | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch44.md] |
| 模型版本 | 权重、Runtime镜像和关键启动参数的可发布组合 | 不等同于Git代码版本 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch44.md] |
| LLM 网关 | 统一LLM API入口与治理执行点 | 不同于API Gateway全站入口 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch45-llm.md] |
| 租户 | 资源与配额隔离单元 | 不同于K8s Namespace本身 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch45-llm.md] |
| 路由规则 | 决定backend的匹配逻辑 | 不同于模型服务发布（ch44） | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch45-llm.md] |
| 降级 | 主backend失败时的备用路径 | 不同于模型Canary（ch44） | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch45-llm.md] |
| GitOps | 用Git驱动部署与变更，关注集群期望状态和实际状态的持续对齐 | 不等同于只负责构建镜像的CI | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch46-gitops-iac.md] |
| IaC | 用代码描述基础设施 | 不等同于配置管理脚本 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch46-gitops-iac.md] |
| 漂移 | 集群实际状态与Git声明不一致 | 和Canary流量比例是两个问题：前者说明事实源被绕过，后者说明发布策略仍在执行 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch46-gitops-iac.md] |

## 六、平台架构与工作流

### 分层架构（自底向上）
```
L6 观测    ─── OTel、Langfuse (ch38)
L5 平台应用 ─── Agent Runtime、DataAgent
L4 网关    ─── LiteLLM (ch45) ── 统一入口、租户路由、配额
L3 模型服务 ─── KServe InferenceService (ch44) ── 金丝雀、回滚、Revision
L2 存储    ─── OSS模型桶、PVC
L1 算力    ─── GPU节点池、Device Plugin (ch43) ── Volcano/Kueue队列
L0 网络    ─── VPC、子网、安全组
```
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch46-gitops-iac.md 表46-2 + ch43图43-1]

### 交付工作流（GitOps Promotion）
```
开发者提交PR → CI校验(terraform plan + helm template + 命名一致性检查)
→ merge到main → dev自动sync → staging自动sync + 集成测试(smoke四路径)
→ 打tag prod-vX.Y.Z → Platform Owner审批 → prod Manual Sync(ArgoCD)
```
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch46-gitops-iac.md]

### 模型发布状态机（ch44）
```
Draft → Staging(离线评测通过) → Canary(5%流量) → Expanding(20%→50%→100%)
→ Stable(全量无告警) → Deprecated(保留N版回滚)
Rollback: 从Canary/Expanding一键回到旧Stable Revision
```
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch44.md 表44-6]

### 网关路由决策链（ch45）
```
租户认证 → 合规硬截断(finance禁止云端) → task_type匹配
→ 成本/延迟偏好 → backend选择 → 显式fallback_chain
```
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch45-llm.md 图45-3]

### 边缘推理拓扑（ch46）
```
门店(llama.cpp 7B Q4) ── FAQ本地处理，复杂问题回传中心网关→llm-general-32b
工厂(ONNX 视觉质检) ── 初判边缘，边界样本回传32B复核
物流(MLC 3B手持) ── 简单查询本地，复杂理赔回传中心
```
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch46-gitops-iac.md 表46-4+46-5]

## 七、方案权衡

### ch43 调度器选型
| 方案 | 优势 | 代价 | 本书建议 |
|------|------|------|----------|
| 默认kube-scheduler | 零额外组件 | 无队列、无Gang、无AI优先级 | 在线推理+Affinity |
| Volcano | Gang、队列、批作业原生 | CRD与运维学习成本 | gpu-train/gpu-batch |
| Kueue | 轻量、多租户配额 | Gang弱于Volcano | 配额治理 |
| **Volcano+Kueue并存** | 各取所长 | 运维成本 | **推荐（规模化）** |
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch43-gpu-kubernetes.md 表43-4+43-9]

### ch43 节点池策略
| 方案 | 优势 | 代价 | 适用场景 |
|------|------|------|----------|
| 专用节点池 | SLO可预测、故障域清晰 | 利用率可能偏低 | 规模化企业推荐 |
| 统一GPU池 | 采购决策简单 | 争抢、OOM | 仅试点 |
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch43-gpu-kubernetes.md 表43-8]

### ch44 部署框架
| 类型 | 代表 | 为什么用 | 不适合什么 |
|------|------|----------|------------|
| K8s原生Serving | KServe | InferenceService CRD、流量分割、自动扩缩 | 非K8s环境 |
| 模型打包 | BentoML | 多框架打包、本地到生产一致 | 复杂多模型流量治理 |
| 高性能推理 | Triton | Embedding、多模型同进程、动态批处理 | 快速原型 |
| 分布式编排 | Ray Serve | Python原生、与Ray生态集成 | 运维团队无Ray经验 |
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch44.md 表44-3]

### ch45 网关产品
| 方案 | 优势 | 代价 | 本书建议 |
|------|------|------|----------|
| LiteLLM | 100+后端、OpenAI兼容、社区活跃 | 深度企业特性需扩展 | 初版推荐 |
| Portkey | 可观测、路由、缓存成熟 | SaaS/许可 | 对标参考 |
| Higress/Kong AI | 与API网关生态集成 | LLM语义需插件 | 混合架构 |
| 自研 | 完全定制 | 工程量大 | 长期选项 |
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch45-llm.md 表45-5]

### ch45 缓存策略
| 方案 | 优势 | 代价 | 适用场景 |
|------|------|------|----------|
| 网关semantic cache | 跨backend、可租户隔离 | 一致性难 | 重复问答多 |
| 后端prefix cache | 延迟最低 | 绑定单引擎 | vLLM长system prompt |
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch45-llm.md 表45-6]

## 八、失败模式

| 编号 | 失败模式 | 现象 | 根因 | 恢复策略 | 来源 |
|------|----------|------|------|----------|------|
| F1 | Device Plugin升级致全集群GPU不可见 | 所有推理Pod Pending | Plugin版本与驱动不匹配 | 节点池滚动升级，先在staging验证 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch43-gpu-kubernetes.md] |
| F2 | Volcano PodGroup minMember>可用GPU | 3个Pod Running、1个Pending，服务永久Unschedulable | Gang永远无法满足 | 降minMember或扩容 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch43-gpu-kubernetes.md] |
| F3 | Time-Slicing与在线推理混部 | P99 TTFT从600ms升至1.8s，GPU利用率却"达标" | Embedding批Job与在线推理共享时间片 | 在线推理改回独占整卡，批处理迁入独立池 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch43-gpu-kubernetes.md] |
| F4 | readiness探针过短致"半启动"接流量 | 32B模型加载需8分钟，Pod在2分钟被标记Ready，返回503 | /health不反映模型加载状态 | 改用/ready端点，initialDelay≥120s | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch44.md] |
| F5 | Canary 5%通过但全量后OOM | Canary阶段指标正常，扩到50%后多节点同时OOM | 5%流量中长context样本不足 | 强制离线长上下文压测，分池部署 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch44.md] |
| F6 | Embedding版本静默失败 | 向量维度变化，RAG检索"能答但引用文档全错" | storageUri指向latest/被覆盖 | URI含不可变版本号 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch44.md] |
| F7 | fallback链配置错误致无限重试 | gateway QPS翻倍，下游更快崩溃 | fallback指向同一KServe Service | 限制每请求最多1次fallback，CI静态检查DAG | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch45-llm.md] |
| F8 | 伪造Header致finance访问云端 | finance Agent请求出现在OpenAI账单 | 网关信任客户端X-Tenant-Id | tenant仅来自Key注册表 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch45-llm.md] |
| F9 | 缓存串扰 | 零售客服返回含金融术语的缓存片段 | cache key未含tenant | key=hash(tenant_id+model+prompt) | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch45-llm.md] |
| F10 | ArgoCD prod误开automated sync | 合并main后prod 5分钟内变更 | Application复制dev配置 | prod必须automated:null+Manual Sync | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch46-gitops-iac.md] |
| F11 | Terraform state未锁致节点池混乱 | max_size被覆盖，HPA扩Pod无节点 | 本地state文件，两人同时apply | Terraform Cloud/OSS backend+锁 | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch46-gitops-iac.md] |
| F12 | 边缘OTA断点续传未校验sha256 | 门店模型文件损坏，少200MB仍加载 | 弱网中断后加载不完整文件 | .staging下载+sha256校验+atomic rename | [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch46-gitops-iac.md] |

## 九、实践清单

### ch43 GPU调度
- [ ] 定义4类节点池：gpu-inference、gpu-train、gpu-batch、gpu-burst
- [ ] 推理节点打污点workload=online-infer:NoSchedule
- [ ] 推理Pod带toleration+nodeAffinity(nodepool=gpu-inference)
- [ ] 训练/批Job进Volcano PodGroup，设置minMember
- [ ] 跨事业部配额进Kueue ClusterQueue
- [ ] GPU监控覆盖DCGM_UTIL、Pod Pending、Volcano队列、节点NotReady
- [ ] 可预测尖峰前预热gpu-burst缓冲池（8-25min冷启动）
- [ ] 月度量复审：模型副本、配额、队列、节点标签、实验命名空间
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch43-gpu-kubernetes.md]

### ch44 模型部署
- [ ] InferenceService名与网关model_name 1:1对应
- [ ] readiness探针指向/ready而非/health，initialDelay≥120s
- [ ] storageURI含不可变版本号，不用latest/
- [ ] 金丝雀发布：5%→20%→50%→100%，每阶段观察P99 TTFT+错误率
- [ ] Rollback一键：patch canaryTrafficPercent:0，保留旧Revision≥7天
- [ ] 季度回滚演练：验证旧权重可拉取、旧Revision可Ready、网关可重新发现
- [ ] Embedding与LLM拆成独立InferenceService，独立HPA
- [ ] 模型服务目录记录：Revision、权重URI、镜像digest、Runtime参数、owner
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch44.md]

### ch45 LLM 网关
- [ ] 所有Runtime/DataAgent/Console只对接网关，禁止直连KServe
- [ ] tenant仅从API Key注册表映射，不信任客户端X-Tenant-Id Header
- [ ] finance租户allowed_models白名单只有本地模型，云端backend不可见
- [ ] fallback DAG静态分析无环，每请求最多1次fallback
- [ ] 缓存key=hash(tenant_id+model+prompt)，finance禁用跨会话缓存
- [ ] 错误码统一：401/403不可重试，429按Retry-After退避，502/503可fallback
- [ ] CI探测：finance请求云端应403，retail超预算应429
- [ ] 网关多副本+Redis限流，限流失效策略finance选fail-closed
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch45-llm.md]

### ch46 GitOps/IaC
- [ ] Terraform管云资源(GPU节点池/VPC/OSS/IAM)，Helm管应用，ArgoCD管同步
- [ ] prod Application automated:null + Manual Sync，追固定tag不追main
- [ ] Terraform state用remote backend+锁，禁止本地apply prod
- [ ] 密钥通过External Secrets从Vault注入，Git中不存明文
- [ ] Promotion tag模型与网关同版本，同一PR原子变更
- [ ] staging smoke覆盖四路径：默认对话、NL2SQL代码、finance拒绝云端、fallback
- [ ] 边缘OTA：下载到.staging→sha256校验→atomic rename→版本上报inventory
- [ ] 三联审计：Git PR + ArgoCD Sync + Terraform apply 字段对齐
[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part08-deployment/ch/ch46-gitops-iac.md]

## 十、跨Part关系

| 目标章节 | 关系 | 说明 |
|----------|------|------|
| Part II ch6 | ch44上游 | ch6讲单机vLLM/SGLang跑起来；ch44讲生产服务化封装 |
| Part II ch7 | ch43互补 | ch7讲KV Cache优化；ch43讲GPU算力分配，互相引用容量约束 |
| Part II ch9 | ch43下游 | ch9的LoRA微调依赖ch43的gpu-train队列与Gang Scheduling |
| ch38 可观测 | ch43/44/45/46下游 | 所有层的指标、Trace、告警接入ch38观测链路 |
| ch39 评测 | ch44下游 | ch44的Canary门禁依赖ch39离线评测通过 |
| ch41 成本 | ch43/45下游 | ch43的FinOps标签+ch45的token计量 → ch41成本核算 |
| ch42 SLO | ch45上游 | ch42定义限流SLO，ch45网关执行 |
| ch50 安全 | ch45下游 | ch50的安全策略在ch45网关执行"finance禁止云端"类规则 |

[COMPREHENSIVE: 基于4章全文交叉引用]

## 十一、开放问题

1. **GPU调度与网关联动的工程实现**：ch43提到GPU调度层和网关配额需要联动（容量信号、租户信号、请求信号），但具体信号格式和联动机制细节留待实施 [INFERRED ch43.3.2]
2. **Slurm-K8s统一台账的具体方案**：ch43提出"CMDB记录物理GPU总量，Slurm分区与K8s节点池分别占配"，但统一台账的数据模型和同步机制未展开 [INFERRED ch43.2.2]
3. **LiteLLM深度企业特性扩展范围**：ch45指出LiteLLM的"复杂RBAC、合规路由和tenant白名单需要二次加固"，但加固的具体架构（custom callback vs 外部策略服务）未定 [INFERRED ch45]
4. **边缘推理的离线降级策略**：门店弱网时"本地降级答'请稍后联系人工'"，但什么条件触发降级、回传中心的重试策略细节未详述 [INFERRED ch46]
5. **GitOps monorepo vs 多仓库的长期演进路径**：ch46建议"长期可按terraform/与helm/拆库"，但跨库Promotion tag对齐机制需要统一release manifest [INFERRED ch46]
6. **GPU配额复审的自动化程度**：ch43.3.3提出月度量复审（检查owner/budget/到期时间），但复审如何自动化、与Kueue配额如何联动未详细说明 [INFERRED ch43]

## 十二、suggested_tags

```
enterprise-agent-platform/deployment, kubernetes/gpu, gpu/scheduling,
llm/serving, llm/gateway, gitops, iac, terraform, argocd,
model-deployment, canary-release, multi-tenant, edge-inference,
volcano, kueue, kserve, litellm
```

## 十三、source_summary

Part VIII「部署与基础设施」由4章组成，以自底向上的层次递进构建企业Agent平台的部署底座。ch43从GPU调度入手，定义节点池、队列（Volcano/Kueue）、优先级，将稀缺算力变成可承诺的资源池；ch44承接算力，通过KServe InferenceService将推理引擎包装成稳定可发布的模型服务，引入金丝雀发布状态机、Revision管理和四类服务形态（在线/批/Embedding/多模型路由）；ch45在此之上建立LLM 网关（LiteLLM）作为统一入口，实现多租户认证/授权/配额/路由/降级/缓存；ch46以GitOps和IaC（Terraform+Helm+ArgoCD）将前三层纳入声明式交付流程，并扩展边缘推理（门店/工厂/物流）的OTA同步策略。全书贯穿"承诺→发布→治理→可重现"的主线，每个章节明确职责边界、失败模式和恢复路径。关键跨Part引用包括ch38（可观测）、ch39（评测门禁）、ch41（成本）、ch42（SLO）、ch50（安全）。

---

## 附录：读取方式报告

| 章节 | 文件路径 | 字节数 | 覆盖状态 | 读取方式 |
|------|----------|--------|----------|----------|
| ch43 | /tmp/eap_ch43.md | 7,721 | ✅ 完整（ax截断后web_extract补齐） | ① ax --md --budget 300000 → 截断(432 more hidden) ② web_extract char_limit=60000 → 完整 |
| ch44 | /tmp/eap_ch44.md | 7,795 | ✅ 完整（ax截断后web_extract补齐） | ① ax --md --budget 300000 → 截断(441 more hidden) ② web_extract char_limit=60000 → 完整 |
| ch45 | /tmp/eap_ch45.md | 6,165 | ✅ 完整（ax截断后web_extract补齐） | ① ax --md --budget 300000 → 截断(443 more hidden) ② web_extract char_limit=60000 → 完整 |
| ch46 | /tmp/eap_ch46.md | 6,242 | ✅ 完整（ax截断后web_extract补齐） | ① ax --md --budget 300000 → 截断(424 more hidden) ② web_extract char_limit=60000 → 完整 |

**最终文件**：`/tmp/eap_part08_summary.md`
