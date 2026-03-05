package deferred

import (
	"context"
	"sync"
	"sync/atomic"
)

type simplePromise[T any] struct {
	val T
	err error

	handlers   []func(T, error)
	isResolved atomic.Bool

	cond *sync.Cond

	zeroVal T
}

func (p *simplePromise[T]) IsResolved() bool {
	return p.isResolved.Load()
}

func (p *simplePromise[T]) OnResolve(h func(val T, err error)) {
	p.cond.L.Lock()
	defer p.cond.L.Unlock()

	if p.IsResolved() {
		go h(p.val, p.err)
		return
	}

	p.handlers = append(p.handlers, h)
}

func (p *simplePromise[T]) OnSuccess(h func(val T)) {
	p.OnResolve(func(val T, err error) {
		if err != nil {
			return
		}
		h(val)
	})
}

func (p *simplePromise[T]) OnFail(h func(err error)) {
	p.OnResolve(func(val T, err error) {
		if err == nil {
			return
		}
		h(err)
	})
}

func (p *simplePromise[T]) Wait(ctx context.Context) (T, error) {
	p.cond.L.Lock()
	defer p.cond.L.Unlock()

	if p.IsResolved() {
		return p.val, p.err
	}

	innerCtx, innerCancelFn := context.WithCancel(ctx)
	defer innerCancelFn()

	go func() {
		<-innerCtx.Done()
		if !p.IsResolved() {
			p.resolve(p.zeroVal, innerCtx.Err())
		}
	}()

	p.cond.Wait()

	return p.val, p.err
}

func (p *simplePromise[T]) success(val T) error {
	return p.resolve(val, nil)
}

func (p *simplePromise[T]) fail(err error) error {
	return p.resolve(p.zeroVal, err)
}

func (p *simplePromise[T]) resolve(val T, err error) error {
	if p.IsResolved() {
		return ErrPromiseResolved
	}

	p.cond.L.Lock()
	defer func() {
		p.cond.L.Unlock()
		p.cond.Broadcast()
	}()
	p.isResolved.Store(true)

	if err == nil {
		p.val = val
	} else {
		p.err = err
	}

	handlers := p.handlers
	p.handlers = nil

	go func() {
		for _, h := range handlers {
			h(val, err)
		}
	}()

	return nil
}

func createSimplePromise[T any]() *simplePromise[T] {
	return &simplePromise[T]{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}
