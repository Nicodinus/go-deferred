package deferred

type Deferred[T any] struct {
	promise *Promise[T]
}

func (d *Deferred[T]) Go(handler func() (*T, error)) *Promise[T] {
	go func() {
		result, err := handler()
		if err != nil {
			d.Reject(err)
		} else {
			d.Resolve(result)
		}
	}()
	return d.Promise()
}

func (d *Deferred[T]) IsResolved() bool {
	return d.promise.IsResolved()
}

func (d *Deferred[T]) Resolve(val *T) error {
	return d.promise.resolve(val, nil)
}

func (d *Deferred[T]) Reject(err error) error {
	return d.promise.resolve(nil, err)
}

func (d *Deferred[T]) Cancel() error {
	return d.promise.resolve(nil, ErrPromiseCancelled)
}

func (d *Deferred[T]) Promise() *Promise[T] {
	return d.promise
}
