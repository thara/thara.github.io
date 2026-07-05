---
title: モノレポのディレクトリ構成のいけすかなさ
date: '2026-07-05'
published: '2026-07-05'
---

最近は特に、アプリケーション作るときにモノレポを採用することが増えているように思う。

バックエンド実装やTerraformなどのIaC、Webフロントエンドまで。

それ自体は個人的に望ましいと思っているのだけれど、困るのはそのディレクトリ構成。
普通にやると、`backend/infrastructure/frontend` とか、`server/unity/web` とか、みたいな構成になるかもしれない。

でも個人的に[Screaming Architecture](https://blog.cleancoder.com/uncle-bob/2011/09/30/Screaming-Architecture.html)が望ましいと思っている。

ので、`backend/infrastructure/frontend` みたいな技術領域で分けられたディレクトリ構成は「叫んでない！」と感じてしまう。

Go単体であれば、シンプルなアプリケーションなら[Flat Application Structure](https://www.calhoun.io/flat-application-structure/)でいいとして、もう少し複雑なアプリケーションなら[こんな感じのコードアーキテクチャ](https://medium.com/sellerapp/golang-project-structuring-ben-johnson-way-2a11035f94bc)(domain packageがrepo root) でいいだろうけど。
<sub><sup> (流石に[golang-standards/project-layout](https://github.com/golang-standards/project-layout) を新規に採用しようとする人はもういないよね)</sup></sub>

これにGo以外のも含めてモノレポにしようとした時には、異なる技術領域が絡み合ってしまう。

じゃあ技術領域ごとに分けよう、となると、`backend/infrastructure/frontend` みたいな分け方になって、「叫び」を感じられなくなってしまう。

ソフトウェアに対する認識というか観点といったものは一つではないので、ファイルシステムという1つの木表現に押し込むこと自体がそもそも無理、という話ではある。
それ前提だと、特定の技術領域は拡張子やパッケージマネージャが求めるディレクトリ構成に強制的に従わざるを得ないので、その強制力が働かないディレクトリ構造においては、技術領域以外の観点で構成した方が良いのではないか、と思うのだが、開発体験とかも含めるとそんなに単純な話ではないよなぁと思ったり。

ここら辺に、モノレポ採用時のディレクトリ構成のいけすかなさ、みたいなのをどうしても感じてしまう。

論理的なファイル構成を表現できるツールがあればいいかも？
でもそれってEclipseとかのIDEが搭載していて、個人的にあんまり嬉しくなかったんだよな...みたいなのもあり、いまだに自分の中の答えはない。
