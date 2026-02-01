---
title: Protocol-Driven Development
date: '2026-02-01'
published: '2026-02-01'
---

10年弱前にファミコンエミュレータをSwiftで書き始めようと思った時、どこから始めれば良いのか全く分からなかった。

[NESdev Wiki](https://www.nesdev.org/wiki/NES_reference_guide) を見ても内部構造についての説明はあるが、コードを書き始めるとっかかりが見当もつかなかった。
なるべく自分の手で作りたくて、他のエミュレータのコードを読むことも避けていたので、暗礁に乗り上げた。

ファミコンエミュレータは作ったことはなかったが、CPUというものが何かは知っている。そこで、実装をどう進めていけばいいかはわからなかったが、以下のようなprotocolを定義してみた。

```swift
struct OpCode {}
struct Instruction {}

protocol Cpu {
    func fetch() -> OpCode
    func decode(_ opcode: OpCode) -> Instruction
    func execute(_ instruction: Instruction)
}

extension Cpu {
    func run() {
        let opcode = fetch()
        let instruction = decode(opcode)
        execute(instruction)
    }
}
```

via [Cpu.swift](https://github.com/thara/SwiftNES/blob/b309be2ee243a215ec34fc8e7be796dee79ca232/Sources/SwiftNES/CPU/Cpu.swift)

このコードによって動作するものは何もないが、コンパイルは通る。そして、これはCPUの基本的な動作を表現している。

これを書いた時に、目の前がパッと明るくなった気がした。

自分がこれから何を実装していけばいいか、どのように進めていけばいいかがわかった。今までふんわりとしていたCPUエミュレーションを、fetch-decode-executeという3つのステップに分解できたことで、実装も小さなタスクに分割できるようになった。

この体験から、自分が問題をどのように理解しているのかが見えてきた。

### Protocol-Driven Development

ここで言う protocol はSwiftの言語機構に限定されない。
問題の構造と責務を明示する「契約」を指している。struct や enum のようなデータ構造も、ユニットテストも、この意味での protocol に含まれる。

Protocol-Driven Development とは、問題を理解するために protocol を定義し、 その protocol を実装することで問題を解決していく開発手法である。

SwiftにおけるProtocol Oriented Programming(POP)はよく知られているが、この Protocol-Driven Development は自分の造語であり、POPとは全く異なる。

この手法のポイントは、protocol を定義すること自体が問題の理解を深める行為ということ。protocol を定義しコンパイルを通すことで、その問題の構造が明確になり、実装すべき詳細が見えてくる。

この protocol はSwiftの言語機構から出発点としているが、そのprotocolを構成するためのデータ構造(struct や enum)も、この protocol と考えられる。

Haskellなどの静的型付け関数型言語の界隈で言われる "型駆動開発 (Type-Driven Development)" に近い。型駆動開発が型の精緻化によって仕様を表現するのに対し、Protocol-Driven Development はインターフェースの分割と責務の明示に重点を置く。

### Agentic Coding時代への適用

Protocol-Driven Development は、AIアシスタントと協働する Agentic Coding 時代にも有効ではないかと考えている。ただし、これはまだ検証していない仮説に過ぎない。

想定しているのはこういう流れ:

1. 人間が問題を理解するために protocol を定義する
2. その protocol をAIアシスタントに渡し、実装させる
3. コンパイラとテストによって、意図通りの実装かを検証する

人間はコンパイラによって問題の構造や制約が正しいかを検証でき、AIアシスタントは人間が定義した protocol をガイドとして実装を行うことで、意図どおりのコードを生成しやすくなる——という期待。

その言語のコンパイラでは表現できない制約やルールは、ユニットテストやプロパティベースドテストとして protocol に付随させることで、AIアシスタントがフィードバックループを回しやすくなり、より正確な実装が期待できる。少なくとも、意図から大きく外れた実装を早期に検出できる。

このユニットテストやプロパティベースドテストも、仕様を表現する一種の protocol と考えてもいいだろう。

依存型を持つ言語であれば、ユニットテストやプロパティベースドテストの多くを型として表現できるため、より強力な Protocol-Driven Development が可能になるかもしれない。

### 動作環境での動的な検証

ここから発展して、コンパイラによって静的に検証されたprotocol、テストによって動的に検証するprotocolに加えて、更に実際の動作環境でprotocol通りの動作をしているかを検証する仕組みがあっても良いかもしれない。

それはおそらく、可観測性(Observability)の仕組みと結びつくことになるだろう。

SLI/SLOのような概念を取り入れ、protocolに基づいた動作が期待通りに行われているかをモニタリングし、逸脱があればアラートを上げたり、自動修正のトリガーにしたりする仕組みだ。

AIアシスタントは現状自然言語を主な入力としているが、Protocol-Driven Developmentはそこにprotocolというガイドを提供した。さらに、SLI/SLOにも自然言語ではなく検証済みのprotocolを提供することで、AIアシスタントがより正確にシステムの状態を把握し、適切な対応を行えるようになるかもしれない。

### 抽象レイヤーと信頼性レイヤー

[過去の記事](./trying-cline.html) では、"生成AI時代の信頼性設計" として信頼性のレイヤーについて触れた。

Protocol-Driven Development における protocol の抽象レイヤーと、信頼性設計における信頼性レイヤーは一致するのではないか。

テストは対象とするprotocolと同じ抽象レイヤーに基づいて設計される。その抽象レイヤー自体がテストレベルを決定づける。信頼性においても同様で、各抽象レイヤーに応じた信頼性対策を講じるはずだ。SLI/SLOは、最もエンドユーザに近いレイヤーのprotocolに基づいて定義された指標/目標値、ということになる。

テストレベルと信頼性レイヤー、どちらもprotocolを中心に据えることで、システム全体の検証アプローチが一貫する。

今後、AIエージェントと協働してソフトウェアを開発していくことが主流になっていくと考えられるが、そのさらに先、AIエージェントがシステムそのものに組み込まれる時代こそ、この一貫性が重要になるのではないか。

これらはまだ直感の域を出ないが、今後検証していきたい仮説の一つだ。

### まとめ

Protocol-Driven Development は当然ながら完成した手法ではない。

だが、AIと共にソフトウェアを作る時代において、人間が「何を考え、何を信じているのか」をコードとして残す一つの方法にはなりうるのではないかと考えている。
