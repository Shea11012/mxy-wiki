---
title: "MCP协议与自定义服务器"
created: 2026-06-11
updated: 2026-06-11
status: stable
type: concept
tags: [agent/mcp, agent/hello-agents]
source: "raw/2026/06/10/第十章 智能体通信协议.md"
confidence: EXTRACTED
contested: false
contradictions: []
related:
  - "智能体通信协议（总览）"
  - "A2A协议"
  - "ANP协议"
---

# MCP协议与自定义服务器

> 摘录：MCP（Model Context Protocol）由 Anthropic 团队提出 <sup>[1]</sup> ，其核心设计理念是 **标准化智能体与外部工具/资源的通信方式** 。

## MCP 架构

MCP 协议采用 Host、Client、Servers 三层架构：

- **Host（宿主层）**：用户直接交互的界面（如 Claude Desktop），管理整个对话流程
- **Client（客户端层）**：Host 内置的 MCP Client，负责与 Server 建立连接和通信
- **Server（服务器层）**：执行实际操作的端点（如文件系统、数据库、GitHub）

> 摘录：这种架构设计的优势在于 **关注点分离** ：Host 专注于用户体验，Client 专注于协议通信，Server 专注于具体功能实现。

## 三大核心能力

| 能力 | 类比 | 说明 |
|------|------|------|
| **Tools**（工具） | API 端点 | 主动执行操作（读文件、查数据库、发消息），由 LLM 决定是否调用 |
| **Resources**（资源） | GET 端点 | 被动提供数据（文件内容、代码库结构、提交历史），按需读取 |
| **Prompts**（提示模板） | 模板 | 预定义的提示模板，引导 LLM 处理特定类型的任务 |

> 摘录：这三种能力的区别在于： **Tools 是主动的** （执行操作）， **Resources 是被动的** （提供数据）， **Prompts 是指导性的** （提供模板）。

## MCP 工作流程

完整的 5 步流程：工具发现 → 上下文构建 → 模型推理 → 工具执行 → 结果整合。整个过程完全自动化，LLM 根据工具描述质量决定是否使用以及如何使用工具。

> 摘录：这个过程是 **完全自动化** 的，LLM 会根据工具描述的质量来决定是否使用以及如何使用工具。

## 五种传输方式

MCP 协议核心特性是**传输层无关性**。HelloAgents 基于 FastMCP 2.0 支持：

1. **Memory** — 内存传输，用于测试和演示
2. **Stdio** — 标准输入输出传输，本地开发首选
3. **Stdio with Args** — 带额外参数的命令传输
4. **npx Stdio** — 启动社区 MCP 服务器
5. **HTTP/SSE/StreamableHTTP** — 远程传输，建议直接使用 `MCPClient`

> 摘录：MCP 协议的一个重要特性是 **传输层无关性** （Transport Agnostic）。

各传输方式的详细使用示例见 [[raw/2026/06/10/第十章 智能体通信协议#10.2.3 MCP 传输方式详解|源文]]。

## MCP 与 Function Calling 的差异

| 维度 | Function Calling | MCP |
|------|-----------------|-----|
| 集成方式 | 代码内嵌（JSON Schema） | 独立协议（网络通信） |
| 跨模型 | 各厂商实现不同 | 统一规范 |
| 动态发现 | 不支持 | 支持（list_tools） |
| 资源/Prompts | 仅有 Tools | 全面（Tools/Resources/Prompts） |

## 构建自定义 MCP 服务器

主要动机：封装业务逻辑、访问私有数据、性能专项优化、功能定制扩展。

以天气查询 MCP 服务器为例，完整实现详见 [[raw/2026/06/10/第十章 智能体通信协议#10.5.1 创建你的第一个 MCP 服务器|源文]]。

**上传社区**：开发完成后可发布到 npm（Node.js 实现）或 PyPI（Python 实现），或提交到 [Awesome MCP Servers](https://github.com/punkpeye/awesome-mcp-servers) 列表。

> 摘录：将企业内部特有的业务流程或复杂操作封装为标准化的 MCP 工具，供智能体统一调用。
