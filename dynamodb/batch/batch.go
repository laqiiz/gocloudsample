package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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
	w1 := Item{MyHashKey: "00001", MyRangeKey: 1, MyText: "some text1..."}
	w2 := Item{MyHashKey: "00001", MyRangeKey: 2, MyText: "some text2..."}
	w3 := Item{MyHashKey: "00001", MyRangeKey: 3, MyText: "some text3..."}
	if err := coll.Actions().Create(&w1).Create(&w2).Create(&w3).Do(ctx); err != nil {
		log.Fatalf("actions: %v", err)
	}

	items := []Item{
		{MyHashKey: "00001", MyRangeKey: 1},
		{MyHashKey: "00001", MyRangeKey: 2},
		{MyHashKey: "00001", MyRangeKey: 3},
	}
	for _, read := range items {
		if err := coll.Get(ctx, &read); err != nil {
			log.Fatalf("get: %v", err)
		}
		fmt.Printf("got: %+v\n", read)
	}

}
