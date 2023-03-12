package main

import (
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
)

func main() {
	publisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   []string{"localhost:9092"},
			Marshaler: kafka.DefaultMarshaler{},
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		panic(err)
	}

	publishMessages(publisher)
}
func publishMessages(publisher message.Publisher) {
	for {
		msg := message.NewMessage(watermill.NewUUID(), []byte("Hello, world!"))

		if err := publisher.Publish("notification.topic", msg); err != nil {
			panic(err)
		}

		time.Sleep(time.Millisecond * 10)
	}
}
