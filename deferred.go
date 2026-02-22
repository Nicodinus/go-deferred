package deferred

import "context"

type Deferred[T any] struct {
	promise internalPromise[T]
}

func (d *Deferred[T]) Go(handler func(ctx context.Context) (T, error)) (Promise[T], context.CancelFunc) {
	ctx, cancelFn := context.WithCancel(context.Background())
	go func() {
		defer cancelFn()
		result, err := handler(ctx)
		if err != nil {
			d.Reject(err)
		} else {
			d.Resolve(result)
		}
	}()
	return d.Promise(), func() {
		d.Cancel()
		cancelFn()
	}
}

func (d *Deferred[T]) IsResolved() bool {
	return d.promise.IsResolved()
}

func (d *Deferred[T]) Resolve(val T) error {
	return d.promise.success(val)
}

func (d *Deferred[T]) Reject(err error) error {
	return d.promise.fail(err)
}

func (d *Deferred[T]) Cancel() error {
	return d.Reject(ErrPromiseCancelled)
}

func (d *Deferred[T]) Promise() Promise[T] {
	return d.promise
}

func CreateDeferred[T any]() *Deferred[T] {
	return &Deferred[T]{
		promise: createSimplePromise[T](),
	}
}
