---
title: golang.org/x/oauth2 package
date: '2026-08-16'
published: '2026-08-16'
---

Go言語には [golang.org/x](https://pkg.go.dev/golang.org/x) 配下で提供されている準標準ライブラリがある。

この中に [golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2) があって、OAuth2クライアント周りの実装が置かれている。

準、とはいえ、OAuth2の実装が標準として提供されているのはGoの良さだよね。

最近、その中の [golang.org/x/oauth2/clientcredentials](https://pkg.go.dev/golang.org/x/oauth2@v0.36.0/clientcredentials) に気づいた。

OAuth2.0 Client Credentials Grants typeのクライアント周りのあれこれをやってくれて、めちゃ便利。

```go
config := clientcredentials.Config{
    ClientID:     "my-client-id",
    ClientSecret: "my-client-secret",
    TokenURL:     "https://auth.example.com/oauth2/token",
}

client := config.Client(ctx)

resp, err := client.Get("https://api.example.com/resources")
```


`golang.org/x/oauth2/clientcredentials` パッケージ以外も見てみたら、

- [amazon](https://pkg.go.dev/golang.org/x/oauth2@v0.36.0/amazon)
- [facebook](pkg.go.dev/golang.org/x/oauth2@v0.36.0/facebook)
- [gitlab](https://pkg.go.dev/golang.org/x/oauth2@v0.36.0/gitlab)
- [fitbit](https://pkg.go.dev/golang.org/x/oauth2@v0.36.0/fitbit)
- [paypal](https://pkg.go.dev/golang.org/x/oauth2@v0.36.0/paypal)

などなどのOAuthプロバイダごとのpackageがあって、それぞれエンドポイントのURLを[oauth2.Endpoint](https://pkg.go.dev/golang.org/x/oauth2@v0.36.0#Endpoint) で定数定義してあるだけ。   
(var定義なので、やろうと思えば書き換えられるので定数ですらないんだけど...)

と思ったら、[golang.org/x/oauth2/endpoints](https://pkg.go.dev/golang.org/x/oauth2@v0.36.0/endpoints) ってのもあって、そこにはサービスごとのoauth2.Endpoint定数が鬼のように並んでいた...

[oauth2/endpoints/endpoints.go at master · golang/oauth2](https://github.com/golang/oauth2/blob/master/endpoints/endpoints.go)

いくらbatteries includedっていっても限度があるやろ...
今は新規にOAuth2プロバイダごとのpackageは作らずにこのoauth2/endpoints packageに追記していく方針らしいけど...

どんな基準で追加されてるんだろう...みたいなのを眺めてたら、[golang.org/x/oauth2/cern](https://pkg.go.dev/golang.org/x/oauth2@v0.36.0/cern) を見つけた。
そう、あのCERN(欧州原子核研究機構)。

CERNの[研究論文とか?オープンデータとか?にアクセスできるAPI](https://repository.cern/docs/reference/reference/)がpublicに提供されているみたい。知らなかった。
まぁ知っていたとして使う機会はなかったと思うけども。

どういう経緯でこのようにOAuth2プロバイダ固有のエンドポイントが定数として提供されているかわからないけど、こういう発見がたまにあるから、golang.org/x packageは面白い。
