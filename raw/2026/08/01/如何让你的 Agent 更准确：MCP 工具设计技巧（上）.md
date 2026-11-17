---
source_url: https://bytedance.larkoffice.com/wiki/HeXswIiAni9blLkvB6eclI8mnD0
ingested: 2026-08-01
extracted_to:
  - "MCP 工具设计：让 Agent 更准确"
source_document_id: Kzp1ddmkKo8DmmxIfpicmCEznZg
source_revision: 34
---

# 如何让你的 Agent 更准确：MCP 工具设计技巧（上）

<readonly-block type="isv"></readonly-block>

---

作者：小夏

TRAE 技术专家



越来越多的开发者开始为 AI Agent 开发工具，无论是通过 MCP（Model Context Protocol）、Skills 脚本、还是直接使用 OpenAI / Claude 的 function calling。但很快大家发现了一个令人困惑的现象：**技术上实现完全正确的工具，Agent 却用不好。** 工具能跑通，schema 定义正确，API 调用成功，但是 Agent 总是选错工具，传错参数，或者在明明应该调用工具的时候却回复「我无法完成这个任务」。

**问题出在哪里？问题在于：我们用写 API 的思维在写 Agent 工具。** 

当你为人类设计 API 时，你可以假设他们会阅读文档、理解上下文、在出错后调试代码，但 Agent 不一样：

- 它只能通过工具的名称、描述和参数 schema 来「理解」这个工具能做什么
- 它「试错」的代价很高，每次调用都消耗 token，都可能影响用户体验
- 它需要在可能几十上百个工具中，瞬间做出选择
- 它有一定的试错能力，但是成本高且不稳定

这意味着，给 Agent 开发工具的真正挑战**不是技术实现，而是设计出 Agent 能用好的工具接口**，这也是我们在开发 TRAE 过程中一直在思考和解决的一个问题。

## **核心理念：Agent 工具是 Agent 的用户界面**

**这里有一个关键的思维转换：Agent 工具是 AI Agent 的用户界面（User Interface），不是已有 REST API 的封装**

传统 REST API 是为人类开发者设计的，我们假设开发者会阅读文档、理解上下文、在出错后调试代码，但 Agent 是完全不同的「用户」，它不会主动查阅文档，不擅长从上下文中推断隐含信息，每次调用都需要从头开始理解工具的用途。

**换句话说：你不是在写 API，你是在教会一个智能体如何与这个世界交互。**

这个智能体（LLM）有着独特的长处与局限性：

- 它**擅长**理解自然语言、推理意图、组合信息
- 它**不擅长**精确计算、记住长上下文、从模糊描述中猜测正确参数
- 它**看不到**你的代码实现，只能看到你暴露的 schema 和描述
- 它**只具备有限的上下文**，并且随着上下文被打满工具调用性能会明显下降

<callout emoji="🌳">
只有理解这个智能体的特性，你才能设计出它真正能用好的工具。本文将以 MCP 为主要切入点，因为它正在成为 Agent 工具开发的主流方式，但文中的设计原则适用于所有 Agent 工具开发场景。那接下来让我们先从工具调用的工作机制开始。
</callout>

---

## LLM Tool Calling：完整的调用链路

要设计好 Agent 工具，首先需要理解它是如何被 Agent 调用的，这条调用链路决定了你的设计将如何被「消费」。

### LLM 原生的 Tool Calling 机制

让我们从最底层开始：LLM 本身是如何调用工具的？

**一个关键的认知：LLM 本身不会「执行」任何函数。** 它只做一件事：生成文本，所谓的「function calling」或「tool calling」本质上是 LLM 与应用程序之间的一个多轮对话协议：

```Plain Text
┌─────────────┐                              ┌─────────────┐
│             │  ① 发送请求 + 工具定义         │             │
│             │ ─────────────────────────────▶│             │
│             │                              │             │
│             │  ③ 返回「工具调用请求」         │             │
│   应用程序   │ ◀─────────────────────────────│     LLM     │
│             │                              │             │
│             │  ⑤ 返回工具执行结果            │             │
│             │ ─────────────────────────────▶│             │
│             │                              │             │
│             │  ⑥ 生成最终回复（或继续调用）    │             │
│             │ ◀─────────────────────────────│             │
└──────┬──────┘                              └─────────────┘
       │
       │ ④ 执行实际的函数调用
       ▼
┌─────────────┐
│  外部工具    │
│ (API/DB/..) │
└─────────────┘

② LLM 分析用户请求，决定是否需要调用工具
```

#### 第一步：定义工具

以 OpenAI API 为例，工具通过 `tools` 参数传递给模型。每个工具定义包含三个核心部分：

```Python
tools = [
    {
        "type": "function",
        "name": "get_weather",                    # 工具名称
        "description": "获取指定城市的当前天气",    # 工具描述 - LLM 理解工具用途的关键
        "parameters": {                           # 参数的 JSON Schema
            "type": "object",
            "properties": {
                "location": {
                    "type": "string",
                    "description": "城市名称，如：深圳、北京"
                },
                "unit": {
                    "type": "string",
                    "enum": ["celsius", "fahrenheit"],
                    "description": "温度单位"
                }
            },
            "required": ["location"]
        }
    }
]
```

