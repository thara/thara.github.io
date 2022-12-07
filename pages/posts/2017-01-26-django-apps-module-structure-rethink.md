---
title: Djangoアプリケーションのモジュール構成再考
tags: [python, django]
---

[以前のDjangoモジュール構成のpost]({% post_url 2016-12-09-module-configuration-of-django-app %}) では、論理凝集度が高過すぎる感があって、もっと機能凝集度を高めたモジュール構成を考えてみた。   
もちろんDjangoアプリケーションは細かく分けることを前提。

以下、全体像を簡単に図示。

![Django App Structure]({{ site.url }}/images/django-app-structure.png)

各要素を説明すると...

- View
    - DjangoのView定義
- Model
    - Django Model定義
    - 他のModelへのクエリ/更新をしない
    - Fat Modelにしない
- Domain
    - 機能単位のモジュール。機能を満たす関数やクラスを定義する。
    - モジュール内には、さらに以下を定義する
      - Command
        - 更新系ユースケースを表すアプリケーションサービス。Viewから使用。
      - Query
        - 参照系ユースケースを表すアプリケーションサービス。Viewから使用。
- Infrastructure
    - Domainの実装のうち、Redisやメール送信、Celeryなどの外部サービスのAPIを使用する実装を置く。
    - DjangoのModelはInfrastructureの責務には **含めない**
        - フルスタックフレームワークを使用している以上、Django Modelへの依存を避けようとすると、フレームワークの良さを殺してしまう


以下、ファイル構造の一例。

```
social_game_api
├── friend
│   ├── models.py  (Model)
│   ├── request.py (Domain)
│   ├── request_tasks.py (Infrastructure: Celery)
│   ├── urls.py
│   └── views.py (View)
└── user
    ├── login.py (Domain)
    ├── login_repos.py (Infrastructure: Redis)
    ├── loginbonus.py (Domain)
    ├── models.py (Model)
    ├── profile.py (Domain)
    ├── signup.py (Domain)
    ├── urls.py
    └── views.py (View)
```

モジュール名は `{Domain}_{Infrastructure type}` という法則にすることで、どのDomainのどのようなInfrastructureを使用するモジュールかが分かりやすくなる。   
また、他のDjangoアプリケーションからアクセスできるのは、ModelsとDomainに限定する。   
Domainでは、他のDomainを参照することにより循環参照になることがある。その場合は、関数やメソッド内でimportするか、別モジュールに切り出すなどして回避する。

やっぱり、Simple is best.
