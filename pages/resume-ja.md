---
title: 職務経歴書
---

原 知愛
Tomochika Hara

- GitHub: [thara](https://github.com/thara)
- Blog: [Tomochika Hara's Blog](https://blog.thara.jp)
- X: [@zetta1985](https://x.com/zetta1985)
- Bluesky: [@thara.jp](https://bsky.app/profile/thara.jp)
- LinkedIn: [Tomochika Hara](https://www.linkedin.com/in/tomochikahara/)

## 概要

バックエンドエンジニアとして10年以上の経験があり、クラウドプラットフォーム上のスケーラブルなアーキテクチャの設計とアプリケーション開発に携わってきた。
特にモバイルゲーム開発においては、サーバーチームリードエンジニアとしてアーキテクチャ設計およびパフォーマンスチューニングを行い、DAU10万・総ダウンロード数1000万を超えるタイトルへの成長に貢献した。

フルリモートワークでのモバイルアプリ開発を経て、現在はメタバースプラットフォームのサーバーサイドエンジニア/ソフトウェアアーキテクトに従事。

## スキル概要

- AWSやGCP上でのバックエンドサービスの開発
- 拡張性/保守性/高負荷を考慮したアーキテクチャ/アプリケーション設計
- 投資対効果を考慮したテストの設計と実装
- フルリモートワーク/非同期コミュニケーションを前提とした効率的な開発
- 開発チーム外を巻き込んだプロジェクト遂行
- テクニカルライティング/技術文書の作成

## スキル詳細

- ★★★★★: 深い専門性を持ち、設計・実装からチームへの指導や最適化が可能
- ★★★★☆: 実務で主力として使用、 複雑な課題を自立して解決可能
- ★★★☆☆: 実務で継続的な使用経験あり、標準的なタスクを安定して遂行可能
- ★★☆☆☆: 実務での使用経験あり、基本的な課題を解決可能
- ★☆☆☆☆: 小規模タスクで使用可能

### プログラミング言語
- **★★★★★:** Go
- **★★★★☆:** Python, Ruby, JavaScript/TypeScript  
- **★★☆☆☆:** Kotlin, Swift, C++, C#, Scala
- **★☆☆☆☆:** Rust

### フレームワーク・ライブラリ
- **★★★★☆:** Goa, Django
- **★★★☆☆:** Ruby on Rails
- **★★☆☆☆:** React, Vue.js  
- **★☆☆☆☆:** RxSwift, UniRx  

### インフラ・ツール
- **★★★★☆:** MySQL, Redis, Gradle, CircleCI, GitHub Actions  
- **★★★☆☆:** AWS, GCP, nginx, Docker, Terraform, Jenkins
- **★★☆☆☆:** Photon Engine, Apache Airflow

### プロトコル・フォーマット
- **★★★★★:** RESTful, Protocol Buffers, gRPC  
- **★★★☆☆:** MQTT  
- **★☆☆☆☆:** glTF, VRM  

## 経歴

### クラスター株式会社 (2020-09 〜 2025-02)

- 役割: シニアソフトウェアエンジニア、ソフトウェアアーキテクト(サーバーサイド)
- 利用技術: Go, AWS, MySQL, Redis, Protocol Buffers, gRPC, MQTT, Docker, GitHub Actions, Circle CI, Gradle, Unity, C#, glTF, VRM

メタバースプラットフォーム cluster のバックエンドの設計・開発・保守に従事。

- 空間上のプレイヤーのレベルやスキルを保存できる「セーブ機能」を実装した。
    - [バーチャルSNS「cluster」プレイヤーのレベルやスキルを保存できる「セーブ機能」をリリース | クラスター株式会社のプレスリリース](https://prtimes.jp/main/html/rd/p/000000096.000017626.html)
- 3D空間の複数ユーザー間の同期処理を担うサーバー(room server)の内部アーキテクチャの設計を行なった。
    - [Go Conference 2022 Spring | メタバースを支える技術 ～UGCに溢れる3D空間のリアルタイム同期を支えるGo〜](https://gocon.jp/2022spring/sessions/a6-c/)
- 既存VerneMQから上記room serverへの移行をメンテナンス期間なしで実施した。
    - [clusterのリアルタイム通信サーバーの漸進的な進化 - Cluster Tech Blog](https://tech-blog.cluster.mu/entry/2022/04/13/143058)
- room serverからRDBに依存する箇所をマイクロサービスとして切り出した。
- Node.js/C#で書かれたUGCパイプラインをGoでリライトし、GitHub ActionsによるCI/CD環境を構築した。
- GoによるglTF/VRMのバリデーション、圧縮処理、3Dモデルに含まれるテクスチャのGPU native imageへの変換を実装した。
- アプリ上で空間上のワールドの状態を保存/公開する機能の概念設計および実装を行なった。
    - [新機能「ワールドクラフト」リリース！スマホで・誰でも・友達と・簡単にメタバース空間が作れるように | クラスター株式会社のプレスリリース](https://prtimes.jp/main/html/rd/p/000000126.000017626.html)
- 上記ワールドクラフトで使用可能なアイテムのUGCストア機能の設計および実装を行なった。
    - [自らの手で作り上げていく、cluster発のバーチャル経済圏の確立に向けてユーザー待望の新機能「ワールドクラフトストア」ついにリリース！ | クラスター株式会社のプレスリリース](https://prtimes.jp/main/html/rd/p/000000144.000017626.html)
- [ProtocolBuffersスキーマ運用の改善](https://tech-blog.cluster.mu/entry/2023/04/27) を行なった
- モノリスアプリケーションをマイクロサービスに分割するプロジェクトを主導した。
    - [無停止で機能開発を継続した、clusterのシステム分割事例 - Cluster Tech Blog](https://tech-blog.cluster.mu/entry/2023/12/26)
- 開発チームの設計/テクニカルライティングスキルの向上を目的に委員会を立ち上げ、design docのテンプレート改善やガイドラインの策定を行なった。
- プラットフォーム運用向けWebアプリケーションのフレームワーク設計/実装を行なった。

### 弥生株式会社 (2020-07 〜 2020-08)

株式会社Misoca吸収合併による転籍

### 株式会社Misoca (2018-07 〜 2020-06) サーバーサイドエンジニア/プロダクトマネージャ(Android/iOS)

- 役割: サーバーサイドエンジニア, プロダクトマネージャ(Android/iOS)
- 利用技術: Ruby, Ruby on Rails, AWS, Kotlin, Swift, RxSwift

請求書作成サービスMisocaのモバイルアプリのバックエンドAPIの開発, Android/iOS向けのプロダクトマネージメントおよびライブラリ開発に従事。

- Android/iOSに向けたKotlin/Swift製の税率計算ライブラリを実装し、サーバーサイドの実装と同じテストデータを用いたCIを構築した。
- 2019/10の軽減税率制度開始に向け、複数プラットフォームに向けて古いアプリバージョンとの後方互換性を保ちつつ機能変更をするための、モバイルアプリ向けのフィーチャーフラグや段階的なデータマイグレーションを設計した。
- tech blogの執筆
    - [Misocaに必要なことは全て受入プロジェクトで学んだ](https://tech-blog.yayoi-kk.co.jp/entry/2018/08/17/185049)

### WonderPlanet (2013-10 〜 2018-06) サーバーサイドエンジニア/リードエンジニア/エンジニアリングマネージャ

- 役割: サーバーサイドエンジニア, リードエンジニア, エンジニアリングマネージャ
- 利用技術: Python, Django, Falcon, PHP, FuelPHP, C++, Photon Server, AWS, GCP, nginx, MySQL, Redis, Airflow, Jenkins, Cocos2d-x, CircleCI, C#, Unity

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
- Unityエンジニア・デザイナー向けのUnity3dエディタ拡張を実装した。
- Unity製ゲームのCI/CD環境を構築した。

### ニューソン株式会社 (2008-04 〜 2013-09)

- 役割: ソフトウェアエンジニア
- 利用技術: Java, Struts, Struts2, Spring framework, Oracle Database, JavaScript, jQuery, Apache Tomcat, Apache HTTP Server

業務用Webアプリケーション設計/開発を担当した。

- 独自JavaScriptフレームワークを開発し、AjaxによるインタラクティブなWebアプリケーションを開発した。
- Javaで開発した多くの入力パターンを受け付ける化学化合物の計算ライブラリを実装し、データ駆動テストによって高品質を実現した。

## 自己PR

### 複雑な問題解決への意欲

複雑で未知の問題に挑戦し、それを解決する過程で成長することに意欲を感じている。たとえば、メタバースプラットフォーム開発では、仮想空間の永続化表現の設計やリアルタイム同期処理の実装に携わった。このプロジェクトでは、複数ユーザー間の同期を担うサーバーアーキテクチャの刷新に取り組み、VerneMQからの移行を無停止で成功させたことで、システムの拡張性と安定性を向上させた。

また、過去の失敗から学びを得て、それを次の挑戦に活かす姿勢も重要視している。モバイルゲームのリリース後に高負荷対応が間に合わず、長期メンテナンスを余儀なくされた経験が、以後のプロジェクトにおいて安定したサービス運用を実現する礎となった。

これらの経験を通じて、効率的かつ柔軟に複雑な課題を解決する力を磨き、プロダクト全体の成果に貢献している。

### チーム貢献と技術的リーダーシップ

技術的リーダーシップを発揮し、チームの成長と持続可能な組織作りに寄与してきた。具体的には、サーバーチームのリーダーとして、メンバーへの教育や技術的指導を行い、リーダーがSPOF（Single Point of Failure）とならない体制を整備。この結果、リーダー不在時でもチームが機能し、自律的にプロジェクトを進められる環境を構築した。

さらに、リモートワーク環境においても、非同期コミュニケーション文化の醸成に貢献。たとえば、設計や技術的知識を効率的に共有するためにテクニカルライティングの重要性を提唱し、design docのテンプレート改善やガイドライン策定を主導。この取り組みにより、メンバー間のスキル向上と生産性の向上を実現した。

これらの取り組みを通じて、チーム全体のスキルアップと効率的なコラボレーションを推進し、組織の長期的な成長に貢献している。

### 継続的な学習と成長意欲

継続的な学びを通じて、未知の技術や視点を実務に取り入れることを重視している。メタバースプラットフォームにおいて、リアルタイム同期処理や仮想空間の永続化といった未経験の技術分野に挑戦し、これらの技術を活用して新たなプロダクトの方向性を生み出した。

業務外でも技術書の執筆やOSS活動に積極的に参加し、得た知識やアイデアを日常業務に還元。未知の技術領域への挑戦を楽しみながら、「無駄」と思われる試行錯誤からも新たな価値を見出すことを心がけている。

これらの活動を通じて、自身の成長にとどまらず、プロジェクトやチームに新たな視点を提供している。


## 個人プロジェクト

- [SwiftNES](https://github.com/thara/SwiftNES): Swift製のクロスプラットフォーム向けNESエミュレータ
- [gorones](https://github.com/thara/gorones): Go製のクロスプラットフォーム向けNESエミュレータ
- [rust_aliasmethod](https://github.com/thara/rust_aliasmethod): Walker's Alias MethodのRust実装
- [erajp](https://github.com/thara/erajp): Rust製 日本の元号変換ライブラリ
- [SoundIO](https://github.com/thara/SoundIO): クロスプラットフォーム向け音声入出力ライブラリlibsoundioのSwiftバインディング

## その他の活動

- [JUnit実践入門 体系的に学ぶユニットテストの技法](http://gihyo.jp/book/2012/978-4-7741-5377-3)の出版前レビュー
- [SpriteKitではじめる2Dゲームプログラミング Swift対応](http://www.shoeisha.co.jp/book/detail/9784798139517) の執筆
  - 第2章 Swiftの基本, 第8章リバーシを作ろう--AIの作り方 を担当
- [iOSDC Japan 2019: Swiftでつくるファミコンエミュレータのススメ](https://fortee.jp/iosdc-japan-2019/proposal/92904657-beda-46fe-8ecb-b27c75ee0f16) (LT発表)
- [Go Conference 2022 Spring | メタバースを支える技術 ～UGCに溢れる3D空間のリアルタイム同期を支えるGo〜](https://gocon.jp/2022spring/sessions/a6-c/) (スポンサーセッション登壇)

## 資格

- 2019-10 End-to-End Machine Learning with TensorFlow on GCP(Coursera) 修了
    - https://www.coursera.org/account/accomplishments/verify/EBWAMSPYDF9B
- 2019-10 Google Cloud Platform Big Data and Machine Learning Fundamentals 日本語版(Coursera) 修了
    - https://www.coursera.org/account/accomplishments/verify/28Q2L33PSUEB
- 2019-09 How Google does Machine Learning 日本語版(Coursera) 修了
    - https://www.coursera.org/account/accomplishments/verify/3HSV6VGSYNEX
- 2019-09 Launching into Machine Learning 日本語版(Coursera) 修了
    - https://www.coursera.org/account/accomplishments/verify/UH9SSFT7WCTA
- 2012-04 エンベデッドシステムスペシャリスト試験 合格
- 2011-09 情報セキュリティスペシャリスト試験 合格
- 2011-04 応用情報技術者試験 合格
