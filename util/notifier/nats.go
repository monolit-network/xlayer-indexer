package notifier

import (
	"context"
	"fmt"
	"os"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type NATSTransport struct {
	conn *nats.Conn
}

func NewNATSFromEnv(logger *zap.Logger, topics ...Topic) (*Notifier, error) {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}

	conn, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to nats at %s: %w", url, err)
	}

	logger.Info("nats transport initialized successfully", zap.String("url", url))

	return NewNotifier(logger, &NATSTransport{conn: conn}, topics...)
}

func (b *NATSTransport) Publish(_ context.Context, topic string, data []byte) error {
	return b.conn.Publish(topic, data)
}

func (b *NATSTransport) Subscribe(_ context.Context, topic string, handler func(data []byte)) (func() error, error) {
	subscription, err := b.conn.Subscribe(topic, func(msg *nats.Msg) {
		handler(msg.Data)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to nats topic %s: %w", topic, err)
	}

	if err := b.conn.Flush(); err != nil {
		_ = subscription.Unsubscribe()
		return nil, fmt.Errorf("failed to confirm nats subscription for topic %s: %w", topic, err)
	}

	return func() error {
		return subscription.Unsubscribe()
	}, nil
}

func (b *NATSTransport) Close() error {
	if b == nil || b.conn == nil {
		return nil
	}

	b.conn.Close()
	return nil
}