这三个部分：**name、description、parameters**，就是 LLM「看到」的工具的全部信息。它看不到你的代码实现，不知道函数内部做了什么。

#### 第二步：LLM 决策与返回工具调用

当用户说「深圳今天天气怎么样？」时，LLM 会分析这个请求，发现需要调用 `get_weather` 工具。但它不会执行任何代码，而是返回一个结构化的「工具调用请求」：

```JSON
{
  "id": "fc_12345xyz",
  "type": "function_call",
  "name": "get_weather",
  "arguments": "{\"location\": \"深圳\", \"unit\": \"celsius\"}"
}
```

注意 `arguments` 是一个 JSON 字符串，LLM 本质上只是在「生成文本」，只不过这段文本遵循了特定的结构化格式，**存在返回非法 JSON 格式的可能**。

#### 第三步：应用程序执行函数

应用程序解析 LLM 返回的工具调用请求，执行实际的函数，这里以 Python 代码进行示例：

```Python
import json

# 解析 LLM 返回的工具调用
tool_call = response.output[0]  # 获取第一个工具调用
args = json.loads(tool_call.arguments)

# 执行实际的函数（这是你的代码，不是 LLM 执行的）
weather_result = get_weather(args["location"], args.get("unit", "celsius"))
# 返回: {"temperature": 14, "condition": "晴", "humidity": 65}
```

#### 第四步：将结果返回给 LLM

执行结果需要通过 `function_call_output` 类型的消息返回给 LLM：

```Python
# 将工具执行结果添加到对话中
input_messages.append({
    "type": "function_call_output",
    "call_id": tool_call.call_id,  # 关联到具体的工具调用
    "output": json.dumps(weather_result)
})

# 再次调用 LLM，让它基于结果生成最终回复
final_response = client.responses.create(
    model="gpt-4",
    tools=tools,
    input=input_messages
)
```

#### 第五步：LLM 生成最终回复

LLM 收到工具执行结果后，会生成用户可读的最终回复：

```Plain Text
"深圳今天天气晴朗，当前气温 14°C，湿度 65%。"
```

#### 工具定义如何被 LLM「看到」？

这是一个容易被忽视但非常重要的细节，要理解工具设计的约束，我们需要从 LLM 实现的角度来看工具调用是如何工作的。

**1. JSON 只是中间格式，不是 LLM 真正「看到」的东西**

当你通过 API 传入 JSON 格式的工具定义时，LLM 提供商通常会将其转换为一种内部优化的格式。这是因为 **JSON 对 LLM 来说并不是一个友好的格式**：

- **边界模糊**：JSON 使用 `{}`、`[]`、`"` 等通用符号标记结构，这些符号在普通文本中也会频繁出现，容易产生歧义
- **严格的语法要求**：少一个逗号、多一个引号就会导致解析失败，而 LLM 生成文本时很容易犯这类错误
- **字符串转义的噩梦**：JSON 字符串中的引号需要转义为 `\"`，反斜杠需要转义为 `\\`，换行需要转义为 `\n`。当参数内容包含代码片段时（这在 coding agent 中极为常见），LLM 需要正确处理代码中的所有引号、反斜杠和换行符，这是一个极易出错的环节
- **远距离依赖**：嵌套结构中，匹配的括号可能相隔很远，LLM 需要「记住」开始标记才能正确闭合
- **缺乏显式结束标记**：JSON 只依赖括号匹配，没有像 `</function>` 这样语义明确的结束信号

相比之下，许多 LLM 提供商内部使用**类 XML 的格式**来表示工具调用：

```XML
<function_calls>
<invoke name="get_weather">
<parameter name="location">北京</parameter>
<parameter name="unit">celsius</parameter>
</invoke>
</function_calls>
```

类 XML 格式的优势在于：

- **明确的边界**：开始标签 `<invoke>` 和结束标签 `</invoke>` 清晰地标记了工具调用的范围
- **自描述性**：标签名本身携带语义信息，`<parameter name="location">` 比 `"location":` 更不容易与内容混淆
- **训练数据丰富**：LLM 在预训练时见过大量 HTML/XML 文档，对这种格式更「熟悉」
- **容错性更好**：即使内容中包含类似符号，也不容易与结构标记产生冲突

实际上，不同提供商采用了不同的内部格式和特殊 token。这些格式在模型训练时就被专门优化过，使模型能够更准确地识别「何时应该调用工具」以及「如何正确构造调用参数」。

当 LLM 决定调用工具时，它实际上是在生成一系列遵循特定模式的 token。这些 token 随后被 API 层解析还原成 JSON 格式返回给开发者。这个解析过程本身存在一定的失败率，模型可能生成格式不完整或不合法的输出，导致工具调用失败。

一些 LLM 提供商（如 OpenAI 的 Strict Mode）使用了 **Constrained Decoding** 技术来保证输出一定是合法的 JSON 结构。这种技术在解码时动态限制下一个 token 的候选集，确保生成的序列符合预定义的 schema。但这种约束并非没有代价：它可能影响生成速度，在某些边界情况下也可能影响模型的表达能力。

**2. 工具定义是 System Prompt 的一部分**

