---
title: さくらのクラウドの「シンプルMQ」で始めるイベント駆動アプリケーション
date: '2025-12-23'
published: '2025-12-23'
---

この記事は [さくらインターネット Advent Calendar 2025](https://qiita.com/advent-calendar/2025/sakura) 23日目の記事です。   
22日目は[kiyokiyo](https://x.com/kiyokiyo)さんの [Containerlabを使ってネットワーク機器のCIを考えてみた](https://qiita.com/kiyokiyo/items/47b779d5601921a4dead) と [SONiCのsystemdを読みながら好きなツールを動かしてみる](https://qiita.com/kiyokiyo/items/3e80799ff9e430282a4b) でした。2つともネットワークエンジニアにとって非常に興味深い内容でしたね。自分には全然わかりませんが、ネットワーク領域にもCI/CD、可観測、IaCの波が来ているのだなと感じました。

------------

[さくらのクラウド シンプルMQ](https://manual.sakura.ad.jp/cloud/appliance/simplemq/index.html)は、マネージドなメッセージキューサービスです。この記事では、シンプルMQを使ったProducer-Consumerパターンの実装例を紹介します。

特に以下の点に焦点を当てます：

- **Protocol Buffers + Deflate圧縮**によるメッセージの効率的な取り扱い
- **Terraform**によるインフラのコード化
- 実践的な実装パターンと運用のポイント

### サンプルアプリケーションの全体像

![sample application overview](../../images/simplemq-sample.png)

- **Producer CLI**: コマンドラインから1メッセージを送信
- **Producer API Server**: HTTP APIでメッセージを受け付け、シンプルMQに転送
- **Consumer Worker**: シンプルMQをポーリングしてメッセージを処理

## 【入門編】基本的な使い方

まずはシンプルMQの基本的な使い方を体験してみましょう。

### リポジトリのクローン

```bash
git clone https://github.com/thara/sakura-simplemq-sample.git
cd sakura-simplemq-sample
```

### TerraformでシンプルMQを構築

Infrastructure as Codeの実践として、TerraformでシンプルMQのキューを作成します。

まず、さくらのクラウドのAPIキーのアクセストークンとアクセストークンシークレットを環境変数に設定します：

```bash
export SAKURACLOUD_ACCESS_TOKEN="your-access-token"
export SAKURACLOUD_ACCESS_TOKEN_SECRET="your-access-token-secret"
```

次にTerraformでキューを作成：

```bash
cd terraform
terraform init
terraform apply
cd ..
```

`terraform/main.tf`では、以下のようにシンプルMQリソースを定義しています：

```hcl
resource "sakura_simple_mq" "playground" {
  name = "playground-mq"
  description = "This is a playground message queue."
  tags = ["playground"]

  visibility_timeout_seconds = 10  # メッセージの可視性タイムアウト（秒）
  expire_seconds             = 100 # メッセージの有効期限（秒）
}
```

作成後、さくらのクラウドのコントロールパネルからシンプルMQのキューのAPIキーをローテートし、環境変数に設定します：

```bash
export SIMPLEMQ_API_KEY="your-simplemq-api-key"
```

2025年12月現在terraform apply時にAPIキーを設定 or 取得する方法は提供されていないため、**キューを作成後に一度APIキーをローテートする必要がある** ことに注意してください。

## 動かしてみる

それでは、3つのコンポーネントを実際に動かしてみましょう。

**Consumer Workerを起動（別ターミナル）:**

```bash
go run cmd/consumer-worker/main.go -queue=playground-mq
```

このWorkerは1秒ごとにシンプルMQをポーリングし、メッセージを処理します。

**Producer CLIでメッセージを送信:**

```bash
go run cmd/producer/main.go -queue=playground-mq -message="Hello from CLI!"
```

Consumer側のターミナルに以下のようなログが表示されるはずです：

```
yyyy/MM/dd HH:mm:ss INFO messages received count=1
yyyy/MM/dd HH:mm:ss INFO message received id=<message-id>
yyyy/MM/dd HH:mm:ss INFO message deleted id=<message-id>
yyyy/MM/dd HH:mm:ss INFO notification received message="Hello from CLI!"
```

**Producer API Serverでの送信（別ターミナル）:**

```bash
# サーバーを起動
go run cmd/producer-api-server/main.go -addr=:8080 -queue=playground-mq

# 別ターミナルからHTTP POSTで送信
curl -X POST http://localhost:8080 -d "Hello from API Server!"
```

同様にConsumer側でメッセージが受信されます。これで非同期メッセージングの基本的な流れを体験できました。

## 【実践編】本番運用に向けた工夫

ここからは、本番運用を想定した実装の工夫について解説します。基本的な使い方を理解した方向けの内容です。

### シンプルMQの特性と設計への影響

[シンプルMQ](https://manual.sakura.ad.jp/cloud/appliance/simplemq/index.html)には以下の特性があり、これらを理解した上で設計する必要があります：

**Pull型 - Consumer側がポーリングする**

Consumer側が能動的にメッセージを取得します。今回のサンプルでは1秒ごとのポーリングを実装しています。

**順序は保証されない**

受信時刻の古いメッセージから配信されますが、[可視性タイムアウト延長などにより順序が入れ替わる](https://manual.sakura.ad.jp/cloud/appliance/simplemq/index.html#:~:text=%E2%80%BB%E5%BE%8C%E8%BF%B0%E3%81%AE-,%E5%8F%AF%E8%A6%96%E6%80%A7%E3%82%BF%E3%82%A4%E3%83%A0%E3%82%A2%E3%82%A6%E3%83%88%E5%BB%B6%E9%95%B7%E3%81%AA%E3%81%A9%E3%81%AB%E3%82%88%E3%82%8A%E3%80%81%E6%96%B0%E3%81%97%E3%81%84%E3%83%A1%E3%83%83%E3%82%BB%E3%83%BC%E3%82%B8%E3%81%8C%E5%85%88%E3%81%AB%E9%85%8D%E4%BF%A1%E3%81%95%E3%82%8C%E3%82%8B%E5%A0%B4%E5%90%88%E3%82%82%E3%81%82%E3%82%8A%E3%81%BE%E3%81%99%E3%80%82,-%E4%BB%95%E6%A7%98%EF%83%81)場合があります。順序が重要な処理では、アプリケーション側で制御してください。

**重複配信を前提に冪等性を確保**

[各メッセージは少なくとも1回の配信が保証](https://manual.sakura.ad.jp/cloud/appliance/simplemq/index.html#:~:text=%E9%85%8D%E4%BF%A1%E3%81%95%E3%82%8C%E3%82%8B%E3%80%82-,%E5%90%84%E3%83%A1%E3%83%83%E3%82%BB%E3%83%BC%E3%82%B8%E3%81%AF%E5%B0%91%E3%81%AA%E3%81%8F%E3%81%A8%E3%82%821%E5%9B%9E%E3%81%AE%E9%85%8D%E4%BF%A1%E3%81%8C%E4%BF%9D%E8%A8%BC%E3%81%95%E3%82%8C%E3%82%8B,-%E3%80%82)されますが、2回以上配信される可能性があります。Consumer側で冪等性を確保する必要があります。

冪等性の実装例：
- メッセージIDをDBに記録し、処理済みなら何もしない
- 処理済みフラグを外部ストア（Redis等）に置く
- 副作用を冪等なAPI（PUT、べき等な更新）に寄せる

**可視性タイムアウトは再配信の原因になる**

メッセージを受信してから他のConsumerに見えなくなる時間（可視性タイムアウト）があります。処理時間より短いと、処理中のメッセージが再配信されます。Terraformでの `visibility_timeout_seconds` 設定時は、想定される処理時間より長めに設定してください。

### メッセージの効率的な取り扱い

このサンプルアプリケーションでは、メッセージの効率性と信頼性を高めるために、**Protocol Buffers**と**Deflate圧縮**を組み合わせています。

#### なぜProtocol Buffersなのか

[Protocol Buffers（protobuf）](https://protobuf.dev/)は、Googleが開発したシリアライズフォーマットです。JSONと比較して以下の利点があります：

1. **スキーマ定義**: `.proto`ファイルでメッセージ構造を明確に定義
2. **型安全性**: コンパイル時に型チェックが行われる
3. **バイナリフォーマット**: JSONよりコンパクトで高速

今回のサンプルでは、シンプルなNotificationメッセージを定義しています：

```protobuf
edition = "2023";

option go_package = "github.com/thara/sakura-simplemq-sample/samplepb";

message Notification {
  string message = 1;
}
```

このシンプルな例でも、将来的にフィールドを追加する際に後方互換性を保ちやすいという利点があります。

**なぜここまでやるのか**

この設計は「キューを長期利用する」前提です。メッセージは将来、他言語・他サービスからも利用される可能性があるため、以下の点を重視しています：

- **スキーマの安定性**: `.proto`ファイルで明確に定義され、サービス間で共有できる
- **後方互換性**: フィールド追加時も既存のConsumerが動き続ける
- **サイズ制限耐性**: シンプルMQの256000文字制限に対して余裕を持たせる

**適用の判断基準**

- メッセージサイズが数十〜数百byte中心なら、圧縮の効果は限定的
- 他言語対応・長期運用を想定するなら Protocol Buffers が有効
- デバッグ容易性を優先するなら JSON も選択肢に入る

#### Deflate圧縮 + Base64エンコーディング

さらに、メッセージサイズを削減するためにDeflate圧縮を適用し、Base64エンコーディングでテキスト化してからシンプルMQに送信しています。Go標準ライブラリの`compress/flate`だけで完結させるためDeflateを採用していますが、gzipでも同様に実装できます。

エンコード処理の流れ（`internal/encoding.go`）：

```
Notification (protobuf)
  → Marshal (バイナリ化)
  → Deflate圧縮
  → Base64エンコード
  → シンプルMQに送信するメッセージ（文字列）
```

実装を見てみましょう（`internal/encoding.go`）：

```go
func encodeProtoMessage(msg proto.Message) (string, error) {
    // 1. Protocol Buffersでバイナリ化
    b, err := proto.Marshal(msg)
    if err != nil {
        return "", fmt.Errorf("failed to marshal proto message: %w", err)
    }

    // 2. Deflate圧縮
    compressed, err := compress(b)
    if err != nil {
        return "", fmt.Errorf("failed to compress proto message: %w", err)
    }

    // 3. Base64エンコード
    return base64.StdEncoding.EncodeToString(compressed), nil
}
```

デコード側は逆の順序で処理を行います：

```go
func decodeProtoMessage(src string, msg proto.Message) error {
    // 1. Base64デコード
    data, err := base64.StdEncoding.DecodeString(src)
    if err != nil {
        return fmt.Errorf("failed to decode base64 string: %w", err)
    }

    // 2. Deflate解凍
    decompressed, err := decompress(data)
    if err != nil {
        return fmt.Errorf("failed to decompress proto message: %w", err)
    }

    // 3. Protocol Buffersでデシリアライズ
    if err := proto.Unmarshal(decompressed, msg); err != nil {
        return fmt.Errorf("failed to unmarshal proto message: %w", err)
    }
    return nil
}
```

#### 圧縮の効果

この方式により、以下のようなメリットが得られます：

- **メッセージサイズ制限への対応**: シンプルMQには1メッセージ[256000文字以内](https://manual.sakura.ad.jp/api/cloud/simplemq/#operation/sendMessage)という制限があります。圧縮により、より多くのデータを送信できます
- **転送サイズの削減**: 特に繰り返しの多いデータで圧縮効果が高い
- **ネットワーク効率**: 帯域幅の節約

ただし、メッセージが数十バイト程度の場合、圧縮によって逆にサイズが増えることもあります。ユースケースに応じて判断してください。

### Producer-Consumerパターンの実装

#### Producer実装のポイント

ProducerとConsumerは共通の`internal`パッケージを使用することで、エンコード/デコード処理の一貫性を保っています。

**CLI Producer**（`cmd/producer/main.go`）:

```go
notification := &samplepb.Notification{
    Message: proto.String(message),
}
if err := internal.SendNotification(ctx, messageOp, notification); err != nil {
    return fmt.Errorf("failed to send notification: %v", err)
}
```

**API Server Producer**（`cmd/producer-api-server/main.go`）:

```go
handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "failed to read request body", http.StatusBadRequest)
        return
    }

    notification := &samplepb.Notification{
        Message: proto.String(string(body)),
    }
    if err := internal.SendNotification(r.Context(), messageOp, notification); err != nil {
        http.Error(w, fmt.Sprintf("failed to send notification: %w", err), http.StatusInternalServerError)
        return
    }
})
```

どちらも`internal.SendNotification`関数を呼び出すだけで、エンコード処理は内部で行われます：

```go
func SendNotification(ctx context.Context, messageOp simplemq.MessageAPI, notification *samplepb.Notification) error {
    content, err := encodeProtoMessage(notification)
    if err != nil {
        return fmt.Errorf("failed to encode message: %v", err)
    }

    resSend, err := messageOp.Send(ctx, content)
    if err != nil {
        return fmt.Errorf("failed to send message: %v", err)
    }

    slog.Info("Message Sent", slog.String("ID", string(resSend.ID)))
    return nil
}
```

#### Consumer実装のポイント

Consumer Worker（`cmd/consumer-worker/main.go`）は、1秒ごとにシンプルMQをポーリングします：

```go
ticker := time.NewTicker(1 * time.Second)
defer ticker.Stop()

for {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-ticker.C:
        if err := receiveMessages(ctx, client, queueName); err != nil {
            slog.ErrorContext(ctx, "failed to receive messages", slog.Any("error", err))
        }
    }
}
```

メッセージ受信処理では、`internal.ReceiveNotifications`を使用してデコードと削除（acknowledge）を行います：

```go
func ReceiveNotifications(ctx context.Context, messageOp simplemq.MessageAPI) ([]*samplepb.Notification, error) {
    // 1. シンプルMQからメッセージを受信
    messages, err := messageOp.Receive(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to receive messages: %v", err)
    }
    slog.Info("messages received", slog.Int("count", len(messages)))

    var notifications []*samplepb.Notification

    // 2. 各メッセージを処理
    for _, msg := range messages {
        if err := func() error {
            // 2-1. デコード
            var notification samplepb.Notification
            if err := decodeProtoMessage(string(msg.Content), &notification); err != nil {
                return fmt.Errorf("failed to decode message ID %s: %v", msg.ID, err)
            }
            slog.Info("message received", slog.String("id", string(msg.ID)), slog.String("content", notification.GetMessage()))

            notifications = append(notifications, &notification)

            // 2-2. メッセージを削除してacknowledge
            if err := messageOp.Delete(ctx, string(msg.ID)); err != nil {
                return fmt.Errorf("failed to delete message ID %s: %v", msg.ID, err)
            }
            slog.Info("message deleted", slog.String("id", string(msg.ID)))

            return nil
        }(); err != nil {
            slog.Error("failed to process message", slog.String("id", string(msg.ID)), slog.Any("error", err))
        }
    }
    return notifications, nil
}
```

### 実践的なユースケースと運用Tips

#### よくあるユースケース

- **バックグラウンド処理**: 画像処理、メール送信、レポート生成など
- **マイクロサービス連携**: サービス間の非同期イベント通知
- **データパイプライン**: ログ収集・加工・格納の各ステージを疎結合に接続

#### 運用上の考慮事項

**メッセージ設定の調整**

Terraformで定義した以下のパラメータは、用途に応じて調整が必要です：

- `visibility_timeout_seconds`: メッセージを受信してから他のConsumerに見えなくなる時間。処理時間より長く設定する
- `expire_seconds`: メッセージの有効期限。処理されなかったメッセージが自動削除されるまでの時間

**複数Consumerによるスケーリング**

処理量が増えた場合、Consumer Workerを複数起動することで並列処理が可能です：

```bash
# ターミナル1
go run cmd/consumer-worker/main.go -queue=playground-mq

# ターミナル2
go run cmd/consumer-worker/main.go -queue=playground-mq

# ターミナル3（さらに追加）
go run cmd/consumer-worker/main.go -queue=playground-mq
```

それぞれのWorkerが異なるメッセージを処理するため、処理能力が向上します。

**エラーハンドリングとリトライ**

このサンプル実装には、本番運用で危険な点があります。デコード失敗などの恒久的に成功しない処理が発生した場合、そのメッセージは削除されないため、可視性タイムアウトが切れると再配信され続けます。
このようなメッセージを[Poison message](https://en.wikipedia.org/wiki/Poison_message)と呼びます。

本番環境ではこのPoison messageに対して以下のような対策を講じることを推奨します：

- リトライ回数の保持（メッセージ属性にカウンタを持たせる）
- 一定回数失敗したメッセージの削除、またはエラーログへの記録
- アラート通知（異常なエラー率の検知）

**モニタリング**

以下の指標を監視することで、システムの健全性を把握できます：

- キュー内のメッセージ数（滞留していないか）
- メッセージの処理時間
- エラー率
- Consumer数とその稼働状況

## まとめ

この記事では、さくらのクラウドの「シンプルMQ」を使ったProducer-Consumerパターンの実装例を紹介しました。Protocol Buffers + Deflate圧縮により、型安全性とメッセージサイズの削減を両立しながら、将来的な拡張性も確保できます。

### 次のステップ

このサンプルをベースに、以下のような拡張も考えられます：

- **複雑なメッセージ**: Protocol Buffersのスキーマを拡張し、より多くの情報を含める
- **エラーハンドリングの強化**: アプリケーション側でのリトライやPoison message対策の実装
- **モニタリング**: [モニタリングスイート](https://manual.sakura.ad.jp/cloud/appliance/monitoring-suite/index.html)を使ってキューの状態やメッセージ処理状況を可視化

サンプルコードは[GitHub](https://github.com/thara/sakura-simplemq-sample)で公開していますので、ぜひ試してみてください。

さくらのクラウド シンプルMQを活用して、スケーラブルで堅牢なイベント駆動アプリケーションを構築しましょう！

---

明日は[@linyows](https://github.com/linyows)さんの「Stalwartについて書く」です。またしても何もわからないのですが、このようにさくらインターネットでは様々な技術分野で活躍されている方が多いので、とても勉強になりますね。
ぜひ他の記事もチェックしてみてください！
