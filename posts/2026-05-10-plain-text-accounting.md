---
title: plain-text-accounting
date: '2026-05-10'
published: '2026-05-10'
---

資産管理アプリとしてマネーフォワードを使っているんだけど、[最近のアレコレ](https://www.itmedia.co.jp/news/articles/2605/07/news096.html) があって、代替手段を探している。

ネットバンクやクレカのパスワードを第三者に預けたくない、という気持ちが強かったのでDeep Researchで調べていたが、今はそのようなサービスが存在しない、という事実がわかり絶望した。

セキュリティを犠牲にせずに資産管理をするには、自前でやるしかなさそう、ということで [Plain Text Accounting](https://plaintextaccounting.org/) というアプローチがあることを知ったので、それにチャレンジしてみたい。

Plain Text Accountingは、plain textで複式簿記の台帳を書く、というアプローチ。

[Ledger](https://ledger-cli.org/doc/ledger3.html) を紐切りに、Haskell製の[hledger](https://hledger.org/1.52/hledger.html) や [Beancount](https://github.com/beancount/beancount) などのツールがある。

Beancountを選択する人が多いっぽいのだけれど、CLIコマンドの体系があまり好みではなかったのでやめた。
Ledgerでも良かったのだけれど、機能が絞り込まれていてとっかかりとして使いやすそう、Goで書かれていて内部の実装が読みやすい、マニュアルが読みやすい、などの理由から[howeyc/ledger](https://github.com/howeyc/ledger) を使ってみようと思う。

自分が使っているネットバンクやクレカ会社はCSVダウンロードをサポートしているので、howeyc/ledgerのimport機能がサポートしているフォーマットに変換するスクリプト用意すれば良さそう。

手で集めるのは流石に手間だけど、月1回に支出を振り返るタイミングでちょっと手間をかければいいだけなので、まぁ許容できるかな...
ChatGPT AgentとかClaude Computer useとか使ってもいいかもしれない。（ちょっと抵抗はあるが...）

仮にCSVダウンロードできない金融機関があったとしても、コピペしてLLMエージェントに渡せばいいし、LLMエージェント自体にGUIを直接操作させてもいいので、だいぶ省力化できる余地が広がった感がある。

plain textに台帳あれば、読み込んで分析したりグラフ作ったりもやろうと思えばできるし、LLMコーディングエージェントのおかげで片手間でも作れそう。

資産管理アプリのベンダーロックインをこういった形で解決できるのはいいね。

代わりにLLMコーディングエージェントにロックインされているだけかもしれんけど...
