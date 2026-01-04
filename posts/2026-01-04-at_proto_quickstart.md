---
title: ATProtocol Quick start guideをやった
date: '2026-01-04'
published: '2026-01-04'
---

最近、Xの代わりに [Bluesky](https://bsky.app/) を使ってるのだけれど、あんまりBlueskyを支えるAT Protocolについて知らなかったので、[AT Protocol Quick start guide](https://atproto.com/guides/applications)をやってみた。

[bluesky-social/statusphere-example-app](https://github.com/bluesky-social/statusphere-example-app) のサンプルアプリを動かす感じ。

ATProtocol全然関係ないんだけど、このサンプルアプリが採用している[uhtml](https://www.npmjs.com/package/uhtml) というシンプルなHTMLテンプレートエンジンは初めて知った。あんまりSPAフレームワーク使いたくない場合には良さそう。

## DID

[DID - AT Protocol](https://atproto.com/specs/did)

Decentralized identifierの略で、[W3Cによって標準化](https://www.w3.org/TR/did-1.0/)されているらしい。知らなかった...

標準の文書全部読むのは別の機会にするとして、雑に言うと分散型の識別子で中央集権的な管理者がいないのが特徴らしい。
中身はURIの形式で `did:method:identifier` みたいな感じ。

![A simple example of a decentralized identifier (DID) ](https://www.w3.org/TR/did-1.0/diagrams/parts-of-a-did.svg)   
<sub><sup>via https://www.w3.org/TR/did-1.0/#a-simple-example</sup></sub>

ATProtocolでは `did:web` と `did:plc` の2つのDIDメソッドが使われているけど、前者はDNSと同じ仕組みなのでドメイン名の所有者が移管された時にマイグレーションしなくちゃいけないんだけどそのメカニズムは提供されていないので、後者の `did:plc` を使うのが良さそう？

これが具体的にどうATProtocolやBluesky上で関係しているかはまだよく分かっていない。

[did-method-plc/did-method-plc: Public Ledger of Credentials: a cryptographic, strongly-consistent, and recoverable DID method](https://github.com/did-method-plc/did-method-plc)

## Repository

チュートリアル上では[Repository](https://atproto.com/specs/repository)の詳細はあんまり説明されていなかったけど、
ユーザーの操作結果がRepositoryに保存され、それがfirehose(イベントのアグリゲーターみたいなもの?)というもので他のユーザーにイベントログとして配信される仕組みらしい。

個人的に、このRepositoryがユーザーデータの保存場所として個人の制御化に置ける(んだよね?)点がおもしろいと思っているんだけれど、チュートリアル上ではどこに保存されているかわからないし、Blueskyでも Repository がどこにあるのかよくわかっていない。

とはいえ、イベントログとして配信される仕組みをベースとしたイベント駆動アーキテクチャは面白くて、それを活用したアプリケーションの設計は楽しそう。

## まとめ

大体雰囲気はわかったけれど、まだ DID とか Repository とか firehose とかの詳細がよくわかってないので、もう少し調べてみたい。
