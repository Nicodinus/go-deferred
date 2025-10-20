package deferred

import (
	"context"
	"sync"
	"sync/atomic"
)

type Promise[T any] struct {
	val *T
	err error

	waitChannels []chan struct{}
	handlers     []func(*T, error)

	isResolved atomic.Bool

	m sync.RWMutex
}

//

func (p *Promise[T]) OnResolve(h func(val *T, err error)) {
	if p.IsResolved() {
		p.m.RLock()
		val := p.val
		err := p.err
		p.m.RUnlock()
		go h(val, err)
		return
	}

	p.m.Lock()
	p.handlers = append(p.handlers, h)
	p.m.Unlock()
}

func (p *Promise[T]) OnSuccess(h func(*T)) {
	p.OnResolve(func(val *T, err error) {
		if err != nil {
			return
		}
		h(val)
	})
}

func (p *Promise[T]) OnFail(h func(error)) {
	p.OnResolve(func(val *T, err error) {
		if err == nil {
			return
		}
		h(err)
	})
}

func (p *Promise[T]) Wait(ctx context.Context) (*T, error) {
	if p.IsResolved() {
		p.m.RLock()
		defer p.m.RUnlock()
		return p.val, p.err
	}

	p.m.Lock()
	ch := make(chan struct{}, 1)
	p.waitChannels = append(p.waitChannels, ch)
	p.m.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-ch:
		p.m.RLock()
		defer p.m.RUnlock()
		return p.val, p.err
	}
}

func (p *Promise[T]) IsResolved() bool {
	return p.isResolved.Load()
}

//

func (p *Promise[T]) resolve(val *T, err error) error {
	if p.IsResolved() {
		return ErrPromiseResolved
	}
	p.isResolved.Store(true)

	p.m.Lock()

	p.val = val
	p.err = err

	waitChannels := p.waitChannels
	p.waitChannels = nil

	handlers := p.handlers
	p.handlers = nil

	p.m.Unlock()

	go func() {
		for _, ch := range waitChannels {
			ch <- struct{}{}
			close(ch)
		}
	}()

	go func() {
		for _, h := range handlers {
			h(val, err)
		}
	}()

	return nil
}
