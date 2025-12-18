package deferred

import (
	"context"
	"sync"
	"sync/atomic"
)

type EmptyPromise struct {
	err error

	waitChannels []chan struct{}
	handlers     []func(error)

	isResolved atomic.Bool

	m sync.RWMutex
}

//

func (p *EmptyPromise) OnResolve(h func(err error)) {
	if p.IsResolved() {
		p.m.RLock()
		err := p.err
		p.m.RUnlock()
		go h(err)
		return
	}

	p.m.Lock()
	p.handlers = append(p.handlers, h)
	p.m.Unlock()
}

func (p *EmptyPromise) OnSuccess(h func()) {
	p.OnResolve(func(err error) {
		if err != nil {
			return
		}
		h()
	})
}

func (p *EmptyPromise) OnFail(h func(error)) {
	p.OnResolve(func(err error) {
		if err == nil {
			return
		}
		h(err)
	})
}

func (p *EmptyPromise) Wait(ctx context.Context) error {
	if p.IsResolved() {
		p.m.RLock()
		defer p.m.RUnlock()
		return p.err
	}

	p.m.Lock()
	ch := make(chan struct{}, 1)
	p.waitChannels = append(p.waitChannels, ch)
	p.m.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ch:
		p.m.RLock()
		defer p.m.RUnlock()
		return p.err
	}
}

func (p *EmptyPromise) IsResolved() bool {
	return p.isResolved.Load()
}

//

func (p *EmptyPromise) resolve(err error) error {
	if p.IsResolved() {
		return ErrPromiseResolved
	}
	p.isResolved.Store(true)

	p.m.Lock()

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
			h(err)
		}
	}()

	return nil
}
