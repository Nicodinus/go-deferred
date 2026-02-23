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

	m            sync.RWMutex
	resolveMutex sync.RWMutex

	zeroVal T
}

func (p *simplePromise[T]) IsResolved() bool {
	return p.isResolved.Load()
}

func (p *simplePromise[T]) OnResolve(h func(val T, err error)) {
	if p.IsResolved() {
		p.m.RLock()
		defer p.m.RUnlock()
		// if p.err != nil {
		// 	var zeroVal T
		// 	go h(zeroVal, p.err)
		// } else {
		// 	go h(p.val, p.err)
		// }
		go h(p.val, p.err)
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
	if p.IsResolved() {
		p.m.RLock()
		defer p.m.RUnlock()

		// if p.err != nil {
		// 	var zeroVal T
		// 	return zeroVal, p.err
		// }

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

	p.resolveMutex.RLock()
	defer p.resolveMutex.RUnlock()

	p.m.RLock()
	defer p.m.RUnlock()

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
	p.m.Lock()
	defer p.m.Unlock()

	if p.IsResolved() {
		return ErrPromiseResolved
	}
	p.isResolved.Store(true)

	if err == nil {
		p.val = val
	}
	p.err = err

	defer p.resolveMutex.Unlock()

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
	p := simplePromise[T]{}
	p.resolveMutex.Lock()
	return &p
}
