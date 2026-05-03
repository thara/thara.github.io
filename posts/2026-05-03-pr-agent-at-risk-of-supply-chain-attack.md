---
title: PR Agentにはサプライチェーン攻撃のリスクがある
date: '2026-05-03'
published: '2026-05-03'
---

[PR-Agent](https://github.com/The-PR-Agent/pr-agent)は、LLMを使用してGitHubのPull Requestのレビューを自動化するOSSのGitHub Actionsであり、いくつかの企業で採用している例を見かける。

しかし、このPR-Agentの実装にはサプライチェーン攻撃のリスクがあり、そのリスクを認識した上で採用する必要がある。

1. Dockerベースイメージをpinningしていない
2. pip installをHash-checking Modeで実行していない
3. Dockerイメージのビルドプロセスが不明瞭


## 1. Docker ベースイメージをpinningしていない

PR-AgentはDockerコンテナを使用したActionだが、そのDockerfileのベースイメージが以下のようにpinningされていない。

```
FROM python:3.12.10-slim AS base
```

via [Dockerfile.github_action#L1](https://github.com/The-PR-Agent/pr-agent/blob/009ba5a116c4d3273368a6dc53a4efdb7904d519/Dockerfile.github_action#L1)

Dockerレジストリのタグは上書き可能であり、悪意あるバージョンのDocker ベースイメージが使用される可能性がある。

ref. [Building best practices | Docker Docs](https://docs.docker.com/build/building/best-practices/#pin-base-image-versions)

## 2. pip installをHash-checking Modeで実行していない

前述のDockerfileではpip installしているが、それに `--require-hashes` オプションが付与されておらず、requirements.txtの各依存ライブラリにも `--hash` オプションが付与されていないため、[Axiosのnpmパッケージ侵害](https://www.trendmicro.com/ja_jp/research/26/d/axios-npm-package-compromised.html)のようなパッケージ侵害の影響を受ける可能性がある。

via [Dockerfile.github_action#L6-L8](https://github.com/The-PR-Agent/pr-agent/blob/009ba5a116c4d3273368a6dc53a4efdb7904d519/Dockerfile.github_action#L6-L8)


## 3. Dockerイメージのビルドプロセスが不明瞭

1,2の対策としてPR-Agentの特定のバージョンを使用する方法がドキュメントに記載されている。

```yaml
steps:
  - name: PR Agent action step
    id: pragent
    uses: docker://pragent/pr-agent@sha256:a0b36966ca3a197ca739fa1e65c16703076fc1c744cd423ca203b8c21707d71c
```

via [SECURITY.md](https://github.com/The-PR-Agent/pr-agent/blob/009ba5a116c4d3273368a6dc53a4efdb7904d519/SECURITY.md#specific-release-version)

Dockerイメージのversionはpinningされており一見問題なさそうに見えるが、このDockerイメージのビルドプロセスが不明瞭であり、その妥当性を第三者が検証できない他、image自体の検証もしづらい。

(同様の問題提起を[issue comment](https://github.com/The-PR-Agent/pr-agent/issues/2306#issuecomment-4304993299)でされている方もいる）

## 3rd-party GitHub Actionsを採用する際の注意点

[pinact](https://github.com/suzuki-shunsuke/pinact)やGitHubのEnforce SHA Pinning機能によって、3rd-partyのGitHub Actionsを使用する際にpinningすることの重要性は広く普及してきたように思える。

しかし、GitHub ActionsをpinningしたとしてもそのGitHub Actions自体のリスク評価についてはまだ認識が広がっていないように思える。

3rd-partyのGitHub Actionsを採用する際には、自分は以下のように注意している。

- Actions内でプログラムを動的にダウンロードしている場合は、ダウンロードしたプログラムのチェックサム検証をしているか確認する
    - 実行バイナリやインストールスクリプトなど、悪意ある第三者に上書きされる可能性がある
- Actions内でGit リモートリポジトリからチェックアウトしている場合は、コミットハッシュ指定しているか確認する
    - ブランチやタグ指定されている場合、悪意ある第三者に上書きされる可能性がある
- Actions内でパッケージインストールしている場合は、パッケージがpinningされていることを確認する
    - 最近のパッケージマネージャにはpinningする機能があるはず
    - 逆説的に言えば、pinning機能を持っていないパッケージマネージャを使用してはいけない
- Dockerコンテナを使用したActionの場合は、Dockerイメージのビルドプロセスに透明性があるかを確認する
    - プライベートな環境でのビルドの場合、そのビルドが侵害されてないことを証明するのが難しい

## まとめ

3rd-partyのGitHub Actionsは気軽に使えて開発生産性に大きく寄与する一方、LLMの発展によりサプライチェーン攻撃などのセキュリティリスクが高まっている現代においては、そのリスクを適切に評価した上でその採用の是非を決定する必要がある。

...嫌な世の中になってしまったな...
