package kafka

import (
	"context"

	"github.com/razedwell/traxex/shared/config"
	kafkago "github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafkago.Writer
}

func NewProducer(cfg config.KafkaConfig, topic string) *Producer {
	return &Producer{
		writer: &kafkago.Writer{
			Addr:                   kafkago.TCP(cfg.Brokers...),
			Topic:                  topic,
			Balancer:               &kafkago.LeastBytes{},
			AllowAutoTopicCreation: true,
		},
	}
}

func (p *Producer) Publish(ctx context.Context, key, value []byte) error {
	return p.writer.WriteMessages(ctx, kafkago.Message{
		Key:   key,
		Value: value,
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
