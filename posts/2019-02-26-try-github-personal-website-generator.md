---
title: GitHubのPersonal website generatorを使ってみた
description: 先日、GitHubがdevドメインで `Personal website generator` なるものを公開したので、試しに使ってみた。
published: '2019-02-26'
tags: [github, jekyll]
---

先日、GitHubがdevドメインで `Personal website generator` なるものを公開したので、試しに使ってみた。

[GitHub personal website generator](https://github.dev/)

仰々しく（？）ブランディングされているけれど要は単なるJekyllテーマで、 GitHub APIに使ってGitHubユーザーのプロフィールやリポジトリの一覧を取得し、それを元に開発者向けの簡易なポートフォリオサイトを作れる、というもの。

作り方は簡単、 [github/personal-website](https://github.com/github/personal-website) をforkしてGitHub Pagesとして公開するだけ。 カスタムドメインなどの設定は、既存のGitHub Pagesと全く同じ。

とりあえず、forkした直後でローカルで `jekyll serve` するとこんな感じ。

![GitHub personal website sample]({{ site.url }}/images/github-personal-website-sample.png)

`My Projects` の一覧はリポジトリをアルファベット順に9つ表示されている。   
`My Interests` は `_config.yml` で固定で設定されているのを表示しているだけ。[^1]


## _config.ymlだけでできること

`_config.yml` を編集するだけで、レイアウトとスタイル、`My Interests`を変更できる。
逆にこれら以外の内容、例えば表示されるリポジトリなどを選択するには、テンプレートを直接変更する必要がある。

### Stacked Layout

[参考](https://github.com/github/personal-website#layout)

![GitHub personal website sample - Stacked Layout]({{ site.url }}/images/github-personal-website-sample-stacked.png)


### Dark Theme

[参考](https://github.com/github/personal-website#style)

![GitHub personal website sample - Dark Theme]({{ site.url }}/images/github-personal-website-sample-dark.png)

### Custom My Interests

[参考](https://github.com/github/personal-website#topics)

(デフォルトのテンプレートでは、4つまでしか表示されない)

![GitHub personal website sample - Custom My Interests]({{ site.url }}/images/github-personal-website-sample-my-interests.png)


## その他

ポートフォリオのページだけでなく、他のページを作成したり、ブログの記事を追加したりもできる。これはJekyllの機能そのまま利用しているだけなので、それほど特別なことはしてないはず。

## 感想

さくっといい感じのデザインのポートフォリオサイトが作れるのは、Webデザインが苦手な自分にとってはありがたい。
テンプレートを自由に編集することができる人であれば必要十分なポートフォリオサイトが作ることができそう。

一方で、リポジトリ一覧は大抵の人がそのままでは満足しないだろうし、My Interestsも手作業で設定する必要があるので、手軽さがもう一歩改善するとより多くの開発者に使ってもらえるんじゃないかと思う。[^2]

まぁ自分にはもう [thara.jp](https://thara.jp) というミニマムなサイトがあるので、新たにポートフォリオサイト作る意味はあんまりないけど。

GitHub上でつくったPersonal websiteをhiring有効にしている人とかMy Interestsとかの条件で検索できるようになると、エコシステムとして機能していくかもしれない。   

GitHub社がどこまで考えてこのPersonal websiteを公開したのかはまだわからない。今後の動向に期待したい。

[^1]: CSSとかWeb DesignとかSassとかはデフォルト値で、自分自身が興味あるわけじゃない。
[^2]: リポジトリの一覧をpinned repositoriesにすればよいのでは、と考えたけど、そもそも[github-metadata gem](https://github.com/jekyll/github-metadata)がpinned repositoriesに対応していないっぽくて詰んでる
