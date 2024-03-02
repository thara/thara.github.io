---
title: Djangoプロジェクト構成 再考
tags: [python, django]
---

[過去のポスト](../my-django-project-structure) で、リポジトリルートにDjangoアプリケーションを配置する構成を挙げたけれど、
その方法はネストが浅くなりコード全体を見渡しやすくなる反面、Djangoアプリケーションを機能単位で作成していると大量のディレクトリがリポジトリルートに配置されることになり、この構成のメリットであるコードの見通しが悪化するという現象にあった。

特に、DjagnoプロジェクトとDjangoアプリケーションのルートディレクトリがリポジトリルートに混在することになってしまう点が一番のデメリットだった。

この反省を踏まえ、今現在では、新規プロジェクトは以下のような構成で構築している。

```
├── config
│   ├── settings.py
│   ├── urls.py
│   └── wsgi.py
├── sample
│   ├── foo
│   │   ├── models.py
│   │   ├── urls.py
│   │   └── views.py
│   ├── bar
│   │   ├── models.py
│   │   ├── urls.py
│   │   └── views.py
│   └── utils.py
├── constraints.txt
├── manage.py
└── requirements.txt
```

## Djangoアプリケーションの配置

上記の `sample` ディレクトリは、プロジェクト名を示すが、 **DjangoプロジェクトでもDjangoアプリケーションでもない** 通常のPythonパッケージを表すディレクトリ。
`sample/for` や `sample/bar` が実際のDjangoアプリケーションのルートとなる。必然的に、settings.py に記述する `INSTALLED_APPS` には `sample.for` や `sample.bar` と指定することになる。
`sample` 直下には、複数のアプリケーションが共通で使用するモジュール、例えば、ModelやViewの基底クラス・ユーティリティ関数などを定義したモジュールを配置する。

この配置は、[`Two Scoops of Django`](https://www.amazon.co.jp/dp/B076D5FKFX) を参考にした。

## Djangoプロジェクトの配置

上記の `config` ディレクトリが Djangoプロジェクトのルートとなる。これも、`Two Scoops of Django` を参考にしている。

プロジェクトによって複数のエンドポイントを用意する必要があるかもしれない。その際は、以下のように `config` 内にサブディレクトリを作成する。

```
└── config
    ├── api
    │   ├── settings.py
    │   ├── urls.py
    │   └── wsgi.py
    └── adminsite
        ├── settings.py
        ├── urls.py
        └── wsgi.py
```

必要であれば、`DJANGO_SETTINGS_MODULE` のデフォルト値を変えた manage.py を各Djangoプロジェクトルートに配置しても良いかもしれない。

## Djangoアプリケーションの粒度について

再利用を意図しないDjangoアプリケーションの粒度については、モノリシックにするか、マイクロにするかでいろんな意見がある。
自分はマイクロ派、特に機能ごとに分割する派であったが、一定規模のDjangoプロジェクトだと、マイクロなDjangoアプリケーションを大量に作成することで以下のような問題に直面した。

- 必ずしも依存関係が単方向であるとは限らないため、modelsやそれを利用するモジュールが循環参照となることがある
    - 最初は依存方向を加味していても仕様変更で崩れる
- 循環参照を避けようと signals で解決しようとすると、1度のsignals送信で発火される複数アプリケーションのreceiverの実行順を意識する必要がある
    - receiverをapp configで登録していた場合、INSTALL_APPSで定義された順にreceiver登録されるが、わかりづらかった
- 機能追加時にどのアプリケーションに入れるのか、または新しいアプリケーションを作成するのか、迷うことがある
    - この判断は難しく、また、人によって異なる

これを踏まえて、今では 「DB分割可能な機能単位」 で Django アプリケーションを作成することにしている。
つまり、それぞれの Django アプリケーションで異なるDB（シャーディングは別）を参照する。

高負荷環境において、DBの垂直/水平分割は避けられない。また、コードベースの変更やデータ移行のコストが膨大であるため、ローンチ後は機能よりも永続化されたデータやDB構成の方が変更が困難であることが多い。
Djangoアプリケーションを作成する段階で、あらかじめDB分割可能な機能単位を分析・設計するのは、悪くないタイミングだと思う。

モバイル向けソーシャルゲームのバックエンドAPIを例に、「DB分割可能な機能単位」 で Djangoアプリケーションを定義すると以下のようになる。

```
└── gamebackend
    ├── auth    -> auth db に接続
    ├── friend  -> friend db に接続
    ├── player  -> player data db に接続
    └── battle  -> battle db に接続
```

それぞれのDjango アプリケーションで、urls.py を定義してもよいし、エンドポイント定義用のモジュールを別途定義しても良いと思う。
後者の場合、 Django アプリケーション間が疎結合であれば、後に負荷対策として Django プロジェクトとして独立したサービスとして稼働させることも可能になる。   
（マイクロサービス化、と言いたいところだが、この粒度だと組織構造と一致しないし、単一リポジトリであるため、スムーズに移行できるとは思えない）

## Summary

- リポジトリルート直下の通常の Python パッケージに Django アプリケーションとアプリケーション共通で使用するモジュールを配置する
- `config` またはそのサブディレクトリに Django プロジェクトのルートディレクトリを配置する
- Django アプリケーションの粒度は、 DB分割可能な機能単位とする

今まで、様々なプロジェクト構成を試してきたが、一周回って元に戻ってきた感があるが、この構成は今までで一番しっくり来ている。   
`Two Scoops of Django` は昔読んだことがあったのだけれど、すっかり内容が吹っ飛んでいたな...   
このブログより、よっぽど有用なことが書いてある。Djangoを使う人は必読ですよ。

---

最近は Golang がバックエンドで人気だが、DB マイグレーションやWeb アプリケーション、restframework のRESTfull API など、Django を使用することで設計から実装・検証まで時間がかかりがちなフルスタックのアプリケーションを素早く構築することができる。Python は Golang に比べてパフォーマンスが劣ることは衆知の事実だが、プロジェクト早期から MVP としてアーキテクチャ検証可能なポテンシャルがある。   
流行りの機械学習だけではなく、 Web アプリケーションに採用される言語/フレームワークとしても、もっと注目を集めてほしい・・・