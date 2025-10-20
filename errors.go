package deferred

import "errors"

var ErrPromiseResolved = errors.New("promise resolved")
var ErrPromiseCancelled = errors.New("promise cancelled")
