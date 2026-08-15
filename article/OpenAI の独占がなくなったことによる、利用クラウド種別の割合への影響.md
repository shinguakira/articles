---
title: OpenAI の独占がなくなったことによる、利用クラウド種別の割合への影響
tags:
  - AWS
  - Azure
  - GoogleCloud
  - Cloud
  - OpenAI
private: false
qiita_url: https://qiita.com/ShinguAkira/items/fb04e65ef707dac3969a
---

## OpenAI の Azure 独占が解消された(今更ですが。)

2026年4月27日、Microsoft と OpenAI が独占パートナーシップの解消（契約の改定）を発表し、OpenAI が AWS や Google Cloud でも自社モデルを提供できるようになりました（[VentureBeat](https://venturebeat.com/technology/microsoft-and-openai-gut-their-exclusive-deal-freeing-openai-to-sell-on-aws-and-google-cloud)）。翌日には **AWS Bedrock** 上で OpenAI モデルが利用可能になっています（[AWS 公式](https://aws.amazon.com/bedrock/openai/)）。

## これまでの縛り：GPT のためだけに Azure を併用していた

GPT クラスのモデルを会社のプロダクトなどで使うには、実質 **Azure OpenAI Service 一択**。
各会社の技術スタックなどやアーキテクチャー図を見ると、AWSやGCPが基本にも関わらず、なぜかAIのところだけAzureという構成が度々見られた。**GPT の利用のためだけに Azure を別途契約する**、必要があった状態だったかと今になって思います。

## 解消後：Azure に分散していた依存を基盤に一本化できる

独占が外れたことで、この分散が不要になります。AWS 基盤の企業は、Azure サブスクを持たずに **Bedrock 上で OpenAI モデルを呼べる**ようになりました（一般提供済み。[AWS 公式](https://aws.amazon.com/bedrock/openai/)）。GPT のためだけの Azure 契約をする必要がなくなり、依存を基盤クラウドにまとめられます。

なお GCP については契約上の障壁が外れただけで、本稿時点では Google Cloud（Vertex AI）での提供条件は未発表です。

### 例外：新モデルだけは今も Azure が先行

新しいフロンティアモデルはまず Azure に約4か月先行提供され、その後で他クラウドへ展開されます（[Windows News](https://windowsnews.ai/article/openai-breaks-cloud-exclusivity-microsoft-and-aws-reshape-enterprise-ai-leverage.415898)）。Azure を使っていない企業でも、**最新モデルの Day-1 アクセスが要る間だけは Azure を併用せざるを得ない**場面が残ります。

## まとめ：基盤の選定そのものは変わらない

GPT のための Azure 併用は不要になり、局所的なAzure利用は減るかと思いますが、AWS・Azure・GCP のどれを基盤にするかは、従来どおりです。
AWSかAzureか、すでに導入している会社はどちらかに移行というのはあまりないかと思います。そのため特に全体のクラウド利用シェアには今回の話は大した影響はないと考えています。

# その他
独占とは違いますが、似たようなところでAWS,AzureでもBigQueryを使っているところを見るので、BigQueryも似たような別クラウド利用の例かと思います。
