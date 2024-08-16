package application

import (
	"context"
	"fmt"
	"log"

	"github.com/milicns/company-manager/company-service/internal/utils"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewProducer(addr string) (*KafkaProducer, func()) {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(addr),
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	return &KafkaProducer{writer: writer}, func() {
		if err := writer.Close(); err != nil {
			log.Fatal("failed to close kafka writer:", err)
		}
		log.Println("closing kafka writer")
	}
}

func (producer *KafkaProducer) Produce(event utils.Event) {
	err := producer.writer.WriteMessages(context.Background(), kafka.Message{
		Topic: fmt.Sprintf("topic-%s", event.Topic),
		Key:   []byte(event.EventType),
		Value: []byte(event.Value),
	})
	if err != nil {
		log.Println("failed to write message:", err)
	}
}
