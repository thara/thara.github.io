---
title: herdrをセットアップした
date: '2026-11-01'
published: '2026-11-01'
---

最近、以下のようなポストをちらほら見て、最近自分の課題感もあって、自分もherdr真面目にセットアップするか、と重い腰を上げた。

<blockquote class="twitter-tweet"><p lang="ja" dir="ltr">ここ1ヶ月くらいはChief of Staff というか、自分が常に話しかけるエージェントを1個おいて、これはずっと同じセッションで動いている（たまに /compact する）、必要な仕事は herdr で他のエージェントを立ち上げて元のエージェントは全部それを把握してる、みたいにしててとてもよい #rebuildfm https://t.co/8bkLPyL519</p>&mdash; Tatsuhiko Miyagawa (@miyagawa) <a href="https://x.com/miyagawa/status/2113515792948158772?ref_src=twsrc%5Etfw">October 23, 2026</a></blockquote>
<script async src="https://platform.x.com/widgets.js" charset="utf-8"></script>

<blockquote class="bluesky-embed" data-bluesky-uri="at://did:plc:fuqczowdfxvncehn5iu65hus/app.bsky.feed.post/3myu377v3qv2i" data-bluesky-cid="bafyreiflyowtzggdcrnzlnjet7ugeyquskgfr4uxl6ldvvxd2vmndgtdwe" data-bluesky-embed-color-mode="system"><p lang="">ツールが働き方を変える、ということを改めて実感している。GitとGitHubでも働き方は大きく変わってた訳だし、結局今の形にすぐ慣れていくもんだなー。herdrに慣れてからエージェントの並列数が大きく変化したし、合わせて同時進行で複数作業するのにも割と慣れてきた。</p>&mdash; joker1007 (<a href="https://bsky.app/profile/did:plc:fuqczowdfxvncehn5iu65hus?ref_src=embed">@joker1007.net</a>) <a href="https://bsky.app/profile/did:plc:fuqczowdfxvncehn5iu65hus/post/3myu377v3qv2i?ref_src=embed">October 27, 2026 at 7:55 PM</a></blockquote><script async src="https://embed.bsky.app/static/embed.js" charset="utf-8"></script>

## 課題感

普段、tmux + OpenCode使って複数タスクをそれぞれでgit worktree立ててAgentic Codingしてる。
大体1タスクごとに、以下のwindow構成になる。

- OpenCode
- bash
- vim (コードレビュー用途）

(pane分けるのは必要に応じて一時的に作るだけで、普段は1 window 1 paneでやってる）

で、今の自分だと3並列ぐらいタスクを回してるんだけど、それぞれのエージェントの状態を把握するのにwindow切り替えて様子見なくちゃいけない。   
OpenCode -> bash -> vim -> OpenCode -> bash -> vim -> OpenCode ... ってタブ切り替えするのが煩わしい。

tmuxのセッションごと1タスクにしてもいいんだけど、タスク以外の仕事もあって、それと並ぶとさらに混乱を及ぼしそう。

## herdrでの解決策

ということで、herd使ってworkspace = 1タスクにしておいて、さらにエージェントの状態も可視化できるようになるし、ある程度なんとかなるんじゃないか、と思って以前使ってみたんだけど、tmuxがあまりにも手に馴染みすぎてしまってそこで手が止まってた。

上述のpost見て、今度はちゃんとセットアップしてみよ、と思って自分のtmuxの使い勝手をherdrで再現するようにセットアップしてみた。

## herdr設定

自分の[tmux.conf](https://github.com/thara/dotfiles/blob/master/config/.tmux.conf)では色々設定してあるけど、herdrでネックになってたのは以下の設定。

- prefix key
    - `ctrl+q` に変更
- pane分割方法
    - 左右分割を `prefix+|` に
    - 上下分割はデフォルト(`prefix+-`) のまま
- タブ作成時に元のペインと同じディレクトリで開く
    - `[terminal]` の `new_cwd` に `"follow"` を設定
- タブ作成時にタブの名前を聞かない
    - めっちゃうざかったやつ
    - `[ui]` の `prompt_new_tab_name` に `false` を設定
- 設定ファイルのリロード
    - tmuxでは `prefix+r`
    - herdrは `prefix+r` がリサイズモードなので `prefix+alt+r` に退避（多分使わない）した上で、設定リロードに `prefix+r` を設定

あんまり大したものはなかった...

これでAgentic Codingしてみて、細かいところを調節していく感じかな。

最初にherdr試した時はtmuxからの置き換えを考えて評価してたんだけど、普段はtmuxで、Agentic Codingするときにherdr使う、という方向が良さそうな気がしてきた。

まぁ、色々試してみる。
