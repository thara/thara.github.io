---
title: モバイルアプリのコアにRustを使う
date: '2026-02-08'
published: '2026-02-08'
---

[Next-generation Proton Mail mobile apps: more than meets the eye | Proton](https://proton.me/blog/next-generation-proton-mail-mobile-apps) を読んで。

ここには、マルチプラットフォームに展開するモバイルアプリのメンテナンス性の問題を解決するために、UI部分はAndroid, iOSのplatform nativeな技術で、共通部分はRustを採用することに決めたことが書いてある。
Flutter, React-Native, Kotlin Multiplatformが採用しないと決めた理由も書いてあって、モバイルアプリの実装はほぼ組み込みだと考えている自分には納得性が高かった。

意外だったのは、思ったよりRustが担当する領域の広さ。UIコンポーネントも含まれていた。

以前から複数プラットフォームに展開するモバイルアプリの共通部分にRustを使用する試みがあることは知っていた。

ニッチなユースケースと思って放置してたんだけど、Proton Mailのような自分が使っているアプリがこのアプローチを取っていたのを知ったのと、前述の通り思った以上にRustの領域が広い点に興味が湧いたので雑に試してみているところ。

[thara-playground/rust-core-mobile-app]( https://github.com/thara-playground/rust-core-mobile-app)
(まだrust-coreとios-appの部分しかない)

正直、Counter app程度だと利点はほぼ皆無だけど、どういう技術要素使えば実現できるかはわかってきたかな。

- [mozilla/uniffi-rs](https://github.com/mozilla/uniffi-rs) でbinding生成
- Swift向けにframework生成すれば、Swiftコード上は割とシンプルに使える
    - 個人的にこれがSwiftの強みと感じる

今は、アプリの状態をSwiftで持っていてそれをRustに渡す、みたいなことしてるけど、これをRust内で閉じる良い方法がないかを模索中。
