---
title: "Part IX：前端、交互与多模态"
created: 2026-07-18
updated: 2026-07-22
status: snapshot
type: source-summary
tags: [ai/enterprise-agent-platform]
collection: "企业级 Agent 平台工程"
collection_kind: book
collection_order: 9
collection_home: "（总览）"
source: "datagallery-lab/enterprise_agent_platform@e5d97a66c467045413fff692d28fd2c401503029"
sources:
  - "[[sources/ch47-ui]]"
  - "[[sources/ch48-generative-ui]]"
  - "[[sources/ch49-agent]]"
source_kind: third-party
as_of: 2026-07-18
confidence: SYNTHESIZED
contested: false
contradictions: []
related:
  - "（总览）"
---
# Part IX：前端、交互与多模态

## 章节覆盖表

| 章节 | 标题 | 核心主题 | 读取方式 |
|------|------|----------|----------|
| Ch47 | 对话 UI 与流式输出 | 企业 Agent UI 七类界面单元、SSE 事件协议、消息模型三层、流式状态机、前端框架选型、可观测 | `web_extract`（ax --budget 300000 被截断，313行隐藏） |
| Ch48 | Generative UI 与富交互 | 受控组件渲染、Artifact 生命周期、审批卡片、组件白名单、渲染契约、XSS 防护、版本治理 | `web_extract`（ax --budget 300000 被截断，314行隐藏） |
| Ch49 | 多模态输入与语音 Agent | 异步文件解析、Context Ref、WebRTC/WebSocket 语音链路、VAD/ASR/TTS 编排、多模态权限与审计 | `web_extract`（ax --budget 300000 被截断，331行隐藏） |

## Part 核心问题

Part IX 回答企业 Agent 前端领域三个递进问题：

1. **对话 UI 如何承载企业任务而非单纯聊天？** — 聊天框只是入口，真正的企业 Agent UI 必须展示工具调用状态、等待审批、可编辑证据和任务进度 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch47-ui.md — "企业 Agent 的前端要展示工具调用状态、等待审批、可编辑证据和任务进度，聊天框只是入口之一"]
2. **富交互如何在不引入安全风险的前提下将工具结果变成可操作界面？** — Generative UI 的关键是 Agent 在受控协议内选择、填充和编排平台登记过的组件，而非让模型随手生成界面 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch48-generative-ui.md — "让 Agent 在受控协议内选择、填充和编排平台登记过的组件，而不是让模型随手生成界面"]
3. **多模态输入如何在进入任务链路前具备治理属性？** — 非文本材料进入上下文之前，先经过上传与解析链路，产出结构化文本、版面、表格、元数据和引用 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch49-agent.md — "非文本材料进入上下文之前，要先经过上传与解析链路，产出结构化文本、版面、表格、元数据和引用"]

## 章节递进

Part IX 三章形成严格的能力递进关系：

1. **Ch47 → 建立事件流基础**：从聊天界面出发，定义企业 Agent UI 七类界面单元（会话入口、消息流、工具进度、上下文面板、业务控件、反馈入口、观测标识），然后深入到流式交互协议（SSE、增量渲染、事件 reducer），最终建立消息模型三层（会话层、消息层、事件层）和流式状态机。这一步把前端从"展示文本"升级为"表达任务过程"。

2. **Ch48 → 在事件流上叠加可操作对象**：基于 Ch47 的事件流和消息模型，将工具结果映射为受控组件（图表、表格、表单、Artifact、审批卡片）。核心递进在于：工具结果不再是"被推送的 JSON"，而是进入白名单组件、绑定字段权限、数据引用和服务端策略。表 48-1（能力递进表）明确把这个关系量化：Vercel AI SDK 的流式能力 → 加上白名单组件映射；assistant-ui 的对话组件 → 加上业务组件注册和审批证据链。

3. **Ch49 → 扩展输入模态并前移治理**：在 Ch47-48 建立的交互框架上，打开多模态输入通道（文件、图片、语音、视频）。关键递进在于治理前移：输入侧每多一种模态，就多一条需要审计的链路。异步文件解析、Context Ref、实时语音确认卡片、打断与取消语义，都是在前两章基础上的输入层治理扩展。

