---
layout: post
title: RealmSwiftでまだ存在しないオブジェクトだけを追加する
published: '2018-10-22'
tags: [ios, swift, realm]
---

RDBを使用したアプリケーションでよくあるパターンなんだけど、Realmではあまり例がなさそうだったので。

以下のモデルを例に取り上げる。

```swift
class SampleObject: Object {
    // 略
    override static func primaryKey() -> String? {
        return "identifier"
    }
}

// すでに保存されている/まだ保存されていないモデルが混在
let objects: [SampleObject] = ....
```

## 問題
- すでに保存されているRealmモデルとまだ保存されていないRealmモデルが混在するArrayがある
- そのArrayのうち、まだ保存されていないRealmモデルだけをRealmに新規追加したい
- すでに保存されているRealmモデルのプロパティを更新したくない
- 少ないデータアクセス回数で抑えたい

`Realm#add` には複数オブジェクトを渡すことができるが、引数`update`をtrueにするとobjects内の「すでに保存されているRealmモデル」が更新されてしまい、falseにするとidentifierの値が重複するためエラーが発生する。

```
realm.add(objects, update: true)  // すでに保存されているRealmモデルが更新される
realm.add(objects, update: false) // identifierの値が重複するためエラー
```


## 解決方法
ポイントは「すでに保存されているRealmモデルのキー」がわかれば「まだ保存されていないRealmモデル」が分かる、ということ。

```swift
// すべてのprimaryKeyのArray
let keys = objects.map { $0.identifier }

// すでに保存されているRealmモデルを取得して、そのprimaryKeyのSetを得る
let storedKeys = Set(
    realm.objects(SampleObject.self)
      .filter("identifier IN %@", keys).map { $0.identifier })

// すでに保存されているRealmモデルのprimaryKeyに一致しないオブジェクトを選択する
let newObjects = objects.filter { !storedKeys.contains($0.identifier) }

if newObjects.count != 0 {
    try realm.write {
        // newObjectsはまだ保存されていないRealmモデルのみなので update: false
        realm.add(newObjects, update: false)
    }
}
```

効率的にデータアクセスするために以下の2点の工夫を入れている。

- 「すでに保存されているRealmモデルのキー」を一度のデータアクセスで取得する
- 「すでに保存されているRealmモデルのキー」を `Set`にすることで`contains`の計算量をO(1)にする

「すでに保存されているRealmモデル」の数が多い場合は、より効果的になる。
