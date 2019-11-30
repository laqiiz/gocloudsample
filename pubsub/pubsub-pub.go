package main

import (
	"context"
	"fmt"
	"gocloud.dev/gcp"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/gcppubsub"
	_ "gocloud.dev/pubsub/gcppubsub"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"log"
	"os"
)

var endPoint = "pubsub.googleapis.com:443"

func main() {
	ctx := context.Background()

	conn, err := dial(ctx)
	if err != nil{
		log.Fatalf("dial error: %v", err)
	}
	defer conn.Close()

	pubClient, err := gcppubsub.PublisherClient(ctx, conn)
	if err != nil {
		log.Fatalf("publisher client error: %v", err)
	}
	defer pubClient.Close()

	topic, err := gcppubsub.OpenTopicByPath(pubClient, "projects/dummy/topics/future-example", nil)
	if err != nil {
		log.Fatalf("open topic error: %v", err)
	}
	defer topic.Shutdown(ctx)

	if err = topic.Send(ctx, &pubsub.Message{Body: []byte("Hello, World!")}); err != nil {
		log.Fatalf("publish error: %v", err)
	}
	fmt.Println("done")

	//TODO subscribe
	subClient, err := gcppubsub.SubscriberClient(ctx, conn)
	if err != nil {
		log.Fatalf("subscriber client error: %v", err)
	}
	defer subClient.Close()

	sub, err := gcppubsub.OpenSubscriptionByPath(subClient, "projects/dummy/subscriptions/gocdk-example1", nil)
	if err != nil {
		log.Fatalf("subscription error: %v", err)
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

func dial(ctx context.Context) (*grpc.ClientConn, error) {
	emulatorEndPoint := os.Getenv("PUBSUB_EMULATOR_HOST")
	if emulatorEndPoint != "" {
		endPoint = emulatorEndPoint
		return grpc.DialContext(ctx, endPoint,
			grpc.WithInsecure(), // ★追加
			grpc.WithUserAgent(fmt.Sprintf("%s/%s/%s", "pubsub", "go-cloud", "0.1.0")),
		)
	}

	creds, err := gcp.DefaultCredentials(ctx)
	if err != nil {
		return nil, err
	}

	return grpc.DialContext(ctx, endPoint,
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")),
		grpc.WithPerRPCCredentials(oauth.TokenSource{TokenSource: creds.TokenSource}),
		grpc.WithUserAgent(fmt.Sprintf("%s/%s/%s", "pubsub", "go-cdk", "0.1.0")),
}
