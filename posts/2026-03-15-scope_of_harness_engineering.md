---
title: Harness Engineeringのスコープ
date: '2026-03-15'
published: '2026-03-15'
---

[Claude Code / Codex ユーザーのための誰でもわかるHarness Engineeringベストプラクティス](https://nyosegawa.github.io/posts/harness-engineering-best-practices-2026/#harness-engineering%E3%81%A8%E3%81%AF%E4%BD%95%E3%81%8B) を読んで。

Harness Engineeringはコーディングエージェントに適用されるものだけど、自分は "Harness" はコーディングエージェントそのものよりエージェントが作り出したソフトウェアに対しても適用すべきだと考えている。

先の記事では主に、以下が "Harness" として挙げられていた。

- AGENTS.md/CLAUDE.md
- ADRなどのドキュメント
- リンター/型チェッカー/テストスイートなどの決定論的ツール
- Hooksによる品質フィードバックループ
- E2Eテスト

これに対して、自分の考える "Harness" は、かなり抽象的だが以下のようなもの。

- 現実世界とのインタフェースの規約
- 業界標準のセキュリティ基準（PCI DSSなど）
- コスト
- コンプライアンス
- KPI
- SLO
- 監査統制（ISMAPなど）

などなど。

これらをエージェントに対して自然言語などの曖昧さがある表現ではなく、先の記事中に述べられているような「決定論的なツール」で制約を加えたい。例えば、CELとかRegoとかOpen Policy Agentとか。

つまり自分が実現したい "Harness" とは「ソフトウェアの外側から課される制約を、決定論的なツールで表現したもの」と言えそう。

先の記事では、

> Harness Engineeringが数カ月後にはとくに重要でない分野になっている可能性はあります。Coding Agent自体に還元され、個々の開発者や組織が意識しなくても良い状態になっているかもしれませんし、LLMの能力が向上しHarness群自体が(あるいはHarness群の一部かもしれませんが)不要になっている可能性もあります。

と述べられていたけれど、自分の考える "Harness" はコーディングエージェントとシステムが同一視される未来になったとしても有用だし、そもそもAI関係なく人間だけで開発されたシステムにも活用できる。

自分の考える "Harness" から見れば、コーディングエージェントだろうがまだ見ぬAGIだろうが人間だろうが、曖昧性があるという点は変わらない。

---

さて、この自分なりの "Harness" はもはや "Harness" = 引き具、馬具、とは言いづらい。

以前の記事で書いた [Protocol-Driven Development](./protocol-driven-development.html) は対象がソースコード自体になっていたが、考え方としては同じだ。

well-definedな点も同様なので、今まで述べてきた「自分が実現したい "Harness"」は "Protocol" と呼称するようにしたい。

つまり、自分が真に実現したいのは "Protocol" の設計と運用を中心に据えた、AIの発展状況に依らず現実問題への影響を統制するソフトウェア開発論 **Protocol Engineering** ということ。

この Protocol Engineering についてはまた、後でじっくりとまとめたい。
