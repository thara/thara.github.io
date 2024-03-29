---
title: 職務経歴書
---

原 知愛
Tomochika Hara

- Github: [thara](https://github.com/thara)
- Blog: [Tomochika Hara's Blog](https://blog.thara.jp)
- Twitter: [@zetta1985](https://twitter.com/zetta1985)
- Twitch: [tharadev](https://twitch.tv/tharadev)
- LinkedIn: [Tomochika Hara](https://www.linkedin.com/in/tomochikahara/)

## 概要

バックエンドエンジニアとして10年以上の経験があり、クラウドプラットフォーム上のスケーラブルなアーキテクチャの設計とアプリケーション開発に携わってきた。
特にモバイルゲーム開発においては、サーバーチームリードエンジニアとしてアーキテクチャ設計およびパフォーマンスチューニングを行い、DAU10万・総ダウンロード数1000万を超えるタイトルへの成長に貢献した。

フルリモートワークでのモバイルアプリ開発を経て、現在はメタバースプラットフォームのサーバーサイドエンジニアに従事。

## スキル概要

- AWSやGCP上でのバックエンドサービスの開発
- 拡張性/保守性/高負荷を考慮したアーキテクチャ/アプリケーション設計
- 投資対効果を考慮した自動化テストの設計と実装
- フルリモートワーク/非同期コミュニケーションを前提とした効率的なプロダクト開発
- 開発チーム外を巻き込んだプロジェクト遂行

## スキル詳細

- プログラミング言語
    - Go, Swift, Rust, Python, C/C++, Ruby, Java, Kotlin, Scala, JavaScript/TypeScript
- インフラストラクチャ
    - Amazon Web Services, Google Cloud Platfrom
    - MySQL, Redis
    - Jenkins, Docker, Circle CI, GitHub Actions
- フレームワークその他
    - Django, Ruby on Rails, Spring framework
    - RESTful, Protocol Buffers, gRPC

## 経歴

### クラスター株式会社 (2020-09 〜 現在)

メタバースプラットフォームのバックエンドの設計・開発・保守に従事。

- サーバーリリースワークフローのSlack workflowによる半自動化を実施した。
- GoによるglTF/VRMのバリデーションおよび圧縮処理を実装した。
- NodeJS/C#で書かれたUGCパイプラインをGoでリライトし、GitHub ActionsによるCI/CD環境を構築した。
- 3D空間の複数ユーザー間の同期処理を担うサーバーの内部アーキテクチャの設計と既存サービスからの移行を行なった。

### 弥生株式会社 (2020-07 〜 2020-08)

株式会社Misoca吸収合併による転籍

### 株式会社Misoca (2018-07 〜 2020-06)

請求書作成サービスMisocaのモバイルアプリのバックエンドAPIの開発およびAndroid/iOS向けのライブラリ開発に従事。

- Android/iOSに向けたKotlin/Swift製の税率計算ライブラリを実装し、サーバーサイドの実装と同じテストデータを用いたCIを構築した。
- 2019/10の軽減税率制度開始に向け、複数プラットフォームに向けて古いアプリバージョンとの後方互換性を保ちつつ機能変更をするための、モバイルアプリ向けのフィーチャーフラグや段階的なデータマイグレーションを設計した。

### WonderPlanet (2013-10 〜 2018-06)

モバイルゲーム向けサーバーサイドアプリケーションの設計・開発・保守、およびAWSやGCPなどのクラウドプラットフォームを使用した高負荷に対応したシステム設計・構築・運用に従事し、
リリース後3ヶ月で100万ダウンロードを達成したタイトルのバックエンドサービスの開発を主導した。
複数プロジェクトでサーバーチームのリーダーを務め、メンバーへの教育や技術的指導を行い、自発的に行動することでリーダーがSPOFにならないようなチーム体制を整えた。

- バックエンドサービスのアーキテクチャ設計と構築、API開発, その他ツール開発を主導した。
- 差分更新やロールバックをサポートした非エンジニア向けのマスターデータ管理ワークフロー設計を行った。
- リリース直後、想定外の負荷により長期メンテナンスに入ったが、Amazon RDS for MySQLのパフォーマンスチューニングやDB構成の変更などの改善を施し安定稼働させた。
- MySQLのインデクシングがボトルネックになるような高頻度の更新処理への負荷対策のために、Amazon DynamoDBを部分的に使用し、DBへの負荷を低減した。
- ガチャのアルゴリズムにWalker's Alias Methodを採用し、定数時間での復元抽出を実現した。
- Redisの集合演算を用いてフレンドとその他プレイヤーを混合したリアルタイムのスコアランキング生成を実現した。
- 数時間かかると考えられていた全プレイヤーに対するグルーピング処理を、Apache AirflowとRedisを用いたPython製のバッチ処理を実装し、30分前後の処理時間に収めた。
- リリース2ヶ月前というタイミングでPhoton Serverのマッチング処理ではより良いゲーム体験を提供できないと判明したため、Python製の独自マッチングサービスを1ヶ月で実装した。
- Jenkins サーバーの構築と運用および、モバイルアプリ向け自動ビルド環境を構築した。
- Go 製 Google PlayStore Voided purchase 集計サービス/ Slack ボットを開発した。
- cocos2d-x アプリのビルド時間のプロファイルと改善を行なった。
- Photon Serverの複数クラスタの負荷分散のために、pre-fork形式のWebサーバーを参考にしたモニタリングサービスをC++で実装した。
- Photon Serverの1クラスタ向けの複数台のAmazon EC2インスタンスをワーカーとするベンチマークツールをC#で実装した。

### ニューソン株式会社 (2008-04 〜 2013-09)

業務用Webアプリケーション設計/開発を担当した。

- 独自JavaScriptフレームワークを開発し、AjaxによるインタラクティブなWebアプリケーションを開発した。
- Javaで開発した多くの入力パターンを受け付ける化学化合物の計算ライブラリを実装し、データ駆動テストによって高品質を実現した。

## 個人プロジェクト

- [SwiftNES](https://github.com/thara/SwiftNES): Swift製のクロスプラットフォーム向けNESエミュレータ
- [erajp](https://github.com/thara/erajp): Rust製 日本の元号変換ライブラリ
- [rust_aliasmethod](https://github.com/thara/rust_aliasmethod): Walker's Alias MethodのRust実装

## その他の活動

- [JUnit実践入門 体系的に学ぶユニットテストの技法](http://gihyo.jp/book/2012/978-4-7741-5377-3)の出版前レビュー
- [SpriteKitではじめる2Dゲームプログラミング Swift対応](http://www.shoeisha.co.jp/book/detail/9784798139517) の執筆
  - 第2章 Swiftの基本, 第8章リバーシを作ろう--AIの作り方 を担当
- [iOSDC Japan 2019: Swiftでつくるファミコンエミュレータのススメ](https://fortee.jp/iosdc-japan-2019/proposal/92904657-beda-46fe-8ecb-b27c75ee0f16) (LT発表)
