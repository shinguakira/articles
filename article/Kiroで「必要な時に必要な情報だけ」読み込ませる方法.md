---
title: Kiroで「必要な時に必要な情報だけ」読み込ませる方法
tags:
  - AWS
  - AI駆動開発
  - Kiro
  - 仕様駆動開発
private: false
qiita_url: https://qiita.com/ShinguAkira/items/bfecef60c1615d5761a6
---

---
Kiroで「必要な時に必要な情報だけ」読み込ませる方法
---

ファイルの先頭に特定の記述をすることで、特定のファイル読み込み時に特定のsteeringを読み込ませることができる。公式にも記載があったが、探しづらかったので、備忘として記録。
初期は、always,manual,fileMatchだった気がするが、autoも追加されていた。

## 参考

https://kiro.dev/docs/steering/

## inclusionモード一覧

**`always`**（デフォルト）は全やり取りで常時読み込まれる。フロントマターを省略した場合もこれ。

**`manual`** はチャットで `#ファイル名` と入力したときだけ読み込まれる。大きなリファレンスや稀にしか使わない手順書向け。

## `fileMatch` — 特定ファイルを触ったときだけ読み込む
```yaml
---
inclusion: fileMatch
fileMatchPattern: "**/*.test.ts"
---
```

`src/*.ts` を開いた時だけAPIルールを読み込む、テストファイルを編集中だけテスト規約を注入する、といった使い方ができる。複数パターンも配列で指定可能。
```yaml
fileMatchPattern:
  - "src/**/*.ts"
  - "lib/**/*.js"
```

## `auto` — リクエスト内容に応じて自動判断
```yaml
---
inclusion: auto
name: api-design
description: REST API design patterns and conventions. Use when creating or modifying API endpoints.
---
```

[Skills](https://kiro.dev/docs/skills) に近い仕組みで、`description` に書いた内容をKiroがリクエストと照合し、関連すると判断した時だけ読み込む。`fileMatch` はファイル種別での自動化だが、`auto` はAIが意味を理解して判断する。`/` のスラッシュコマンドから手動で呼び出すことも可能。

| フィールド | 必須 | 説明 |
|-----------|------|------|
| `name` | ○ | ステアリングの識別子。表示や照合に使われる |
| `description` | ○ | いつ読み込むかの説明。Kiroがリクエストと照合する |

## 整理すると

| モード | いつ読む | 向いているもの |
|--------|----------|---------------|
| `always` | 毎回 | スタック・規約・ポリシー |
| `fileMatch` | 該当ファイル編集時 | ドメイン特有のルール |
| `auto` | リクエストが合致した時 | 専門的なワークフロー・複雑なルール |
| `manual` | 明示的に呼び出した時 | 手順書・重いリファレンス |

## ファイルの置き場所
```
.kiro/steering/
  tech.md          # always
  api-standards.md # fileMatch: src/api/**
  auth-flow.md     # auto: 認証周りの実装時に自動判断
  migration.md     # manual
```

## **その他**
以前の状態や、その時のアップデートで変わるが、以前specsをspecs直下ではなく、特定のフォルダを作ってその中に作らせていたが、fileMatchで、.kiro/specs全体をmatchパターンに設定して、requirements作成前にフォルダを作るような構成にしていたが、うまくいかず、結局その内容をalways読み込みにしたり、AGENTS.mdに記載することで動作するようになった。まだ細かいところの動作が安定しないため、場合によってはその際に対応が必要かもしれません。

## **まとめ** 
デフォルトや、AGENTS.mdに書きすぎると、関係ないcontextを消費してしまうため、
単純に、記載内容を減らすのも良いが、必要なファイルだけど必要な時に読み込ませたい。
例:requimrentsを作成するときは、コーディングが技術関連の情報は基本必要ないはずなので読み込ませない等。
