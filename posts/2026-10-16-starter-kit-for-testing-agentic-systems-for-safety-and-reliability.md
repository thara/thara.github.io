---
title: Starter Kit for Testing Agentic Systems for Safety and Reliability
date: '2026-10-16'
published: '2026-10-16'
---

10/14に、シンガポールIMDA(情報通信メディア開発局) が、エージェントシステムの安全性・信頼性のテストのためのスターターキットとやらを公開していた。

[Singapore launches toolkit to help firms test AI agents | The Straits Times](https://www.straitstimes.com/singapore/singapore-launches-new-tool-kit-to-help-firms-test-ai-agents-for-safety-and-reliability?utm_source=chatgpt.com)

IMDAが立ち上げたAI Verify Foundationというオープンソースコンソーシアムから、「Starter Kit for Testing Agentic Systems for Safety and Reliability」と題して配布されていた。

[Starter Kit for Testing Agentic Systems for Safety and Reliability - Version 1.0](https://www.imda.gov.sg/assets/603572eb-724a-46cf-ac12-14fe38e1558c.pdf) (PDF)
([Resources – AI Verify Foundation](https://aiverifyfoundation.sg/resources/) に導線がある)

エージェントを使ったシステムのうち、以下の代表的な2つのタイプに対して、どう安全性・信頼性のテストをすれば良いかが記載されている。

- Workflow automation
    - あらかじめ定められた手順に沿ってエージェントがタスクを実行するモデル
- Goal-oriented
    - エージェントが計画を立てて、柔軟にタスクを実行するモデル

(これ以上の柔軟性、自律性を持ったエージェントシステムはスコープ外)

エージェントシステムの安全性・信頼性テスト、といっても基本はソフトウェアのそれと大枠は変わらないが、以下の4ステップのうち特徴的なのは3, 4あたりだろうか。

1. リスクの特定
2. テスト設計
3. テストの実行とログ分析
4. 継続監視・再テスト

エージェントシステムは複数のコンポーネントからなるので、安全性・信頼性テストにおいても、そのコンポーネント単位とE2Eでテストした方が良いとのこと。マルチエージェントシステムの場合においては、個々のエージェントとエージェントクラスタ全体も同様。

ログ(Trajectory っていうらしい)を分析することで、エージェントの行動要因を分析し、ツールやハーネスなどのどこに不具合があったかを特定することができる。また、インシデント発生時のエビデンスにもなる。

また、モデルの更新やツールやハーネスなどの環境によって挙動が変化する可能性があるため、従来の安全性・信頼性テストよりも、定期的に行うことの重要性が高い。
先述のTrajectoryやトレーシングを通じて、可観測性を高く保つことも同様。

しかし、当然、限界はある。

「Starter Kit for Testing Agentic Systems for Safety and Reliability - Version 1.0」では、課題として以下が挙げられている。

- 自律性が高いエージェントでは、事前に全実行経路をテストできない
    - エッジケース、リスクシナリオと障害パターンの組み合わせが膨大
- 大量のTrajectoryを自動評価する技術は、まだ十分に成熟していない
    - 膨大 & 複雑なTrajectoryの評価に課題
- 大規模なテストにコストがかかる
    - 多くの計算リソース、トークン消費量、ツール使用コスト
- 長時間実行するエージェントの挙動の予測が困難
    - 施行を重ね、コンテキストが拡大するにつれてTrajectoryが膨大になり可観測性の難易度も上がる
- 人間による監視の限界
    - レビューコスト、決断疲れ

大枠の考え方は似てるにしても、今までのソフトウェアの安全性・信頼性テストにはなかった課題が少なくはない。

モデルの進化もいいけど、より自律的なエージェントを社会的に実現させるためには、こういう安全性・信頼性とかハーネスとかポリシーとか、そこらへんの問題は避けて通れない。

まさに今、研究が進められている領域で面白い、ね。
