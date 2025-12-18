package deferred

// static package methods

func Create[T any]() Deferred[T] {
	return Deferred[T]{
		promise: &Promise[T]{},
	}
}

func CreateEmpty() EmptyDeferred {
	return EmptyDeferred{
		promise: &EmptyPromise{},
	}
}

func Go[T any](handler func() (*T, error)) *Promise[T] {
	d := Create[T]()
	return d.Go(handler)
}

func GoEmpty(handler func() error) *EmptyPromise {
	d := CreateEmpty()
	return d.Go(handler)
}
