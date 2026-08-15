---
title: Node.jsのファイル間共有グローバル変数について
tags:
  - Node.js
  - 状態管理
  - グローバル変数
  - ファイル間共有変数
private: false
qiita_url: https://qiita.com/ShinguAkira/items/2b1d7567c7e62de45f03
---

# Node.jsのファイル間共有グローバル変数について

## 基本的な仕組み

Node.jsでは、ファイル間でグローバル変数を簡単に共有できます。これはNode.jsのモジュールキャッシュの仕組みによるもので、一度読み込まれたモジュールは同じインスタンスが再利用されます。

## 最小限の実装例

### グローバル変数を定義するファイル

```javascript
// global-store.js
// 共有したいグローバル変数
let counter = 0;

// 変数を増加させる関数
export function incrementCounter() {
  counter++;
  return counter;
}

// 現在の値を取得する関数
export function getCounter() {
  return counter;
}
```

### 別々のファイルからのアクセス

```javascript
// file1.js
import { incrementCounter, getCounter } from './global-store.js';

console.log('File1の初期カウンター値:', getCounter()); // 0
incrementCounter();
console.log('File1の処理後カウンター値:', getCounter()); // 1
```

```javascript
// file2.js
import { incrementCounter, getCounter } from './global-store.js';

console.log('File2の初期カウンター値:', getCounter()); // 1 (file1.jsで増加済み)
incrementCounter();
console.log('File2の処理後カウンター値:', getCounter()); // 2
```

```javascript
// main.js
import { getCounter } from './global-store.js';
import './file1.js';
import './file2.js';

console.log('すべての処理後のカウンター値:', getCounter()); // 2
```

## なぜこれが可能なのか？

Node.jsがモジュールを読み込むとき、以下のプロセスが発生します：

1. モジュールが最初に`import`されると、そのコードが実行され、変数が初期化されます
2. モジュールの内容はメモリにキャッシュされます
3. 同じモジュールが他の場所から`import`されると、キャッシュされたインスタンスが返されます
4. このため、モジュール内の変数は事実上「シングルトン」として機能します

モジュール内で宣言された変数（この例では`counter`）は、モジュール自体がメモリ上に存在する限り保持されます。関数（`incrementCounter`と`getCounter`）はこの変数へのクロージャとなり、どこからでもアクセスして変更できます。

## 使用上の注意

この方法は簡単ですが、大規模なアプリケーションでは変数の変更を追跡しづらくなる可能性があります。状態管理が複雑になる場合は、専用のライブラリやパターンの使用を検討することをお勧めします。
