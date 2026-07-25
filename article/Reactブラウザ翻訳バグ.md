---
title: 多言語対応アプリで「特定のユーザーだけ画面が真っ白」になる — 原因はブラウザ翻訳
tags:
  - React
  - ブラウザ翻訳
  - i18n
  - フロントエンド
  - JavaScript
private: false
---

## 現象

多言語対応の React アプリで、一部のユーザーだけ画面が真っ白になることがある。開発者の手元では再現せず、コンソールに次の例外が1つ出るだけ。

```
NotFoundError: Failed to execute 'removeChild' on 'Node':
The node to be removed is not a child of this node.
```

原因はブラウザの翻訳機能で、ブラウザ翻訳を使うのはUIと母語が違うユーザーなので、特定のユーザーだけで起き、再現しにくいバグになる。

## 参考

この問題についていちばん詳しいのが Martijn Hols の記事で、原因・各対策の効果・限界まで書かれているのでまずここを読む。

- [Everything about Google Translate crashing React - Martijn Hols](https://martijnhols.nl/blog/everything-about-google-translate-crashing-react)

Zenn の記事はこの Martijn Hols 記事を参照しており、実際の障害事例と ESLint での対策についてまとめている。

- [ReactのページでGoogle翻訳するとエラーになる事象 - Zenn](https://zenn.dev/chot/articles/google_translate_crashing_react)

その他

- [Make React resilient to DOM mutations from Google Translate · Issue #11538 · facebook/react](https://github.com/facebook/react/issues/11538)
- [Why React App May Be Broken By Google Translate Extension - DEV](https://dev.to/ivanturashov/preventing-react-crashes-handling-google-translate-5bi0)

## 原因

ブラウザ翻訳は、ページ内のテキストノードを翻訳済みのテキストに置き換えるが、Google 翻訳の場合は元のテキストノードを `<font>` 要素に差し替える。

```html
<!-- 元 -->
<p>合計 100 円</p>

<!-- 翻訳後: テキストノードが font に差し替わる -->
<p><font>Total: 100 yen</font></p>
```

React はレンダー時に生成した DOM ノードへの参照を持っていて、再レンダーのたびにそのノードを更新・削除・挿入する。翻訳で元のテキストノードが差し替えられると参照しているノードはもう親の子ではなくなるので、`removeChild` や `insertBefore` が例外を投げて描画が止まり、白画面になる。

これは React 特有の問題ではなく、DOM ノードへの参照を持って更新する仕組み（Vue・Angular・Svelte や自前の DOM 操作）なら同様に起こる。React で報告が多いだけ。

## 対策

### translate="no" で翻訳対象から外す

動的に書き換わる部分を `translate="no"` にすると、クラッシュも翻訳のズレも防げる。

```jsx
<span translate="no">{count}件</span>
```

ページ全体なら `<html translate="no">` でも無効化できるが、その部分が翻訳されなくなるので多言語ユーザーの利便性と引き換えになる。

### span で包む（部分的な緩和）

動的なテキストを `<span>` で包むとクラッシュは減るが、全部は直らないので確実な解決策ではない（詳細は参考リンク）。

### ESLint での予防（不十分）

この問題を検出する ESLint プラグイン。

**eslint-plugin-react-google-translate**

- `react-google-translate/no-conditional-text-nodes-with-siblings`
- `react-google-translate/no-return-text-nodes`

**eslint-plugin-sayari**

- `no-unwrapped-jsx-text` — JSX内でラップされていないテキストノードを検出するルール

危険なパターンは検出できるが、完全には防げない。

## Reactの対応状況

Issue #11538 は2017年から open のままで、`removeChild` / `insertBefore` を上書きして握りつぶす回避策もあるが根本解決ではない。ブラウザ翻訳が DOM を書き換えるのが原因なので、React 単体では直せない。
