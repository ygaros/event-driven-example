package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

const (
	kafkaUrl   string = "localhost:9092"
	topic      string = "notification.topic"
	topicRetry string = "notification.topic.retry"
)

//using watermill router API

func main() {
	logger := watermill.NewStdLogger(false, false)
	subClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	sub, err := redisstream.NewSubscriber(
		redisstream.SubscriberConfig{
			Client:        subClient,
			ConsumerGroup: "test_consumer_group",
		},
		logger,
	)
	if err != nil {
		panic(err)
	}
	messages, err := sub.Subscribe(context.Background(), topic)
	if err != nil {
		panic(err)
	}
	process(messages)
}
func printMessages(msg *message.Message) error {
	fmt.Printf(
		"\n> Received message: %s\n> %s\n> metadata: %v\n\n",
		msg.UUID, string(msg.Payload), msg.Metadata,
	)
	// return fmt.Errorf("error occured")
	return nil
}

func process(messages <-chan *message.Message) {
	counter := 2
	for msg := range messages {

		// counter += 1
		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.

		if counter%3 == 0 {
			log.Printf("FAILED processing message: %s, payload: %s", msg.UUID, string(msg.Payload))
			// if err := publisher.Publish(topicRetry, msg); err != nil {
			// 	panic(err)
			// }
		} else {
			log.Printf("SUCCESS received message: %s, payload: %s", msg.UUID, string(msg.Payload))
		}
		msg.Ack()
	}
}
