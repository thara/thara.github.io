---
title: 個人的によく採用するDjangoプロジェクト構成
date: '2016-10-07 17:13:43'
last_modified_at: '2018-01-14'
tags: [python, django]
---

## 2018-01-14追記
この記事の意見は、2018-01-14 現在の見解とは異なります。   
最新の見解は [こちら]({{ site.baseurl }}{% post_url 2018-01-14-rethink-django-project-structure %}) を参照してください。

---

[Django](https://www.djangoproject.com/) を実際にプロダクトで使用するとなると、まずプロジェクト構成をどうするか、という問題にぶち当たる。

ここでいう プロジェクト構成とは以下を指す。

- Djangoプロジェクトのディレクトリ構成
- Djangoアプリケーションの単位・ファイル構成
- settings管理方法
- 各Djangoアプリケーションごとの内部構成

現在の自分の意見をまとめてみると、こんな感じ


```
├── sampleapp
│   ├── settings.py
│   ├── urls.py
│   └── wsgi.py
├── sampleapp_foo
│   ├── models.py
│   ├── urls.py
│   └── views.py
├── sampleapp_bar
│   ├── models.py
│   ├── urls.py
│   └── views.py
├── manage.py
├── requirements.txt
└── setup.py

```

ルートは作業ツリーのルートと同じ = READMEを置くところ。   
`sampleapp` がいわゆるプロジェクト名で、その配下にプロジェクトの設定をおいておく。   
settings.pyは複数環境に対応できるよう、データベースの接続先URLなどの環境ごとに異なる値は環境変数から取得するようにする。

`sampleapp/urls.py` にURLパターンを記述していくが、基本的に他のDjangoアプリケーションの `urls.py` を読み込むだけにしておく。

`sampleapp_foo` がいわゆるDjangoアプリケーション。
モジュール名の重複を避けるため`{プロジェクト名}_XXXX` みたいな感じにしている。基本的に、**機能**単位で作成し、機能的凝集を高める。（自然とbase URLのサブディレクトリごとになることが多い）

Djangoアプリケーションの配置場所には、
`sampleapp/foo` とか `apps/foo` とか `sampleapp/apps/foo` とか、いろいろ流儀はあるようだが、

- Djangoアプリケーションはそれぞれで機能的に独立していることが望ましい
- Pythonではネストの深いパッケージはあまり見かけない
- Djangoアプリケーションは機能追加によりどんどん増えていく

という観点から `sampleapp_foo` にしている。
次点で、`sampleapp/foo`。

最初は作業ツリー直下に `sampleapp_XXX` が増えていくのに抵抗があったが、今は慣れた。なにより目的のファイルを見つけやすいことが嬉しい。

models.pyは、大きくなったら `models/XXX.py` みたいにmodelsディレクトリ作って、分割。`models/__init__.py` 内でmodelのクラスをimportして、Djangoから認識されるようにする。

だいたいこんな感じ。
これから更に、Djangoアプリケーション内でどのようにもジョールを定義していくか、という問題があるけれど、それはまた今度・・・
