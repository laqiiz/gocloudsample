package main

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"gocloud.dev/docstore"
	"gocloud.dev/docstore/awsdynamodb"
	"gocloud.dev/gcerrors"
	"log"
	"time"
)

type Item struct {
	MyHashKey        string `docstore:"MyHashKey"`
	MyRangeKey       int    `docstore:"MyRangeKey"`
	MyText           string `docstore:"MyText"`
	DocstoreRevision interface{}
}

func main() {
	// 別のgoroutineでも無限書き込み
	go UpdateLoop()

	// メインスレッドでも無限書き込み
	UpdateLoop()
}

func UpdateLoop() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	db := dynamodb.New(sess, &aws.Config{Endpoint: aws.String("http://localhost:8000")})

	coll, err := awsdynamodb.OpenCollection(db, "MyFirstTable", "MyHashKey", "MyRangeKey", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer coll.Close()

	ctx := context.Background()

	for {
		read := Item{MyHashKey: "00001", MyRangeKey: 1}
		if err := coll.Get(ctx, &read); err != nil {
			log.Fatalf("get: %v", err)
		}
		if err := coll.Update(ctx, &read, docstore.Mods{"MyText": "update text: " + time.Now().String()}); err != nil {
			if gcerrors.Code(err) == gcerrors.FailedPrecondition {
				log.Fatalf("optimistic locking: %v", err)
			}
			log.Fatalf("update: %v", err)
		}
	}
}