---
title: Incident Management System
date: '2026-08-23'
published: '2026-08-23'
---

SRE本で紹介されたIncident Command System(ICS)は、今ではもはや当たり前になった（よね？）。

が、肝心の元になった[FEMA National Incident Management System](https://www.fema.gov/emergency-managers/nims) (NIMS)について1ミリも知らなかったので、ちょっと調べてみた。

今まで、システム障害としてインシデント管理にしか触れてこなかったが、改めて[The National Incident Management System guides (PDF)](https://www.fema.gov/sites/default/files/2020-07/fema_nims_doctrine-2017.pdf) を読んで印象に残ったこととかシステム障害においてどう考えるかを雑に挙げてみる。

## FEMAだけの取り組みじゃない

FEMA = アメリカ合衆国連邦緊急事態管理庁、なんだけど、てっきり消防だけの取り組みだと思っていた。
NIMSはそうじゃなくて、警察や医療、政府や自治体、NGOや民間企業などが緊急事態にどう強調するかの標準仕様みたいなものだった。

システム障害って文脈でICSを取り扱ってきたけど、組織内での緊急事態における標準的なフレームワークとして、もっと普及されてもいいのかもしれない。

## インシデント管理の原則

インシデント担当者は、以下の原則にしたがってNIMSを実施するとのこと。

- Flexibility
- Standardization
- Unity of Effort

システム障害対応としては、こういうのは社内に閉じてしまいがちだけど、ソフトウェア業界全体のスキルアップ的なことを考えると、ある程度標準化されていた方が良いかもしれない。

そういえば、こういう企業跨いだ標準化みたいなの、技術領域ではあるけど、ソフトウェア開発・運用組織運営とかの標準みたいなものって、ソフトウェア業界ってあんまないよねぇ...

## Incident Command Systemは一部

NIMSはおおよそ、以下の構造を取っている。   

- Resource Management
    - Resource Management Preparedness
    - Resource Management During an Incident
    - Mutual Aid
- Command and Coordination
    - NIMS Management Characteristics
    - Incident Command System (ICS) <- ここ
    - Emergency Operations Centers (EOC)
    - Multiagency Coordination Group (MAC Group)
    - Joint Information System (JIS)
    - Interconnectivity of NIMS Command and Coordination Structures
- Communications and Information Management
    - Communications Management
    - Incident Information
    - Communications Standards and Formats

見ればわかるように、NIMSの中でIncident Command Systemは一部でしかない。

もちろん、システム障害対応に置いてあまり意味をなさないものもあるかもしれないが、インシデント対応フローなどを社内で規定する際に、これらの何をどう無視 or 追加して、何をどこまでカバーするかは考えてみる価値はあるかもしれない。

というか、そもそもインシデント対応 "フロー" を考える前に、インシデント対応に何が必要なのかとか、どういう原則で取り組むのかとか、そこらへんの定義がまず第一に必要であると感じた。

## Incident Commander配下のCommand Staff

ICSには、Incident Commanderの指揮機能を支援するために、以下のようなCommand Staffも定義されていた。

- Public Information Officer
    - 広報担当官
- Safety Officer
    - 安全担当官 (インシデント対応メンバーの "安全")
- Liaison Officer
    - 連絡調整官

Incident Commanderは、必要に応じてさらに新規の補佐を任命することもある。

システム障害フローにおいても、広報担当官的な役割は設けられることがよくありそう。
レピュテーションリスク対応としてPdMとかが入って、カスタマーサポートとかマーケとかに指示するイメージ。

一方で、安全担当官みたいなのはあんまり見かけないかも。直接的な身の危険はあまりないかもしれないけど、労務管理的な観点でそういう担当がいても良いかもしれない。

連絡調整官は、EMとかが務めるイメージ。（そもそもIncident CommanderがEMのことも少なくないと思うけど）

## Unity of Command と Unity of Effort

インシデント対応メンバーが単一の指揮官に対して報告するのが、Unity of Command。   
複数の指揮系統があると、メンバーの中で優先順位をつけなくてはいけなくなるので、重要。

インシデント対応メンバーが共通の目標を達成するために調整するのが、Unity of Effort。   
こちらも、その時点でのインシデント対応自体の目標が明確に定義され、それが指揮系統に浸透されていることが、重要。

## Joint Information System

外部に対して、適時かつ正確で、入手しやすく、かつ行動につながる情報を提供するために整理されたプロセスや手順、ツールから成るのがJoint Information System(JIS)。

システム障害だと、サービスステータスページとその更新方法とか、対外発信方法などになるだろうか。
Incident Commanderからどのような経路でそういったアウトプットがなされるかがあらかじめ定められているとよさそう。
複数の経路から発信する場合に、その整合性を保つための仕組みは必要そう。

---

パラパラと読んでみたけど、単に「インシデント対応」といっても抽象的なドメインであることがわかった。

元のFEMA National Incident Management Systemは、それこそ災害とかテロ対応とか、そういったレベルのものだけれど、システム障害以外の組織内での緊急事態の対応にも広く展開できそうではある。

それこそ、元のNIMSは組織の枠を超えた枠組みを提供しているわけで、企業を超えた緊急事態対応の標準フレームワークとして、同じようなものが定義されると良いのではないか。

どうやってそういうのを進めるのかわからんけど。業界団体的なもの？とかが進めるんかな...

NIMSでは[トレーニングプログラム](https://www.fema.gov/sites/default/files/documents/fema_nims-training-program-2020.pdf)とかもあるみたい。   
少なくともICSぐらいは、業界団体的なもの？とかがこういった教育プログラムを通じて、ソフトウェア業界全体に普及していくのがいいのかもしれない。