三章共享的核心架构假设：**所有交互事件最终进入同一条 Trace**，前端行为和后端 Run 状态必须在同一事件链上可排序、可恢复、可审计。

## 关键结论

### Ch47 关键结论

1. **企业 Agent UI 包含七类界面单元**，缺一不可：会话入口、消息流、工具进度、上下文面板、业务控件、反馈入口、观测标识 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch47-ui.md — 表 47-2]
2. **SSE 作为文本 DataAgent 默认传输协议**，WebSocket 留给语音和强实时控制，WebRTC 留给第49章 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch47-ui.md — "企业 DataAgent 默认可以选择 SSE 作为文本任务的主传输协议"]
3. **消息模型拆为三层**：会话层（conversation_id/租户/权限）、消息层（用户/助手/工具/审批/错误消息）、事件层（增量事件，支持流式、恢复、幂等） [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch47-ui.md — 表 47-4]
4. **流式协议必须包含 run_id/message_id/event_id/seq/事件类型/可恢复标记/trace 关联** [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch47-ui.md — "事件协议要比 token 流多几层信息"]
5. **最终状态必须以 Runtime 事件为准，而非模型文字为准** [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch47-ui.md — "最终状态必须以 Runtime 事件为准，而非以模型文字为准"]
6. **前端 reducer 采用事件折叠模式**，而非客户端拼接 token 或全量消息刷新 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch47-ui.md — §47.4.3]
7. **企业自研方向是"可演进的事件协议 + 受控组件 + 可观测链路"**，框架可换，协议和治理边界不能丢 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch47-ui.md — §47.4.2]

### Ch48 关键结论

1. **Generative UI 的价值不在视觉效果，而在工作界面** [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch48-generative-ui.md — 表 48-1 下方判断]
2. **工具结果到 UI 的最小链路**：工具执行 → Runtime 附加权限/证据 → Render Contract 决定组件 → 前端按白名单和权限渲染 → 用户动作回 Policy/Observability [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch48-generative-ui.md — §48.3]
3. **Artifact 是业务产物，不是消息附件**，需要独立生命周期（生成中→已生成→编辑中→已提交→已归档）和版本管理 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch48-generative-ui.md — 图 48-6 + 表 48-6]
4. **组件白名单是最小安全边界**：平台默认最小化暴露组件类型、字段、版本和动作集合 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch48-generative-ui.md — 表 48-3]
5. **敏感字段不进入消息流**：小表格随消息返回，大结果走 `data_ref`，服务端在展开/筛选/导出时重新做权限和脱敏 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch48-generative-ui.md — §48.5.1]
6. **Artifact 权限再校验**：产物保存后权限变化时，元数据保留给审计，图表数据/明细表按当前权限重新判断 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch48-generative-ui.md — §48.12]
7. **Artifact 需撤回机制而非简单删除**：区分隐藏内容/标记过期/撤销发布/替换新版本/通知接收者 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch48-generative-ui.md — §48.15]

### Ch49 关键结论

1. **多模态输入能力拆成四层**：异步文件解析层、实时媒体会话层、语音能力编排层、浏览器采集与传输层 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch49-agent.md — 表 49-1]
2. **文件上传采用"对象存储 + 解析任务 + Context Ref"模式**，Agent 不直接消费原始文件 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch49-agent.md — §49.3.1]
3. **Context Ref 至少包含租户、权限、来源、解析版本、质量评分和保留策略** [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch49-agent.md — §49.2.1]
4. **浏览器实时语音默认选择 WebRTC**，后台或非浏览器链路用 WebSocket，录音上传适合异步分析 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch49-agent.md — §49.6.1]
5. **语音确认不等于语音 Agent**：还需轮次管理、工具确认、上下文恢复和事故回放 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch49-agent.md — 表 49-1 下方说明]
6. **敏感动作必须在转写后再次确认**，不能把实时识别结果直接当作最终意图 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch49-agent.md — §49.15]
7. **多模态证据链必须跨模态回放**：经营报告引用截图→回到截图文件/解析版本/识别区域/确认记录 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch49-agent.md — §49.7]

## 概念边界

以下概念边界是 Part IX 的核心区分，横跨三章：

