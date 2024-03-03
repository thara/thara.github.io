---
title: RustでWalker's alias method実装した
date: '2016-08-29 15:26:00'
published: '2016-08-29'
tags: [rust]
---

[Cargoで公開した](https://crates.io/crates/aliasmethod)。

実装的には、加減算したi64型の値をvectorの添字にするときにusizeにキャストするのが非常にめんどくさかった。   
usizeのまま加減算できないものか・・・

手続き型で書いているので、もう少し関数型プログラミングスタイルで書きたいところ。特にalias tableを生成している部分。


[creates.io](https://creates.io)の各パッケージのページからはダウンロード数の推移をみることができる。   
このような簡単なライブラリでもダウンロードされているとちょっとうれしい。

ドキュメントは・・・もっとちゃんと書かなきゃなぁ・・・

