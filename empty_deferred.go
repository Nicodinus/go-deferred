package deferred

type EmptyDeferred struct {
	promise *EmptyPromise
}

func (d *EmptyDeferred) Go(handler func() error) *EmptyPromise {
	go func() {
		err := handler()
		if err != nil {
			d.Reject(err)
		} else {
			d.Resolve()
		}
	}()
	return d.Promise()
}

func (d *EmptyDeferred) IsResolved() bool {
	return d.promise.IsResolved()
}

func (d *EmptyDeferred) Resolve() error {
	return d.promise.resolve(nil)
}

func (d *EmptyDeferred) Reject(err error) error {
	return d.promise.resolve(err)
}

func (d *EmptyDeferred) Cancel() error {
	return d.promise.resolve(ErrPromiseCancelled)
}

func (d *EmptyDeferred) Promise() *EmptyPromise {
	return d.promise
}
