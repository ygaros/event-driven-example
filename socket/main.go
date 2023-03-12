package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
)

const (
	kafkaUrl   string = "localhost:9092"
	topic      string = "notification.topic"
	topicRetry string = "notification.topic.retry"
)

//using watermill router API

func main() {
	logger := watermill.NewStdLogger(false, false)
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}
	router.AddPlugin(plugin.SignalsHandler)

	router.AddMiddleware(
		middleware.CorrelationID,

		middleware.Retry{
			MaxRetries:      3,
			InitialInterval: time.Millisecond * 100,
			Logger:          logger,
		}.Middleware,
		middleware.Recoverer,
	)

	saramaSubscriberConfig := kafka.DefaultSaramaSubscriberConfig()
	// equivalent of auto.offset.reset: earliest
	saramaSubscriberConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	sub, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:               []string{kafkaUrl},
			Unmarshaler:           kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: saramaSubscriberConfig,
			ConsumerGroup:         "test_consumer_group",
		},
		logger,
	)
	if err != nil {
		panic(err)
	}
	// pub, err := kafka.NewPublisher(
	// 	kafka.PublisherConfig{
	// 		Brokers:   []string{kafkaUrl},
	// 		Marshaler: kafka.DefaultMarshaler{},
	// 	},
	// 	logger,
	// )
	if err != nil {
		panic(err)
	}
	// router.AddHandler(
	// 	"process_messages",
	// 	topic,
	// 	sub,
	// 	topicRetry,
	// 	pub,
	// 	structHandler{}.Handler,
	// )
	router.AddNoPublisherHandler(
		"print_incoming_messages",
		topic,
		sub,
		printMessages,
	)
	// router.AddNoPublisherHandler(
	// 	"print_incoming_messages",
	// 	topicRetry,
	// 	sub,
	// 	printMessages,
	// )
	messages, err := sub.Subscribe(context.Background(), topic)
	if err != nil {
		panic(err)
	}
	process(messages)
	// ctx := context.Background()
	// if err := router.Run(ctx); err != nil {
	// 	panic(err)
	// }
}
func printMessages(msg *message.Message) error {
	fmt.Printf(
		"\n> Received message: %s\n> %s\n> metadata: %v\n\n",
		msg.UUID, string(msg.Payload), msg.Metadata,
	)
	// return fmt.Errorf("error occured")
	return nil
}

// type structHandler struct {
// 	// we can add some dependencies here
// 	counter int
// }

// func (s structHandler) Handler(msg *message.Message) ([]*message.Message, error) {
// 	log.Println("structHandler received message", msg.UUID)
// 	s.counter += 1
// 	if s.counter%3 == 0 {
// 		log.Printf("FAILED processing message: %s, payload: %s", msg.UUID, string(msg.Payload))
// 		return nil, fmt.Errorf("FAILED processing message: %s, payload: %s", msg.UUID, string(msg.Payload))
// 	}
// 	log.Printf("SUCCESS processing message: %s, payload: %s", msg.UUID, string(msg.Payload))
// 	msg = message.NewMessage(watermill.NewUUID(), []byte("message produced by structHandler"))
// 	return message.Messages{msg}, nil

// }

// implementation of low-level PUB SUB watermill interfaces
// var publisher *kafka.Publisher

// func main() {

// 	saramaSubscriberConfig := kafka.DefaultSaramaSubscriberConfig()
// 	// equivalent of auto.offset.reset: earliest
// 	saramaSubscriberConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

// 	subscriber, err := kafka.NewSubscriber(
// 		kafka.SubscriberConfig{
// 			Brokers:               []string{kafkaUrl},
// 			Unmarshaler:           kafka.DefaultMarshaler{},
// 			OverwriteSaramaConfig: saramaSubscriberConfig,
// 			ConsumerGroup:         "test_consumer_group",
// 		},
// 		watermill.NewStdLogger(false, false),
// 	)
// 	if err != nil {
// 		panic(err)
// 	}
// 	publisher, err = kafka.NewPublisher(
// 		kafka.PublisherConfig{
// 			Brokers:   []string{kafkaUrl},
// 			Marshaler: kafka.DefaultMarshaler{},
// 		},
// 		watermill.NewStdLogger(false, false),
// 	)
// 	if err != nil {
// 		panic(err)
// 	}
// 	messages, err := subscriber.Subscribe(context.Background(), topic)
// 	if err != nil {
// 		panic(err)
// 	}
// 	process(messages)

// }
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
