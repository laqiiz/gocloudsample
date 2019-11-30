package main

import (
	vkit "cloud.google.com/go/firestore/apiv1"
	"context"
	"fmt"
	"gocloud.dev/docstore/gcpfirestore"
	"gocloud.dev/gcp"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"log"
	"os"
)

type Record struct {
	ID   string `docstore:"id"`
	Name string `docstore:"name"`
}

var endPoint = "firesotore.googleapis.com:443"

func main() {
	ctx := context.Background()

	client, err := dial(ctx)
	if err != nil {
		log.Fatalf("dial error: %v", err)
	}

	resourceID := gcpfirestore.CollectionResourceID("dummy", "example-collection")
	coll, err := gcpfirestore.OpenCollection(client, resourceID, "id", nil)
	if err != nil {
		log.Fatalf("open collection error: %v", err)
	}
	defer coll.Close()

	// 書き込み
	row := Record{
		ID:   "1",
		Name: "hoge",
	}
	if err := coll.Create(ctx, &row); err != nil {
		log.Fatalf("create error: %v", err)
	}

	// 読み込み
	rowToRead := Record{
		ID: "1",
	}

	if err := coll.Get(ctx, &rowToRead); err != nil {
		log.Fatalf("get error: %v", err)
	}
	fmt.Printf("get: %+v\n", rowToRead)
}

func dial(ctx context.Context) (*vkit.Client, error) {
	emulatorEndPoint := os.Getenv("FIRESTORE_EMULATOR_HOST")
	if emulatorEndPoint != "" {
		endPoint = emulatorEndPoint

		conn, err := grpc.DialContext(ctx, endPoint, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		return vkit.NewClient(ctx,
			option.WithEndpoint(endPoint),
			option.WithGRPCConn(conn),
			option.WithUserAgent(fmt.Sprintf("%s/%s/%s", "firestore", "go-cloud", "0.1.0")),
		)
	}

	creds, err := gcp.DefaultCredentials(ctx)
	if err != nil {
		return nil, err
	}

	client, _, err := gcpfirestore.Dial(ctx, creds.TokenSource)
	return client, err
}
