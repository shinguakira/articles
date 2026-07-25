- UIで閉じ込めて、日付の関わるものはカレンダーもしくは表示ラベルに限定するか
- 登録系は、行きと帰りの考慮。登録時と表示時
- frontend backendで時刻の使用関数をまとめる。一旦無理な共通化は行わず、現状把握のため集約飲みさせる。
UIコンポーネントの強制が必要かどうか。


首位う約するイメージ

まず、frontend / backend / shared に散らばっている時刻関連関数を調査します。共通化できる（完全一致する）ものを見極める必要があるので、実物を確認します。


Searched code, found files
時刻専用ファイルが両方にあります。中心となる3ファイルを読みます。


Read 3 files
3ファイルを把握しました。frontend と backend で完全一致する関数を特定するため、他の重複候補と各関数の使用箇所を確認します。


Ran a command, searched code
root に datetime-format-design.md / datetime-type-format-final.md という未追跡の設計ドキュメントがあります。この作業に直結しそうなので確認します。


Read 2 files
設計ドキュメントは型/フォーマット方針全体の話で、今回の「関数のshared集約」とは別スコープです。作業対象を絞ります。

toDateString が FE/BE 双方に存在し、実質同一です。使用箇所と import 元を確認します。

Searched
from ["'].*(date-utils|time-utils)["']
明確になりました。FE/BE 両方に存在する完全一致関数は toDateString のみです。他はどちらか片側にしか存在しません。

shared パッケージの構成（ビルド方式・index・package名）を確認します。


Read 2 files
スコープを確定します。調査結果はこうです。

FE/BE で完全一致する関数は toDateString だけです（実装は実質同一、FEがundefinedも受ける上位互換）。それ以外は片側にしか存在しません:

場所	関数
FE date-utils.ts	toDate, formatDate, formatRelativeDate, formatTime, formatMessageDate, normalizeDate
FE time-utils.ts	formatRelativeTime, isWithinHours, isSameDay, WEEKDAYS
BE date-utils.ts	toIsoString, toIsoDateString, toISOString, buildDayRange
両方（完全一致）	toDateString
shared は dist ビルド参照なので、追加後は shared の再ビルドが要ります。

ここで方針が2通りに割れます。ご指示の「完全一致だけ共通化」と「少なくとも時刻関連関数はsharedを絶対参照」が指す範囲が異なるため、確認させてください。