| 概念 | 是 | 不是 | 来源 |
|------|----|------|------|
| 对话 UI | 承载会话、工具进度、业务动作和反馈的交互层 | 聊天气泡组件 | Ch47 表47-3 |
| 流式输出 | 服务端把生成过程拆成事件或增量片段发送 | token 打字机效果 | Ch47 表47-3 |
| SSE | 浏览器通过 HTTP 接收服务端单向事件 | WebSocket 替代品 | Ch47 表47-3 |
| 增量渲染 | 前端按事件更新局部消息、工具卡和状态（需幂等和回滚） | 字符串拼接 | Ch47 表47-3 |
| Generative UI | LLM 根据任务上下文选择并填充受控 UI 组件 | 模型生成任意前端代码 | Ch48 表48-3 |
| 工具调用渲染 | 将 Tool Call 状态和输出映射为卡片/图表/表格/表单 | 把工具日志原样展示 | Ch48 表48-3 |
| Artifact | 对话外可编辑/可保存/可复用的产物（有生命周期和版本） | 消息附件 | Ch48 表48-3 |
| 组件白名单 | 平台允许 Agent 引用的组件类型、字段、版本和动作集合 | 前端组件库全集 | Ch48 表48-3 |
| 审批卡片 | 高风险动作前展示影响范围、证据和确认入口（需留痕+策略绑定） | 普通确认弹窗 | Ch48 表48-3 |
| 多模态输入 | 文本之外的文件、图像、音频、视频等输入 | 多模态模型专属（也可用解析器+单模态模型组合） | Ch49 表49-4 |
| Context Ref | 指向解析后安全上下文的引用（带权限、版本、质量） | 原始文件 URL | Ch49 表49-4 |
| ASR | 自动语音识别，将音频转为文本 | 语义理解 | Ch49 表49-4 |
| TTS | 文本转语音，将文字回答合成为音频 | 语音 Agent（后者还需打断/轮次/工具控制） | Ch49 表49-4 |
| WebRTC | 浏览器实时音视频通信技术 | WebSocket 替代（更贴近媒体传输） | Ch49 表49-4 |
| VAD | 语音活动检测，判断用户何时开始/停止说话 | 唤醒词 | Ch49 表49-4 |

## 平台架构/工作流

### 对话 UI 在企业平台中的位置（Ch47 图 47-1）

```
┌─────────────────────────────────────┐
│         Agent Console (对话 UI)       │
│  消息展示 | 状态感知 | 审批交互 | 证据展示 │
├─────────────────────────────────────┤
│       Conversation API               │
│        (统一事件流转换)                │
├─────────────────────────────────────┤
│   Runtime | LLM 网关（LLM Gateway） |            │
│   Tool Registry | Observability      │
└─────────────────────────────────────┘
```

三条关键边界：
1. Agent Console 不直接调用模型，所有模型交互经 Conversation API 和 Runtime 统一鉴权、trace、限流
2. 工具执行结果不直接写入前端；经 Runtime 整理、权限过滤和渲染契约转换
3. 前端事件回流到观测系统，用户停止生成、重试、展开工具卡、提交差评都与后端 trace 关联

### Generative UI 在企业平台中的位置（Ch48 图 48-4）

Generative UI 位于 Tool Registry、Policy、Runtime、Observability 与 Console 的交汇处：
- Tool Registry 返回工具能力和结构化结果，不是前端组件
- Policy 管后端工具执行 + 前端动作（导出/提交审批/保存等）
- Observability 记录用户看到的界面状态（哪张图/哪张表/改了什么/批准了什么）
- Render Gateway 把内部事件转换成受控渲染契约

### 多模态输入层在企业平台中的位置（Ch49 图 49-4）

多模态输入层位于前端 Console 与 Agent Runtime 之间：
- 上传入口和 Agent Runtime 解耦：文件先入对象存储和解析任务，解析结果通过 Context Store 暴露给 Runtime
- 实时媒体和业务动作解耦：语音链路负责采集/传输/转写/播放/打断，业务动作仍经 Tool Registry 和 Policy
- 多模态输入接入 Observability：上传失败/解析警告/转写置信度/用户修正/打断/确认/工具调用同一条 trace

