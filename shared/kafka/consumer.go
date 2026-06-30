package kafka

import (
	"context"

	"github.com/razedwell/traxex/shared/config"
	kafkago "github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafkago.Reader
}

func NewConsumer(cfg config.KafkaConfig, topic string) *Consumer {
	return &Consumer{
		reader: kafkago.NewReader(kafkago.ReaderConfig{
			Brokers: cfg.Brokers,
			Topic:   topic,
			GroupID: cfg.GroupID,
		}),
	}
}

func (c *Consumer) Read(ctx context.Context) (kafkago.Message, error) {
	return c.reader.ReadMessage(ctx)
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
