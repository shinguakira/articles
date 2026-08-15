---
title: ElectronアプリでPythonを実行する際、stderr出力されるもの特異なもの
tags:
  - Python
  - Electron
  - tqdm
  - エラーハンドリング
  - stderr
private: false
qiita_url: https://qiita.com/ShinguAkira/items/d471fd6353f23da315b0
---

# ElectronアプリでPythonを実行する際、stderr出力されるもの特異なもの


ElectronアプリからPythonスクリプトを実行する際、進捗表示のために`tqdm`ライブラリを使用すると予期せぬエラーが発生することがあります。本記事では、tqdmがstderr（標準エラー出力）を使用する理由と、Electronアプリでこの問題を適切に処理する方法について解説します。

## Electronアプリでの問題点

Electronアプリから`child_process.spawn`を使用してPythonプロセスを起動する際、次のような問題が発生します：

1. **プロセス終了の誤検出**: 多くの場合、stderrに何かが出力されると「エラーが発生した」と解釈され、プロセスが強制終了されることがあります。

2. **エラー判別の複雑さ**: tqdmの進捗表示とアプリケーションの実際のエラーメッセージが同じstderrに混在するため、真のエラーを判別するのが難しくなります。

3. **UIへの反映**: stderrからの出力をUIに適切に反映させる方法が必要です。

## 解決策

### 1. Pythonコード修正が可能な場合(stderrを適切に処理する)
もしくは、outputに | #など、tqdmのみだけフィルターできそうな文字にすることで、他のモジュール実行を同コード内で行っても問題なくなる。

```javascript
// Electronアプリ内のコード
const { spawn } = require('child_process');

const pythonProcess = spawn('python', ['script_with_tqdm.py']);

// stderrのデータを処理
pythonProcess.stderr.on('data', (data) => {
  const output = data.toString();
  
  // tqdmの出力かエラーかを判別する
  if (output.includes('\r')) {
    // カーソルを行の先頭に戻す文字（\r）があれば、おそらくtqdmの進捗表示
    updateProgressBar(output);
  } else {
    // それ以外は実際のエラーとして処理
    handleError(output);
  }
});

// プロセス終了時の処理
pythonProcess.on('exit', (code) => {
  if (code !== 0) {
    console.error(`プロセスがエラーコード ${code} で終了しました`);
  }
});
```

### 2. Pythonコード修正が不可能な場合(tqdmの出力先を変更する)

Pythonスクリプト側でtqdmの出力先を制御する方法もあります：

```python
from tqdm import tqdm
import sys

# ファイルに出力
with open('progress.log', 'w') as f:
    for i in tqdm(range(100), file=f):
        pass

# stdoutに出力（あまり推奨されない）
for i in tqdm(range(100), file=sys.stdout):
    pass
```

### 3. 専用の通信チャネルを使用する

より高度な解決策として、WebSocketやZMQなどの双方向通信を使用する方法があります：

```python
# Pythonスクリプト側
import zmq
import time

context = zmq.Context()
socket = context.socket(zmq.PUSH)
socket.connect("tcp://127.0.0.1:5555")

total = 100
for i in range(total):
    time.sleep(0.1)
    # 進捗情報をJSONとして送信
    socket.send_json({
        "type": "progress",
        "current": i + 1,
        "total": total,
        "percent": (i + 1) / total * 100
    })
```

```javascript
// Electron側
const zmq = require('zeromq');

async function startServer() {
  const sock = new zmq.Pull();
  await sock.bind('tcp://127.0.0.1:5555');
  
  for await (const msg of sock) {
    const data = JSON.parse(msg.toString());
    if (data.type === "progress") {
      // UIの進捗バーを更新
      updateProgressUI(data.percent);
    }
  }
}

startServer();
```

## tqdmがstderrを使用する理由

tqdmライブラリは、デフォルトで**stderr**に出力します。これには重要な理由があります：

1. **標準出力の保護**: Pythonスクリプトの主な出力結果（データ処理結果など）は通常stdoutに送られます。tqdmがstdoutを使うと、進捗バーと実際の出力が混在して、パイプラインやログ記録が破壊される可能性があります。

2. **非ブロッキング動作**: stderrはバッファリングされないため、リアルタイムで進捗状況が更新されます。

3. **分離の原則**: エラーや補助情報（進捗状況など）はstderrに、実際の処理結果はstdoutに出力することで、論理的に出力を分離できます。

```python
from tqdm import tqdm
import time

# 標準的なtqdmの使用例
for i in tqdm(range(100)):
    time.sleep(0.1)
    # 処理結果はstdoutへ
    print(f"処理結果: {i}")
```

## Electronアプリでの問題点

Electronアプリから`child_process.spawn`を使用してPythonプロセスを起動する際、次のような問題が発生します：

1. **プロセス終了の誤検出**: 多くの場合、stderrに何かが出力されると「エラーが発生した」と解釈され、プロセスが強制終了されることがあります。

2. **エラー判別の複雑さ**: tqdmの進捗表示とアプリケーションの実際のエラーメッセージが同じstderrに混在するため、真のエラーを判別するのが難しくなります。

3. **UIへの反映**: stderrからの出力をUIに適切に反映させる方法が必要です。


## まとめ

Electronアプリ内でPythonのtqdmを使用する際は、以下の点に注意しましょう：

1. tqdmはデフォルトでstderrに出力するため、spawn処理でstderrを適切に処理する必要がある
2. stderrに出力があってもプロセスを強制終了しないよう設定する
3. tqdmの進捗表示と実際のエラーを区別する仕組みを実装する
4. より複雑なアプリケーションでは、専用の通信チャネルを検討する

これらの対策により、ElectronとPythonを組み合わせたアプリケーションでも、ユーザーフレンドリーな進捗表示を実現できます。

## 参考リンク

- [tqdmの使い方を自分なりにまとめてみた](https://qiita.com/kazuki__003/items/cb993f9cde9bc0e88123)
- [Why is tqdm output directed to sys.stderr and not to sys.stdout?](https://stackoverflow.com/questions/75580592/why-is-tqdm-output-directed-to-sys-stderr-and-not-to-sys-stdout)
- [Electron child_process API](https://www.electronjs.org/docs/latest/api/child-process)
- [ZeroMQを使ったプロセス間通信](https://zeromq.org/)
