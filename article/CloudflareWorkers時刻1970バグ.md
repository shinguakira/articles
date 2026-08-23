---
title: Cloudflare Workers で new Date() が 1970 年になる
tags:
  - Cloudflare
  - CloudflareWorkers
  - JavaScript
  - Date
  - タイムゾーン
private: false
---

## 現象

Cloudflare Workers（や Pages Functions）で、`new Date()` や `Date.now()` が現在時刻ではなく **1970-01-01（Unix エポック）** を返すことがあります。

```js
// モジュールのトップレベル（グローバルスコープ）
const bootTime = new Date();
console.log(bootTime); // → 1970-01-01T00:00:00.000Z
```

ローカルの `wrangler dev` や Node では現在時刻が出るのに、デプロイ後だけ 1970 年になるので、原因が見えにくいバグです。

## 対策

`new Date()` や `Date.now()` は、リクエストハンドラの中で呼べば現在時刻が取れます。

```js
export default {
  async fetch(request, env, ctx) {
    console.log(new Date()); // 現在時刻でOK
    return new Response(new Date().toISOString());
  },
};
```

グローバルで起動時刻のような値を持ちたいときは、最初のリクエストで確定させます。

```js
let bootTime = null;

export default {
  async fetch(request) {
    if (bootTime === null) {
      bootTime = new Date(); // 最初のリクエストで確定
    }
    return new Response(`boot: ${bootTime.toISOString()}`);
  },
};
```

## 原因

バグではなく、Workers ランタイムの時刻の扱いによるものです。

### グローバルスコープでは時刻が 0 になる

モジュールのトップレベル、つまりリクエストを受ける前のグローバルスコープでは、`Date.now()` が **0** を返します。

グローバルスコープのコードは、デプロイ時やスナップショット復元時、アイソレートの起動時など、Cloudflare 側の都合で実行されるタイミングが一定しません。ここで実時間を返すと、いつ呼ばれたかで結果が変わってしまいます。そのため時刻を常に 0（1970-01-01）に固定しています。

### 時計は I/O のたびにしか進まない（Spectre 対策）

もう一つの理由が、Spectre 系のタイミング攻撃への対策です。Workers では時計が同期処理の間は止まっていて、ネットワーク受信などの I/O が発生したときにだけ進みます。

「現在時刻」は直近の I/O の時点で更新された値なので、まだ一度も I/O が起きていないグローバルスコープでは、初期値の 0 のままになります。

## タイムゾーンにも注意

時刻が取れるようになっても、Workers は UTC で動きます。JST 前提のコードはズレます。

- `new Date()` が持つ値は UTC。表示は `toLocaleString('ja-JP', { timeZone: 'Asia/Tokyo' })` のように `timeZone` を明示する。
- Cron Triggers（`wrangler.toml` の `crons`）も UTC 基準。日本時間のつもりで組むと9時間ずれる。JST 9:00 に動かしたいなら UTC 0:00、つまり `0 0 * * *` で書く。

## 実際に踏んだケース

zod のスキーマをフロント・バックの共有箇所に置いていて、未来日を弾くバリデーションで踏みました。

```ts
// 共有スキーマ
const schema = z.object({
  date: z.date().max(new Date(), { message: "未来日は指定できません" }),
});
```

`new Date()` がスキーマ定義（モジュールのトップレベル）で評価されると、Workers 上では 1970-01-01 になります。その結果、上限が 1970 年に固定されてしまい、現在の日付まで「未来日」として弾かれました。共有スキーマなのでフロントでは正常に動き、Workers 側だけで再現するのも気づきにくかった点です。

`new Date()` はスキーマ定義時ではなく、バリデーションが走るタイミング（リクエストハンドラの中）で評価されるようにすれば直ります。

```ts
const schema = z.object({
  date: z.date().max(new Date(), { message: "未来日は指定できません" }),
});
// ↑ 呼び出しをリクエスト内に置く。あるいは refine で実行時に new Date() を取る
const schema = z.object({
  date: z.date().refine((d) => d <= new Date(), { message: "未来日は指定できません" }),
});
```

## まとめ

- グローバルスコープでは `Date.now()` は 0（1970-01-01）を返す。決定性と Spectre 対策のため。
- 時刻はリクエストハンドラの中で取る。
- Workers は UTC 基準。表示も Cron も `timeZone` を明示する。

## 参考

一番の情報源は、Cloudflare のエンジニア Kenton Varda による Stack Overflow の回答です。`Date.now()` が最後の I/O の時刻を返す挙動をここで説明しています。

- [Cloudflare Workers: Date.now() の挙動（Kenton Varda の回答）](https://stackoverflow.com/questions/50611329/cloudflare-workers-date-now)
- [Fetch handler - Runtime APIs · Cloudflare Workers docs](https://developers.cloudflare.com/workers/runtime-apis/handlers/fetch/)
- [JavaScript | Web standards · Cloudflare Workers docs](https://developers.cloudflare.com/workers/runtime-apis/web-standards/)
