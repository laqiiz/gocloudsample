package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"gocloud.dev/blob"
	"gocloud.dev/blob/s3blob"
	_ "gocloud.dev/blob/s3blob"
	"log"
	"net/url"
	"os"
)

func main() {
	s, err := readAll("s3://future-example", "test.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("read: %v", s)
}

func readAll(bucketURL, key string) (string, error) {
	bucket, err := openBucket(bucketURL)
	if err != nil {
		return "", err
	}

	b, err := bucket.ReadAll(context.Background(), key)
	if err != nil {
		return "", err
	}
	return string(b), err
}

func openBucket(bucketURL string) (*blob.Bucket, error) {
	endpoint := os.Getenv("ENDPOINT_URL") // ★環境変数化する
	if len(endpoint) == 0 {
		return blob.OpenBucket(context.Background(), bucketURL)
	}

	u, err := url.Parse(bucketURL)
	if err != nil {
		return nil, err
	}

	fmt.Println(u.Path)

	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(true),
	}))
	return s3blob.OpenBucket(context.Background(), sess, u.Host, nil) // ★hostnameに切り替える
}
