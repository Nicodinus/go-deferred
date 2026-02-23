package deferred

import "context"

type EmptyPromise interface {
	IsResolved() bool

	OnResolve(func(error))

	OnSuccess(func())

	OnFail(func(error))

	Wait(context.Context) error
}
