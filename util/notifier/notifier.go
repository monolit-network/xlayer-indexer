package notifier

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Notifier struct {
	Logger *zap.Logger

	transport Transport
	topics    []Topic
	active    map[string]*topicState
	mu        sync.RWMutex
}

type Transport interface {
	Publish(ctx context.Context, topic string, data []byte) error
	Subscribe(ctx context.Context, topic string, handler func(data []byte)) (func() error, error)
	Close() error
}

type Message[T any] struct {
	Value T
	Raw   []byte
}

type typedSubscriber interface {
	deliver(value any, raw []byte)
	close()
}

type subscriber[T any] struct {
	mu     sync.RWMutex
	ch     chan Message[T]
	closed bool
}

func newSubscriber[T any](buffer int) *subscriber[T] {
	return &subscriber[T]{
		ch: make(chan Message[T], buffer),
	}
}

func (s *subscriber[T]) deliver(value any, raw []byte) {
	typedValue, ok := value.(T)
	if !ok {
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return
	}

	select {
	case s.ch <- Message[T]{
		Value: typedValue,
		Raw:   raw,
	}:
	default:
	}
}

func (s *subscriber[T]) close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return
	}

	s.closed = true
	close(s.ch)
}

func NewNotifier(logger *zap.Logger, transport Transport, topics ...Topic) (*Notifier, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	if transport == nil {
		return nil, fmt.Errorf("notifier transport is required")
	}

	registry := make([]Topic, 0, len(topics))
	for _, topic := range topics {
		registry = append(registry, topic)
	}

	return &Notifier{
		Logger:    logger,
		transport: transport,
		topics:    registry,
		active:    make(map[string]*topicState),
	}, nil
}

const defaultSubscriberBuffer = 100

func Subscribe[T any](n *Notifier, ctx context.Context, topic string) (<-chan Message[T], string, error) {
	if n == nil {
		return nil, "", fmt.Errorf("notifier is nil")
	}
	if _, err := n.getTopic(topic, typeOf[T]()); err != nil {
		return nil, "", err
	}

	sessionID := uuid.New().String()
	sub := newSubscriber[T](defaultSubscriberBuffer)

	n.mu.Lock()
	defer n.mu.Unlock()

	state, err := n.getOrCreateTopicStateLocked(ctx, topic)
	if err != nil {
		return nil, "", err
	}
	state.addSubscriber(sessionID, sub)

	return sub.ch, sessionID, nil
}

func Publish[T any](n *Notifier, ctx context.Context, topic string, value T) error {
	if n == nil {
		return fmt.Errorf("notifier is nil")
	}

	topicData, err := n.getTopic(topic, typeOf[T]())
	if err != nil {
		return err
	}

	data, err := topicData.encode(value)
	if err != nil {
		return err
	}

	return n.transport.Publish(ctx, topic, data)
}

func (n *Notifier) Close() error {
	if n == nil || n.transport == nil {
		return nil
	}

	return n.transport.Close()
}

func Unsubscribe(n *Notifier, topic string, sessionID string) {
	if n == nil {
		return
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	state, ok := n.active[topic]
	if !ok {
		return
	}

	state.removeSubscriber(sessionID)
	if len(state.subscribers) == 0 { // можно без лока тк state не может поменять количество подписчиков после n.mu.Lock()
		delete(n.active, topic)
	}
}
