---
title: ElixirでWalker's alias method実装した
date: '2021-10-12 02:00:00'
tags: [elixir]
---

[Hexで公開した](https://hex.pm/packages/alias_method)。   
初Hexパッケージ公開。

Walkes's alias methodは [Rust](https://github.com/thara/rust_aliasmethod) とか [Go](https://github.com/thara/go-aliasmethod) とかですでに書いていて、   
わりと書き慣れているアルゴリズムのはずだったけど関数型プログラミング言語で書くと割と手間取った。

富豪的に書いた方が抽象度上がって全体の見通しが良かった。

あと久しぶりにElixir触った感想。

- map/filter/reduce最高
- 再帰関数の終了条件をpattern matchで書けるのはわかりやすい
- Rubyのモジュールシステムを踏襲しているっぽいけど、ちょっとわかりづらい（自分が慣れてないだけ）
- ExUnitでparameterized testする標準的な方法はないっぽい
- typespec書くのめんどい
    - 書かねば、というモチベーションはあるが、dialyzerで警告出すだけではなくて修正内容を提案してほしい
- hex package公開時にドキュメンテーション必須なのは良い
