---
title: モバイルアプリのコアにRustを使う(2) - Android対応
date: '2026-02-11'
published: '2026-02-11'
---

[モバイルアプリのコアにRustを使う](./2026-02-08-use-rust-in-mobile-app.md)ではiOSしか対応してなかったのでAndroidにも対応した。

[thara-playground/rust-core-mobile-app](https://github.com/thara-playground/rust-core-mobile-app/tree/main/android-app)

AndroidのNDKまわり、昔cocos2d-xのビルドパフォーマンス改善をした時に触ったんだけど結構大変だった。   
が、Rustのエコシステムが強くてそこら辺はほとんど気にしなくて良い感じになってた。

- [cargo-ndk](https://crates.io/crates/cargo-ndk) でAndroid向けのnativeバイナリを作成
- [mozilla/uniffi-rs](https://github.com/mozilla/uniffi-rs) でbinding生成

Rustからのcallbackも自然に書ける。エコシステムの勝利、って感じ。

関係ないけど、数年前と比べてgradleのKotlinサポートも問題なさそう。
