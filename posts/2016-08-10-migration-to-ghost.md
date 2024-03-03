---
title: ブログをGhostに移行した
date: '2016-08-10 17:11:17'
published: '2016-08-10'
tags:
- ghost-tag
---

最近買ったChromebookから手軽にブログ投稿することを目指して、[Ghost](https://ghost.org/)（Self-Hosting)に移行した。

移行した、といってもコンテンツはめんどくさいのであんまり頑張って移行しない。
Themeは凝りだすといつまでも探し続けてしまい、結局自作することになるのは目に見えているので、ひとまずデフォルトのままでいくことにした。

公式ドキュメントどおりにすればインストールは簡単だが、以下の点を工夫した。

- nginxの設定には[Mozilla SSL Configuration Generator](https://mozilla.github.io/server-side-tls/ssl-config-generator/)をベースにした
  - limit_rateなどのDDoS対応も少々
  - Content-Security-Policyも設定した
    - 厳しくしすぎてJSが動かなかったり、Web Fontを取ってこれなかったりしたけど解決した
- [Let's Encrypt](https://letsencrypt.org/)でHTTPS対応した
  - [itamae](https://github.com/itamae-kitchen/itamae)でのプロビジョニングの際にACMEを効かせるのに一苦労
    - プロビジョニング中はACME用に80ポートでlistenしといて、後々443ポートの設定ファイルを書く、とかしないといけない
  - 自動更新はまだ仕掛けてない
  - 自前メールサーバーも対応させたい
- データベースにsqliteではなく[MariaDB](https://mariadb.org/)を設定した
  - sqliteはなんか不安だし、MySQLは仕事でよく使うので。
  - GhostからはMySQL扱い

こんな感じ。Let's Encryptで生成した証明書を使って、WebブラウザがHTTPSでサイトにアクセスしてくれたときは胸熱だった。

ChromebookでWebに繋がっていればどこでもブログを手軽に書けるようになったので、ちまちま更新していきたい。
