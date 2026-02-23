package deferred

import (
	"context"
	"sync"
	"sync/atomic"
)

type simpleEmptyPromise struct {
	err error

	handlers   []func(error)
	isResolved atomic.Bool

	m            sync.RWMutex
	resolveMutex sync.RWMutex
}

func (p *simpleEmptyPromise) IsResolved() bool {
	return p.isResolved.Load()
}

func (p *simpleEmptyPromise) OnResolve(h func(err error)) {
	if p.IsResolved() {
		p.m.RLock()
		defer p.m.RUnlock()
		go h(p.err)
		return
	}

	p.m.Lock()
	defer p.m.Unlock()

	if p.IsResolved() {
		p.OnResolve(h)
		return
	}

	p.handlers = append(p.handlers, h)
}

func (p *simpleEmptyPromise) OnSuccess(h func()) {
	p.OnResolve(func(err error) {
		if err != nil {
			return
		}
		h()
	})
}

func (p *simpleEmptyPromise) OnFail(h func(err error)) {
	p.OnResolve(func(err error) {
		if err == nil {
			return
		}
		h(err)
	})
}

func (p *simpleEmptyPromise) Wait(ctx context.Context) error {
	if p.IsResolved() {
		p.m.RLock()
		defer p.m.RUnlock()

		return p.err
	}

	innerCtx, innerCancelFn := context.WithCancel(ctx)
	defer innerCancelFn()

	go func() {
		<-innerCtx.Done()
		if !p.IsResolved() {
			p.resolve(innerCtx.Err())
		}
	}()

	p.resolveMutex.RLock()
	defer p.resolveMutex.RUnlock()

	p.m.RLock()
	defer p.m.RUnlock()

	return p.err
}

func (p *simpleEmptyPromise) success() error {
	return p.resolve(nil)
}

func (p *simpleEmptyPromise) fail(err error) error {
	return p.resolve(err)
}

func (p *simpleEmptyPromise) resolve(err error) error {
	if p.IsResolved() {
		return ErrPromiseResolved
	}
	p.m.Lock()
	defer p.m.Unlock()

	if p.IsResolved() {
		return ErrPromiseResolved
	}
	p.isResolved.Store(true)

	p.err = err

	defer p.resolveMutex.Unlock()

	handlers := p.handlers
	p.handlers = nil

	go func() {
		for _, h := range handlers {
			h(err)
		}
	}()

	return nil
}

func createSimpleEmptyPromise() *simpleEmptyPromise {
	p := simpleEmptyPromise{}
	p.resolveMutex.Lock()
	return &p
}