### 流式对话时序（Ch47 图 47-2）

```
用户提问 → 服务端推送: state → token_delta → tool_call_started → tool_call_delta → tool_call_completed → done
前端逐步渲染: 占位消息 → 文本增量 → 工具卡片(loading) → 工具卡片(进度) → 工具卡片(结果) → 完成态
```

### 文件上传异步解析流水线（Ch49 图 49-5）

```
前端上传 → 对象存储 → 触发异步解析任务 → OCR/文档解析 → 分块入库 → 返回解析状态
                                              ↓ (独立进行，不阻塞 UI)
                                         对话可并行进行
```

### Artifact 生命周期状态机（Ch48 图 48-6）

```
generating → generated → editing → committed → archived
                ↓            ↓           ↓
             (可重新生成)  (可放弃编辑)  (可撤回/标记过期)
```

## 方案权衡

### Ch47：传输协议选择

| 方案 | 适用 | 代价 | 企业建议 |
|------|------|------|----------|
| SSE | 文本对话、工具进度、报告生成 | 双向控制较弱，长连接受代理/超时影响 | **DataAgent 默认选择** [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch47-ui.md — "企业 DataAgent 默认可以选择 SSE 作为文本任务的主传输协议"] |
| WebSocket | 语音、多人协作、强控制 | 网关/鉴权/心跳/扩容复杂 | 留给 Ch49 多模态和实时交互 |
| HTTP 分块响应 | 内部极简示例/模型流适配 | 事件语义弱、恢复困难 | 不推荐作为企业默认协议 |
| WebRTC | 低延迟音视频 | 部署和调试复杂 | 浏览器实时语音默认方案（Ch49） |

### Ch48：渲染策略选择

| 方案 | 优势 | 风险 | 企业建议 |
|------|------|------|----------|
| 模型生成 HTML/JS | 表达力强，原型快 | XSS、越权动作、审计缺口 | **禁止在生产使用**，仅限隔离沙箱 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch48-generative-ui.md — "生产系统不应让模型直接生成任意 HTML/JS"] |
| 组件白名单 + 模型生成配置 | 安全可控，可审计 | 表达力受限 | **企业默认路径**，通过扩展组件库解决表达力不足 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch48-generative-ui.md — "企业默认路径应是组件白名单与模型生成配置结合"] |
| 对话内卡片 | 上下文连续，理解成本低 | 不适合大产物 | 小图表/短表格/审批提示 |
| 独立 Artifact 工作区 | 可编辑/审阅/导出/版本 | 上下文割裂 | 报告/方案/代码/长分析 |
| 前端直接渲染数据 | 首屏速度快 | 大结果卡浏览器，敏感字段暴露 | 小摘要表可，大结果必须走 `data_ref` |
| 数据引用渲染 (`data_ref`) | 动态权限/脱敏/分页 | 多一次数据获取 | 大结果和敏感数据的默认方式 |
| 静态截图 | 适合归档和外部分享 | 不可交互/难审计 | 不应替代数据引用 |
| 前端工具调用 | 读页面状态/切换筛选器 | 绕过治理 | 仅低风险 UI 动作；涉及数据查询/导出/写入须后端 |
| 后端工具渲染 | 完整审计和控制 | 灵活性略低 | 高风险动作默认方式 |

### Ch49：多模态处理策略

| 方案 | 适用 | 代价 | 企业建议 |
|------|------|------|----------|
| 异步解析 | 文件/表格/长录音 | 用户需等待 | **企业默认**，大文件和复杂版面必须异步 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch49-agent.md — "文件上传默认应走异步解析和上下文引用"] |
| 同步问答 | 小图片/小文本/低风险原型 | 大文件拖垮体验和治理 | 仅限轻量场景 |
| WebRTC 实时语音 | 浏览器端低延迟交互 | 部署复杂度高 | **浏览器实时语音默认** |
| WebSocket 语音 | 服务端到服务端/非浏览器客户端 | 媒体处理弱于 WebRTC | 后台或非浏览器链路 |
| 录音上传 + 异步 ASR | 会议纪要/客服分析/离线质检 | 不支持实时打断 | 不适合实时助理体验 |
| 直接多模态模型 | 复杂视觉语义（截图理解/图片问答） | 成本高，证据结构弱 | 仅用于无法提前结构化的内容 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch49-agent.md — "直接多模态模型适合处理复杂视觉语义……代价是成本高，证据结构弱"] |
| 解析器流水线 | 文档/表格/票据/录音 | 不支持非结构化视觉语义 | **高价值流程默认选择**，产文本块+时间戳+置信度 |
| 混合模式 | 多模态高价值流程 | 架构复杂 | 先用解析器拿证据，再用多模态模型处理非结构化 |
| 保存原始音频 | 客服质检/合规/争议复盘 | 隐私和存储风险 | 配套加密/访问审批/保留期限/删除机制 |
| 只保存转写 | 普通语音助手 | 削弱争议复盘能力 | 需提前说明产品边界 |

