---
title: setup-trivyの侵害コミットで自分の名前が詐称されていた件
date: '2026-03-22'
published: '2026-03-22'
---

[2026年3月19日の Trivy 再侵害の概要と対応指針 – やっていく気持ち](https://diary.shift-js.info/trivy-compromise/) を読んで。

この記事では、 [aquasecurity/setup-trivy](https://github.com/aquasecurity/setup-trivy) と [aquasecurity/trivy-action](https://github.com/aquasecurity/trivy-action) にクレデンシャル窃取コードが注入するためにgit commitが以下のように偽装されたと書かれている。

> 著者名の詐称: setup-trivy側のcommit 8afa9b9 は実在のコントリビュータ “thara” を、trivy-action側の ddb9da4 はAquaチームメンバー “DmitriyLewen” を author として詐称している

**はい、このsetup-trivy側のコントリビュータ「thara」とは自分のことです。**

## コントリビュート経緯

2026-01-13、setup-trivyの存在を知った自分はそれを使用しても問題ないかを検証するためにリポジトリ内のaction.yamlを確認しにいった。   
そこでsetup-trivyが[aquasecurity/trivy](https://github.com/aquasecurity/trivy)からcheckoutしてくるインストールスクリプト( `contrib/install.sh`)が常に最新のものを使用しており、それによってサプライチェーンアタックの可能性があることに気づき、以下のissueを立てた。

[Pin Trivy install script checkout to commit hash · Issue #27 · aquasecurity/setup-trivy](https://github.com/aquasecurity/setup-trivy/issues/27)

提案した修正方針が既存のコントリビュータに認められたため、続けてPRを作成した。

[Pin Trivy install script checkout to a specific commit by thara · Pull Request #28 · aquasecurity/setup-trivy](https://github.com/aquasecurity/setup-trivy/pull/28)

自分の貢献はこの程度だが、この時のcommit logの情報を悪用されたようだ。

## 詐称されたコミット

先の記事では、以下のように自分が詐称されたコミットが記載されている。

```
// GET /repos/aquasecurity/setup-trivy/commits/8afa9b9f9183b4e00c46e2b82d34047e3c177bd0
{
  "sha": "8afa9b9f9183b4e00c46e2b82d34047e3c177bd0",
  "commit": {
    "author": {"name": "Tomochika Hara", "email": "github@thara.dev", "date": "2026-01-15T10:21:20Z"},
    "committer": {"name": "GitHub", "email": "noreply@github.com", "date": "2026-01-15T10:21:20Z"},
    "message": "Pin Trivy install script checkout to a specific commit (#28)",
    "verification": {"verified": false, "reason": "unsigned"}
  },
  "author": {"login": "thara"},
  "committer": {"login": "web-flow"}
}
```

比べて、自分が作成したPRに含まれるコミット(の一つ)は以下。

```
// GET /repos/aquasecurity/setup-trivy/commits/0fa8dc73b42c00a5983466410dcdab3a4071f21e より抜粋
{
    "sha": "0fa8dc73b42c00a5983466410dcdab3a4071f21e",
    "commit": {
        "author": {"name": "Tomochika Hara", "email": "git@thara.dev",  "date": "2026-01-13T08:27:30Z",},
        "committer": {"name": "Tomochika Hara", "email": "git@thara.dev", "date": "2026-01-14T08:48:52Z"},
        "message": "fix: pin Trivy install script to specific commit hash...(略)",
        "verification": {
            "payload": "(略)",
            "reason": "valid",
            "signature": "(略)",
            "verified": true,
            "verified_at": "2026-01-14T08:48:57Z"
        }
    },
    "author": {"login": "thara"},
    "committer": {"login": "thara"}
}
```

自分の正規のコミットと攻撃者が詐称したコミットを比べてみる。

- commit.{author|commiter}.emailが `git@thara.dev` ではなく `github@thara.dev` になっている
    - `github@thara.dev` は自分のGitHubアカウントの[プライマリメールアドレス](https://docs.github.com/ja/account-and-profile/how-tos/email-preferences/changing-your-primary-email-address)
    - 先の記事中でも触れられているけれど、GitHubのcommit authorはメールアドレスさえ合わせれば誰でも詐称できる（署名なしの場合）
- commit.verification.verifiedがfalseになっている
    - 自分は常に[signed commit](https://docs.github.com/ja/authentication/managing-commit-signature-verification/signing-commits)している
    - 自分視点では自分が作成したcommitではないと明らかだけど、第3者から見たら確実なことは言えない、か...
- committer.login が `thara` ではなく `web-flow` になっている
    - GitHub上からマージした際に付与されるもの


## 自分への影響

直接的な実害はないが、個人のレピュテーション上の影響はあった。

一応確認したが、GitHubアカウントやthara.devのメールアカウント（Google Workspaceで独自ドメイン設定しているやつ）にも不正なログインの兆候はなかった。

あえて言えば、[詐称コミット](https://github.com/aquasecurity/setup-trivy/commit/8afa9b9f9183b4e00c46e2b82d34047e3c177bd0)を見ると自分が攻撃したみたいで気分が悪い、ぐらいだろうか。

ひとまず、プライベートや業務で使用しているメールアドレスをgit commitに使ってなくてよかった。

## 感想

trivyが3週間前に侵害された時に「これって収束したと言えるんかな〜」と思っていたが、やっぱり次があった...   
今まで儀式のようにsigned commitの設定をしていたけど、今回の件で改めて重要性を実感した。今後も続けていきたい。

## まとめ

[自分がsetup-trivyに入れた変更](https://github.com/aquasecurity/setup-trivy/pull/28) によって、v0.2.6、厳密にはcommit hash [3fb12ec](https://github.com/aquasecurity/setup-trivy/commit/3fb12ec12f41e471780db15c232d5dd185dcb514) は安全に使えるので、setup-trivyを使う場合はそれを使ってください。

setup-trivy以外にも、サードパーティのバイナリやGitHubアクションを使う場合は以下に気をつけましょう。

- バイナリダウンロード後にchecksumを検証する
    - checksumは動的にダウンロードするのではなく、あらかじめchecksumファイルをダウンロードする or CIの設定に直値を埋め込む
- setup-xxx を使用する際はGitHubアクションの実装を確認する
    - バイナリダウンロード後に検証しているか、インストールスクリプトの正当性が保証されているか、など

