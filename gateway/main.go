package main

import (
	"gateway/server"
)

func main() {
	server.StartNewHttpServer(8080, "localhost:9092")
}

// func main() {
// 	saramaSubscriberConfig := kafka.DefaultSaramaSubscriberConfig()
// 	// equivalent of auto.offset.reset: earliest
// 	saramaSubscriberConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

// 	subscriber, err := kafka.NewSubscriber(
// 		kafka.SubscriberConfig{
// 			Brokers:               []string{"localhost:9092"},
// 			Unmarshaler:           kafka.DefaultMarshaler{},
// 			OverwriteSaramaConfig: saramaSubscriberConfig,
// 			ConsumerGroup:         "test_consumer_group",
// 		},
// 		watermill.NewStdLogger(false, false),
// 	)
// 	if err != nil {
// 		panic(err)
// 	}

// 	messages, err := subscriber.Subscribe(context.Background(), "example.topic")
// 	if err != nil {
// 		panic(err)
// 	}
// 	go process(messages)

// 	//producer
// 	publisher, err := kafka.NewPublisher(
// 		kafka.PublisherConfig{
// 			Brokers:   []string{"localhost:9092"},
// 			Marshaler: kafka.DefaultMarshaler{},
// 		},
// 		watermill.NewStdLogger(false, false),
// 	)
// 	if err != nil {
// 		panic(err)
// 	}

// 	publishMessages(publisher)
// }
// func publishMessages(publisher message.Publisher) {
// 	for {
// 		msg := message.NewMessage(watermill.NewUUID(), []byte("Hello, world!"))

// 		if err := publisher.Publish("example.topic", msg); err != nil {
// 			panic(err)
// 		}

// 		time.Sleep(time.Second)
// 	}
// }
// func process(messages <-chan *message.Message) {
// 	for msg := range messages {
// 		log.Printf("received message: %s, payload: %s", msg.UUID, string(msg.Payload))

// 		// we need to Acknowledge that we received and processed the message,
// 		// otherwise, it will be resent over and over again.
// 		msg.Ack()
// 	}
// }
