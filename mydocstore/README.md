# MyDocstore

## DynamoDBの起動

```sh
$ docker-compose up -d

# AWS CLIのセットアップ
$ aws configure set aws_access_key_id dummy     --profile local
$ aws configure set aws_secret_access_key dummy --profile local
$ aws configure set region ap-northeast-1       --profile local

# テーブル作成
$ aws dynamodb create-table --endpoint-url http://localhost:8000 --cli-input-json file://my_table.json
```

