package deferred

type internalEmptyPromise interface {
	EmptyPromise

	success() error
	fail(err error) error
	resolve(err error) error
}
