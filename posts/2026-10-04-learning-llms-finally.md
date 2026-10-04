---
title: 今更ながらLLMの中身を学習し始めた
date: '2026-10-04'
published: '2026-10-04'
---

見ようと思って忘れてた↓の動画を、この前の休みに見た。

<iframe width="560" height="315" src="https://www.youtube.com/embed/gmBo-C11we4?si=3UQfuynUnDgLV7Gv" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>

(新たなインサイトを得られて良かったけど、全体的な感想はここでは書かない)

冒頭に [Transformer Explainer](https://transformer-explainer.explorable-explanations.com/) と [LLM Architecture Gallery](https://sebastianraschka.com/llm-architecture-gallery/) が紹介されていた。

機械学習のビジュアライズは、以前CourseraのML系のシリーズをやった時に確認していたのだけれど、Transformerに関してはイメージ持ってなかったのでTransformer Explainerでの表現方法が新鮮だった。

それよりも何よりも、LLM Architecture Gallery見てそれぞれのモデルの「アーキテクチャ上の差」を評価したくなった。

正直、モデルの性能差みたいなものはあまり興味がない。
どんどん性能が良いものが出てくるから基本的にはそれを使っていくだけだし、その性能も普段の使い方が非決定的なので体感的に評価できない。評価フレームワークはあると思うが、それが実体験に伴わない。

しかし、アーキテクチャの差は、何かしらの課題があってそれに対する解答の現れ（だと思う）ので、そっちの方が気になる。
そこらへんの認識は、自分の中では、ソフトウェアアーキテクチャとなんら変わりがない。   
（ちなみに、LLMのプロンプトエンジニアリングとか、LLMアプリケーション開発とかも、あんまり興味がない。
LLMを使った自律型エージェント開発なら興味はあるけど。分散システムになると思うし）

というわけで、オライリーから2026年6月に出た [ゼロから作るDeep Learning ❻ ―LLM編](https://www.oreilly.co.jp/books/9784814401611/) を読む...

前に、その前提となる [ゼロから作るDeep Learning](https://www.oreilly.co.jp/books/9784873117584/) を読んでる。
2016年に出たらしいけど、そういえば読んでなかった。
知識としては知ってるけど、自分で書いたことはなかったので良い機会かな。

ちょうどいいタイミングで2026-10-19に [Transformer実践ガイド](https://www.oreilly.co.jp/books/9784814401833/) が出るので、それの理解を一つの目標にしておく。

まぁ今更すぎるんだけど。

普段Coding Harness使いまくってるけど、概要の理解だけどそこまで深く理解しようとしていなかった。
やっと、LLMの中で自分が興味を持てたものが見つかったので、まぁしょうがない。
自分が何を面白いと思うかまではコントロールできないし、面白いと思わないと深く理解できないし。

LLMアーキテクチャの評価にはまだまだ遠いけど、理解はオフロードできないので、じっくり取り組んでいきたい。
