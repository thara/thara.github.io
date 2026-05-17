---
title: ledger向けスクリプト作り始めた
date: '2026-05-17'
published: '2026-05-17'
---

[先週書いた記事](../2026-05-10-plain-text-accounting.md)で書いた通り、[Plain Text Accounting](https://plaintextaccounting.org/)のために、[howeyc/ledger](https://github.com/howeyc/ledger)のimport機能がサポートしているフォーマットに変換するスクリプトをぼちぼち書き始めた。

howeyc/ledgerがサポートしているフォーマットはこういうの。

```csv
Transaction Date,Description,Amount
01/12/22,Dominoes Pizza     HOUSTON TX,12.34
01/23/22,Dominoes Pizza     PEARLAND TX,14.34
01/02/22,Half Price Books   AUSTIN TX,5.24
```

各金融機関のWebサイトからエクスポートしたCSVファイルをこれに変換するためのスクリプトを書いた。

```bash
./import-ufj.py --description-map descriptions.csv /path/to/exported.csv > import.csv
```

descriptions.csvは、元データの「取引概要」みたいなのを、ledgerに入れたいDescriptionに変換するためのマッピング情報。

ledgerには↓のようにimportする。

```bash
ledger -f ledger.bat import UFJ import.csv >> ledger.bat
```

コマンドライン引数に `Account` を設定しないといけないんだけど、元のledgerファイルにその `Account` がないと失敗するので、↓のように初期状態を入れておいた。

```
2026/04/01 Opening Balances
    Assets:Bank:UFJ           50000.00
    Equity:Opening Balances
```

importしたら全部 Equity:Opening Balances に紐づいてしまったので、まだなんか間違ってそうな気がしないでもないが、とりあえず最初の一歩ということで...

スクリプトは公開しても良いかな、って思ってたけど、自分が使っている金融機関がモロにバレるので、とりあえず非公開。

いろんな金融機関のCSVファイルのフォーマットが公開されてれば、片っ端から実装するだけで済むんだけど。

API連携といい、こういうところでオープンじゃないのが辛いねぇ...
