---
title: ブログをGithub Pagesに移行した
---

昨年、ブログをGhost(Self-Hosting)に移行したけど、どうしてもプレーンテキストの安心感が忘れられず、Github PagesもといJekyllに戻してしまった。。。   
それほど手間はかからなかったが、一応メモ。


## なぜGithub pagesか

Github pagesであれば、わざわざ静的サイトをCIなどでgenerateせずともサイトに反映できるし、Github自体がなくなったとしても、プレーンテキストなのでどうにでもなる。
Githubのprivate repositoryでもGithub pagesは公開されるので、公開サイトの元データを外部に晒すこともなし。   
運用面でもとても楽。

Ghostに移行した理由はChromebookから編集できない、という点であったが、よくよく考えるとGithub上からもファイル追加・編集・削除できたのだった・・・   
半年後に気づくとは・・・ 😫

## Github pages用Gemfile

Github pagesは、Githubサーバー上でのJekyllやそのプラグインのバージョンに制限があるので、ローカルマシンで確認する場合はgemのバージョンをGithub上の環境に合わせておいたほうがよい。
ありがたいことに[github-pages](https://rubygems.org/gems/github-pages) というgemで依存関係がまとめられているので、これを使用する。

Gemfileはこんな感じ。

```ruby
source "https://rubygems.org"

gem "jekyll", "3.3.1"
gem 'github-pages', "115", group: :jekyll_plugins

gem "minima"
gem 'jekyll-compose', group: [:jekyll_plugins]

group :jekyll_plugins do
  gem "jekyll-feed", "~> 0.6"
end
```

jekyll-compose は便利コマンドがほしかったので入れといた。   
テーマは minima。
テーマはこだわっても仕方がないし、Webデザインをアピールしたいわけでもないので、このままでいいや。


## Ghostからの移行

Ghostはブログのデータをexportする機能がlabに含まれているので、それを使用。[参考](https://help.ghost.org/hc/en-us/articles/224112927-Import-Export-Data)   
Content-Security-Policy の設定が厳しいからかダウンロードダイアログが出なかったので、URLをコピって直接ダウンロードした。

GhostがexportするフォーマットはJSONであるため、Markdownに変換する必要がある。   
[jekyll_ghost_importer](https://github.com/eloyesp/jekyll_ghost_importer) を使うと、Ghostがexportしたデータから_postディレクトリ内にMarkdownファイルを生成してくれるので、それをありがたく使わせていただいた。   
ただ、バージョン0.4.0では、ブログエントリか静的ページかを問わず、全てのコンテンツを_postsに出力してしまうため、Jekyllで静的ページをブログエントリにしたくない場合は手作業でファイルを移動したりlayout設定を変更したりしなければならない。   
（エントリ内画像はどのようにexportされるかはexportしたブログエントリ内に画像が含まれていないので分からない）


## カスタムドメインのHTTPSについて

以前は、Ghostを適当な非公開portでlistenさせておいて、Nginxでリバースプロキシ & HTTPS化させていた。
Github pageにも https://blog.thara.jp でアクセスできるように以下のようにNginxを設定。

```
server {
    listen       443;
    server_name  blog.thara.jp;
(略)
    location / {
        proxy_pass https://{Github page URL};
        proxy_intercept_errors on;
        expires off;
    }
}
```

だいたいブログの構築して満足してしまって全然エントリを書けていない・・・

2017年は週1で更新できるように頑張る ✊
