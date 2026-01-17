---
title: "AT Protocolアプリにdisplay name表示機能を追加する"
date: '2026-01-11'
published: '2026-01-11'
---


[前回の記事](./at_proto_quickstart.html)でAT Protocolのチュートリアルを一通り完了したので、今度は[Quick start guide](https://atproto.com/guides/applications)のNext Stepsセクションにある「Sync the profile records of all users so that you can show their display names instead of their handles」を実装してみた。

実装内容は[これ](https://github.com/thara/statusphere-example-app/pull/6)。

## 解決策: Write-time caching

今回は、ステータス投稿（書き込み）時に、その投稿者のプロフィール情報を取得してローカルに保存するアプローチを採用した。

表示の即時性を優先しつつ、実装の複雑さを最小限に抑えるのが狙い。

### 実装の流れ

src/routes.ts - POST /status ハンドラー内

```typescript
// ステータス投稿後、プロフィールを取得してキャッシュ
const profileResponse = await agent.com.atproto.repo.getRecord({
  repo: agent.assertDid,
  collection: 'app.bsky.actor.profile',
  rkey: 'self',
})

// 取得したプロフィール情報をSQLiteにキャッシュ
await ctx.db
  .insertInto('profile')
  .values({
    did: agent.assertDid,
    displayName: record.displayName || null,
    // ...
  })
  .onConflict((oc) =>
    oc.column('did').doUpdateSet({
      displayName: record.displayName || null,
      // ...
    })
  )
  .execute()
```

これにより、投稿直後にすぐに表示名が表示されるようになる。エラーが起きてもステータス投稿自体は失敗させないようにエラーハンドリングも入っている。

### 実装の詳細

#### データベーススキーマ


新しくprofileテーブルを追加。マイグレーションは起動時に自動実行される:

src/db.ts

```typescript
export type Profile = {
  did: string              // Primary key
  displayName: string | null
  description: string | null
  avatarCid: string | null
  avatarMimeType: string | null
  bannerCid: string | null
  bannerMimeType: string | null
  indexedAt: string
}

export type DatabaseSchema = {
  status: Status
  profile: Profile  // 追加
  auth_session: AuthSession
  auth_state: AuthState
}
```


[Kysely](https://kysely.dev/) を使ったのは初めてだったけど、まぁ他のORMと似た感じなので特に問題なく使えた。


#### 表示ロジック

homeページではstatusテーブルとprofileテーブルをLEFT JOINして、表示名を取得:

src/routes.ts
```typescript
const statusesWithProfiles = await ctx.db
  .selectFrom('status')
  .leftJoin('profile', 'status.authorDid', 'profile.did')
  .select([
    'status.uri',
    'status.authorDid',
    'status.status',
    'profile.displayName',
    // ...
  ])
  .execute()

const statuses = statusesWithProfiles.map(row => ({
  uri: row.uri,
  authorDid: row.authorDid,
  status: row.status,
  displayName: row.displayName || null,
  // ...
}))
```

フロントエンドでは、表示名がある場合は「DisplayName (@handle)」、ない場合は「@handle」という形式で表示:

src/pages/home.ts
```typescript
const displayName = status?.displayName

const authorDisplay = displayName
  ? html`${displayName} <span class="handle">(@${handle})</span>`
  : html`@${handle}`
```

## Lexiconとスキーマ定義

AT Protocolでは、データ構造を [Lexicon](https://atproto.com/specs/lexicon) というスキーマ定義言語で定義する。今回使った [app.bsky.actor.profile](https://github.com/bluesky-social/atproto/blob/main/lexicons/app/bsky/actor/profile.json) もLexiconで定義されている。

チュートリアルアプリのxyz.statusphere.statusは独自のLexiconスキーマで、lexicons/ディレクトリに定義されている。npm run lexgenコマンドでTypeScriptの型定義が自動生成される:

lexicons/status.json

```json
{
  "lexicon": 1,
  "id": "xyz.statusphere.status",
  "defs": {
    "main": {
      "type": "record",
      "key": "tid",
      "record": {
        "type": "object",
        "required": ["status", "createdAt"],
        "properties": {
          "status": {
            "type": "string",
            "minLength": 1,
            "maxGraphemes": 1,
            "maxLength": 32
          },
          "createdAt": { "type": "string", "format": "datetime" }
        }
      }
    }
  }
}
```

src/routes.ts
```typescript
import * as Status from '#/lexicon/types/xyz/statusphere/status'

const record = {
  $type: 'xyz.statusphere.status',
  status: req.body?.status,
  createdAt: new Date().toISOString(),
}

// バリデーション
if (!Status.validateRecord(record).success) {
  ...
}
```

この仕組みにより、型安全性を保ちながら柔軟なデータ構造を定義できる。   
ここら辺は自分で開発する際もよくProtocolBuffersスキーマとか使うので、似たような感覚で使えそう。

## 分散型プロトコルならではの課題

今回の実装を通して、「分散型プロトコル上でアプリケーションを作るときの前提」が、これまで慣れてきたWebアプリケーションとは大きく異なると感じた。

特に印象的だったのは、**自分が管理していないユーザーのデータが、自然に流れ込んでくる** こと。

今回の実装では、一旦UI上の課題に対応する目的で複雑なキャッシュ戦略ではなく「投稿時にプロフィールをフェッチする」というシンプルなアプローチを採用した。しかし、これだと以下の問題がある。

- 投稿していないユーザーのdisplay nameの更新が反映されない
- DBに保存されていないユーザーのdisplay nameが表示されない

後者については説明が必要かもしれない。

チュートリアルアプリでは `xyz.statusphere.status` コレクションでステータス投稿しているが、このコレクションはAT Protocol上の全ユーザーが投稿できる。つまり、他の人が投稿したステータスが自分のローカルで起動しているチュートリアルアプリに届く可能性がある。

今回の実装では、**ステータス投稿時にそのユーザーのプロフィールをキャッシュする**仕組みを入れたため、他の人の手元で投稿されたステータスの投稿者プロフィールは自分のローカルDBにキャッシュされない。

今までのアプリケーション開発の発想だと当たり前なのだが、分散型プロトコル上で動作するアプリケーションでは、**自分が管理しているデータ以外にも他のユーザーのデータが届く可能性がある** という点が新しい発想だった。

なので、次回はこの課題を解決に挑戦したい。

実は、初期の実装ではFirehoseでプロフィール更新イベント(`app.bsky.actor.profile`) を監視してキャッシュする実装にしていた。しかし、これだとAT Protocol上の全てのユーザーのプロフィール更新が送られてくるため、負荷が高くなる可能性があったのでやめた。

これを踏まえて、次は以下の2つのアプローチを検証してみる予定。

- `xyz.statusphere.status` に display name を含めるようにスキーマを拡張する
- ステータス投稿イベント監視時に、その投稿者のプロフィールを取得してキャッシュする

## 参考リンク

- https://atproto.com/guides/applications
- https://atproto.com/specs/repository
- https://atproto.com/specs/lexicon
- https://github.com/thara/statusphere-example-app/pull/6
