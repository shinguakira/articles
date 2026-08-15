---
title: SELECT list is not in GROUP BY clause and contains nonaggregated column MyBatis使用時のSQLエラー解消
tags:
  - Java
  - MySQL
  - SQL
  - MyBatis
private: false
qiita_url: https://qiita.com/ShinguAkira/items/c53194b1ba4cf88f901d
---

# SELECT list is not in GROUP BY clause and contains nonaggregated column MyBatis使用時のSQLエラー
Java,Mybatisが生成したMapperを使用していて、当エラーが出ている場合、以下SQLモードを使用しているmySQLで設定することで解消できます。

このコマンドを実行するのは、mySQLにログインして、コマンドを実行します。
または、A5:SQL Mk-2を使って、データベースに接続し、クエリ実行できます。
ダウンロード→https://a5m2.mmatsubara.com/
A5:SQL Mk-2でのコマンドの実行の仕方→https://qiita.com/kanfutrooper/items/a4db929306efcc91b610#:~:text=%E3%81%AE%E3%81%BF%E5%AE%9F%E8%A1%8C%E5%8F%AF%E8%83%BD%E3%80%82-,%E3%82%B3%E3%83%9E%E3%83%B3%E3%83%89%EF%BC%88Ctl%20%2B%20Enter%EF%BC%89%E3%81%A7%E5%AE%9F%E8%A1%8C,%E3%81%B0SQL%E3%82%92%E5%AE%9F%E8%A1%8C%E5%8F%AF%E8%83%BD%E3%80%82

```
SET GLOBAL sql_mode = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';
```

ちなみに
```
SELECT @@GLOBAL.sql_mode;
```
で現在設定されているSQLモードを確認できます。
NO_ZERO_IN_DATE,NO_ZERO_DATEもしくはONLY_FULL_GROUP_BYが含まれるとこのエラーが発生する気がします。

ちなみに、A5:SQL Mk-2にてSET GLOBAL sql_mode実行後、SELECT @@GLOBAL.sql_modeでモードの確認をされる場合、アプリ自体を再起動もしくは、再接続を行う必要があるかもしれません。

## 余談
MyBatisで生成される一部の条件設定でのクエリのコードにdistinctが含まれるのですが、その「distinctを消すことで解消できました！」なんて話もあったんですが、自動生成するたびにそれを行う必要がある上に、テーブルが多い場合、いちいち修正するのは論外かと思うので、SQLモードについて、プロジェクトで縛りがなければ、上記設定することで、解消できるかと思います。
