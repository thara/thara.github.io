---
title: Go ジェネリクスでインターフェース制約を満たしつつ zero valueで初期化する
date: '2025-07-21'
published: '2025-07-21'
---

例えば、proto.Message を型制約に持つgenericな関数を書こうとすると、次のような落とし穴がある：

```go
func DecodeProtoMessage[T proto.Message](data []byte) (*T, error) {
	...
	var m T
	if err := proto.Unmarshal(b, m); err != nil {
		...
	}
	return &m, nil
}
```

このとき、`T` に `*Pet` のようなポインタ型を指定すると、`var m T` は nilのままとなり、`proto.Unmarshal` はpanicしてしまう。
Go のジェネリクスでは、ポインタ型に対するゼロ値の `new(T)` や `var m T` は nilなので、non-nil前提の処理ではこうなる。

## 解決策：値型 T とポインタ型 *T を分離して制約する

そこで、以下のように値型 `T` に対するポインタ型 `*T` に制約をかけることで、安全な初期化ができる：

```go
type ProtoMessagePtr[T any] interface {
	*T
	proto.Message
}

func DecodeProtoMessage3[T any, P ProtoMessagePtr[T]](data []byte) (*T, error) {
	...
	var m T
	var p P = &m
	if err := proto.Unmarshal(b, p); err != nil {
		...
	}
	return &m, nil
}
```

`T` が `any` にもかかわらず、 `P` の型制約によって`*T`が `proto.Message` を満たすことを関数定義が要求していることがポイント。
この方法は、`proto.Message` に限らず、任意のinterfaceに応用できる。

使用例:
```go
	dst, err := DecodeProtoMessage[petstorev1.Pet](encodedData)
```

このように `P` は第1型パラメータによって自明なので、関数を利用する際には `ProtoMessagePtr` を指定する必要もない。
reflectionやfactoryなどを使わずに、型安全にzero valueを生成できていい感じ。
