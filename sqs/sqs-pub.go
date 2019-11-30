package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/awssnssqs"
	"log"
)

func main() {
	ctx := context.Background()

	sess, err := session.NewSession(&aws.Config{
		Endpoint: aws.String("http://localhost:4576"),
	})
	if err != nil {
		log.Fatal(err)
	}

	topic := awssnssqs.OpenSQSTopic(ctx, sess, "http://localhost:4576/queue/test-queue", nil)
	defer topic.Shutdown(ctx)

	err = topic.Send(ctx, &pubsub.Message{
		Body:     []byte("Hello, World!\n"),
		Metadata: map[string]string{"Env": "test"},
	})
	if err != nil {
		log.Fatal(err)
	}

	sub := awssnssqs.OpenSubscription(ctx, sess, "http://localhost:4576/queue/test-queue", nil)
	defer sub.Shutdown(ctx)

	for {
		msg, err := sub.Receive(ctx)
		if err != nil {
			log.Printf("Receiving message: %v", err)
			break
		}
		fmt.Printf("Got message: %q\n", msg.Body)
		msg.Ack()
	}

}
