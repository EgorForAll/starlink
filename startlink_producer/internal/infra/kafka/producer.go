package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"

)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(brokers...),
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: true,
		},
	}
}

// Publish отправляет сообщение в топик, равный event_type (например "user.created").
func (p *Producer) Publish(ctx context.Context, topic string, payload []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Value: payload,
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
