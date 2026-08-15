---
title: Electron における IPC 通信の型安全性の問題と対策
tags:
  - TypeScript
  - Electron
  - IPC
  - 型安全
private: false
qiita_url: https://qiita.com/ShinguAkira/items/8fe6bf510aec1cf761a3
---

# Electron における IPC 通信の型安全性の問題と対策

## 問題：型チェックの欠如

Electron の IPC 通信において、`ipcRenderer.invoke` と `ipcMain.handle` の間では型チェックが行われません。これは異なるプロセス間で通信が行われるため、TypeScript の静的型チェックが適用されないことが原因です。

### 具体的な問題例

#### 例1：引数の型の不一致

```typescript
// メインプロセス (main.ts)
ipcMain.handle('get-user', async (_event, userId: number) => {
  // userId が数値であることを期待
  console.log(`数値として処理: ${userId + 1}`);
  return await database.getUserById(userId);
});

// プリロードスクリプト (preload.ts)
contextBridge.exposeInMainWorld('api', {
  // ここで文字列を渡してしまっても型チェックエラーは発生しない
  getUser: (userId: string) => ipcRenderer.invoke('get-user', userId)
});

// レンダラープロセス
// 文字列が渡され、メインプロセスではそれを数値として処理しようとする
const user = await window.api.getUser("1");  // ランタイムエラーの可能性！
```

#### 例2：複雑なオブジェクトの不一致

```typescript
// 共通の型定義（本来あるべき姿）
interface UserFilter {
  ageRange: { min: number; max: number };
  active: boolean;
}

// メインプロセス
ipcMain.handle('search-users', async (_event, filter: UserFilter) => {
  // filter.ageRangeの存在とその型を前提としたコード
  const { ageRange, active } = filter;
  return await database.findUsers({
    minAge: ageRange.min,  // filter.ageRangeが存在しない場合はエラー
    maxAge: ageRange.max,
    isActive: active
  });
});

// プリロードスクリプト
contextBridge.exposeInMainWorld('api', {
  // 不完全なオブジェクトが渡せてしまう
  searchUsers: (activeOnly: boolean) => ipcRenderer.invoke(
    'search-users', 
    { active: activeOnly }  // ageRangeが欠けている！
  )
});
```

#### 例3：返り値の型の不一致

```typescript
// メインプロセス
ipcMain.handle('get-settings', () => {
  // 文字列を返す
  return "設定文字列";
});

// プリロードスクリプト
contextBridge.exposeInMainWorld('api', {
  // 数値を期待しているが、実際には文字列が返る
  getSettings: (): Promise<number> => ipcRenderer.invoke('get-settings')
});

// レンダラープロセス
const settings = await window.api.getSettings();
console.log(settings * 2);  // 文字列に数値演算を適用してエラー
```

## 対処法

### 1. 共通の型定義ファイルの作成

```typescript
// shared-types.ts（メインプロセスとレンダラー両方からアクセス可能な場所に配置）
export interface User {
  id: number;
  name: string;
  email: string;
}

export interface UserFilter {
  ageRange: { min: number; max: number };
  active: boolean;
}

// IPC通信用の型定義
export interface IpcChannels {
  'get-user': {
    params: [userId: number];
    result: User | null;
  };
  'search-users': {
    params: [filter: UserFilter];
    result: User[];
  };
  'get-settings': {
    params: [];
    result: string;
  };
}
```

### 2. 型安全なラッパー関数の作成

```typescript
// preload.ts
import { IpcChannels } from './shared-types';

// 型安全なinvoke関数
function typedInvoke<C extends keyof IpcChannels>(
  channel: C,
  ...args: IpcChannels[C]['params']
): Promise<IpcChannels[C]['result']> {
  return ipcRenderer.invoke(channel, ...args);
}

// 型安全なAPI
contextBridge.exposeInMainWorld('api', {
  getUser: (userId: number) => typedInvoke('get-user', userId),
  searchUsers: (filter: UserFilter) => typedInvoke('search-users', filter),
  getSettings: () => typedInvoke('get-settings')
});
```