工具定义会被注入到 LLM 的 system prompt 中，占用宝贵的 context window。当你定义了 10 个工具，每个工具有详细的描述和参数 schema，这些信息都会序列化后添加到每次请求的 prompt 中，这带来两个重要影响：

- **上下文占用**：工具定义占用的 token 越多，留给实际对话内容的空间就越少。在长对话或需要处理大量代码的场景中，这个问题尤为突出。
- **Prompt Caching**：现代 LLM API 通常会缓存 system prompt 的 KV cache 来加速推理。如果你动态修改工具列表（比如根据用户状态添加或移除工具），就会导致缓存失效，每次请求都需要重新计算，显著增加延迟和成本。

**3. 为什么要用原生的 Function Calling，而不是自己定义格式？**

你可能会想：既然工具定义最终也是放在 prompt 里，我能不能自己在 system prompt 中定义一套工具调用的格式，让 LLM 按照我的格式输出？

技术上可以，之前也有很多 Agent 实现这样做，但效果会差很多。原因在于：**LLM 在训练过程中已经对原生的工具调用格式进行了专门的优化**。模型见过大量使用这种格式的训练数据，对这些特殊 token 有更强的「注意力」和「遵循力」。使用原生格式，模型更容易：

- 准确识别何时应该调用工具
- 正确选择要调用的工具
- 生成符合 schema 的参数

相比之下，自定义格式需要模型「临时学习」你定义的规则，效果和稳定性都会打折扣。这也是为什么现在主流的 Agent 框架都直接使用各 LLM 提供商原生的 function calling 机制，而不是自己发明一套。

**4. 工具数量对模型效果的影响**

工具数量不仅影响 token 消耗，更直接影响模型的决策质量。这是一个容易被低估的问题。

**Token 消耗的具体数据**

每个 MCP 工具都带有 schema，描述它做什么以及如何使用。这些 schema 被注入到 system prompt 中。假设每个工具定义平均约 250-300 tokens：

| 工具数量 | System Prompt 占用 | 对 200K 上下文的影响 |
|-|-|-|
| 5 个 | \~1,500 tokens | 约占 0.75% |
| 10 个 | \~3,000 tokens | 约占 1.5% |
| 20 个 | \~6,000 tokens | 约占 3% |
| 50 个 | \~15,000 tokens | 约占 7.5% |

这些数字看起来不大，但有两个关键问题：

1. **工具定义占据「特权空间」**：System prompt 是上下文中最重要的部分，直接与指令、策略和任务框架竞争位置
2. **每轮对话都存在**：这些定义在每一轮对话中都会被发送，即使本轮根本不需要使用任何工具

**认知过载：模型变「笨」了**

当可用工具较少时（比如 5-10 个），模型可以较容易地「记住」每个工具的用途，做出准确的选择。但当工具数量增加到几十甚至上百个时，问题就变得复杂：

- **选择困难**：模型需要在大量相似的工具中做出选择，容易混淆或选错
- **注意力稀释**：模型对每个工具描述的「关注度」会下降，可能忽略关键细节
- **Prompt 拥挤**：大量工具定义挤占上下文空间，影响模型对用户实际请求的理解

**成本也在增加**

所有这些工具开销都要付费，更长的 prompt、更多的调用、更多的步骤都意味着更高的 API 成本， 不仅变慢了，还在为更差的结果花更多的钱。

OpenAI 官方建议：**尽量将工具数量控制在 20 个以内**。这虽然是一个软性建议，但背后反映的是真实的性能瓶颈。在实践中，如果你的 MCP Server 需要暴露大量功能，就应该慎重考虑一下。

### MCP 的定位：标准化的工具协议层

MCP（Model Context Protocol）并没有改变上述的 tool calling 机制，它解决的是另一个问题：**如何标准化地定义和暴露工具**。

在 MCP 出现之前，如果你想让 Agent 调用外部工具，你需要：

- 为 OpenAI 写一套 function schema
- 为 Anthropic 写一套（格式略有不同）
- 为其他 LLM 再写一套……
- 每个工具、每个 LLM 都要单独适配

这就是经典的 **N×M 问题**：N 个工具 × M 个 LLM = N×M 个适配器。

MCP 的解决方案是引入一个**标准化的中间层**：

```Plain Text
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   LLM A     │     │   LLM B     │     │   LLM C     │
└──────┬──────┘     └──────┬──────┘     └──────┬──────┘
       │                   │                   │
       └───────────────────┼───────────────────┘
                           │
                    MCP Client（适配层）
                           │
                    MCP Protocol（标准协议）
                           │
       ┌───────────────────┼───────────────────┐
       │                   │                   │
┌──────┴──────┐     ┌──────┴──────┐     ┌──────┴──────┐
│ MCP Server  │     │ MCP Server  │     │ MCP Server  │
│  (GitHub)   │     │  (Slack)    │     │  (Database) │
└─────────────┘     └─────────────┘     └─────────────┘
```

MCP Server 只需要按照 MCP 协议暴露工具，MCP Client 负责将这些工具转换成各个 LLM 能理解的格式。

### 从 MCP Tool 到 LLM Tool Call 的转换

当你在 MCP Server 中定义一个工具时，你实际上定义的是：

