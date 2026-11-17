---
author: Avi Chawla & Akshay Pachaur
tags:
  - agent
  - ai
created: 2026-06-07
---

### AI Agent的定义
AI Agent是能够自主推理、思考、规划、确定相关来源并在需要时从中提取信息、采取行动，甚至在出现问题时自我纠正的自治系统。

### 与LLM和RAG的区别
| 概念 | 描述 |
|------|------|
| **LLM (大型语言模型)** | 相当于"大脑"，基于海量文本数据训练，能推理、生成、总结，但仅限于其训练数据内的知识，无法主动访问外部信息或调用API |
| **RAG (检索增强生成)** | 相当于给"大脑"提供新鲜信息，通过从向量数据库、搜索引擎等外部源检索文档并将其作为上下文输入LLM来增强其能力，使LLM无需重新训练即可了解更新的相关信息 |
| **AI Agent** | 相当于"决策者"，利用LLM作为大脑，结合工具（如API、数据库）进行规划和行动，不仅能回答问题，还能决定采取哪些步骤，实现工作流程的自动化 |

### AI Agent的核心构建模块
AI Agent有六个关键构建模块：
1. **角色扮演 (Role-playing)**: 为Agent赋予明确具体的角色，使其输出更精准相关
2. **专注/任务 (Focus/Tasks)**: 聚焦特定任务，避免过载导致的混乱和不一致
3. **工具 (Tools)**: 允许Agent使用适当的工具（如网络搜索、代码执行）来增强能力
4. **协作 (Cooperation)**: 多个专业化的Agent协同工作，交换反馈，提高结果质量
5. **护栏 (Guardrails)**: 设置约束条件，防止Agent偏离轨道、产生幻觉或做出错误决策
6. **记忆 (Memory)**: 使Agent能够记住过去交互的上下文，实现持续学习和更连贯的响应

## 3. 五大Agentic AI设计模式

| 编号  | 设计模式                         | 特点与应用                                                                           |
| --- | ---------------------------- | ------------------------------------------------------------------------------- |
| 1   | 反思模式 (Reflection pattern)    | Agent审查自身工作以发现错误并迭代，直到产生最终响应。适用于需要自我验证和改进的任务                                    |
| 2   | 工具使用模式 (Tool use pattern)    | Agent调用外部工具（如查询向量数据库、执行Python脚本、调用API）来获取更多信息，不再仅依赖内部知识                         |
| 3   | ReAct (Reason and Act) 模式    | 结合反思和工具使用模式，Agent在Thought→Action→Observation循环中操作，类似于人类解决问题的方式。CrewAI等框架默认使用此模式 |
| 4   | 规划模式 (Planning pattern)      | AI创建路线图，通过细分任务和概述目标来战略性地解决问题，而不是一次性解决                                           |
| 5   | 多智能体模式 (Multi-Agent pattern) | 多个具有特定角色和任务的Agent协同工作，必要时可将任务委派给其他Agent，共同交付最终成果                                |

## 4. 五个级别的Agentic AI系统

| 等级 | 名称 | 特点描述 |
|------|-----|---------|
| 1 | 基础响应器 (Basic Responder) | 人类引导整个流程，LLM只是接收输入并产生输出的通用响应器，对程序流几乎没有控制权 |
| 2 | 路由模式 (Router Pattern) | 人类定义存在的路径/功能，LLM根据情况做出基本决策，选择应采取哪个函数或路径 |
| 3 | 工具调用 (Tool Calling) | 人类定义LLM可访问的一组工具以完成任务，LLM决定何时使用这些工具以及执行参数 |
| 4 | 多智能体模式 (Multi-agent Pattern) | 一个管理Agent协调多个子Agent，并迭代地决定下一步行动。人类设定代理间的层次结构、角色、工具等 |
| 5 | 自主模式 (Autonomous Pattern) | 最高级模式，LLM独立生成和执行新代码，有效地充当独立的AI开发者 |

## 5. 12个AI Agents项目案例

以下表格总结了12个AI Agents项目案例的技术栈、工作流程和核心功能：

