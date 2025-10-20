package deferred

// static package methods

func Create[T any]() Deferred[T] {
	return Deferred[T]{}
}
