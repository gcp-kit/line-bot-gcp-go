# line-bot-gcp-go
GCPine(GC🍍)はローカル開発とクラウド運用が共存できるツールキットです。

## 概要
Google Cloud Platform上で動作するLINE Botを運用することができます。

## 特徴
GCPineはローカル開発時にはWebHookを受けて動作し  
クラウド開発によるデプロイ待ちなどを気にせずに開発効率を維持できるようにしました。

## 必須

* go1.13以上
* 外部依存モジュールはgo.modから取得

## インストール
```shell script
go get github.com/gcp-kit/line-bot-gcp-go
```

## Usage
Google App Engine用のサンプル  
[gcp-kit/gcpine-gae-example](https://github.com/gcp-kit/gcpine-gae-example)

TODO
- Google Cloud Functions用のサンプル
- Local用のサンプル

## License
[MIT license](../LICENSE).
