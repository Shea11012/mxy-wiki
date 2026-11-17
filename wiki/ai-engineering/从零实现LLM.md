---
title: "从零实现LLM"
created: 2026-07-15
updated: 2026-07-16
status: stable
type: concept
tags: [ai/llm, ai/transformer, ai/training, ai/moe, ai/sft, ai/rlhf, ai/grpo]
source: "workspace/agents/AI-Engineer-from-scrach/07-llm_from_scrach/"
confidence: SYNTHESIZED
related:
  - "AI工程师从零开始（总览）"
  - "Agent范式详解"
---

# 从零实现 LLM

> 课程覆盖从 decoder-only Transformer 到偏好/RL 后训练的教学链路。它适合建立机制理解，但不等于训练生产级基础模型所需的完整数据、分布式系统和评测体系。

## 一、完整生命周期

```text
数据治理与 tokenizer
 → 预训练（next-token prediction）
 → 评测与 checkpoint 选择
 → SFT / continued pretraining
 → 偏好优化或可验证奖励 RL
 → 安全评测、部署优化与持续监控
```

“预训练 → SFT → Reward Model → PPO/GRPO”只是若干常见路径之一。现代后训练还包括 DPO/IPO/ORPO 等直接偏好优化、蒸馏、拒绝采样、过程/结果监督以及 RLVR；并非所有模型都训练独立 Reward Model，也并非所有对齐都使用 RLHF。

## 二、Decoder-only Transformer

```text
Token IDs → Embedding
 → N × [Norm → Causal Self-Attention → Residual
        Norm → MLP/SwiGLU → Residual]
 → Final Norm → LM Head → Next-token logits
```

核心概念：

- causal mask 保证 token 只能观察当前位置及之前的 token
- multi-head attention 在不同子空间建立 token 交互
- residual path 与 normalization 维持深层优化稳定性
- MLP 通常占大量参数和计算
- 训练目标是对序列的条件概率做最大似然估计

教学中的单头、多头、LayerNorm、GELU 和绝对位置编码是理解起点，不代表所有现代模型的统一实现。

## 三、Tokenization 与数据

字节级 256 词表简单且无 OOV，但序列较长。BPE、Unigram 或 byte-level BPE 在词表大小、压缩率、多语言公平性和代码处理间取舍。

生产训练真正困难的部分包括：

- 来源授权、隐私、版权与可删除性
- 质量过滤、去重、污染和 benchmark leakage
- 语言/领域配比和 curriculum
- tokenizer 训练与 chat template 一致性
- 数据版本、lineage 和可复现采样

仅验证 loss 不足以评价模型；还应测试下游能力、记忆/污染、安全和不同数据切片。

## 四、现代架构组件

| 组件 | 正确定位 |
|------|----------|
| RMSNorm | 不做均值中心化，计算更简单；是否“更稳定”取决于整体架构和训练配置 |
| RoPE | 将位置相关旋转施加到 Q/K，产生相对位置结构；长度外推仍需缩放或训练策略 |
| SwiGLU | 门控 MLP，常与不同中间维度配合，公式不是简单的 `Swish(xW1) ⊙ xW2` 后结束，还需输出投影 |
| GQA/MQA | 减少 KV head，降低 KV cache 内存与带宽；课程原文缺失的重要推理组件 |
| KV Cache | 避免自回归解码重复计算历史 K/V；加速生成但占显存，不是训练加速机制 |
| FlashAttention | 精确 attention 的 IO-aware 实现，降低显存读写；不等于近似/滑窗 attention |
| Sliding Window | 限制局部注意范围，单层复杂度下降；通常与全局层或其它机制组合 |
| Attention Sink | 流式长文本中保留少量初始 token 的 cache 技术，不是普遍必需组件 |

## 五、训练工程

原文“FP16/BF16 前向 + FP32 权重”过于笼统。实践中需要区分：参数、计算、梯度、optimizer states 和 master weights 的 dtype；BF16 通常比 FP16 有更大指数范围，FP16 常配 loss scaling。

规模化训练还包括：

