# Go / Java スタンドアロンAPI比較

## 比較結果

| API | Go成果物 | Java成果物 | 成果物削減率 | Go実行時 | Java実行時 | 実行時差分 |
|---|---:|---:|---:|---:|---:|---:|
| 外部API呼び出し API | 9.24 MiB | 18.30 MiB | 49.5% | 12.58 MiB Working Set | 150.74 MiB Working Set | 138.16 MiB / 91.7%削減 |
| Repair Chatbot Dashboard | 9.72 MiB | 54.58 MiB | 82.2% | 11.43 MiB Working Set | 480.75-611.63 MiB RSS | 469.32-600.20 MiB / 97.6-98.1%削減（参考値） |
| 外部API呼び出し API | 10.64 MiB | 35.39 MiB | 69.9% | 14.33 MiB Working Set | 189.67 MiB Working Set | 175.34 MiB / 92.4%削減 |
| データ分析 API | 17.05 MiB | 43.73 MiB | 61.0% | 15.32 MiB Working Set | 181.07 MiB Working Set | 165.75 MiB / 91.5%削減 |
| チャットボット API | 16.37 MiB | 60.85 MiB | 73.1% | 14.57 MiB Working Set | 201.82 MiB Working Set | 187.26 MiB / 92.8%削減 |
| チャットボットデータ処理 API | 15.70 MiB | 39.06 MiB | 59.8% | 12.63 MiB Working Set | 170.30 MiB Working Set | 157.67 MiB / 92.6%削減 |
| **比較可能分の合計** | **78.72 MiB（6件）** | **251.91 MiB（6件）** | **173.19 MiB / 68.8%削減** | **69.43 MiB Working Set（5件）** | **893.60 MiB Working Set（5件）** | **824.17 MiB / 92.2%削減** |

## 外部API呼び出しAPI 1プロセスの内部メモリ概算

外部API呼び出しAPI 1プロセスを、runtime・web framework・applicationの所有単位で概算したスタック。

```text
Java版                                      Go版
Working Set 150.74 MiB                      Working Set 12.58 MiB

┌──────────────────────────────┐            ┌──────────────────────────────┐
│ アプリ固有コード application          │            │ アプリ固有コード  application          │
│ uses ≈ 22.61 MiB（15%）       │            │ uses ≈ 1.88 MiB（15%）        │
├──────────────────────────────┤            ├──────────────────────────────┤
│ Spring Boot / Tomcat         │            │ HTTP layer                   │
│ uses ≈ 75.37 MiB（50%）       │            │ uses ≈ 3.15 MiB（25%）        │
├──────────────────────────────┤            ├──────────────────────────────┤
│ JVM / JDK runtime            │            │ Go runtime / standard library│
│ uses ≈ 52.76 MiB（35%）       │            │ uses ≈ 7.55 MiB（60%）        │
└──────────────────────────────┘            └──────────────────────────────┘
```

合計Working Setは測定値、各layerのMiBと割合は構成差を把握するための概算モデル。