```TypeScript
// MCP Server 中的工具定义
{
  name: "create_issue",
  description: "Create a new issue in a GitHub repository",
  inputSchema: {
    type: "object",
    properties: {
      repo: { type: "string", description: "Repository name (owner/repo)" },
      title: { type: "string", description: "Issue title" },
      body: { type: "string", description: "Issue body content" }
    },
    required: ["repo", "title"]
  }
}
```

当 Agent 作为 Client 连接到这个 MCP Server 时，会发生以下转换：

1. **工具发现**：Client 通过 MCP 协议获取所有可用工具的列表和 schema
2. **工具合并**：MCP 工具被添加到 Agent 的工具列表中，与 Agent 自带的工具成为「一等公民」
3. **注入 System Prompt**：所有工具信息被添加到发送给 LLM 的请求中
4. **工具调用**：当 LLM 决定调用工具时，Agent 根据工具来源将请求转发给对应的处理程序
5. **结果返回**：执行结果返回给 Agent，再返回给 LLM

#### MCP 工具的命名约定与潜在问题

由于 MCP 工具需要与 Agent 自带的工具共存，大多数 Agent 实现会给 MCP 工具添加前缀以避免命名冲突。常见的命名模式是：

```Plain Text
mcp_<server-name>_<tool-name>
```

例如，一个名为 `github` 的 MCP Server 中的 `create_issue` 工具，最终呈现给 LLM 的名称可能是：

```Plain Text
mcp_github_create_issue
```

这种机制看似简单，但会带来一系列问题：

**1. 与自带工具的冲突或歧义**

如果 Agent 自带了一个 `read_file` 工具，而你的 MCP Server 也暴露了 `read_file`，即使加了前缀变成 `mcp_myserver_read_file`，LLM 仍然可能在两个功能相似的工具之间产生困惑。更糟的是，如果描述不够清晰，LLM 可能会选错工具，或者在应该用自带工具时调用了 MCP 工具（反之亦然）。

**2. 工具名过长与长度限制**

当 server name 和 tool name 都较长时，最终的工具名可能变得冗长：

```Plain Text
mcp_my-awesome-productivity-server_create_calendar_event_with_reminder
```

过长的名称不仅占用更多 token，还可能影响 LLM 对工具的「记忆」和选择准确性。一些研究表明，LLM 对简短、直观的名称有更好的响应。

更重要的是，**某些 MCP Client 有硬性的长度限制**。

例如 TRAE 对 server name + tool name 的总长度限制为 60 个字符。超过限制的工具名会被截断，可能导致工具无法正常工作或产生命名冲突。这意味着在设计 MCP Server 时，你需要为工具名预留足够的「前缀空间」。

**3. 工具数量爆炸**

当用户连接多个 MCP Server 时，工具数量会快速累加。假设：

- Agent 自带 15 个工具
- MCP Server A 暴露 10 个工具
- MCP Server B 暴露 12 个工具
- MCP Server C 暴露 8 个工具

总工具数已达 45 个，超过了前文提到的 20 个工具的建议上限。这还没考虑一些「全功能」的 MCP Server 可能一次性暴露几十个工具的情况。

面对工具数量爆炸的问题，一些 Agent 实现开始探索**动态工具发现**（Dynamic Tool Discovery）机制：不在会话开始时注册所有工具，而是让 LLM 在需要时主动「查询/搜索」可用的工具。

这种方式的优势很明显：大幅减少了 system prompt 中的工具定义数量，避免了上下文污染。但它也有明显的局限：

- **增加调用轮次**：原本一次工具调用能完成的任务，现在可能需要「查询 → 选择 → 调用」多轮交互
- **对 LLM 能力要求更高**：模型需要理解「先查询再使用」这种元认知模式，而不是直接从可用工具中选择
- **当前模型支持有限**：大多数 LLM 对这种动态发现模式的训练还不充分，效果可能不如静态注册稳定

尽管如此，随着 Agent 场景对工具扩展性需求的增长，以及 LLM 在 agentic 能力上的持续进化，动态工具发现很可能成为未来的主流范式。在设计 MCP Server 时，你可以为此做一些准备，比如提供清晰的工具分类和摘要描述，方便未来被动态发现机制索引。

**4. 命名空间污染**

不同的 MCP Server 可能提供功能相似但实现不同的工具。例如

- `mcp_lark_search` - 搜索飞书文档
- `mcp_google-drive_search` - 搜索 Google Drive

对于「帮我搜索 XXX」这样的请求，LLM 需要在多个 search 工具中做选择，而这些工具的描述可能都很相似。

**5. 不同 LLM 厂商 API 的 Schema 兼容性问题**

这是一个容易被忽视但影响很大的问题。前面我们提到，MCP Client 会将工具定义转换为各个 LLM 原生的 tool calling 格式。问题在于：**不同 LLM 厂商的 API 对 JSON Schema 的支持程度差异很大**。

这不是 MCP 协议本身的限制，而是底层 LLM API 的限制。MCP Server 的开发者必须注意这个差异，在不同的 Client 和 LLM 组合下测试

| **LLM 厂商/产品** | **JSON Schema 支持限制** |
|-|-|
| **OpenAI** | 只支持 `anyOf`（不支持 `allOf` 或 `oneOf`）；不能在 schema 根级别使用 `anyOf` |
| **Anthropic** | 支持根级别的 `anyOf` |
| **部分模型** | 可能将对象或数组参数错误地序列化为 JSON 字符串 |

