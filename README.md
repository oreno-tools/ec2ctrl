# ec2ctrl

## これなに

EC2 一覧の取得, 起動, 停止を行うワンバイナリツールでごわす.

## 使い方

### インストゥール

1. お使いの OS に応じて, osx, linux, win ディレクトリから ec2ctrl のバイナリをダウンロードする
2. ec2ctrl をパスの通ったディレクトリにコピー又は移動する

### ヘルプ

```sh
$ ec2ctrl -help
Usage of ec2ctrl:
  -csv
        CSV 形式で出力する
  -endpoint string
        AWS API のエンドポイントを指定.
  -instances string
        Instance ID 又は Instance Tag 名を指定.
  -json
        JSON 形式で出力する
  -profile string
        Profile 名を指定.
  -region string
        Region 名を指定. (default "ap-northeast-1")
  -start
        Instance を起動.
  -state
        Instance の状態を出力.
  -stop
        Instance を停止.
  -tags string
        Tag Key 及び Tag Value を指定.
  -version
        バージョンを出力.
```

### EC2 一覧の取得

```sh
$ ec2Ctrl -profile=your-profile # -profile を指定しない場合, 環境変数又は default のプロファイルを読み込む
```

`-csv` オプションを付与した場合, カンマ区切りで出力する.

```sh
$ ec2Ctrl -csv
```

`-json` オプションを付与した場合, JSON フォーマットで出力する.

```sh
$ ec2Ctrl -json
```

[jq](https://stedolan.github.io/jq/) 等との併用で対象の絞り込みも可能.

```sh
$ ec2ctrl -json | jq '.instances[]|select(.name == "your-instance-tag-name")'
```

### EC2 の起動

```sh
$ ec2Ctrl -start -instances=foo,bar
上記のインスタンスを操作しますか?(y/n): y
EC2 を起動します.
i-xxxxxxxxxxxxxxxxx を起動しました.
i-yyyyyyyyyyyyyyyyy を起動しました.
```

### EC2 の停止

```sh
$ ec2Ctrl -stop -instances=foo,bar
上記のインスタンスを操作しますか?(y/n): y
EC2 を停止します.
i-xxxxxxxxxxxxxxxxx を停止しました.
i-yyyyyyyyyyyyyyyyy を停止しました.
```

## Todo

* 力技で書いているような部分をエレガントに書き換える
* テストを書く