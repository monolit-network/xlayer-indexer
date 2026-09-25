package promise

import (
	"context"
	"sync"
)

type Promise[T any] struct {
	state *promiseState[T]
}

type promiseState[T any] struct {
	once   sync.Once
	doneCh chan struct{}
	val    T
	err    error
}

func NewPromise[T any]() Promise[T] {
	return Promise[T]{
		state: &promiseState[T]{
			doneCh: make(chan struct{}),
		},
	}
}

func (p Promise[T]) Resolve(val T) {
	p.state.once.Do(func() {
		p.state.val = val
		close(p.state.doneCh)
	})
}

func (p Promise[T]) Reject(err error) {
	p.state.once.Do(func() {
		p.state.err = err
		close(p.state.doneCh)
	})
}

func (p Promise[T]) Wait() (T, error) {
	return p.WaitContext(context.Background())
}

func (p Promise[T]) WaitContext(ctx context.Context) (T, error) {
	select {
	case <-ctx.Done():
		select {
		case <-p.state.doneCh:
			return p.state.val, p.state.err
		default:
			var zero T
			return zero, ctx.Err()
		}
	case <-p.state.doneCh:
		return p.state.val, p.state.err
	}
}
