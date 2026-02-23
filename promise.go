package deferred

import "context"

type Promise[T any] interface {
	IsResolved() bool

	OnResolve(func(T, error))

	OnSuccess(func(T))

	OnFail(func(error))

	Wait(context.Context) (T, error)
}
