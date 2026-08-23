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

Cloudflare Workers や Pages Functions で、`new Date()` や `Date.now()` が現在時刻ではなく **1970-01-01（Unix エポック）** を返すことがあります。

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

Spectre 系のタイミング攻撃を防ぐため、Workers の時計は同期処理の間は止まっていて、I/O（ネットワーク受信など）が発生したときにだけ進みます。`Date.now()` は「最後の I/O の時刻」を返すので、まだ一度も I/O が起きていないグローバルスコープ（モジュールのトップレベル）では初期値の 0、`new Date()` は 1970-01-01 になります。ローカルの `wrangler dev` は I/O と関係なく時計が進むので、ここだけ現在時刻が出て気づきにくくなります。

## タイムゾーンにも注意

時刻が取れるようになっても、Workers は UTC で動きます。JST 前提のコードはズレます。

- `new Date()` が持つ値は UTC。表示は `toLocaleString('ja-JP', { timeZone: 'Asia/Tokyo' })` のように `timeZone` を明示する。
- Cron Triggers（`wrangler.toml` の `crons`）も UTC 基準。日本時間のつもりで組むと9時間ずれる。JST 9:00 に動かしたいなら UTC 0:00（`0 0 * * *`）で書く。

## 実際に踏んだケース

zod のスキーマをフロント・バックの共有箇所に置いていて、未来日を弾くバリデーションで踏みました。

```ts
// 共有スキーマ
const schema = z.object({
  date: z.date().max(new Date(), { message: "未来日は指定できません" }),
});
```

`.max(new Date())` の `new Date()` は、スキーマを構築した時点で一度だけ評価されます。この共有スキーマはモジュールの読み込み時（グローバルスコープ）に構築されるので、Workers 上ではその値が 1970-01-01 になり、上限が 1970 年に固定されて、現在の日付まで「未来日」として弾かれました。共有スキーマなのでフロントでは正常に動き、Workers 側だけで再現するのも気づきにくかった点です。

`refine` に変えると、コールバックはバリデーションが走るたびに呼ばれるので、`new Date()` もその時点で評価されて直ります。

```ts
const schema = z.object({
  date: z.date().refine((d) => d <= new Date(), { message: "未来日は指定できません" }),
});
```

## 参考

- [Performance and timers · Cloudflare Workers docs](https://developers.cloudflare.com/workers/runtime-apis/performance/)
- [JavaScript and web standards · Cloudflare Workers docs](https://developers.cloudflare.com/workers/runtime-apis/web-standards/) — `Date.now()` が最後の I/O の時刻を返すことに触れている
- [Mitigating Spectre and Other Security Threats（Kenton Varda / Cloudflare Blog）](https://blog.cloudflare.com/mitigating-spectre-and-other-security-threats-the-cloudflare-workers-security-model/) — なぜ実行中は時計を止めるのか、その背景
