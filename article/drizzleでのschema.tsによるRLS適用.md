---
title: drizzleでのschema.tsによるRLS適用
tags:
  - PostgreSQL
  - TypeScript
  - Drizzle
  - RLS
  - RowLevelSecurity
private: false
qiita_url: https://qiita.com/ShinguAkira/items/6bd740ec0e201525c877
---

# Drizzle ORM で schema.ts に RLS を適用する

https://orm.drizzle.team/docs/rls

## .enableRLS()（〜v0.x）

```ts:schema.ts
import { integer, pgTable } from 'drizzle-orm/pg-core';

export const users = pgTable('users', {
  id: integer(),
}).enableRLS();
```

## pgTable.withRLS()（v1.0.0-beta.1〜）

`.enableRLS()` は v1.0.0-beta.1 で非推奨になり、`pgTable.withRLS()` に変わった。

```ts:schema.ts
import { integer, pgTable } from 'drizzle-orm/pg-core';

export const users = pgTable.withRLS('users', {
  id: integer(),
});
```

## 生成されるマイグレーション

```sql
ALTER TABLE "users" ENABLE ROW LEVEL SECURITY;
```
