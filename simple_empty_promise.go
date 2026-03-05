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

	cond *sync.Cond
}

func (p *simpleEmptyPromise) IsResolved() bool {
	return p.isResolved.Load()
}

func (p *simpleEmptyPromise) OnResolve(h func(err error)) {
	p.cond.L.Lock()
	defer p.cond.L.Unlock()

	if p.IsResolved() {
		go h(p.err)
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
	p.cond.L.Lock()
	defer p.cond.L.Unlock()

	if p.IsResolved() {
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

	p.cond.Wait()

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

	p.cond.L.Lock()
	defer func() {
		p.isResolved.Store(true)
		p.cond.L.Unlock()
		p.cond.Broadcast()
	}()

	if err != nil {
		p.err = err
	}

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
	return &simpleEmptyPromise{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}
