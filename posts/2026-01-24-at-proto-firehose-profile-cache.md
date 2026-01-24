---
title: "AT Protocol: Firehoseイベントからプロフィールをキャッシュする"
date: '2026-01-24'
published: '2026-01-24'
---

[前回の記事](./at-proto-statusphere-display-name.html)で、ステータス投稿時に投稿者のプロフィールをキャッシュする実装を行った。しかし、これだと他ユーザーが投稿したステータス（Firehose経由で受信）の投稿者プロフィールがキャッシュされず、display nameが表示されないという課題があった。

今回は、前回の記事の最後で挙げた「ステータス投稿イベント監視時に、その投稿者のプロフィールを取得してキャッシュする」アプローチを実装した。

実装内容は[これ](https://github.com/thara/statusphere-example-app/pull/7)。

## 分散プロトコルにおけるプロフィール取得

AT Protocolは分散プロトコルであり、各ユーザーのデータは異なるPDS (Personal Data Server) に保存されている。そのため、あるユーザーのプロフィールを取得するには、まず「そのユーザーのデータがどこにあるか」を知る必要がある。

### DID と DID Document

AT Protocolでは、ユーザーは [DID (Decentralized Identifier)](https://atproto.com/specs/did) で識別される。例えば `did:plc:z72i7hdynmk6r22z27h6tvur` のような形式だ。

DIDを解決すると、DID Documentが得られる。これにはそのユーザーのPDSのURLが含まれている:

```json
{
  "id": "did:plc:z72i7hdynmk6r22z27h6tvur",
  "service": [
    {
      "id": "#atproto_pds",
      "type": "AtprotoPersonalDataServer",
      "serviceEndpoint": "https://bsky.social"
    }
  ],
  ...
}
```

### プロフィール取得の流れ

Firehoseからステータスイベントを受信した時、イベントには `authorDid` しか含まれていない。プロフィールを取得するには以下の手順が必要になる:

```
1. DID解決: did:plc:xxx → DID Document
2. PDS特定: DID Document → serviceEndpoint (例: https://bsky.social)
3. プロフィール取得: PDS に getRecord を呼び出し
```

コードで書くとこんな感じ:

```typescript
// 1. DID Document を解決
const didDoc = await idResolver.did.resolve(did)

// 2. PDS の URL を取得
const pdsService = didDoc.service?.find(
  (s) => s.id === '#atproto_pds'
)
const pdsUrl = pdsService.serviceEndpoint  // "https://bsky.social"

// 3. 認証なしの Agent でプロフィールを取得
const agent = new Agent({ service: pdsUrl })
const profile = await agent.com.atproto.repo.getRecord({
  repo: did,
  collection: 'app.bsky.actor.profile',
  rkey: 'self',
})
```

従来のWebアプリケーション開発では「ユーザーデータは自分のDBにある」という前提だったが、分散プロトコルでは「ユーザーデータは世界中に散らばっている」という前提になる。これが根本的な発想の違いだと感じる。

## 認証なしでプロフィールが取得できる理由

上記のコードでは認証情報を渡していない。これで動くのは、AT Protocolではプロフィールが**公開データ**として設計されているため。

AT Protocolのリポジトリ内のデータは基本的に公開されている:

- `app.bsky.actor.profile` (プロフィール) → 公開
- `app.bsky.feed.post` (投稿) → 公開
- `xyz.statusphere.status` (ステータス) → 公開

一方、ユーザーの設定など一部のデータは非公開で、認証が必要:

- `app.bsky.actor.preferences` → 非公開

`com.atproto.repo.getRecord` は公開レコードに対しては認証不要で呼び出せる。これにより、任意のユーザーのプロフィールを取得できる。

## Firehose処理をブロックしない

今回の実装で気をつけたのは、プロフィール取得がFirehoseのイベント処理をブロックしないようにすること。

Firehoseは大量のイベントをリアルタイムで受信するため、各イベントの処理は高速に完了させる必要がある。プロフィール取得はネットワークI/Oを伴うため、同期的に実行するとFirehose全体の処理が詰まってしまう。

そのため、プロフィール取得は非同期で実行し、完了を待たずに次のイベント処理に進むようにした:

```typescript
// ステータス保存後
if (!profileCached) {
  // await しない = Firehose処理をブロックしない
  fetchAndCacheProfile(did, db, idResolver, logger).catch(() => {})
}
```

プロフィール取得に失敗しても、ステータス自体は正常に保存される。display nameは次回のページ読み込み時に表示されればよいので、この程度の遅延は許容範囲とする。

## まとめ

今回の実装を通じて、分散プロトコル上でのデータ取得の考え方を学んだ:

- **DID解決が必要**: ユーザーのデータがどのサーバーにあるか、DID Documentから特定する
- **公開データモデル**: プロフィールなどの公開データは認証なしで取得可能
- **非同期処理の重要性**: Firehoseのようなリアルタイムストリームでは、ブロッキング処理を避ける

前回の記事で挙げたもう一つのアプローチ「スキーマにdisplay nameを含める」は、AT Protocolの分散性を考えると筋が悪いと判断した。プロフィールは頻繁に更新される可能性があり、ステータスレコードに埋め込むと古いデータが残り続けてしまう。必要な時に最新のプロフィールを取得する今回のアプローチの方が、分散プロトコルの特性に合っていると思う。

## 参考リンク

- [AT Protocol: Identity](https://atproto.com/specs/did)
- [AT Protocol: Repository](https://atproto.com/specs/repository)
- [実装PR](https://github.com/thara/statusphere-example-app/pull/7)
