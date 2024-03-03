---
title: KotlinのIterable.unzipはPairのリストにしか対応していない
description: KotlinのIterable.unzipはPairにしか対応していないのでTriple用の拡張メソッドを作ったときのメモ
published: '2019-05-26'
tags:
  - kotlin
---
タイトルの通り、 [kotlin.collectionsのunzip](https://kotlinlang.org/api/latest/jvm/stdlib/kotlin.collections/unzip.html) はPairのリスト（厳密に言うと `Array` と `Iterable` ）にしか対応していなかった。

なので、 `Triple` 用のを定義した。
まぁ [`Iterable<Pair>#unzip`](https://github.com/JetBrains/kotlin/blob/deb416484c5128a6f4bc76c39a3d9878b38cec8c/libraries/stdlib/src/kotlin/collections/Iterables.kt#L83) の定義をほとんど流用しているので、特に面白いことはない。


```
fun <A, B, C> Iterable<Triple<A, B, C>>.unzip(): Triple<List<A>, List<B>, List<C>> {
    val listA = ArrayList<A>()
    val listB = ArrayList<B>()
    val listC = ArrayList<C>()
    for (triple in this) {
        listA.add(triple.first)
        listB.add(triple.second)
        listC.add(triple.third)
    }
    return Triple(listA, listB, listC)
}
```

`Iterable<Iterable<T>>#unzip` も作れると思うけど、内側のIterableの要素数足りないときにどうするのか決めないといけなくて、今は必要としていなかったので、特に考えてない。

それにしても、なぜ `Iterable<Pair>#unzip` があって `Iterable<Triple>#unzip` が用意されていないのだろう？

`libraries/stdlib/src/kotlin/collections/Iterables.kt` をGit blameして、リファクタリング前のコミットを追ってみると、 `Iterable<Pair<T, R>>.unzip` が最初に実装されたのは [このコミット](https://github.com/JetBrains/kotlin/commit/da3ec891d0b476b26881bd122dcb63c24d42b711) だった。

もとのissueは[KT-5793](https://youtrack.jetbrains.com/issue/KT-5793)。
最初からTripleのことは考えられていなかったようだ。

[このコミット](https://github.com/JetBrains/kotlin/commit/8333448f10ef16de2fcc2868664176206db1fa55#diff-b3b2e95da541fae3f1e3d4139032c03e) でPairとTripleは同時に入ったようなので、KT-5793の対応するときに単に `Triple` のことが忘れ去られていたか、その時必要ではなかったので後回しにされているか、だと思う。

## Next Action

KotlinのSlackにjoinして聞いてみる。
