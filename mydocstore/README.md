# MyDocstore

## DynamoDBの起動

$ docker-compose up -d
$ aws dynamodb create-table --endpoint-url http://localhost:8000 --cli-input-json file://my_table.json

