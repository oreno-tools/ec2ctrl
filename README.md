# ec2ctrl

## これなに

EC2 一覧の取得, 起動, 停止を行うワンバイナリツールでごわす. [jq](https://stedolan.github.io/jq/) や [direnv](https://github.com/direnv/direnv) 等と併用すると貴方の EC2 ライフを少しだけ豊かなものにしてくれるでしょう.

## 使い方

### インストゥール

https://github.com/inokappa/ec2ctrl/releases から環境に応じたバイナリをダウンロードしてください.

```
wget https://github.com/inokappa/ec2ctrl/releases/download/v0.0.7/ec2ctrl_darwin_amd64 -O ~/bin/ec2ctrl
chmod +x ~/bin/ec2ctrl
```

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
  -stop
        Instance を停止.
  -tags string
        Tag Key 及び Tag Value を指定.
  -version
        バージョンを出力.
```

### EC2 一覧の取得

```sh
$ ec2ctrl -profile=your-profile # -profile を指定しない場合, 環境変数又は default のプロファイルを読み込む
```

`-csv` オプションを付与した場合, カンマ区切りで出力する.

```sh
$ ec2ctrl -csv
```

`-json` オプションを付与した場合, JSON フォーマットで出力する.

```sh
$ ec2ctrl -json
```

[jq](https://stedolan.github.io/jq/) 等との併用で対象の絞り込みも可能.

```sh
$ ec2ctrl -json | jq '.instances[]|select(.name == "your-instance-tag-name")'
```

### EC2 状態確認

```sh
$ ec2ctrl -instances=foo,bar
```

`-instances` は EC2 インスタンス ID 又はインスタンスに付与されている `Name` タグを指定可能.

### EC2 の起動

```sh
$ ec2ctrl -start -instances=foo,bar
上記のインスタンスを操作しますか?(y/n): y
EC2 を起動します.
i-xxxxxxxxxxxxxxxxx を起動しました.
i-yyyyyyyyyyyyyyyyy を起動しました.
```

`-instances` は EC2 インスタンス ID 又はインスタンスに付与されている `Name` タグを指定可能.

### EC2 の停止

```sh
$ ec2ctrl -stop -instances=foo,bar
上記のインスタンスを操作しますか?(y/n): y
EC2 を停止します.
i-xxxxxxxxxxxxxxxxx を停止しました.
i-yyyyyyyyyyyyyyyyy を停止しました.
```

`-instances` は EC2 インスタンス ID 又はインスタンスに付与されている `Name` タグを指定可能.

## Todo

* 力技で書いているような部分をエレガントに書き換える
* テストを書く