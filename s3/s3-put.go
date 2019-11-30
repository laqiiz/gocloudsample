package main

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"gocloud.dev/blob/s3blob"
	_ "gocloud.dev/blob/s3blob"
	"io"
	"log"
)

func main() {
	w, err := writer("s3://future-example", "test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	if _, err := w.Write([]byte("1234567890")); err != nil {
		log.Fatal(err)
	}
}

func writer(bucketURL, key string) (io.WriteCloser, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint:         aws.String("http://localhost:4572"),
		S3ForcePathStyle: aws.Bool(true),
	}))
	bucket, err := s3blob.OpenBucket(context.Background(), sess, bucketURL, nil)
	if err != nil {
		return nil, err
	}

	return bucket.NewWriter(context.Background(), key, nil)
}
