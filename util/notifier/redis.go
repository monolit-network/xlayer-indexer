package notifier

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisTransport struct {
	redisClient *redis.Client
}

func NewRedisFromEnv(logger *zap.Logger, topics ...Topic) (*Notifier, error) {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		return nil, fmt.Errorf("REDIS_HOST is not set")
	}
	port := os.Getenv("REDIS_PORT")
	if port == "" {
		return nil, fmt.Errorf("REDIS_PORT is not set")
	}
	password := os.Getenv("REDIS_PASSWORD")

	addr := fmt.Sprintf("%s:%s", host, port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to ping redis at %s: %w", addr, err)
	}

	return NewNotifier(logger, &RedisTransport{redisClient: client}, topics...)
}

func (b *RedisTransport) Publish(ctx context.Context, topic string, data []byte) error {
	_, err := b.redisClient.Publish(ctx, topic, data).Result()
	if err != nil {
		return err
	}

	return nil
}

func (b *RedisTransport) Subscribe(_ context.Context, topic string, handler func(data []byte)) (func() error, error) {
	pubsub := b.redisClient.Subscribe(context.Background(), topic)

	if _, err := pubsub.Receive(context.Background()); err != nil {
		_ = pubsub.Close()
		return nil, fmt.Errorf("redis subscription receive failed for topic %s: %w", topic, err)
	}

	go func() {
		channel := pubsub.Channel()
		for msg := range channel {
			handler([]byte(msg.Payload))
		}
	}()

	return func() error {
		return pubsub.Close()
	}, nil
}

func (b *RedisTransport) Close() error {
	if b == nil || b.redisClient == nil {
		return nil
	}
	return b.redisClient.Close()
}
