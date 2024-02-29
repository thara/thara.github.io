---
title: 中規模Djangoアプリケーションのモジュール構成
date: '2016-12-09 16:42:23'
tags: [python, django]
---

DjangoやRuby on Railsといったフルスタックフレームワークは、シンプルなCRUDアプリケーションを素早く構築することができるが、複雑なアプリケーションの設計をどうすればよいかという問題にはいつも悩まされる。

どのようなプロジェクトによるかはケースバイケースだが、自分の中でおおよそ固まってきたので、それをまとめてみる。

Djangoの場合、フレームワークが必要とするモジュールはsettings.pyとmodels.pyぐらいなので、各Djangoアプリケーション内では割と自由にモジュールを定義できる。

自分は以下のような自作モジュールを配置するようにしている。

- commands
- query
- services
- utils

### commands

まずは、command。このcommandはGoFのCommandパターンのことではなく、Command-Query Responsibility SegregationのCommand。   
つまり、更新系のオペレーションのこと。
このモジュールには、更新用の関数とその関数用の引数/戻り値（入力/出力）用ののDTOが含まれ、viewsから直接使用される。
このモジュールの更新用関数は、DBのトランザクション境界とおおよそ一致する。

大体以下のような感じ。

```python
from collections import namedtuple

UserProfileUpdateCommand = namedtuple(
    "UserUpdateCommand", "user_id user_name email_address")
class UserProfileUpdateResult(namedtuple(
    "UserProfileUpdateResult", "success errors")):
    def to_dict(self):
        return {"success": self.success, "errors": self.errors}


def _validate_command(command: UserProfileUpdateCommand):
    """ commandに指定された値のプロパティ検証 """

def update_user_profile(command: UserProfileUpdateCommand):
    """ 更新系処理 """
    errors = _validate_command(command)
    if errors:
        return UserProfileUpdateResult(success=False, errors=errors)

    with transaction.atomic():
        user_profile = UserProfile.objects.get(command.user_id)
        user_profile.user_name = command.user_name
        user_profile.email_address = command.email_address
        user_profile.save(
            update_fields=['user_name', 'email_address'])

    return UserProfileUpdateResult(success=True, errors=[])    
```

`_validate_command` は、上の例のようにシンプルな更新であれば`UserProfile`のメソッドとして提供してもよいかもしれない。   
が、modelを更新する機能が複数ある場合は検証処理もそれによって異なることがあるので、そのような検証処理はcommandに定義した方がいいだろう。   
また、commandでmodelの`save`メソッドを直接呼ぶのも良いが、複数の機能から同じような更新処理を実行するのであれば、modelにメソッドを実装するのも悪くない。

これらの関数/クラス はViewから以下のように使われる。（RESTｆrameworkを使用を前提としている）

```python

class UserProfileView(APIView):

    def post(self, request):
        s = UserProfileUpdateSerializer(data=request.data)
        s.is_valid(raise_exception=True)
        command = UserProfileUpdateCommand(
            user_id=serializer.validated_data['user_id'],
            user_name=serializer.validated_data['user_name'],
            email_address=serializer.validated_data['email_address'])
        result = update_user_profile(command)
        return Response(result.to_dict(), status=status.HTTP_200_OK)

```

### query

queryは、commandと対照的に、viewsから使用される参照系の関数を定義する。

```
def get_user_profile(user_id):
    return UserProfile.objects.using('read_replica').get(user_id)
```

この例だとあまり利点は見えないが、DBリードレプリカを参照したり、キャッシュからデータを取ってくるなど、パフォーマンス優先の実装を配置するのに最適。

## services / utils

このserviceとは、レイヤードアーキテクチャのサービス層ともDDDのドメインサービスとも関係がない。
services.pyには、viewsからは直接使用されないが、同じDjangoアプリケーション内のcommandsやquery、または異なるDjangoアプリケーションから使用される関数やクラスを定義する。


一方のutilsは、あらゆるアプリケーションから使用される可能性があるユーティリティ関数を定義する。
こちらは極力、Djangoなどのフレームワークに依存しないことが望ましい。

## Djangoアプリケーションの構成

ここまで挙げた、 commands, query, services, utilsを、機能を表すパッケージごとに定義する。

以下は一例。

より実際のファイル構成に近くなるよう、フレームワークの使用にあたって必要になるであろうファイルも配置した。

models.py, views.py, urls.pyはそれぞれDjangoのDBモデル, View, URLディスパッチ設定。serializers.pyはDjango RESTframework用のJSONシリアライズ設定。


```
sns_sample
├── friend
│   ├── request
│   │   ├── commands.py
│   │   ├── query.py
│   │   ├── serializers.py
│   │   └── services.py
│   ├── commands.py
│   ├── models.py
│   ├── query.py
│   ├── serializers.py
│   ├── serivces.py
│   ├── urls.py
│   └── views.py
└── user
    ├── login
    │   ├── commands.py
    │   ├── serializers.py
    │   └── services.py
    ├── signup
    │   ├── commands.py
    │   ├── serializers.py
    │   └── utils.py
    ├── models.py
    ├── serializers.py
    ├── urls.py
    └── views.py
```

だいたいこんな感じ。
それほど複雑になるわけでもなく、かといって不足もないようなちょうどよい塩梅だと思う。もちろん、プロジェクトの規模にもよるのだが、中規模であれば、これぐらいで良いだろう。

実際のモジュール配置場所を決定するために一番考えなければいけないことは、モジュール構成を変更するのはかなりリスキーである、ということを認識すること。

バージョン管理システムによるソースコードのバージョン管理もファイルの配置場所を変えるとうまくいかないことが多いため、今までのソースコードの変更履歴を追うのがめんどくさくなる。

よりシンプルでわかりやすい構成にしておくことを心掛けたい。
