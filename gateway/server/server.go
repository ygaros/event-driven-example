package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
)

type HttpServer interface {
	HandleNotification(w http.ResponseWriter, r *http.Request)
}
type httpServer struct {
	publisher *kafka.Publisher
}
type NotificationRequest struct {
	Message string `json:"message"`
	Sender  string `json:"sender"`
}

func (s *httpServer) HandleNotification(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	defer r.Body.Close()
	request := NotificationRequest{}
	if err := json.Unmarshal(body, &request); err != nil {
		log.Println("Failed to parse json data:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	msg := message.NewMessage(watermill.NewUUID(), body)

	if err := s.publisher.Publish("notification.topic", msg); err != nil {
		panic(err)
	}
}
func StartNewHttpServer(port int, kafkaUrl string) {
	if kafkaUrl == "" {
		kafkaUrl = "localhost:9092"
	}
	publisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   []string{kafkaUrl},
			Marshaler: kafka.DefaultMarshaler{},
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		panic(err)
	}
	server := httpServer{
		publisher: publisher,
	}
	http.HandleFunc("/", server.HandleNotification)

	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
