---
title: Kotlinの.isInitialized呼び出しになぜCallableReferenceの指定が必要なのか
published: '2018-12-07'
tags: [kotlin]
---

Misocaの [thara](https://twitter.com/zetta1985) です。   
この記事は [Misoca+弥生+ALTOA Advent Calendar 2018 - Qiita](https://qiita.com/advent-calendar/2018/misoca-yayoi) 7日目の記事です。

最近は主にモバイルアプリ([Android版](https://play.google.com/store/apps/details?id=jp.misoca.misoca&hl=ja)、[iOS版](https://itunes.apple.com/jp/app/id1026534800))のバックエンドを担当してますが、Java歴が長いのとSwift関連の本を書いたこともあってクライアントサイドのコードもちょこちょこ書いてます。

今回はKotlinで書かれたソースコードをレビューしていて疑問に思ったことを掘り下げて調べてみたので、そのことについて書きます。   
（本文は丁寧語を考えるのが面倒なので、文語体風で書きます）

----

## .isInitializedについての疑問
Kotlinには`lateinit` という修飾子がある。
この修飾子は、非nullのプロパティや変数に対してコンストラクタ外で初期化することを許可する。
初期化する前にそのプロパティや変数にアクセスすると、`UninitializedPropertyAccessException` という例外がスローされる。
よって、`lateinit` がつけられたプロパティや変数に安全にアクセスするためには、`.isInitialized`  が真であることを確認してからアクセスしなければならない。

初期化しない場合

<script src="https://gist.github.com/thara/838f2ee3bcc166ff6bcea00860d38c49.js?file=0-lateinit-not-initialized.kt"></script>

これは以下のような例外が発生する。

```bash
Exception in thread "main" kotlin.UninitializedPropertyAccessException: lateinit property s has not been initialized
    at A.print(Simplest version.kt:8)
```

`.isInitialized` で初期化済みであることを確認

<script src="https://gist.github.com/thara/838f2ee3bcc166ff6bcea00860d38c49.js?file=1-lateinit-isInitialized.kt"></script>

さて、上記のコードでは `.isInitialized` を呼び出す際に`::`という`prefix`をつけている。これは、[Callable Reference](https://kotlinlang.org/docs/reference/reflection.html#callable-references)を省略したもので、たとえば `lateinit` がクラスのインスタンスフィールドにつけられていた場合には`this::`という記述になる。   
これを `s.isInitialized` に変えるとコンパイルが通らない。`String` には `isInitialized`というプロパティがないからだ。

<script src="https://gist.github.com/thara/838f2ee3bcc166ff6bcea00860d38c49.js?file=2-lateinit-unresolved-reference.kt"></script>

また、`lateinit` という修飾子がつけていなくても、当然のようにコンパイルが通らない。

<script src="https://gist.github.com/thara/838f2ee3bcc166ff6bcea00860d38c49.js?file=3-lateinit-only-be-called-on-a-lateinit-ref.kt"></script>

では、この`.isInitialized`はどこからやってくるのだろうか？   
なぜ`::`や`this::`といった、Callable Referenceの指定が必要なのだろうか？

## Kotlinの内部を追う
Android Studio上で `.isInitialized` の定義元に飛ぶと、[kotlin/Lateinit.kt](https://github.com/JetBrains/kotlin/blob/e21c235bbb839d273fd6b419d0b3527e573f5e4e/libraries/stdlib/src/kotlin/util/Lateinit.kt) であった。

```kotlin
/**
 * Returns `true` if this lateinit property has been assigned a value, and `false` otherwise.
 *
 * Cannot be used in an inline function, to avoid binary compatibility issues.
 */
@SinceKotlin("1.2")
@InlineOnly
inline val @receiver:AccessibleLateinitPropertyLiteral KProperty0<*>.isInitialized: Boolean
```

どうやら通常のプロパティアクセスとは異なり、コンパイラからは特別なプロパティリテラルとして認識されるらしい。

`lateinit`が導入された[コミットログ](https://github.com/JetBrains/kotlin/commit/c6263ac8e6cab87d18f92d77c401faf4b6cad431)を参考にコードを追うと `.isInitialized` の実態は [kotlin/LateinitIntrinsics.kt](https://github.com/JetBrains/kotlin/blob/2c4e023cba45b8dea65cd99d86373445eec35336/compiler/backend/src/org/jetbrains/kotlin/codegen/intrinsics/LateinitIntrinsics.kt) で自動生成されたもののようだ。   
また、[kotlin/LateinitIntrinsicApplicabilityChecker.kt](https://github.com/JetBrains/kotlin/blob/0e5544a4919a185f52b5e4294820d6ed45a45d6d/compiler/frontend/src/org/jetbrains/kotlin/resolve/calls/checkers/LateinitIntrinsicApplicabilityChecker.kt) では、コンパイル時の構文チェックを行っている。そこで以下のような記述を見つけた。

```kotlin
if (!referencedProperty.isLateInit) {
    context.trace.report(LATEINIT_INTRINSIC_CALL_ON_NON_LATEINIT.on(reportOn))
}
```

`context.trace.report(LATEINIT_INTRINSIC_CALL_ON_NON_LATEINIT.on(reportOn))` は、おそらくコンパイルエラーのメッセージを指定しているのだろう。[kotlin/DefaultErrorMessages.java](https://github.com/JetBrains/kotlin/blob/2182be82e64cb57c43d711d3696e21b2133ab285/compiler/frontend/src/org/jetbrains/kotlin/diagnostics/rendering/DefaultErrorMessages.java)には、   
この `LATEINIT_INTRINSIC_CALL_ON_NON_LATEINIT` という名前で、さきほどの `lateinit` がつけられていない変数にアクセスした際に表示されたエラーメッセージが指定されていた。

[kotlin/LateinitIntrinsicApplicabilityChecker.kt](https://github.com/JetBrains/kotlin/blob/0e5544a4919a185f52b5e4294820d6ed45a45d6d/compiler/frontend/src/org/jetbrains/kotlin/resolve/calls/checkers/LateinitIntrinsicApplicabilityChecker.kt) をさらに読み進めると、以下の判定式を見つけられた。
```swift
} else if (!isBackingFieldAccessible(referencedProperty, context)) {
    context.trace.report(LATEINIT_INTRINSIC_CALL_ON_NON_ACCESSIBLE_PROPERTY.on(reportOn, referencedProperty))
```

`isBackingFieldAccessible` という名前から、backing fieldにアクセスできない場合は `.isInitialized` を使うことができないらしい。

Kotlinの通常のプロパティアクセスの構文ではgetterやsetterが呼び出されるが、特定のコンテキストではgetterやsetterを介さず直接フィールドにアクセスできる。   
そのとき、そのフィールドをbacking fieldと呼ぶ。

backing fieldにアクセスできるのは、そのプロパティが宣言されたソースファイル内に限られる。
これは、そのソースファイル以外のクラスの`lateinit` が修飾されたプロパティや変数に対して、`.isInitialized`を呼ぶことができないことでもある。

これで、先程挙げた2つの疑問が解消した。

## 疑問の答え
まず、「`.isInitialized`はどこからやってくるのだろうか？」という疑問。

これは、`lateinit`を修飾したプロパティや変数は`.isInitialized` という特殊なプロパティリテラルを使えるようにコンパイラが特別扱いしていたのだった。具体的な実装はインライン化されており自動生成される。[^1]   

そして、次の「なぜ`::`や`this::`といった、Callable Referenceの指定が必要なのだろうか？」という疑問。

これは`.isInitialized`を使用可能なのがbacking fieldに限られるからだ。`lateinit`をつけたプロパティのオブジェクトに対する通常のプロパティアクセス（getter/setterの呼び出し）と区別するためにCallable Referenceが構文上必要だった。

さて、ここで新たな疑問が起こる。   
Kotlinでは、当然のようにCallable Referenceを指定して通常のプロパティアクセスも可能だ。`foo`というプロパティが存在すれば、`this::foo` と呼び出せる。   
では、元々 `lateinit`をつけたプロパティの型に、`.isInitialized`というメソッドが定義されていた場合はどうなるのだろうか？

### .isInitializedメソッドが存在するときの挙動

挙動を確認するために以下のコードを書いた。


<script src="https://gist.github.com/thara/838f2ee3bcc166ff6bcea00860d38c49.js?file=4-lateinit-duplicated-isInitialized.kt"></script>

これを実行すると、 `Initialized` が表示される。つまり、[kotlin/Lateinit.kt](https://github.com/JetBrains/kotlin/blob/e21c235bbb839d273fd6b419d0b3527e573f5e4e/libraries/stdlib/src/kotlin/util/Lateinit.kt) で定義されている `.isInitialized`が実行される。

次は、以下のような変更を加えてみる。Callable Referenceの指定を取ったものだ。

<script src="https://gist.github.com/thara/838f2ee3bcc166ff6bcea00860d38c49.js?file=4-lateinit-duplicated-isInitialized.diff"></script>

これを実行すると、今度は `Not initialized` が表示される。これは、`X`クラスの`isInitialized` プロパティが呼び出されたことを表している。

これらのことから、 **Callable Referenceを指定してisInitializedプロパティを使うとkotlin/Lateinit.ktのisInitializedが優先される** ことがわかった。   
これは2つ目の「なぜ`::`や`this::`といった、Callable Referenceの指定が必要なのだろうか？」という疑問の答えを補完しているようにも思える。  
Callable Referenceを指定すると`isInitialized`の解決に[kotlin/Lateinit.kt](https://github.com/JetBrains/kotlin/blob/e21c235bbb839d273fd6b419d0b3527e573f5e4e/libraries/stdlib/src/kotlin/util/Lateinit.kt) の`.isInitialized`が優先的に扱われるため、コンパイラから見て曖昧さが無くなるのではないだろうか。

## まとめ

- `lateinit`を修飾したプロパティや変数が使用できる`.isInitialized` はコンパイラから特別扱いされた特別なプロパティリテラルである
- `.isInitialized`はbacking fieldに対してのみ使用可能であるため、Callable Referenceを用いる必要がある。
- Callable Referenceの使用により[kotlin/Lateinit.kt](https://github.com/JetBrains/kotlin/blob/e21c235bbb839d273fd6b419d0b3527e573f5e4e/libraries/stdlib/src/kotlin/util/Lateinit.kt) の`.isInitialized`を使用することをコンパイラに伝える。

## 感想
`isInitialized` を最初レビューで見かけたときには、奇妙な構文だと思った。   
通常のプロパティアクセスのような構文を部分的に特殊扱いしている点が洗練されていない印象を持ったが、既存の構文への影響を考えると、実用的なKotlinらしい判断でもあると言えそうだ。

ちなみに、今回の記事を書くにあたって[KotlinのGitHubリポジトリ](https://github.com/JetBrains/kotlin)のコードを読んでいたのだが、途中で[Kotlin/KEEP: Kotlin Evolution and Enhancement Process](https://github.com/Kotlin/KEEP)なるリポジトリを見つけた。   
このリポジトリはKotlinの言語仕様に対するproposalを管理しているのだが、今回の`isInitialized`の件もちゃんとproposalが出ていた。

[KEEP/lateinit-property-isinitialized-intrinsic.md at master · Kotlin/KEEP · GitHub](https://github.com/Kotlin/KEEP/blob/master/proposals/lateinit-property-isinitialized-intrinsic.md)

これを最初に読んでおけば、Kotlin本体を見なくても疑問解決したなぁ・・・

自分が「洗練されていない印象を持った」という感想を持ったと先に書いたが、 現に

> The solution is admittedly very ad-hoc.

と言及されていた。ですよねー。

まぁ、Kotlin本体のソースコードリーディング楽しかったし、なんかコントリビュートしたい気持ちになったので、良しとしよう。

Swiftもそうだが、最近の言語はKEEPのような言語仕様のproposalもGitHub上でオープンに議論される風潮にあり、良い傾向だと思う。

----


明日は [めろたん](https://qiita.com/merotan)さん が「なんかがんばります」だそうです。   
これは期待できますよ・・・！

[^1]: MisocaでPRレビューしたときには、見当違いの指摘をしていた気がする。ごめんなさい…