此外，MCP Client 自身和 LLM API 也可能引入额外限制，限制工具的数量、限制 `server name + tool name` 的最大字符数等等。这些差异带来几个实际影响：

- **复杂 Schema 可能失效**：如果你的工具参数使用了 `$ref` 引用、联合类型（`anyOf`/`oneOf`）或递归结构，在某些 LLM 上可能完全无法工作
- **参数解析错误**：有些 LLM 会把本应是对象的参数序列化成 JSON 字符串传递，导致类型不匹配
- **跨平台兼容困难**：理想情况下你的 MCP Server 应该能在所有 Client 上工作，但现实中可能需要针对不同 LLM 提供不同的 schema 变体

**设计建议**：尽量使用简单、扁平的 schema 结构，避免深层嵌套、递归引用和复杂的联合类型。如果必须使用这些特性，要在目标 LLM 和 Client 上充分测试。

**6. LLM 参数传递的不确定性**

即使你的 schema 定义完全正确，LLM 生成的参数也可能出现问题：

- **类型错误**：期望数字却传来字符串，期望数组却传来单个值
- **格式不符**：日期、URL 等特殊格式的字符串可能不符合预期
- **必填字段缺失**：LLM 可能「忘记」传递某些 required 参数
- **额外字段**：LLM 可能传递 schema 中未定义的字段

这意味着你的 MCP Server 实现需要做好**防御性编程**：验证输入、提供合理的默认值、对错误格式做兼容处理，并返回清晰的错误信息帮助 LLM 自我纠正。

### 理解这条链路的意义

为什么要花这么多篇幅讲这条调用链路？因为它直接影响你的设计决策，理解了这条链路，我们接下来就可以开始讨论：如何站在 Agent 的角度，设计它能用好的工具接口，以及在实现和测试的过程中有哪些需要注意的细节。

---

## Agent 如何「看」工具？

在深入具体的设计原则之前，我们需要先建立一个关键的心智模型：**Agent 究竟是如何「感知」和「理解」你设计的工具的？**理解这一点，是设计出好用工具的前提。

### Agent 眼中的工具：三元组

从 Agent（LLM）的视角来看，每个工具就是一个简单的**三元组**：

```Plain Text
工具 = (名称, 描述, 参数 Schema)
```

就这些，没有代码实现，没有注释，没有文档链接， Agent 对工具的全部认知，就来自这三个元素：

**1. 名称（Name）**

```Plain Text
create_github_issue
```

名称是 Agent 对工具的「第一印象」，一个好的名称应该让 Agent 在看到的瞬间就能大致猜到这个工具是做什么的。

**2. 描述（Description）**

```Plain Text
Create a new issue in a GitHub repository. Use this when you need to 
report bugs, request features, or track tasks.
```

描述是工具的「使用说明书」，告诉 Agent 这个工具能做什么、应该在什么场景下使用。

**3. 参数 Schema（Input Schema）**

```JSON
{
  "type": "object",
  "properties": {
    "repo": { 
      "type": "string", 
      "description": "Repository in owner/repo format, e.g. facebook/react" 
    },
    "title": { 
      "type": "string", 
      "description": "Issue title, should be concise and descriptive" 
    },
    "body": { 
      "type": "string", 
      "description": "Issue body in markdown format" 
    }
  },
  "required": ["repo", "title"]
}
```

参数 Schema 告诉 Agent 调用这个工具需要提供什么信息、每个参数是什么含义、哪些是必填的。

**这三个元素就是 Agent 理解工具的全部信息来源。** 如果名称模糊、描述不清、参数含义不明，Agent 就会困惑、选错工具、传错参数。

### Agent 依赖显式语义，而非隐含上下文

这是传统 API 设计与 Agent 工具设计的**根本差异**。

为什么好的 REST API 不等于好的 MCP Server？因为 REST API 的设计原则（可组合性、灵活性、自我发现）对人类开发者和 Agent 的影响完全不同：

| **REST API 设计原则** | **对人类开发者** | **对 Agent** |
|-|-|-|
| **发现性（Discovery）** | 成本低：读一次文档就记住 | 成本高：每次请求都要带上完整 schema |
| **可组合性（Composability）** | 优势：灵活组合多个小端点 | 劣势：多步调用意味着更多 token、更慢的迭代 |
| **灵活性（Flexibility）** | 更多选项 = 更强大 | 更多选项 = 更容易产生幻觉和选错 |
| **文档丰富度** | 可以查阅外部文档补充理解 | 只能依赖工具定义中的描述 |
| **试错迭代** | 成本低：快速调试修复 | 成本高：每次失败都消耗 token 和上下文 |

人类开发者能从上下文推断意图。当我们看到 `get_user(id)` 时，会「显然」地认为 `id` 是用户的唯一标识符（比如 UUID）。但 Agent 没有这种隐含知识，它可能会尝试用邮箱、用户名甚至随机字符串来调用这个函数。下面我们用 Python 代码来举例，函数名可以理解为工具名，Docstring 可以理解为工具的 `description`：

