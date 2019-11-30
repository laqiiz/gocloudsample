package main

import (
	"context"
	"fmt"
	"gocloud.dev/docstore"
	_ "gocloud.dev/docstore/gcpfirestore"
	"log"
)

type Entity struct {
	ID   string `docstore:"ID"`
	Name string `docstore:"NAME"`
}

func main() {
	ctx := context.Background()

	const url = "firestore://projects/dummy/databases/(default)/documents/example-collection?name_field=userID"
	coll, err := docstore.OpenCollection(ctx, url)
	if err != nil {
		log.Fatalf("open collection error: %v", err)
	}
	defer coll.Close()

	// 書き込み
	row := Entity{
		ID:   "1",
		Name: "hoge",
	}
	if err := coll.Create(ctx, &row); err != nil {
		log.Fatalf("create error: %v", err)
	}

	// 読み込み
	rowToRead := Entity{
		ID: "1",
	}
	if err := coll.Get(ctx, &rowToRead); err != nil {
		log.Fatalf("get error: %v", err)
	}
	fmt.Printf("get: %+v\n", rowToRead)
}