| 编号 | 项目名称 | 技术栈 | 工作流程 | 核心功能 |
|------|--------|-------|--------|---------|
| 1 | Agentic RAG | CrewAI, Firecrawl, Qdrant, LightningAI's LitServe, Qwen 3 | 用户查询 → Retriever Agent调用Firecrawl网络搜索或向量数据库工具获取上下文 → Writer Agent生成响应 | 动态从不同来源（向量数据库和互联网）获取上下文的RAG管道 |
| 2 | Voice RAG Agent | CartesiaAI, AssemblyAI, LlamaIndex, Livekit, Gemma 3 | 实时音频输入 → AssemblyAI转录 → LlamaIndex基于文档检索上下文 → Cartesia合成语音输出 | 实时语音交互，基于文档库进行问答并语音回复 |
| 3 | Multi-agent Flight finder | CrewAI, BrowserbaseHQ's headless browser tool, Ollama, DeepSeek-R1 | 解析自然语言查询生成Kayak搜索URL → 浏览器工具提取航班信息 → Summarization Agent汇总信息 | 解析自然语言查询，抓取Kayak实时航班数据并汇总最佳选项 |
| 4 | Financial Analyst | CrewAI, Ollama, DeepSeek-R1, Cursor MCP | 用户提交查询 → MCP启动财务分析团队 → 研究并创建可执行脚本 → 运行脚本生成分析图表 | 获取、分析并生成股票市场趋势洞察，生成可视化图表 |
| 5 | Brand Monitoring System | Bright Data, CrewAI, Ollama, DeepSeek-R1 | 使用Bright Data抓取X、Instagram、YouTube、网站等平台的品牌提及 → 平台特定团队分析数据 → 合并所有洞察生成最终报告 | 跨多个平台抓取品牌提及，分析并生成综合洞察报告 |
| 6 | Multi-agent Hotel Finder | CrewAI, Browserbase's headless browser tool, Ollama, DeepSeek-R1 | 解析旅行查询 → 生成Kayak酒店搜索URL → 浏览器工具提取酒店信息 → Summarizer Agent汇总最佳选项 | 解析旅行查询，抓取Kayak实时酒店数据并汇总最佳住宿方案 |
| 7 | Multi-agent Deep Researcher | Linkup platform, CrewAI, Ollama, DeepSeek-R1, Cursor MCP | 用户提交查询 → Web Search Agent通过Linkup进行深度网络搜索 → Research Analyst验证和去重结果 → Technical Writer撰写带引用的连贯响应 | 构建本地替代ChatGPT深度研究功能的系统，提供详细洞察和引用 |
| 8 | Human-like Memory for Agents | Zep AI, Microsoft AutoGen, Ollama, Qwen 3 | 用户提交查询 → Agent保存对话并提取事实到内存 → 检索事实和摘要 → 利用历史信息生成知情响应 | 为AI Agent构建类似人类的记忆，实现跨会话的个性化和上下文感知响应 |
| 9 | Multi-agent Book Writer | Firecrawl, CrewAI, Ollama, Qwen 3, LightningAI | Outline Crew使用Firecrawl研究书名并决定章节 → 多个Writer Crew并行撰写各章 → 合并所有章节成书 | 从3-5字书名开始，自动研究、规划大纲、并行撰写，最终生成约2万字的书籍 |
| 10 | Multi-agent Content Creation System | Motia, Firecrawl, Ollama, Deepseek-R1, Typefully | 用户提交博客链接 → Firecrawl抓取内容 → X和LinkedIn Agent并行生成社交媒体帖子 → 通过Typefully自动安排发布 | 将任何URL内容转化为社交媒体帖子并自动安排发布的端到端系统 |
| 11 | Documentation Writer Flow | CrewAI, Ollama, DeepSeek-R1 | 用户指定GitHub仓库URL → Planning Crew创建文档计划 → Documentation Crew根据计划编写文档 → 保存到本地目录 | 从GitHub仓库URL自动生成完整的项目文档 |
| 12 | News Generator | Cohere Command R 7B, CrewAI, SerperDevTool | 用户查询 → Research Analyst Agent使用Serper网络搜索工具获取结果并整合 → Content Writer Agent创作精炼的新闻文章 | 接收用户查询，在网上搜索相关信息，并生成带有引文的精心制作的新闻文章 |

## 6. 实用工具和技术栈

书中提到了多种关键工具、框架和技术，主要可分为以下几类：