```Python
# ❌ 依赖隐含上下文（人类能理解，Agent 容易误解）
def get_user(id):
    """获取用户信息"""
    pass

# ✅ 显式语义化（Agent 友好）
def get_user_by_uuid(user_uuid: str):
    """
    根据 UUID 获取用户信息。
    
    参数：
    - user_uuid: 用户的唯一标识符，格式为 'usr_xxxxxxxx'
    
    返回：用户信息的 JSON 对象，包含 name、email、created_at 等字段
    """
    pass
```

这个差异贯穿整个工具设计过程：

| **人类开发者** | **Agent** |
|-|-|
| 依赖隐含上下文和经验 | 只能理解显式提供的信息 |
| 看到 `id` 就知道是标识符 | 需要明确说明 `user_uuid` 的格式 |
| 会主动查阅文档补充信息 | 只能依赖工具定义中的描述 |
| 出错后能推理修正 | 出错后可能重复同样的错误 |

因此，设计 Agent 工具的核心原则是：**防呆式语义化**，假设 Agent 会完全按字面意义理解你的工具，不会做任何「显然」的推断。

### Agent 的「试错」成本

人类开发者使用 API 时，通常的流程是：

1. 阅读文档，大致了解
2. 写代码调用，看看返回什么
3. 报错了，看错误信息，调整参数
4. 反复试验，直到成功
5. 记住这个经验，下次直接用

但 Agent 不一样，Agent 需要**尽量一次做对**，原因有几个：

**成本高昂**

每次工具调用都会消耗 token，一次失败的调用不仅浪费了调用本身的 token，还需要额外的 token 来处理错误、重新规划、再次尝试。在复杂任务中，这种「试错」的成本会快速累积。

**用户体验差**

想象一下：用户让 Agent「帮我在 GitHub 上创建一个 issue」，Agent 先调用了错误的工具，然后参数传错了，再然后格式不对…… 用户看着 Agent 反复折腾，体验会非常糟糕。

**上下文污染**

每次失败的尝试都会被记录在对话历史中，占用宝贵的上下文空间。随着失败尝试的累积，真正有用的信息反而被挤出了上下文窗口。

**没有「记忆」跨会话复用**

人类开发者踩过的坑会记住，下次不会再犯。但 Agent 的每次会话都是独立的，上一次学到的「这个工具要这样用」的经验，下一次会话就忘了。

这意味着：**你的工具设计必须让 Agent 在第一次看到时就能正确使用**，不能太指望它「试几次就会了」。

### 上下文窗口：稀缺的认知资源

我们在上面讨论过，工具定义会占用 context window，但这里我想从另一个角度来看这个问题：**上下文窗口就像 Agent 的「工作记忆」**。

人类的工作记忆容量有限（著名的 7±2 法则），LLM 也是如此。虽然现代 LLM 的上下文窗口可以达到 200K 甚至更长，但这并不意味着它能同等质量地「关注」窗口中的每一部分内容。研究表明，LLM 对上下文的注意力分布是不均匀的：

- 开头和结尾的内容通常得到更多关注
- 中间部分的信息更容易被「忽略」
- 当上下文过长时，整体的推理质量会下降

**工具定义与实际任务的竞争**

更重要的是，工具定义占据的是上下文中最「特权」的位置：system prompt。当你往 system prompt 塞入大量工具定义时，模型会开始关注工具选择逻辑，而不是用户的实际意图。如果你的 Agent 开始不遵循指令，问题可能不在模型本身，而在你的工具集。这对工具设计的启示是：

**工具数量要克制**

20 个精心设计的工具，比 100 个随意堆砌的工具效果更好。当工具太多时，Agent 的「注意力」会被稀释，选择准确率会下降。实际上，经验丰富的开发者建议：**给模型 1-5 个精心设计的工具**，而不是 20 个「可能有用」的工具。

**描述要精准而简洁**

冗长的描述不仅占用更多 token，还可能让关键信息淹没在文字海洋中。好的描述应该用最少的词传达最关键的信息。

**参数要必要且充分**

每多一个参数，Agent 就多一份「认知负担」。只暴露真正必要的参数，对于可以有合理默认值的参数，考虑不暴露或标记为可选。

### 关键洞察：好的工具设计 = 减少 Agent 的认知负担

综合以上分析，我们可以得出一个核心洞察：**设计工具的本质，是在设计 Agent 的认知体验。好的工具设计，就是不断减少 Agent 的认知负担。**

具体来说：

| **认知负担** | **设计目标** |
|-|-|
| 理解工具是做什么的 | 名称自解释，描述清晰 |
| 判断什么时候该用这个工具 | 描述中明确使用场景 |
| 知道该传什么参数 | 参数名直观，description 有示例 |
| 避免传错参数格式 | Schema 约束严格，有 enum 就用 enum |
| 理解执行结果 | 输出结构化，包含足够的上下文 |
| 处理错误情况 | 错误信息明确，指出如何修正 |

接下来我们将围绕这个核心洞察，逐一展开具体的设计原则：命名、描述、输入、输出、错误处理。每一个原则的目标都是一样的，**让 Agent 更容易理解、更容易用对、更难用错**。

---

## 如何命名：让工具「自解释」

