package notifier

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type Topic struct {
	matches   func(topic string) bool
	valueType reflect.Type
	encode    func(value any) ([]byte, error)
	decode    func(data []byte) (any, error)
}

type topicState struct {
	closeRemote func() error
	subscribers map[string]typedSubscriber
	mu          sync.RWMutex
}

func (s *topicState) makeSnapshot() []typedSubscriber {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.subscribers) == 0 {
		return nil
	}

	typedSubs := make([]typedSubscriber, 0, len(s.subscribers))
	for _, sub := range s.subscribers {
		typedSubs = append(typedSubs, sub)
	}
	return typedSubs
}

func (s *topicState) addSubscriber(sessionID string, sub typedSubscriber) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.subscribers[sessionID] = sub
}

func (s *topicState) removeSubscriber(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if sub, ok := s.subscribers[sessionID]; ok {
		delete(s.subscribers, sessionID)
		sub.close()
	}

	if len(s.subscribers) > 0 {
		return
	}

	if len(s.subscribers) == 0 {
		s.closeRemote()
	}
}

func TopicOf[T any](name string) Topic {
	return newTopic[T](func(topic string) bool {
		return strings.HasPrefix(topic, name)
	})
}

func (n *Notifier) getOrCreateTopicStateLocked(ctx context.Context, topic string) (*topicState, error) {
	if state, ok := n.active[topic]; ok {
		return state, nil
	}

	topicData, ok := n.findTopic(topic)
	if !ok {
		return nil, fmt.Errorf("topic %q is not registered", topic)
	}

	state := &topicState{
		subscribers: make(map[string]typedSubscriber),
	}

	closeRemote, err := n.transport.Subscribe(ctx, topic, func(data []byte) {
		rawData := append([]byte(nil), data...)

		typedSubs := state.makeSnapshot()

		value, err := topicData.decode(rawData)
		if err != nil {
			return
		}

		for _, sub := range typedSubs {
			sub.deliver(value, rawData)
		}
	})
	if err != nil {
		return nil, err
	}

	state.closeRemote = closeRemote
	n.active[topic] = state
	return state, nil
}

func (n *Notifier) getTopic(topic string, expectedType reflect.Type) (Topic, error) {
	spec, ok := n.findTopic(topic)
	if !ok {
		return Topic{}, fmt.Errorf("topic %q is not registered", topic)
	}
	if spec.valueType != expectedType {
		return Topic{}, fmt.Errorf("topic %q registered with type %s, got %s", topic, spec.valueType, expectedType)
	}

	return spec, nil
}

func (n *Notifier) findTopic(topic string) (Topic, bool) {
	for _, spec := range n.topics {
		if spec.matches(topic) {
			return spec, true
		}
	}

	return Topic{}, false
}

func newTopic[T any](matches func(topic string) bool) Topic {
	valueType := typeOf[T]()

	return Topic{
		matches:   matches,
		valueType: valueType,
		encode: func(value any) ([]byte, error) {
			typedValue, ok := value.(T)
			if !ok {
				return nil, fmt.Errorf("expected value of type %s, got %T", valueType, value)
			}

			data, err := json.Marshal(typedValue)
			if err != nil {
				return nil, err
			}
			return data, nil
		},
		decode: func(data []byte) (any, error) {
			var value T
			if err := json.Unmarshal(data, &value); err != nil {
				return nil, err
			}
			return value, nil
		},
	}
}

func typeOf[T any]() reflect.Type {
	var zero *T
	return reflect.TypeOf(zero).Elem()
}
