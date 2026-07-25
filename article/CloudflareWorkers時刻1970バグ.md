---
title: Cloudflare Workers で new Date() が 1970 年になる — バグではなく仕様
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

ローカルの `wrangler dev` や Node 環境では現在時刻が出るのに、デプロイ後だけ 1970 年になる、という形でハマります。

## 原因

これはバグではなく、Workers ランタイムの**意図した仕様**です。ポイントは2つあります。

### 1. グローバルスコープでは時刻が「0」になる

Workers ではモジュールのトップレベル（＝リクエストを受ける前のグローバルスコープ）で `Date.now()` が **0** を返します。理由は決定性（deterministic）の確保です。

グローバルスコープは、デプロイ時・スナップショット復元時・アイソレート起動時など、Cloudflare 側の都合で「いつ実行されるか保証されない」コードです。そこで実時間を返すと、実行タイミングによって結果が変わってしまいます。これを避けるため、グローバルスコープでの時刻は常に 0（= 1970-01-01）に固定されています。

### 2. 時計は I/O のたびにしか進まない（Spectre 対策）

もう1つの背景が、Spectre 系タイミング攻撃への対策です。Workers では、時計は同期実行の間は止まっていて、**I/O（ネットワーク受信などの非同期処理）が発生したタイミングでだけ進みます**。

つまり「現在時刻」は、直近の I/O の時点で更新された値です。まだ一度も I/O が起きていないグローバルスコープでは、更新される前の初期値（0）のままになる、というわけです。

## 正しい使い方

`new Date()` / `Date.now()` は、**リクエストハンドラの中**で呼べば現在時刻が取れます。

```js
export default {
  async fetch(request, env, ctx) {
    console.log(new Date()); // → 現在時刻でOK
    return new Response(new Date().toISOString());
  },
};
```

グローバルで「起動時刻」のようなものを持ちたい場合は、遅延初期化にして最初のリクエスト時に確定させます。

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

## タイムゾーンの注意（あわせてハマる点）

時刻が取れるようになっても、Workers は基本 **UTC** で動きます。JST 前提のコードはズレます。

- `new Date()` の内部値は UTC。`toLocaleString('ja-JP', { timeZone: 'Asia/Tokyo' })` のように `timeZone` を明示する。
- Cron Triggers（`wrangler.toml` の `crons`）も **UTC 基準**。日本時間で組むと9時間ずれる。JST 9:00 に動かしたいなら UTC 0:00（`0 0 * * *`）で書く。

## まとめ

- Workers のグローバルスコープでは `Date.now()` は 0（1970-01-01）を返す。これは決定性と Spectre 対策のための仕様。
- 時計は I/O ごとにしか進まないので、時刻はリクエストハンドラ内で取得する。
- タイムゾーンは UTC 前提。表示・Cron ともに明示的に扱う。

## 参考

- [Fetch handler - Runtime APIs · Cloudflare Workers docs](https://developers.cloudflare.com/workers/runtime-apis/handlers/fetch/)
- [JavaScript | Web standards · Cloudflare Workers docs](https://developers.cloudflare.com/workers/runtime-apis/web-standards/)
- [Cloudflare Workers を支える技術 / Kenton Varda の Stack Overflow 回答（`Date.now()` は最後の I/O 時刻を返す）](https://stackoverflow.com/questions/50611329/cloudflare-workers-date-now)