工具名称是 Agent 对工具的「第一印象」，也是它在几十个工具中快速筛选的主要依据。一个好的名称应该让 Agent 在看到的瞬间就能判断：这个工具是不是我需要的？

### 命名要完整，不要依赖隐含上下文

前面我们强调过，Agent 不会做「显然」的推断。这个原则在命名上尤为重要：

```Python
# ❌ 不好：依赖隐含上下文
send_message      # 发给谁？通过什么渠道？
get_user          # 根据什么获取？返回什么？
delete_item       # 删除什么类型的 item？

# ✅ 好：完整、自解释
slack_send_message           # 明确是 Slack 消息
get_user_by_email           # 明确是通过邮箱查找
delete_project_by_uuid      # 明确是删除项目，通过 UUID
```

命名完整性的几个维度：

| **维度** | **不好** | **好** |
|-|-|-|
| 作用域 | `search` | `github_search_issues` |
| 操作对象 | `delete` | `delete_calendar_event` |
| 操作方式 | `get_user` | `get_user_by_uuid` |
| 数据格式 | `export` | `export_report_as_csv` |

### 动词优先：Action-Oriented 命名

工具本质上是「动作」，命名应该以动词开头，清晰表达这个工具会「做什么」：

```Python
# ✅ 动词优先，清晰表达动作
create_github_issue
send_slack_message  
search_documents
update_user_profile
delete_expired_sessions

# ❌ 名词或模糊命名
github_issue          # 是创建？查询？删除？
slack_message_handler # handler 做什么不清楚
document_search       # 不如 search_documents 直观
```

常用的动词模式：

| **动词** | **语义** | **示例** |
|-|-|-|
| `create` | 创建新资源 | `create_project`, `create_user` |
| `get` | 获取单个资源 | `get_issue_by_id`, `get_current_user` |
| `list` | 获取资源列表 | `list_repositories`, `list_recent_commits` |
| `search` | 搜索/查询 | `search_issues`, `search_documents` |
| `update` | 更新现有资源 | `update_issue_status`, `update_profile` |
| `delete` | 删除资源 | `delete_comment`, `delete_branch` |
| `send` | 发送消息/通知 | `send_email`, `send_slack_message` |
| `run` / `execute` | 执行操作 | `run_query`, `execute_workflow` |

### 命名即分类：帮助 Agent 快速筛选

当 Agent 面对几十个工具时，它需要快速判断哪些工具与当前任务相关。**使用一致的前缀可以帮助 Agent 进行「分类筛选」**：

```Python
# 按服务/领域分组的命名
github_create_issue
github_list_pull_requests
github_merge_pull_request
github_search_code

slack_send_message
slack_list_channels
slack_get_channel_history

calendar_create_event
calendar_list_events
calendar_update_event
```

这种命名模式的好处：

1. **视觉分组**：相关工具在列表中自然聚集
2. **语义关联**：Agent 看到 `github_` 前缀就知道这是 GitHub 相关操作
3. **避免冲突**：不同服务的 `create`、`delete` 不会混淆

### 前缀 vs 后缀：因 LLM 而异

一个有趣的发现是：**选择前缀命名还是后缀命名，对不同 LLM 的工具使用评测有的影响**。

```Python
# 前缀命名风格
github_search_issues
github_create_issue
slack_send_message
slack_list_channels

# 后缀命名风格  
search_issues_github
create_issue_github
send_message_slack
list_channels_slack
```

效果因 LLM 而异，没有绝对的「最佳」选择。Anthropic 的研究发现，在他们的内部工具使用评估中，前缀和后缀的选择会产生可测量的性能差异。

**实践建议**：

- 根据你自己的评估来选择命名方案
- 一旦选定，在整个 Server 中保持一致
- 如果你的工具主要被特定 LLM 使用，可以针对该 LLM 优化

### 长度与清晰度的权衡

前面我们提到，某些 MCP Client 对工具名长度有限制（如 TRAE 的 60 字符限制包含 server name）。这需要在完整性和简洁性之间找到平衡：

```Python
# 太长：可能超出限制，也增加 token 消耗
mcp_productivity_suite_create_calendar_event_with_reminder_and_notification

# 太短：信息不足
create_evt

# 合适：完整但不冗余
calendar_create_event  # 如果需要，reminder 可以是参数而非名称的一部分
```

**实用建议**：

- 保持工具名在 30-50 字符以内
- 使用常见缩写（`repo` 代替 `repository`，`msg` 代替 `message`）但要确保不产生歧义
- 把细节放到参数和描述中，而非全部塞进名称

---

## 描述的艺术：精准的契约

如果说名称是工具的「标题」，那么描述就是工具的「使用手册」。对于 Agent 来说，**描述是理解工具如何使用的主要信息来源**，它会认真「阅读」每一个描述来决定工具的使用方式。

### 描述即上下文：Agent 真的会去读

与人类开发者不同，Agent 不会跳过文档直接看代码。它会仔细分析你写的每一句描述，这意味着：

- 描述中的每个词都可能影响 Agent 的行为
- 遗漏的信息会导致 Agent 猜测（通常猜错）
- 错误的描述比没有描述更糟糕