### 3. ハンドラー登録のラッパー関数

```typescript
// main.ts
import { IpcChannels } from './shared-types';

// 型安全なハンドラー登録関数
function typedHandle<C extends keyof IpcChannels>(
  channel: C,
  handler: (
    event: Electron.IpcMainInvokeEvent,
    ...args: IpcChannels[C]['params']
  ) => Promise<IpcChannels[C]['result']> | IpcChannels[C]['result']
): void {
  ipcMain.handle(channel, handler);
}

// 型安全に登録
typedHandle('get-user', async (_event, userId) => {
  // ここで userId は必ず number 型
  return await database.getUserById(userId);
});

typedHandle('search-users', async (_event, filter) => {
  // filter は必ず UserFilter 型を満たす
  const { ageRange, active } = filter;
  return await database.findUsers({ minAge: ageRange.min, maxAge: ageRange.max, isActive: active });
});
```

### 4. ランタイム検証による二重保護

```typescript
// zod を使用した検証
import { z } from 'zod';
import { UserFilter } from './shared-types';

// 検証スキーマ
const UserFilterSchema = z.object({
  ageRange: z.object({
    min: z.number().min(0).max(120),
    max: z.number().min(0).max(120)
  }).refine(data => data.min <= data.max, {
    message: "Min age must be less than or equal to max age"
  }),
  active: z.boolean()
});

// ハンドラー登録時に検証を追加
ipcMain.handle('search-users', async (_event, filterInput) => {
  try {
    // ランタイム検証
    const filter = UserFilterSchema.parse(filterInput) as UserFilter;
    
    // 検証が通ったら処理を続行
    const { ageRange, active } = filter;
    return await database.findUsers({
      minAge: ageRange.min,
      maxAge: ageRange.max,
      isActive: active
    });
  } catch (error) {
    // 検証エラーを適切に処理
    console.error('Invalid filter object:', error);
    throw new Error('無効なフィルター条件が指定されました');
  }
});
```

### 5. チャネル名の定数化

```typescript
// ipc-constants.ts
export const IPC_CHANNELS = {
  GET_USER: 'get-user',
  SEARCH_USERS: 'search-users',
  GET_SETTINGS: 'get-settings'
} as const;

// main.ts
import { IPC_CHANNELS } from './ipc-constants';

ipcMain.handle(IPC_CHANNELS.GET_USER, async (_event, userId) => {
  // 処理
});

// preload.ts
import { IPC_CHANNELS } from './ipc-constants';

const api = {
  getUser: (userId: number) => ipcRenderer.invoke(IPC_CHANNELS.GET_USER, userId)
};
```

### 6. レスポンスの標準化

```typescript
// APIレスポンスの型
type ApiResponse<T> = 
  | { success: true; data: T }
  | { success: false; error: string };

// メインプロセス
ipcMain.handle('get-user', async (_event, userId) => {
  try {
    const user = await database.getUserById(userId);
    if (!user) {
      return { success: false, error: 'User not found' };
    }
    return { success: true, data: user };
  } catch (error) {
    return { 
      success: false, 
      error: error instanceof Error ? error.message : 'Unknown error' 
    };
  }
});

// プリロード
contextBridge.exposeInMainWorld('api', {
  getUser: async (userId: number): Promise<User> => {
    const response = await ipcRenderer.invoke('get-user', userId) as ApiResponse<User>;
    if (!response.success) {
      throw new Error(response.error);
    }
    return response.data;
  }
});
```

## 実践的なアプローチ

実際のアプリケーション開発では、これらの手法を組み合わせて使用するのが効果的です：

1. **基本**: 共有型定義と型安全なラッパー関数
2. **安全性強化**: ランタイム検証の追加
3. **保守性向上**: チャネル名の定数化とレスポンスの標準化

このアプローチにより、コンパイル時の型チェックとランタイム検証の両方を活用して、より堅牢なIPC通信を実現できます。
