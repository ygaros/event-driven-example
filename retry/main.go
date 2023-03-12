package main

import (
	"context"
	"log"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
)

var cache map[string]int

var publisher *kafka.Publisher
var lock sync.RWMutex

const (
	kafkaUrl   string = "localhost:9092"
	topic      string = "notification.topic"
	topicRetry string = "notification.topic.retry"
)

func main() {

	cache = make(map[string]int)
	lock = sync.RWMutex{}

	saramaSubscriberConfig := kafka.DefaultSaramaSubscriberConfig()
	// equivalent of auto.offset.reset: earliest
	saramaSubscriberConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	subscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:               []string{kafkaUrl},
			Unmarshaler:           kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: saramaSubscriberConfig,
			ConsumerGroup:         "test_consumer_group",
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		panic(err)
	}

	publisher, err = kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   []string{kafkaUrl},
			Marshaler: kafka.DefaultMarshaler{},
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		panic(err)
	}

	messages, err := subscriber.Subscribe(context.Background(), topicRetry)
	if err != nil {
		panic(err)
	}
	process(messages)
}
func process(messages <-chan *message.Message) {
	for msg := range messages {
		log.Printf("received message for retry processing: %s, payload: %s", msg.UUID, string(msg.Payload))
		lock.Lock()
		_, ok := cache[msg.UUID]
		if ok {
			cache[msg.UUID] += 1
		} else {
			cache[msg.UUID] = 1
		}
		if cache[msg.UUID] > 3 {
			log.Printf("message %s processed 3 times without any success abandoning this message", msg.UUID)
			delete(cache, msg.UUID)
		} else {
			// msg := message.NewMessage(msg.UUID, msg.Payload)

			if err := publisher.Publish(topic, msg); err != nil {
				panic(err)
			}
			log.Printf("publishing messege %d time on notification.topic", cache[msg.UUID])
		}
		lock.Unlock()
		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
	}
}
