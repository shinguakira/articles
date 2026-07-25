---
title: OpenAI の独占がなくなったことによる、利用クラウド種別の割合への影響
tags:
  - OpenAI
  - Azure
  - AWS
  - GCP
  - クラウド
  - 生成AI
private: false
---

## OpenAI の Azure 独占が解消された

2026年4月27日、Microsoft と OpenAI が独占パートナーシップの解消（契約の改定）を発表し、OpenAI が AWS や Google Cloud でも自社モデルを提供できるようになりました（[VentureBeat](https://venturebeat.com/technology/microsoft-and-openai-gut-their-exclusive-deal-freeing-openai-to-sell-on-aws-and-google-cloud)）。翌日には **AWS Bedrock** 上で OpenAI モデルが利用可能になっています（[AWS 公式](https://aws.amazon.com/bedrock/openai/)）。

## これまでの縛り：GPT のためだけに Azure を併用していた

独占の実務的な影響はここでした。GPT クラスのモデルをエンタープライズ規模で使うには、実質 **Azure OpenAI Service 一択**。その結果、基盤は AWS や GCP なのに **GPT の利用のためだけに Azure を別途契約する**、という構成が強制されていました。

## 解消後：Azure に分散していた依存を基盤に一本化できる

独占が外れたことで、この分散が不要になります。AWS 基盤の企業は、Azure サブスクを持たずに **Bedrock 上で OpenAI モデルを呼べる**ようになりました（一般提供済み。[AWS 公式](https://aws.amazon.com/bedrock/openai/)）。GPT のためだけの Azure 契約を畳んで、依存を基盤クラウドにまとめられます。調達も「OpenAI＝Azure」の一択交渉から、複数クラウドが同じワークロードを取り合うマルチベンダー構図に変わります（[Windows News](https://windowsnews.ai/article/openai-breaks-cloud-exclusivity-microsoft-and-aws-reshape-enterprise-ai-leverage.415898)）。

なお GCP については契約上の障壁が外れただけで、本稿時点では Google Cloud（Vertex AI）での提供条件は未発表です。「どのクラウドでも呼べる」状態が実際に揃うのはこれからです。

### 例外：新モデルだけは今も Azure が先行

新しいフロンティアモデルはまず Azure に約4か月先行提供され、その後で他クラウドへ展開されます（[Windows News](https://windowsnews.ai/article/openai-breaks-cloud-exclusivity-microsoft-and-aws-reshape-enterprise-ai-leverage.415898)）。Azure を使っていない企業でも、**最新モデルの Day-1 アクセスが要る間だけは Azure を併用せざるを得ない**場面が残ります。

## まとめ：基盤の選定そのものは変わらない

今回変わったのは「GPT を使うためにどのクラウドが要るか」であって、「基盤をどのクラウドにするか」ではありません。GPT のための Azure 併用は不要になりましたが、AWS・Azure・GCP のどれを基盤にするかは、従来どおり既存資産・コスト・運用体制で決まります。

選定を左右するのはモデルの可用性ではなく **データがどこにあるか（データの重力）** です。特に BigQuery にデータを持つ企業は、その近くで分析を回すために GCP を選ぶ動機が残ります。
