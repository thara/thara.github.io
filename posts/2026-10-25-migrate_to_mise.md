---
title: dotfiles管理をmiseに移行した
date: '2026-10-25'
published: '2026-10-25'
---

[私物Macと仕事Macのdotfilesとアプリ管理をmise bootstrapで管理しやすくした - $shibayu36->blog;](https://blog.shibayu36.org/entry/2026/10/13/173000) を参考に、
今まで[GNU Stow](https://www.gnu.org/software/stow/)でやってたdotfiles管理を[mise bootstrap](https://mise.jdx.dev/bootstrap.html)に移行した。

元々、以下のようなStow用に（Stowにおける）パッケージごとにディレクトリ切ってた。

- agents
- config
- etc
- git
- macos
- shell
- vim

miseでもこれを踏襲して、それぞれを[dotfile_groups](https://mise.jdx.dev/dotfiles/groups.html) にした。

こういうの。

```
[dotfile_groups.shell]
root = "shell"
target = "~"
mode = "symlink-each"

[dotfile_groups.vim]
root = "vim"
target = "~"
(略)
```

なので、移行内容としてはmiseの設定ファイルを置いたり、移行用のスクリプトを用意するだけで済んだ。

どちらもLLM任せ。いやー、いい時代ですね。

ツールのバージョンとかもasdf使ってたけど、miseに移行して良いかも。
