package deferred

import "context"

type EmptyDeferred struct {
	promise internalEmptyPromise
}

func (d *EmptyDeferred) Go(handler func(ctx context.Context) error) (EmptyPromise, context.CancelFunc) {
	ctx, cancelFn := context.WithCancel(context.Background())
	go func() {
		defer cancelFn()
		err := handler(ctx)
		if err != nil {
			d.Reject(err)
		} else {
			d.Resolve()
		}
	}()
	return d.Promise(), func() {
		d.Cancel()
		cancelFn()
	}
}

func (d *EmptyDeferred) IsResolved() bool {
	return d.promise.IsResolved()
}

func (d *EmptyDeferred) Resolve() error {
	return d.promise.success()
}

func (d *EmptyDeferred) Reject(err error) error {
	return d.promise.fail(err)
}

func (d *EmptyDeferred) Cancel() error {
	return d.Reject(ErrPromiseCancelled)
}

func (d *EmptyDeferred) Promise() EmptyPromise {
	return d.promise
}

func CreateEmptyDeferred() *EmptyDeferred {
	return &EmptyDeferred{
		promise: createSimpleEmptyPromise(),
	}
}
