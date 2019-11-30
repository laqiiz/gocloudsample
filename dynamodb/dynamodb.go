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

type Entity struct {
	ID   string `docstore:"ID"`
	Name string `docstore:"NAME"`
}

func main() {
	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint:         aws.String("http://localhost:4569"),
	}))
	coll, err := awsdynamodb.OpenCollection(dynamodb.New(sess), "test", "ID", "", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer coll.Close()

	// 書き込み
	row := Entity{
		ID:   "1",
		Name: "hoge",
	}
	coll.Create(context.Background(), &row)

	// 読み込み
	rowToRead := Entity{
		ID: "1",
	}
	coll.Get(context.Background(), &rowToRead)
	fmt.Printf("get: %+v\n", rowToRead)
}