- data/tensor/pipeline/context/expert parallelism
- FSDP/ZeRO 分片参数、梯度和 optimizer state
- activation checkpointing（重计算）与 gradient accumulation（扩大有效 batch）
- fused kernels、FlashAttention、通信计算重叠
- checkpoint 完整状态：模型、optimizer、scheduler、scaler、RNG、data cursor
- MFU、tokens/s、通信、内存、loss spike 和硬件故障监控

Gradient accumulation 不是分布式并行的替代品；checkpointing 一词也应区分“保存训练状态”和“activation checkpointing”。

## 六、MoE

Sparse MoE 通过 router 为每个 token 选择少量专家，在相近激活计算下扩大总参数容量，但“参数多、计算少”不代表系统更便宜。

必须处理：

- top-k routing、capacity 和 token dropping/padding
- load-balancing 或其它 router 正则
- expert parallelism 与 all-to-all 通信
- 热点专家、推理部署、量化和容错
- shared experts、dense/MoE 层布局

比较 MoE 与 Dense 应使用激活参数、每 token FLOPs、通信、内存和实际吞吐，而不只比较总参数量。

## 七、SFT 与偏好学习

SFT 应强调数据质量、chat template、loss mask、长度分布、去重和 eval。并非所有数据都必须只在 assistant response 计算 loss；是否训练 system/user/tool token 取决于目标和模板设计。

偏好学习常见两类：

1. **显式 Reward Model + RL**：成对偏好可用 Bradley–Terry 类目标训练序列奖励，再用 PPO 等优化 policy
2. **直接偏好优化**：DPO 等直接从 chosen/rejected 与 reference policy 优化，无需单独 rollout RL loop

Reward Model 常复用 decoder-only LM 并加 scalar head，不宜笼统写成 Transformer encoder。Reward hacking、长度偏差、标注者分歧和分布外泛化必须单独评测。

## 八、PPO

PPO-based RLHF 通常包含 policy、reference policy、reward model、value function 和 rollout 数据。优化目标由 clipped policy objective、value loss、entropy/KL 控制等组成；“Policy Network = SFT + Value Head”只是某种工程封装，policy 与 value 的参数可共享也可分离。

PPO 的优势是成熟且能利用 token-level credit assignment，代价是组件多、显存/训练稳定性要求高。

## 九、GRPO 及后续修正

GRPO 为同一 prompt 采样一组 responses，用组内奖励统计构造相对 advantage，通常不训练独立 value model，并采用 PPO 风格 clipped objective 与 KL 控制。

原文需纠正：

- “去掉 Value Head，减少一半参数量”不准确。PPO 的 value function 可能共享 backbone，省下的显存/计算取决于实现，不固定为一半
- PPO 不要求“每个 prompt 只生成 1 个”，rollout batch 可以包含多个样本
- KL 在 PPO/GRPO 中的实现位置并非固定，不能用“折叠在 advantage vs 显式 loss”概括所有实现
- 组内标准化在奖励方差接近零时会不稳定，也可能引入难度和长度偏差
- response-level advantage 广播到所有 token，credit assignment 较粗

2025–2026 的 DAPO、Dr.GRPO、GSPO 等工作针对 clipping、长度偏差、零方差组、token/sample importance ratio 和训练稳定性提出修正。因此应把 GRPO 视为仍在快速演进的一类 group-based policy optimization，而不是 PPO 的全面替代。

## 十、评测与部署（原文缺失）

训练完成不等于可用：

- loss/perplexity 与能力 benchmark
- 数据污染和记忆测试
- instruction following、tool use、长上下文
- safety/red-team 和分语言/领域切片
- 多次采样的稳定性和 calibration
- serving：continuous batching、prefix/KV cache、quantization、speculative decoding、并行部署
- 质量—吞吐—延迟—显存的联合权衡

## 核心结论

> 从零实现小模型能理解算法；训练可用 LLM 的难点更多在数据、规模化系统、后训练目标、可靠评测和部署经济性。任何单一架构或 RL 算法都只是完整生命周期中的一个组件。