```Python
# ❌ 描述过于简略
def delete_item(id):
    """删除一个项目"""
    pass

# ✅ 描述完整、语义化
def delete_item_by_uuid(item_uuid: str):
    """
    根据 UUID 永久删除一个项目。
    
    参数：
    - item_uuid: 项目的唯一标识符，格式为 'item_xxxxxxxx'
    
    返回：
    - 成功时返回 "Item deleted successfully"
    - 如果项目不存在，返回描述性错误信息
    
    注意：此操作不可逆，删除前请确认。
    """
    pass
```

### 描述的核心要素

一个好的工具描述应该回答以下问题：

**1. 这个工具做什么？( What)** 

```Plain Text
Create a new issue in a GitHub repository.
```

**2. 什么时候应该使用它？( When)** 

```Plain Text
Use this when you need to report bugs, request features, or track tasks.
```

**3. 有什么限制或前提条件？( Constraints)** 

```Plain Text
Requires authentication. The repository must exist and you must have write access.
```

**4. 会返回什么？( Output)** 

```Plain Text
Returns the created issue object with id, url, and status fields.
```

完整示例：

```Python
{
  "name": "github_create_issue",
  "description": """Create a new issue in a GitHub repository.

Use this when you need to report bugs, request features, or track tasks.
Requires write access to the target repository.

Returns the created issue with id, number, html_url, and state fields.
If the repository doesn't exist or access is denied, returns an error message.""",
  "inputSchema": { ... }
}
```

### 参数描述：示例的价值

参数的 `description` 字段同样重要，**一个好的示例胜过千言万语**：

```JSON
{
  "properties": {
    "repo": {
      "type": "string",
      "description": "Repository in owner/repo format, e.g. 'facebook/react' or 'microsoft/vscode'"
    },
    "labels": {
      "type": "array",
      "items": { "type": "string" },
      "description": "Labels to apply to the issue, e.g. ['bug', 'high-priority']"
    },
    "assignees": {
      "type": "array", 
      "items": { "type": "string" },
      "description": "GitHub usernames to assign, e.g. ['octocat', 'hubot']. Must be valid collaborators."
    }
  }
}
```

示例的作用：

- 明确**格式**：`owner/repo` 而非 `repo` 或完整 URL
- 展示**真实值**：`facebook/react` 比 `<owner>/<repo>` 更直观
- 暗示**边界**：多个示例展示值域范围

### 参数描述的规范

除了示例，参数描述还应该遵循以下规范：

**明确标注必填/可选**

```JSON
{
  "repo": {
    "type": "string",
    "description": "(Required) Repository in owner/repo format"
  },
  "branch": {
    "type": "string", 
    "description": "(Optional) Branch name, defaults to 'main'"
  }
}
```

**说明默认值**

```JSON
{
  "limit": {
    "type": "integer",
    "description": "Maximum results to return (optional, default: 20, max: 100)"
  },
  "format": {
    "type": "string",
    "enum": ["json", "markdown", "text"],
    "description": "Output format (optional, default: 'json')"
  }
}
```

这些信息可以在官方的 MCP Inspector 中查看，帮助 Agent（和开发者）快速理解参数要求。

### 说明失败情况

传统 API 文档往往只描述「成功时会怎样」，但对 Agent 来说，**知道失败时会发生什么同样重要**：

```Python
# ❌ 只描述成功情况
"""
根据用户 ID 获取用户信息。
返回用户的姓名、邮箱和注册时间。
"""

# ✅ 同时描述失败情况
"""
根据用户 ID 获取用户信息。

返回：
- 成功时返回 JSON 对象，包含 name、email、created_at 字段
- 如果用户不存在，返回 "User not found: {id}. Please verify the ID format 
  (should be 'usr_xxx') or try searching by email using find_user_by_email()."
- 如果 ID 格式错误，返回格式说明和正确示例
"""
```

这种描述方式让 Agent 知道：

1. 遇到错误不要惊慌，这是预期内的情况
2. 错误信息本身包含修正指引
3. 有替代方案（`find_user_by_email`）可以尝试

### 引导工具选择顺序

当你有多个功能相似的工具时，可以在描述中**明确指导 Agent 的选择顺序**：

```Python
def get_variable_value(address: str):
    """
    获取指定地址的变量值（推荐首选）。
    
    自动识别变量类型并返回格式化的字符串表示。
    大多数情况下应该优先使用这个函数。
    """
    pass

def read_raw_memory(address: str, size: int):
    """
    读取指定地址的原始内存数据。
    
    ⚠️ 只有当 get_variable_value 失败或需要原始字节时才使用此函数。
    此函数忽略类型信息，返回原始字节数组。
    """
    pass
```

通过在描述中写明「推荐首选」和「只有当 X 失败时才使用」，可以有效引导 Agent 的工具选择策略。



至此，我们建立了 Agent 工具设计的核心认知框架，并深入探讨了命名、描述和输入设计三个最基础的维度。

在下一篇，我们将继续探讨输出设计、错误处理这两个同样关键的维度，以及工具粒度的权衡、跨环境可移植性、Skills 与 MCP 的互补等更高级的实践模式，帮助你将这些原则应用到真实的 MCP Server 开发中。



想学习更多内容，也欢迎直接进入官方线上社区或者直接点击「阅读原文」了解。

官方社区：https://forum.trae.cn/