### Agent框架与编排
- **CrewAI**: 用于多智能体编排和工作流管理的主要框架
- **AutoGen (Microsoft)**: 用于构建多智能体系统的框架，用于实现带记忆的Agent
- **Motia**: 统一的后端框架，用于构建内容创作等工作流

### 数据检索与爬虫工具
- **Firecrawl**: 用于网页搜索和抓取
- **SerperDevTool**: 用于从互联网检索搜索结果
- **Browserbase LoadTool**: 用于运行、管理和监控无头浏览器
- **Linkup platform**: 用于深度网络研究
- **Bright Data**: 用于大规模数据抓取

### 记忆与状态管理
- **Zep AI**: 为AI Agent提供记忆层
- **CrewAI Flows**: 用于管理复杂工作流的状态和持久化

### 语音技术
- **CartesiaAI**: 用于文本到语音合成
- **AssemblyAI**: 用于语音到文本转录

### 部署与基础设施
- **LightningAI's LitServe**: 用于部署AI应用
- **Livekit**: 用于实时音视频通信和编排

### 开发环境集成
- **Model Context Protocol (MCP)**: 用于在Cursor等IDE中集成AI Agent功能

## 7. 主要章节要点

| 章节 | 核心内容总结 |
|------|-------------|
| What is an AI Agent? | 定义了AI Agent为能够自主推理、规划、采取行动并自我纠正的系统，强调其与传统LLM的根本区别在于自主性和行动能力 |
| Agent vs LLM vs RAG | 通过简单类比解释三者关系：LLM是"大脑"，RAG是"喂给大脑的新鲜信息"，而Agent是"使用大脑和工具进行规划和行动的决策者" |
| Building blocks of AI Agents | 详细阐述了构建可靠AI Agent的六个关键原则：角色扮演、专注性、工具使用、协作、护栏和记忆，每个模块都有具体示例说明 |
| 5 Agentic AI Design Patterns | 介绍了五种最流行的AI Agent设计模式：反思、工具使用、ReAct、规划和多智能体模式，解释了每种模式的工作原理和适用场景 |
| 5 Levels of Agentic AI Systems | 提出了从基础到高级的五个AI代理系统级别，帮助理解AI代理能力的演进路径 |
| AI Agents Projects | 通过12个具体的项目案例，展示了如何将理论概念应用于实际场景，涵盖从RAG增强、语音交互到复杂多智能体协作的各种应用 |
| Custom tools via MCP | 介绍了如何创建自定义工具并通过MCP（Model Context Protocol）服务器使其在多个Agent和工作流中复用，实现模块化和可扩展性 |

## 8. 关键实践建议

从书中提取的关键实践建议和最佳实践包括：

- **明确定义角色**: 为Agent分配清晰、具体的角色（如"高级合同律师"），这能显著提升其输出的准确性和专业性。
- **保持专注**: 避免让单个Agent承担过多任务，应使用多个专业化Agent，每个专注于狭窄的职责范围，以减少幻觉并提高准确性。
- **谨慎选择工具**: "更多工具 ≠ 更好结果"，只集成对Agent核心任务真正必要的工具，避免不必要的工具造成混淆和效率降低。
- **促进Agent协作**: 设计工作流使多个专业化的Agent能够协作、交换反馈并相互改进输出，团队合作通常比单一全能Agent效果更好。
- **实施护栏机制**: 为Agent设置约束，如限制工具使用频率、设置验证检查点、建立回退机制，确保其行为符合预期并维持质量标准。
- **利用长期记忆**: 实现短期和长期记忆，使Agent能够记住过去的交互、用户偏好和关键实体，从而提供更连贯和个性化的体验。
- **采用分步模式**: 对于复杂任务，优先考虑使用ReAct（思考-行动-观察）或规划模式，让Agent逐步分解问题并迭代解决。
- **模块化设计**: 创建可复用的自定义工具并通过MCP等协议暴露，便于在不同项目和Agent间共享和集成。
- **本地化部署**: 优先考虑在本地运行LLM（如通过Ollama），以保护数据隐私并完全控制模型和工作流。
- **结构化输出**: 使用Pydantic等工具定义结构化输出模式，确保Agent生成的数据格式正确且易于后续处理。
