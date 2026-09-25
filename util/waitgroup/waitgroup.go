package waitgroup

import (
	"errors"
	"sync"
)

type WaitGroup struct {
	wg  sync.WaitGroup
	mu  sync.Mutex
	err error
}

func (wg *WaitGroup) Go(f func() error) {
	wg.wg.Add(1)
	go func() {
		defer wg.wg.Done()
		if err := f(); err != nil {
			wg.mu.Lock()
			wg.err = errors.Join(wg.err, err)
			wg.mu.Unlock()
		}
	}()
}

func (wg *WaitGroup) Wait() error {
	wg.wg.Wait()
	wg.mu.Lock()
	defer wg.mu.Unlock()
	return wg.err
}

func (wg *WaitGroup) WaitCh() <-chan error {
	ch := make(chan error, 1)
	go func() {
		ch <- wg.Wait()
		close(ch)
	}()
	return ch
}