## 失败模式

### Ch47 流式 UI 高发故障（表 47-7）

| 故障 | 触发条件 | 处理 |
|------|----------|------|
| 连接中断 | SSE 断开或代理超时 | 事件游标恢复；超窗口则展示继续生成 |
| 事件乱序 | 多路工具事件并发返回 | 以 seq 和 event_id 做幂等折叠 |
| 旧流污染 | 用户取消后旧请求继续返回 | 前端按 run_id 丢弃失效事件 |
| 工具结果过大 | 查询返回过多行或大对象 | 消息只展示摘要，表格走数据引用 |
| 敏感字段泄漏 | 工具结果包含未脱敏字段 | 后端字段脱敏，前端只渲染白名单字段 |
| Markdown 注入 | 模型输出包含脚本或危险链接 | 禁用原始 HTML，外链做跳转确认 |
| trace 断链 | 前端事件未关联后端 trace | 所有事件携带 trace，前端交互复用同一 trace |

### Ch48 Generative UI 高发故障（表 48-7）

| 故障 | 触发条件 | 处理 |
|------|----------|------|
| 组件越权 | 模型请求渲染未授权组件或动作 | 拒绝渲染，降级为安全摘要 |
| Schema 漂移 | 工具输出字段与组件版本不匹配 | 版本化契约，前端展示兼容错误 |
| 图表误导 | 模型选择不合适图表或隐藏分母 | 显示图表规格、指标口径和样本范围 |
| 审批绕过 | 前端直接调用业务动作 | 所有动作服务端二次校验 |
| Artifact 污染 | 文档内容中提示注入被写入报告 | 标记生成来源，高风险段落要求人工确认 |
| 数据泄漏 | 表格导出绕过字段脱敏 | 下载走服务端导出任务，按权限重算 |
| 旧版本不可回放 | 组件升级后历史消息渲染失败 | 组件版本保留兼容层，必要时降级为静态摘要 |

### Ch49 多模态输入高发故障（表 49-7）

| 故障 | 触发条件 | 处理 |
|------|----------|------|
| 解析质量低 | OCR 置信度低、表格结构不完整 | 要求用户确认关键字段或改走人工校验 |
| 文件越权 | 用户上传不属于当前租户/项目的文件 | 拒绝进入上下文，记录安全事件 |
| 图片泄漏隐私 | 图片包含人脸/工牌/客户姓名 | 脱敏/裁剪或禁止进入上下文 |
| 语音误识别 | 噪声/口音/多人说话导致转写错误 | 展示实时转写，敏感动作前要求确认文本意图 |
| 实时延迟过高 | 网络抖动/模型响应慢 | 降级为按键发言/文本输入或录音上传 |
| 打断失效 | 模型仍在播放旧回答 | 前端停止播放，服务端取消 response，丢弃旧音频事件 |
| 审计缺失 | 音频/文件/工具调用无统一 trace | 会话创建时生成 trace，所有输入引用和工具调用继承 trace |

### 用户常见误解与错误归因

- 用户看到"正在生成"转 40 秒，不知道系统在等模型、等 SQL、等审批还是连接断开 → 工具进度必须拆成可辨识阶段 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch47-ui.md — 导言段]
- 用户上传附件后以为"立刻成为可信上下文" → 前端应展示 uploaded → parsing → ready_for_context → rejected 状态链 [INFERRED: 综合 Ch49 多处状态分离设计]
- 用户把实时语音"听懂了"理解为"可以执行了" → 高风险动作必须转确认卡片 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch49-agent.md — "实时语音 Agent 的交互速度很容易让用户把'听懂了'理解成'可以执行了'"]

