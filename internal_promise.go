package deferred

type internalPromise[T any] interface {
	Promise[T]

	success(val T) error
	fail(err error) error
	resolve(val T, err error) error
}
