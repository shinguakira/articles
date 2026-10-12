# Java Go移行で使用メモリ効率化しました

おおよそ90%のメモリ改善ができました。特にこの改善によるCPU使用率、処理速度の悪化などはありません。

# 背景
AzureでデプロイしているAPIのメモリ使用量が90%以上張り付きになってしまい、マシンなどのスペックを上げたり等で対応はできるが、今後アプリが増えて、その都度スペックを上げることになってしまうため、今できる部分でできるだけメモリ効率を上げたいと思いJavaからGoに段階的に移行しました。

# ボトルネック
単にJavaがGoに比べてメモリを消費および保持が長い部分はありましたが、マイクロサービス化して、別プロセスや別マシン起動のAPIにした場合に、JVM分加算されてしまい、分割すればするほど1つの巨大なAPIと比べ、積み重ねで合計サイズが肥大化してしまっていました。


# 移行前状況
元々マイクロサービス化の計画はあったので、モジュラーモノリス的に分割していて、一部マイクロサービス化はやっていました。
すでに利用がある機能もあるため、一気に移行はせずに、API単位で段階的に移行を行いました。
e2e,Unit共にある程度揃っており、足りない部分は移行前に書かせて、振る舞いを確認できる状態にして移行しました。


## 以下改善比較結果です。(副作用に成果物サイズ差)

| API | Go成果物 | Java成果物 | 成果物削減率 | Go実行時 | Java実行時 | 実行時差分 |
|---|---:|---:|---:|---:|---:|---:|
| 外部API呼び出し API | 9.24 MiB | 18.30 MiB | 49.5% | 12.58 MiB Working Set | 150.74 MiB Working Set | 138.16 MiB / 91.7%削減 |
| Repair Chatbot Dashboard | 9.72 MiB | 54.58 MiB | 82.2% | 11.43 MiB Working Set | 480.75-611.63 MiB RSS | 469.32-600.20 MiB / 97.6-98.1%削減（参考値） |
| 外部API呼び出し API | 10.64 MiB | 35.39 MiB | 69.9% | 14.33 MiB Working Set | 189.67 MiB Working Set | 175.34 MiB / 92.4%削減 |
| データ分析 API | 17.05 MiB | 43.73 MiB | 61.0% | 15.32 MiB Working Set | 181.07 MiB Working Set | 165.75 MiB / 91.5%削減 |
| チャットボット API | 16.37 MiB | 60.85 MiB | 73.1% | 14.57 MiB Working Set | 201.82 MiB Working Set | 187.25 MiB / 92.8%削減 |
| チャットボットデータ処理 API | 15.70 MiB | 39.06 MiB | 59.8% | 12.63 MiB Working Set | 170.30 MiB Working Set | 157.67 MiB / 92.6%削減 |
| **比較可能分の合計** | **78.72 MiB（6件）** | **251.91 MiB（6件）** | **173.19 MiB / 68.8%削減** | **69.43 MiB Working Set（5件）** | **893.60 MiB Working Set（5件）** | **824.17 MiB / 92.2%削減** |

## 外部API呼び出しAPI 1プロセスの内部メモリ概算

外部API呼び出しAPI 1プロセスを、runtime・web framework・applicationの所有単位で概算したスタック。

<table>
<tr>
<td valign="top">

Java版<br>Working Set 150.74 MiB

<table border="1" cellpadding="6" cellspacing="0">
<tr><td>アプリ固有コード application<br>uses ≈ 22.61 MiB（15%）</td></tr>
<tr><td>Spring Boot / Tomcat<br>uses ≈ 75.37 MiB（50%）</td></tr>
<tr><td>JVM / JDK runtime<br>uses ≈ 52.76 MiB（35%）</td></tr>
</table>

</td>
<td valign="top">

Go版<br>Working Set 12.58 MiB

<table border="1" cellpadding="6" cellspacing="0">
<tr><td>アプリ固有コード application<br>uses ≈ 1.89 MiB（15%）</td></tr>
<tr><td>HTTP layer<br>uses ≈ 3.14 MiB（25%）</td></tr>
<tr><td>Go runtime / standard library<br>uses ≈ 7.55 MiB（60%）</td></tr>
</table>

</td>
</tr>
</table>

合計Working Setは測定値、各layerのMiBと割合は構成差を把握するための概算モデルイメージです。

# その他
メモリなどサーバーリソースは、ユーザーが多いほど増え、対応リージョンなどが増えるほど数も増えるので、リソース最適化の効果は、利用が多ければ多いほどわかりやすいものかと思います。
