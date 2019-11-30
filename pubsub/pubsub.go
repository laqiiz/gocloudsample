package main

import (
	"context"
	"fmt"
	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/gcppubsub"
	"log"
)

func main() {
	ctx := context.Background()

	topicURL := "gcppubsub://projects/dummy/topics/future-example"
	topic, err := pubsub.OpenTopic(ctx, topicURL)
	if err != nil {
		log.Fatal(err)
	}
	defer topic.Shutdown(ctx)

	fmt.Println("open topic")
	if err = topic.Send(ctx, &pubsub.Message{Body: []byte("Hello, World!")}); err != nil {
		log.Fatal(err)
	}
	fmt.Println("send")

	sub, err := pubsub.OpenSubscription(ctx, topicURL)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Shutdown(ctx)

	for {
		msg, err := sub.Receive(ctx)
		if err != nil {
			log.Fatalf("Receiving message: %v", err)
		}
		fmt.Printf("Got message: %q\n", msg.Body)
		msg.Ack()
	}

}