## 实践清单

### 对话 UI 上线前验收（Ch47 §47.6-47.7）

- [ ] 七类界面单元是否齐全（会话入口/消息流/工具进度/上下文面板/业务控件/反馈入口/观测标识）
- [ ] SSE 断线后能否按游标恢复；超窗口后能否展示继续生成
- [ ] 用户取消后前端是否立即停止渲染，后端是否取消 Runtime 任务
- [ ] 工具返回 3000 行 JSON 时前端是否只展示摘要和数据引用
- [ ] 高危 HTML 是否被渲染层拦截
- [ ] 前端事件是否全部关联 trace_id
- [ ] 事件协议是否有版本号、兼容窗口和回放样例
- [ ] 所有异常路径文案是否统一（区分临时等待/可恢复失败/永久拒绝/需人工处理）
- [ ] 不同角色视图是否分离（业务用户/平台工程师/安全合规）
- [ ] 移动端/多标签/刷新页面后任务能否恢复

### Generative UI 上线前验收（Ch48 §48.16-48.17）

- [ ] 组件白名单是否最小化定义（类型/字段/版本/动作）
- [ ] 每条工具结果能否回答"进入哪个组件/展示哪些字段/允许哪些动作/是否需要二次确认"
- [ ] 图表是否绑定指标口径、数据快照和证据引用
- [ ] 导出/审批按钮是否在服务端二次校验
- [ ] Artifact 是否有独立生命周期（draft → reviewing → approved → exported → archived）
- [ ] Artifact 编辑记录是否区分机器生成与人工修改
- [ ] 权限变化后历史 Artifact 敏感内容是否重新鉴权
- [ ] 组件升级后历史消息是否仍可读（降级为静态摘要）
- [ ] Artifact 是否有撤回机制（区分隐藏/标记过期/撤销发布/替换版本/通知接收者）
- [ ] 每类组件是否有统一的空状态/加载态/错误态/权限拒绝态/审批态

### 多模态入口上线前验收（Ch49 §49.12-49.18）

- [ ] 文件上传是否走"对象存储 + 解析任务 + Context Ref"而非直接消费原始文件
- [ ] 前端是否展示文件状态链（uploaded → parsing → ready_for_context → rejected）
- [ ] Context Ref 是否包含租户/权限/来源/解析版本/质量评分/保留策略
- [ ] 实时语音是否有 VAD/打断/取消/重连/降级路径
- [ ] 高风险语音动作是否转为确认卡片（含转写文本/置信度/确认方式）
- [ ] 多模态审计是否回答五问：谁上传/说了什么→如何解析→哪些进入上下文→触发哪些工具→确认哪些高风险动作
- [ ] 隐私策略是否在上传入口生效（人脸/工牌/客户信息检测）
- [ ] 源文件删除后派生产物（OCR文本/embedding/摘要/缓存）是否同步清理
- [ ] 准入数据集是否覆盖三类样本（普通/边界/对抗）
- [ ] 解析失败是否有补救流程（重新选择页码/修正字段/框选区域/改用文本）

## 跨Part关系

| 相关 Part/章节 | 关系 | 说明 |
|---------------|------|------|
| Ch19 OCR、Ch20 知识库 | 上游能力 | 多模态文件解析后进入知识库的切分、索引和引用链路 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch49-agent.md — "稳定知识库则不应每次上传即问，而应进入第19章和第20章的入库、切分、索引和引用链路"] |
| Ch22 Runtime | 核心依赖 | 对话 UI 接收 Runtime 事件；取消/重试/审批状态都依赖 Runtime 状态；前端 reducer 折叠 Runtime 事件 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch47-ui.md — "最终状态必须以 Runtime 事件为准"] |
| Ch23 工具契约 | 关键配合 | Ch48 的渲染契约依赖工具 Schema；Tool Registry 返回结构化结果决定渲染组件 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch48-generative-ui.md — "工具返回什么结构，前端允许渲染什么组件，三者必须匹配"] |
| Ch30 人工审批 (HITL) | 流程联动 | Ch48 审批卡片 + 策略绑定直接对接 HITL 流程；审批状态回写 Run 和 Trace |
| Ch36 表达层 | 产物传递 | Ch48 的 Artifact 生成后进入表达层，报告/图表/代码需要格式化和版本化 |
| Ch38 Trace、Ch39 评测 | 下游消费 | 三章的事件全部进入 Trace；用户点踩/展开证据/修改报告/重试问题形成评测样本；前端可观测四维指标（体验/可靠性/质量/治理）回流评测系统 |
| Ch50 安全 | 持续约束 | Ch48 组件白名单/XSS 防护/审批绕过治理 + Ch49 多模态隐私/脱敏/留存，都是安全边界的前端部分 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch48-generative-ui.md — 小结中明确与 Ch50 安全边界的配合] |
| Ch52 合规 | 证据输出 | Ch48 Artifact 证据链 + Ch49 多模态审计记录，进入合规证据包 |

