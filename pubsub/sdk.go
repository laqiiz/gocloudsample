package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"log"
)

func main() {
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, "dummy")
	if err != nil {
		log.Fatal(err)
	}

	topic := client.Topic("future-example")
	if _, err := topic.Publish(ctx, &pubsub.Message{Data: []byte("12345")}).Get(ctx); err != nil {
		log.Fatalf("publish error: %v", err)
	}

	sub, err := client.CreateSubscription(context.Background(), "sub-name",
		pubsub.SubscriptionConfig{Topic: topic})

	err = sub.Receive(context.Background(), func(ctx context.Context, m *pubsub.Message) {
		log.Printf("Got message: %s", m.Data)
		m.Ack()
	})
	if err != nil {
		log.Fatalf("subscribe error: %v", err)
	}

	fmt.Println("done")
}
