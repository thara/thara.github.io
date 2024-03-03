---
title: Kotlinでは関数呼び出しの次のlambdaリテラルがパースできない
description: Kotlinにおけるパースできないlambdaリテラルについて調べた
published: '2019-05-22'
tags:
  - kotlin
---
Kotlin Testでテスト書いてて、ブロックスコープ作ろうとして `{}` で文を囲ったが、この `{}` はlambda リテラルにパースされるようで、中の文が実行されない。
Swiftでもあったなー `do {...}` とかやったなーとか思いつつ、Kotlinでどうやるか色々探していたが、どうやらキーワードつけるだけでは無理っぽい。

`{...}()` と書いて、lambdaを即時実行させる例を見かけたので、それを使って以下のように書いたが、 今度はコンパイルが通らない。

```kotlin
{
        println("Foo")
}()

{
        println("Bar")
}()
```

> Too many arguments for public abstract operator fun invoke()

というエラーメッセージ。

ひとまず、その場はテスト対象のオブジェクトに対してScope Functionの `run` を使うことで、目的は達成できたのだが、
先のエラーメッセージの意味がわからない。引数、渡してないし。。。

ということをつぶやいていたら [@kokuyouwind](https://twitter.com/kokuyouwind) から「パース狂ってるのでは」「2個目のlambdaが1個目の呼び出しの引数としてパースされてそう」という指摘をもらったので、その線で調べてみたらドンピシャのissueがあった。

[Lambda literal on the next line of a function call is parsed as an argument ("Too many arguments" error) : KT-17884](https://youtrack.jetbrains.com/issue/KT-17884)

さらに同氏より 「`()` だけは多分valid。件のエラーメッセージは型検査レベルのエラーなので、文法エラーではないっぽい」と言われた。

## 2つ目の()が構文エラーにならないのはなぜなのか

手元で確認したところによると `()` だけ書いてもコンパイルエラーになる。
よくわからんなー、と思って、今この文章を書いていたら、理由が思い当たった。

2個目の `()` は、 1個目のlambdaの戻り値がclosureだったときの呼び出し用の `()` としてパースされているのではないか。
構文解析の段階ではvalidであり、まだ型検査中のなので、1個目のlambdaの戻り値が呼び出し可能か判定できていない。
そのため、先に `Too many arguments ...`のエラーとなった、ということなのだろう。

検証するために、 まず `Too many arguments...` をパスさせるために、以下のように1個目のlambdaが引数に関数を受け取るようにした。

```kotlin
{ p: () -> Unit ->
    println("Foo")
    p()
}()

{
    println("Bar")
} // ()  ... ひとまずコンパイルを通すためにコメントアウト
    
// 以下のように出力される
// Foo
// Bar
```

上記のコメントアウトを取ると、 

> Expression .... of Type 'Unit' cannot be invoked as a function. The function `invoke()` is not found

という想定通りのエラーが出た。このコンパイルエラーを通すには、1個目のlambdaが関数を返さなければならない。

ここでは、渡された p をそのまま返してみよう。

```kotlin
{ p: () -> Unit -> 
    println("Foo")
    p()
    p
}()

{
    println("Bar")
}()

// 以下のように出力される
// Foo
// Bar
// Bar
```

「1個目のlambdaの戻り値がclosureだったときの呼び出し用の `()` としてパースされている」という予想は、どうやらあっていたようだ。

## 2つのlambdaを即時実行させるには

最初の問題にもどって2つのlambdaを即時実行させるには、1個目のlambdaの呼び出しと2個目のlamdbaの定義が別れていることをコンパイラに伝えればよいので、セミコロンで区切れば良い。

```
    { 
        println("Foo")
    }();
    
    {
        println("Bar")
    }()
```

普段セミコロンを使わずにコードを書いているので、違和感しかない・・・

パーサが修正されることは期待できないので、その場その場でScope FunctionなどのKotlinらしい解決方法を見つけたほうが良さそう。
