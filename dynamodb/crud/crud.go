package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"gocloud.dev/docstore"
	"gocloud.dev/docstore/awsdynamodb"
	"log"
)

type Item struct {
	MyHashKey  string `docstore:"MyHashKey"`
	MyRangeKey int    `docstore:"MyRangeKey"`
	MyText     string `docstore:"MyText"`
}

func main() {

	// Create session.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	db := dynamodb.New(sess, &aws.Config{Endpoint: aws.String("http://localhost:8000")})

	// Open Collection via Go CDK
	coll, err := awsdynamodb.OpenCollection(db, "MyFirstTable", "MyHashKey", "MyRangeKey", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer coll.Close()

	ctx := context.Background()

	// Create
	//write := Item{MyHashKey: "00001", MyRangeKey: 1, MyText: "some text..."}
	//if err := coll.Create(ctx, &write); err != nil {
	//	log.Fatalf("create: %v", err)
	//}

	// Read
	read := Item{MyHashKey: "00001", MyRangeKey: 1}

	if err := coll.Get(ctx, &read); err != nil {
		log.Fatalf("get: %v", err)
	}
	fmt.Printf("got: %+v\n", read)

	// Update
	updateKey := Item{MyHashKey: "00001", MyRangeKey: 1}
	if err := coll.Update(ctx, &updateKey, docstore.Mods{"MyText": "update text"}); err != nil {
		log.Fatalf("update: %v", err)
	}
	if err := coll.Get(ctx, &read); err != nil {
		log.Fatalf("get: %v", err)
	}
	fmt.Printf("got: %+v\n", read)

	// Replace
	replace := Item{MyHashKey: "00001", MyRangeKey: 1, MyText: "replace"}
	if err := coll.Replace(ctx, &replace); err != nil {
		log.Fatalf("replace: %v", err)
	}
	if err := coll.Get(ctx, &read); err != nil {
		log.Fatalf("get: %v", err)
	}
	fmt.Printf("got: %+v\n", read)


	//// Update: 存在しないキー
	//notFoundKey := Item{MyHashKey: "99999", MyRangeKey: 1}
	//if err := coll.Update(ctx, &notFoundKey, docstore.Mods{"MyText": "update text"}); err != nil {
	//	log.Fatalf("not found: %v", err)
	//}

	// Delete
	deleteKey := Item{MyHashKey: "00001", MyRangeKey: 1}
	if err := coll.Delete(ctx, &deleteKey); err != nil {
		log.Fatalf("delete: %v", err)
	}


}
