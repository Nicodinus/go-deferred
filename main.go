package deferred

import "context"

// static package methods

func Go[T any](handler func(ctx context.Context) (T, error)) (Promise[T], context.CancelFunc) {
	d := CreateDeferred[T]()
	return d.Go(handler)
}

func GoEmpty(handler func(ctx context.Context) error) (EmptyPromise, context.CancelFunc) {
	d := CreateEmptyDeferred()
	return d.Go(handler)
}
