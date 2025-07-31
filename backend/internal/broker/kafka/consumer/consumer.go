package consumer

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	r *kafka.Reader
}

func New(topic string, groupId string, addr ...string) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: addr,
		GroupID: groupId,
		Topic:   topic,
	})
	return &Consumer{r: r}
}

func (c *Consumer) Consume(ctx context.Context) (kafka.Message, error) {
	return c.r.ReadMessage(ctx)
}

func (c *Consumer) Close() error {
	return c.r.Close()
}
