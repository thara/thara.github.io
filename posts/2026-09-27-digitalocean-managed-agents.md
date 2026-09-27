---
title: DigitalOcean Managed Agents
date: '2026-09-27'
published: '2026-09-27'
---

先週、DigitalOceanのManaged AgentsというサービスがPublic Previewになった。

[Introducing DigitalOcean Managed Agents: One AI-native stack to power your intelligence | DigitalOcean](https://www.digitalocean.com/blog/managed-agents-public-preview)

DigitalOcean自体は結構好きなIaaSだけど、AI関連はあまり追っていなかったので、軽くドキュメント読んだ。

DigitalOcean Managed Agentsは、まさに [先週の記事](./ai-control-plane-as-control-architecture.html)で 書いた「AI Control Plane」だった。

DigitalOcean Managed Agentsは以下の2つのコンポーネントで構成される。

- Harness Runtime
- Action Gateway

Harness Runtimeは、FirecrackerでmicroVMを立ち上げて、AIエージェントを隔離された環境で動作させるやつ。
人間のapprove待ちとかのアイドル期間は課金されないのが売り、かな。

Action Gatewayは、サードパーティのAPIをwrapして、認証情報をエージェントから分離したりレートリミット管理やリトライなどをしてくれるやつ。
これは [Control Plane as a Tool: A Scalable Design Pattern for Agentic AI Systems](https://arxiv.org/abs/2505.06817) そのままな感。

これ以外にも、今後追加されるサービスとしては以下のようなものがある。

- Insights: エージェントやツール利用の観測
- Signals: 評価や強化学習のためのフィードバック

[The agent-first cloud: why we built Managed Agents | DigitalOcean](https://www.digitalocean.com/blog/why-we-built-managed-agents) では、AI時代のクラウドとして新たな価値を提供しようと攻勢を取っているのが見て取れる。

今後のDigitalOceanの動向が楽しみ。
