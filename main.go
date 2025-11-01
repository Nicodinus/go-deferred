package deferred

// static package methods

func Create[T any]() Deferred[T] {
	return Deferred[T]{}
}

func Go[T any](handler func() (*T, error)) *Promise[T] {
	d := Deferred[T]{}
	return d.Go(handler)
}
