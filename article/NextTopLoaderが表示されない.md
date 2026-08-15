---
title: NextTopLoaderが表示されない
tags:
  - 画面遷移
  - loading
  - React
  - Next.js
  - loader
private: false
qiita_url: https://qiita.com/ShinguAkira/items/9680bd7839994e1418dd
---

# Next.jsでプログラマティックなナビゲーション時にローディングバーを表示する方法

## はじめに

Next.jsでアプリケーションを開発していると、ページ遷移時にユーザーに視覚的なフィードバックを提供したいケースがあります。特に、Link コンポーネントを使用せずに `router.push()` などでプログラマティックにナビゲーションする場合、ローディングインジケーターが表示されないことがあります。

この記事では、Next.jsアプリケーションで `nextjs-toploader` (NextTopLoader) を使用して、プログラマティックなナビゲーション時にもローディングバーを表示する方法について解説します。

## 問題: Link以外のナビゲーションでローディングバーが表示されない

通常、Next.jsの `<Link>` コンポーネントを使ってページ間を移動する場合、NextTopLoaderは自動的にページ遷移を検知してローディングバーを表示します。しかし、以下のようなケースではローディングバーが表示されないことがあります：

- `router.push()` や `router.replace()` を使用したプログラマティックなナビゲーション
- フォーム送信後のリダイレクト
- ボタンクリックやその他のイベントハンドラーからのナビゲーション

## 解決策: useTopLoader フックの活用

`nextjs-toploader` パッケージが提供する `useTopLoader` フックを使うことで、この問題を解決できます。このフックを使って手動でローディングバーの開始と終了を制御しましょう。
[nextjs-toploader定義](https://www.npmjs.com/package/nextjs-toploader)


## 実装方法


### 1. NextTopLoaderの設定

[layout.tsx](cci:7://file:///e:/workspace/brighty/frontend/src/app/%28tool%29/rooms/layout.tsx:0:0-0:0) ファイルでNextTopLoaderを設定します：

```tsx
'use client';

import { NextTopLoader } from 'nextjs-toploader';
import { ReactNode } from 'react';

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="ja">
      <body>
        <NextTopLoader 
          color="#2299DD"
          initialPosition={0.08}
          height={3}
          showSpinner={false}
        />
        {children}
      </body>
    </html>
  );
}
```

### 2. useTopLoader フックの使用

プログラマティックなナビゲーションを行うコンポーネントで `useTopLoader` フックを使用します：

```tsx
'use client';

import { useTopLoader } from 'nextjs-toploader';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button'; // 例としてのUIコンポーネント

export default function NavigationExample() {
  const router = useRouter();
   const { start, done } = useTopLoader()

  const handleNavigation = async () => {
    // ローディングバーの開始
    start();
    
    try {
      // 何らかの非同期処理（APIリクエストなど）
      await someAsyncFunction();
      
      // 処理が完了したらページ遷移
      router.push('/destination-page');
      
      // Note: router.push()後はコンポーネントがアンマウントされるため、
      // 通常はdone()は不要です
    } catch (error) {
      // エラーが発生した場合はローディングを停止
      done();
      console.error('Navigation error:', error);
    }
  };

  return (
    <div>
      <h1>プログラマティックナビゲーションの例</h1>
      <Button onClick={handleNavigation}>ページ遷移</Button>
    </div>
  );
}
```

## 注意点

1. `useTopLoader` フックは `'use client'` ディレクティブを使用したクライアントコンポーネント内でのみ使用できます。

2. ページ遷移が成功した場合、通常はコンポーネントがアンマウントされるため、`done()` を明示的に呼び出す必要はありません。エラー処理時のみ呼び出すことを推奨します。

3. 非同期処理がない単純なナビゲーションの場合は、以下のように簡略化できます：

```tsx
const handleClick = () => {
  start();
  router.push('/some-page');
};
```
