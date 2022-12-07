---
title: 名古屋Ruby会議03に行ってきた
---

[名古屋Ruby会議03](http://regional.rubykaigi.org/nagoya03/) に行ってきた。   
最初に謝っておかなければならない。

# **僕は普段Pythonをメインに触っているフレンズです。**

itamaeぐらいにしかRuby使ってない... が、それでも予想以上に楽しめた。   
メモを取っていないので、まとまった文章で感想を述べることはできないが、各発表の感想を形に残しておきたい。   
（LT & パネルディスカッションはパス）


### Ruby/Railsはじめてチームの力をメキメキつけた！
by 小芝 敏明さん

非常に身につまされる話であった。
だいたいはできている、と思う...

「テストの失敗をすぐ直す」あたりは、自分のチームでも割りと先延ばしになりがちなので気をつけたいところ。   
「ペア作業」、おそらくペアプログラミングだけを指しているのではないと思う。   
一人に任せるとベクトルが間違っているときに気づきづらいので、体制でカバーする仕組みがあるのは素晴らしい。   


### ぺろぺろ: Github pull request bot framework
by mzpさん

[mzpぺろぺろ](https://github.com/mzp/prpr)

自分のチームで使おうと思ってて、使えてなかったやつ。   
trelloとgemfileプラグイン以外は、実際に使ってみたい。

### 招待講演 Apache ArrowのRubyバインディングをGObject Introspectionで
by 須藤 功平さん

Pythonistaな自分からしたら、全部Pythonにしたらいいのではという印象を持ったが、GObject Introspectionの手軽さには驚いた。   
Rroongaを使っていたが、Pyroongaではだめだったのだろうか。Cythonで書き直したいのだろうか...    

国産全文検索エンジンGroongaは触ったことはないが、自作WebサービスのElasticSearchはGroongaに置き換えるのも面白そうだ。   
明日の[Groonga Meatup名古屋2017](https://misoca.doorkeeper.jp/events/56673) には参加できそうにない。残念。


### mrubyのJIT
by 三浦 英樹さん

ここからmrubyのターンである。   

本業が水道屋さん、というのが驚き。実は自分も昔少し水道屋っぽいことをしていたことを思い出した。   
恥ずかしながら、[Tracing JIT](https://ja.wikipedia.org/wiki/%E3%83%88%E3%83%AC%E3%83%BC%E3%82%B7%E3%83%B3%E3%82%B0%E5%AE%9F%E8%A1%8C%E6%99%82%E3%82%B3%E3%83%B3%E3%83%91%E3%82%A4%E3%83%AB)を初めて知った。   
LuaJITやPyPyのJITの実装がTracing JITとのこと。 勉強になりました。 

### Dynamic certificate internals with ngx_mruby
by 奥村 晃弘さん

ペパボがmrubyをガンガン使っているの事実に驚いた。   
最近話題に上がることが多いと思っていたら、まさかプロダクションに使われていたとは。   
nginxの設定ファイルにmruby書けるのは楽しそうだ。 というかnginxの設定全体をmrubyで書きたい。   

### 未来のサーバ基盤への Haconiwa/mruby の関わり - コンテナ仮想化のその先へ
by 近藤 宇智朗さん

[発表スライド](https://speakerdeck.com/udzura/haconiwa-and-future-os)

haconiwa、コンテナにフック入れられるのが面白い。   
FastContainerは、Docker on Google App Engineのようなものを目指している、という印象。   
mrubyで実装されているのは面白みある。   

### Ruby で TensorFlow by antimon2さん

セグフォで落ちるのは辛い・・・
TensorFlow、どこかでチャレンジしなければ・・・

### Fight with growing data on Rails
by joker1007さん

おそらく自分の仕事上での問題領域に最も近い発表であったと思う。   

中規模以上のデータ分析を扱うならActive Recordは捨てろ、というのが印象的。   
O/Rマッパーは、オブジェクトのCRUDには有効だが、複雑なSQLクエリには適応しないほうがよい、という点で自分も同意見である。   
BigQueryはやはり偉大。   
仕事ではAirflowを使っているが、中間データの生成には使っていない。もし必要となるとしたら、Python製のluigiも検討したい。   
データのクレンジングにはEmbulk。タスクに冪等性を持たせる件は見習いたい。   

「Rubyに頼るところ、そうでないところを明確にする。アプリ基盤がRubyでも剥がせる余地を残しておく。自分が抱えている課題のステージを見極め、適切な道具を選ぶ。」   
大事。


### 招待講演 re: rinda
by 関 将俊さん

並列処理の協調言語Linda。そのRuby版のRinda。双方とも初めて聞いた。   
とてもアイディアが面白い。 特にパターンマッチでwaitするところ。   
n次元の内積問題という例はわかりやすかった。


## 総評

初めてRuby会議に参加したが、普段Rubyをガッツリ触っていない自分でも楽しめた。   
Rubyコミュニティの良さがわかった気がする。自分も何かで貢献したい。
