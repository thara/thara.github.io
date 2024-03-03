---
title: Trelloボード検索できるAlfred 4 workflow作った
description: Trelloボード検索できるAlfred 4 workflow作った
published: '2019-06-30'
tags:
  - rust
---
仕事で [Trello](https://trello.com) をハードに使っているが、TrelloのページをWebブラウザで開く前にシュッと検索したかった。
検索用のURLがあると良かったのだけれど、そんなものはなかった（ページ上からボード検索あるのに何故？）。

しょうがないので、TrelloボードをインクリメンタルサーチできるAlfred 4 workflowを作った。

[thara/alfred_trello_board_search: Trello Board Search on Alfred 4 workflow](https://github.com/thara/alfred_trello_board_search)

Rustで書いたのには特に意味はない。やることはTrello APIでボード情報を検索して、[Script Input Filter](https://www.alfredapp.com/help/workflows/inputs/script-filter/) が読み込めるJSONのフォーマットを標準出力するだけなので、どんな言語でも大して変わりはないと思う。

サクッと作ったので、「ボード全件取得＋フィルタリングはAlfredにお任せ」という設計ゆえに検索するときに若干待たされる感がある。

キャッシュすればよいのだろうけれど、そもそもボード名でインクリメンタルサーチするよりもよく利用するボードをWebブラウザのブックマークに検索しやすい名前で登録しておいて、AlfredからWebブラウザのブックマーク検索をしたほうが格段に使いやすいことに実際に使ってみてから気づいた。   

自分の環境ではボード名に記号が入っていることが多く、それがボード名をそのままインクリメンタルサーチするというソリューションとは相性が悪いのもあって、「ブックマーク検索でいいじゃん」となってしまい、これ以上改善する気力が薄れてしまった。   

とはいえ、よく使うわけではないボードや新規のボードを検索するにはそこそこ使える。   

もうちょっと自分で使ってみて、頻度が高ければアイコンなりインストール方法なりを整備しようかな。   
どれぐらい使っているか、計測用スクリプトを挟んでみるのも面白いかもしれない。
