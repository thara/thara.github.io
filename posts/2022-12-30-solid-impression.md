---
title: Solid所感
date: '2022-12-30'
published: '2022-12-30'
---

ちょっと遅れたけど、[「Web3はWebじゃない」とWebの父が断言](https://www.axion.zone/web3-9/) を読んで。

いわゆるWeb3がブロックチェーンを前提にするのに対し、Web3.0はセマンティック・ウェブを前提としていて、その思想や目的は全く異なる。
Tim Berners-Leeが嘆くのもしょうがない。

ところで、Timが創業した [Inrupt](https://www.inrupt.com/) 、そしてその事業の中心である [Solid](https://solidproject.org/) については初めて知ったので、今の見解を記録に残しておきたい。

記事中では「Inruptが開発するSolid」と説明されていたがこれは不正確。
Solidという分散データストアとそれへのアクセス制御を規定する仕様があり、その個人ごとのデータ保存単位であるPodを管理するサーバーを[エンタープライズレベルで提供する](https://www.inrupt.com/products/enterprise-solid-server)のがInrupt。   
   
Solidを広めるためにTimが事業を立ち上げた、と自分は捉えた。

## Solidとは何か

Solidとは前述のように、分散データストアとそれへのアクセス制御を規定する仕様であり、分散されたデータストア間では各ユーザーやアプリケーションは [WebID](https://www.w3.org/2005/Incubator/webid/spec/identity/) によって識別される。

[GitHub orgがあり](https://github.com/solid/)、以下のようなdescriptionが書かれている。

> Re-decentralizing the web.

もともとWebは非中央集権であったものが、今は特定のサービスプロバイダによる中央集権となってしまったものを、再度非中央集権に戻す、ということか。

Solid自体はHTTPS上で構築できるプロトコルで、2022-12現在はW3Cにworking groupが存在し仕様はオープンになっているがまだEditor's draft。

[Solid Protocol](https://solidproject.org/TR/protocol) 

軽く中身を見ると、懐かしきリソース指向で、RDF(Resource Description Framework)をベースにしているっぽい。
GET/POST/PUT/PATCH/DELETE などのHTTPメソッドが正しく使われる世界。個人的には好み。

だが、令和のこの時代にRDFが普及しきるのか、という疑問がある。
更新のための「N3 Patch」という、[Notation 3](https://w3c.github.io/N3/spec/) を用いたPATCHリクエストをサポートしなければならないが、
RDFとかN3とかを扱ったライブラリが現時点で各言語に備わっているかというとだいぶ怪しい。

RDF自体が仕様が膨大で、一開発者としては、SolidのためだけにRDFライブラリを開発するモチベーションはなかなか湧かないし、Solid以外にRDFを必要とするユースケースが自分の身の回りではあまりない。  

[国立国会図書館でLinked Open Data扱ってる](https://www.ndl.go.jp/jp/dlib/standards/lod/index.html) から、オープンデータを扱う人とかにはそのモチベーションがあるかもしれないけれど。

## Solidの今後

仕様がRDFをベースにしていて、普及の大きな障壁になる、という課題がありそう。

一方で、GDPR/CCPAをはじめとする個人データ保護の強化の流れは変わらず、これにSolidが刺さる可能性がある。
サービスプロバイダとしてもGDPR/CCPA対応には苦心するポイントではあるので、それがSolidによってサービスプロバイダ側の負担が減るのであれば、エンドユーザーとサービスプロバイダ双方にとってメリットがある。

しかし、Solidでは個人ごとにPodを管理するSolidサーバーを持つことで個人データ保護をできると謳っているように見えるが、少なくともここ数年で個人ごとにSolidサーバーを管理するようにはならないように思える。
非中央集権を目的としているのにも関わらず、GoogleやAppleなどが管理するSolidサーバーをエンドユーザーが使用することになる。結果、中央集権であることに変わりはないのではないか。

これまでの流れを見ると、中央集権・非中央集権どちらかで定着するのではなく、螺旋のように時間軸でどちらかに揺れ動く、それ自体に意味や価値があるのかもしれない。

## 自分のスタンス

非中央集権やリソース指向であるSolid protocolなどは自分の好みではあるが、如何せんRDFなどの技術に自分のプライベートの時間をベットする余裕はないので、一旦は静観。
仕事でも利用する可能性があるのはだいぶ先に思えるし。

ただ、GDPR/CCPA周りの流れと合わせて動向は見守っておきたいところ。

<iframe width="560" height="315" src="https://www.youtube.com/embed/qWVTjMsv7AE" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