## 开放问题

1. **组件白名单的动态扩展机制**：平台应如何安全地在生产环境中扩展受控组件？新增图表类型、审批卡片变体或表单控件时，更新流程和兼容策略是什么？[INFERRED: 综合 Ch48 组件治理讨论]
2. **Artifact 跨工作区协作**：多用户协作编辑 Artifact 时，Agent 如何参与协作而不覆盖人工复核段落？锁机制和冲突解决策略需要进一步定义 [INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch48-generative-ui.md — §48.14 多人协作与冲突处理]
3. **多模态成本和证据质量的平衡点**：OCR/ASR/视觉模型调用成本高，按风险分层的解析策略（低风险轻量解析 vs 高风险字段级证据）需要在具体业务中量化 [INFERRED: 综合 Ch49 成本控制和证据治理讨论]
4. **前端事件协议在后端的版本演进**：事件协议如何像后端 API 一样管理版本、兼容窗口和废弃字段？[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch47-ui.md — §47.8 事件协议版本复审]
5. **实时语音确认卡片的最佳交互模式**：语音场景中如何在低延迟体验和高风险动作确认之间取得平衡？确认卡片的时机、形式和降级路径尚未标准化 [INFERRED: Ch49 §49.15 提出确认边界但未给出具体交互规范]
6. **多端入口一致性治理的深度**：Web/移动端/IM 工具/嵌入页共享同一 Runtime 时，核心事件语义一致性和幂等规则如何跨端强制执行？[INFERRED: https://raw.githubusercontent.com/datagallery-lab/enterprise_agent_platform/e5d97a66c467045413fff692d28fd2c401503029/docs/part09-frontend-multimodal/ch/ch47-ui.md — §47.14 多端入口一致性治理]

## suggested_tags

`enterprise-agent-platform`, `frontend`, `conversational-ui`, `streaming`, `sse`, `generative-ui`, `artifact`, `component-whitelist`, `render-contract`, `multimodal`, `voice-agent`, `webrtc`, `context-ref`, `async-parsing`, `observability`, `trace`, `approval-flow`, `security`, `data-ref`, `event-reducer`

## source_summary

Part IX "前端、交互与多模态" 的三章从企业 Agent 前端工程视角出发，系统阐述了对话 UI 从"聊天窗口"到"任务工作台"的演化路径。Ch47 建立事件流基础：定义七类界面单元，选择 SSE 作为文本 DataAgent 默认协议，设计三层消息模型和流式状态机，并明确"可演进的事件协议 + 受控组件 + 可观测链路"的架构方向。Ch48 在事件流上叠加可操作对象：通过受控组件白名单将工具结果映射为图表/表格/表单/Artifact/审批卡片，建立 Artifact 独立生命周期和版本治理，强调"模型表达意图，平台决定执行"。Ch49 扩展输入模态并前移治理：文件走异步解析+Context Ref 链路，语音走 WebRTC 实时交互+确认卡片，所有输入材料在进入模型上下文前必须完成权限/质量/来源/审计治理。三章共享核心命题：**所有交互事件进入同一条 Trace，前端行为和后端 Run 状态可排序、可恢复、可审计**——这是企业 Agent 平台从"能对话"走向"能负责任地执行任务"的关键工程基础。

---
