---
title: Azure でのフロントエンド監視
tags:
  - Azure
  - 監視ログ
  - ApplicationInsights
  - エラー監視
  - TypeScript
private: false
qiita_url: https://qiita.com/ShinguAkira/items/0ab7f408fbffa20d4a87
---

## Azure でのフロントエンド監視

Azure Monitor の **Application Insights** には、ブラウザ用の JavaScript SDK が公式に用意されています。ブラウザで起きた未捕捉の例外や失敗した API 呼び出し、ページ遷移などをこの SDK が送り、Portal で確認できます。やっていることは Sentry とほぼ同じで、集めたデータは自分の Azure サブスクリプションにある Log Analytics ワークスペースに入ります。

実装はこのリポジトリに置いています。(デモ程度ですが)

- [shinguakira/azure-frontend-monitoring](https://github.com/shinguakira/azure-frontend-monitoring)

## 集められるもの

- エラー（未捕捉のエラー、未捕捉の Promise reject、React のレンダー中の例外、API の失敗）
- ページ遷移
- ページの滞在時間
- クリック

ただ、`console.log` や `console.error` は自動では拾ってくれません。そこは [`consoleCapture.ts`](https://github.com/shinguakira/azure-frontend-monitoring/blob/main/src/monitoring/consoleCapture.ts) で `console.error` と `console.warn` を差し替えて送っています。

## 動作確認

以下はリポジトリで実装している画面です。ボタンでエラーやログを起こすと、送った内容が右側の画面ログに出ます。

![動作確認用パネル](https://raw.githubusercontent.com/shinguakira/azure-frontend-monitoring/main/docs/frontend-demo.png)

同じものが Application Insights にも届いていて、Portal で 1 件ずつ見られます。(撮った画像はとった時間の問題で時刻があっていませんが。)

![Application Insights に届いたデータ](https://raw.githubusercontent.com/shinguakira/azure-frontend-monitoring/main/docs/azure-portal.png)
## 似たサービスとの比較

Application Insights と同じ用途のサービスはいくつかあります。

| サービス | どんなものか | データの保存先 |
| --- | --- | --- |
| **Application Insights**（Azure） | Azure Monitor の一部 | 自分の Azure サブスクリプション内の Log Analytics |
| **Sentry** | エラー監視の SaaS。セルフホストもできる | Sentry のクラウド、または自前のサーバー |
| **Datadog RUM** | Datadog の監視サービスの一部 | Datadog のクラウド |
| **Grafana Faro** | OSS。Grafana / Loki / Tempo と組み合わせて使う | 自前、または Grafana Cloud |

## バンドルサイズ

| パッケージ | バージョン | minified | min+gzip |
| --- | --- | --- | --- |
| [`@microsoft/applicationinsights-web`](https://bundlephobia.com/package/@microsoft/applicationinsights-web@3.4.4) | 3.4.4 | 186.5 KB | **72.9 KB** |
| [`@microsoft/applicationinsights-react-js`](https://bundlephobia.com/package/@microsoft/applicationinsights-react-js@19.4.2) | 19.4.2 | 16.2 KB | **6.4 KB** |
| [`@grafana/faro-web-sdk`](https://bundlephobia.com/package/@grafana/faro-web-sdk@2.12.0) | 2.12.0 | 108.4 KB | **36.9 KB** |
| [`@datadog/browser-rum`](https://bundlephobia.com/package/@datadog/browser-rum@7.13.0) | 7.13.0 | 188.3 KB | **64.3 KB** |
| [`@sentry/react`](https://bundlephobia.com/package/@sentry/react@10.75.0) | 10.75.0 | 462.7 KB | **150.7 KB** |

## 参考

- [Application Insights JavaScript SDK Setup in Azure Monitor - Microsoft Learn](https://learn.microsoft.com/en-us/azure/azure-monitor/app/javascript-sdk) — ブラウザ用 SDK の入れ方
