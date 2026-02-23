package tests

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/nicodinus/go-deferred"
	"github.com/stretchr/testify/require"
)

func CreateEmptyDeferredAndResolveAfterWait(t *testing.T, ctx context.Context) {
	t.Helper()

	d := deferred.CreateEmptyDeferred()
	resolveChan := make(chan testChanHelperResult1[any], 1)
	defer close(resolveChan)
	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		time.Sleep(100 * time.Millisecond)
		resolveChan <- testChanHelperResult1[any]{id: testChanHelperResult1ResolveGoroutine, err: d.Resolve()}
		wg.Done()
	}()
	d.Promise().OnSuccess(func() {
		resolveChan <- testChanHelperResult1[any]{id: testChanHelperResult1PromiseOnSuccess}
		wg.Done()
	})
	d.Promise().OnFail(func(err error) {
		resolveChan <- testChanHelperResult1[any]{id: testChanHelperResult1PromiseOnFail, err: err}
		wg.Done()
	})
	d.Promise().OnResolve(func(err error) {
		resolveChan <- testChanHelperResult1[any]{id: testChanHelperResult1PromiseOnResolve, err: err}
		wg.Done()
	})
	go func() {
		time.Sleep(300 * time.Millisecond)
		resolveChan <- testChanHelperResult1[any]{}
	}()

	require.False(t, d.Promise().IsResolved())
	err := d.Promise().Wait(ctx)
	require.True(t, d.Promise().IsResolved())
	require.NoError(t, err)

	for range 4 {
		tv := <-resolveChan
		switch tv.id {
		case testChanHelperResult1ResolveGoroutine:
			require.NoError(t, tv.err)
		case testChanHelperResult1PromiseOnSuccess:
			require.EqualValues(t, nil, tv.val)
		case testChanHelperResult1PromiseOnFail:
			t.Fatal("unexpected call")
		case testChanHelperResult1PromiseOnResolve:
			require.EqualValues(t, nil, tv.val)
			require.NoError(t, tv.err)
		}
	}
}

func CreateResolvedEmptyDeferred(t *testing.T, ctx context.Context) {
	t.Helper()

	d := deferred.CreateEmptyDeferred()
	resolveChan := make(chan testChanHelperResult1[any], 1)
	defer close(resolveChan)
	wg := sync.WaitGroup{}
	wg.Add(3)

	d.Promise().OnSuccess(func() {
		resolveChan <- testChanHelperResult1[any]{id: testChanHelperResult1PromiseOnSuccess}
		wg.Done()
	})
	d.Promise().OnFail(func(err error) {
		resolveChan <- testChanHelperResult1[any]{id: testChanHelperResult1PromiseOnFail, err: err}
		wg.Done()
	})
	d.Promise().OnResolve(func(err error) {
		resolveChan <- testChanHelperResult1[any]{id: testChanHelperResult1PromiseOnResolve, err: err}
		wg.Done()
	})
	go func() {
		time.Sleep(300 * time.Millisecond)
		resolveChan <- testChanHelperResult1[any]{}
	}()

	resolveChan <- testChanHelperResult1[any]{id: testChanHelperResult1ResolveGoroutine, err: d.Resolve()}
	wg.Done()

	require.True(t, d.Promise().IsResolved())
	err := d.Promise().Wait(ctx)
	require.True(t, d.Promise().IsResolved())
	require.NoError(t, err)

	for range 4 {
		tv := <-resolveChan
		switch tv.id {
		case testChanHelperResult1ResolveGoroutine:
			require.NoError(t, tv.err)
		case testChanHelperResult1PromiseOnSuccess:
			require.EqualValues(t, nil, tv.val)
		case testChanHelperResult1PromiseOnFail:
			t.Fatal("unexpected call")
		case testChanHelperResult1PromiseOnResolve:
			require.EqualValues(t, nil, tv.val)
			require.NoError(t, tv.err)
		}
	}
}

func CreateEmptyDeferredAndRejectAfterWait(t *testing.T, rejectVal error, ctx context.Context) {
	t.Helper()

	d := deferred.CreateEmptyDeferred()
	resolveChan := make(chan testChanHelperResult1[any], 1)
	defer close(resolveChan)
	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		time.Sleep(100 * time.Millisecond)
		resolveChan <- testChanHelperResult1[any]{id: testChanHelperResult1ResolveGoroutine, err: d.Reject(rejectVal)}
		wg.Done()
	}()
	d.Promise().OnSuccess(func() {
		resolveChan <- testChanHelperResult1[any]{id: testChanHelperResult1PromiseOnSuccess}
		wg.Done()
	})
	d.Promise().OnFail(func(err error) {
		resolveChan <- testChanHelperResult1[any]{id: testChanHelperResult1PromiseOnFail, err: err}
		wg.Done()
	})
	d.Promise().OnResolve(func(err error) {
		resolveChan <- testChanHelperResult1[any]{id: testChanHelperResult1PromiseOnResolve, err: err}
		wg.Done()
	})
	go func() {
		time.Sleep(300 * time.Millisecond)
		resolveChan <- testChanHelperResult1[any]{}
	}()

	require.False(t, d.Promise().IsResolved())
	err := d.Promise().Wait(ctx)
	require.True(t, d.Promise().IsResolved())
	require.ErrorIs(t, err, rejectVal)

	for range 4 {
		tv := <-resolveChan
		switch tv.id {
		case testChanHelperResult1ResolveGoroutine:
			require.NoError(t, tv.err)
		case testChanHelperResult1PromiseOnSuccess:
			t.Fatal("unexpected call")
		case testChanHelperResult1PromiseOnFail:
			require.ErrorIs(t, tv.err, rejectVal)
		case testChanHelperResult1PromiseOnResolve:
			require.EqualValues(t, nil, tv.val)
			require.ErrorIs(t, tv.err, rejectVal)
		}
	}
}

func CreateRejectedEmptyDeferred(t *testing.T, rejectVal error, ctx context.Context) {
	t.Helper()

	d := deferred.CreateEmptyDeferred()
	resolveChan := make(chan testChanHelperResult1[any], 1)
	defer close(resolveChan)
	wg := sync.WaitGroup{}
	wg.Add(3)

	d.Promise().OnSuccess(func() {
		resolveChan <- testChanHelperResult1[any]{id: testChanHelperResult1PromiseOnSuccess}
		wg.Done()
	})
	d.Promise().OnFail(func(err error) {
		resolveChan <- testChanHelperResult1[any]{id: testChanHelperResult1PromiseOnFail, err: err}
		wg.Done()
	})
	d.Promise().OnResolve(func(err error) {
		resolveChan <- testChanHelperResult1[any]{id: testChanHelperResult1PromiseOnResolve, err: err}
		wg.Done()
	})
	go func() {
		time.Sleep(300 * time.Millisecond)
		resolveChan <- testChanHelperResult1[any]{}
	}()

	resolveChan <- testChanHelperResult1[any]{id: testChanHelperResult1ResolveGoroutine, err: d.Reject(rejectVal)}
	wg.Done()

	require.True(t, d.Promise().IsResolved())
	err := d.Promise().Wait(ctx)
	require.True(t, d.Promise().IsResolved())
	require.ErrorIs(t, err, rejectVal)

	for range 4 {
		tv := <-resolveChan
		switch tv.id {
		case testChanHelperResult1ResolveGoroutine:
			require.NoError(t, tv.err)
		case testChanHelperResult1PromiseOnSuccess:
			t.Fatal("unexpected call")
		case testChanHelperResult1PromiseOnFail:
			require.ErrorIs(t, tv.err, rejectVal)
		case testChanHelperResult1PromiseOnResolve:
			require.EqualValues(t, nil, tv.val)
			require.ErrorIs(t, tv.err, rejectVal)
		}
	}
}

func CreateEmptyDeferredAndCancelAfterWait(t *testing.T, ctx context.Context) {
	t.Helper()
	d := deferred.CreateEmptyDeferred()
	resolveChan := make(chan error, 1)
	go func() {
		time.Sleep(100 * time.Millisecond)
		resolveChan <- d.Cancel()
		close(resolveChan)
	}()

	err := d.Promise().Wait(ctx)
	require.ErrorIs(t, err, deferred.ErrPromiseCancelled)

	require.NoError(t, <-resolveChan)
}

func CreateCancelledEmptyDeferred(t *testing.T, ctx context.Context) {
	t.Helper()
	d := deferred.CreateEmptyDeferred()
	require.NoError(t, d.Cancel())

	err := d.Promise().Wait(ctx)
	require.ErrorIs(t, err, deferred.ErrPromiseCancelled)

}

func CreateEmptyDeferredAndMultipleResolveAfterWait(t *testing.T, ctx context.Context) {
	t.Helper()
	d := deferred.CreateEmptyDeferred()
	resolveChan := make(chan error, 2)
	go func() {
		time.Sleep(100 * time.Millisecond)
		resolveChan <- d.Resolve()
		time.Sleep(100 * time.Millisecond)
		resolveChan <- d.Resolve()
		close(resolveChan)
	}()
	err := d.Promise().Wait(ctx)
	require.NoError(t, err)

	require.NoError(t, <-resolveChan)
	require.ErrorIs(t, <-resolveChan, deferred.ErrPromiseResolved)
}

func CreateEmptyDeferredAndMultipleRejectAfterWait(t *testing.T, rejectVal error, ctx context.Context) {
	t.Helper()
	d := deferred.CreateEmptyDeferred()
	resolveChan := make(chan error, 2)
	go func() {
		time.Sleep(100 * time.Millisecond)
		resolveChan <- d.Reject(rejectVal)
		time.Sleep(100 * time.Millisecond)
		resolveChan <- d.Reject(rejectVal)
		close(resolveChan)
	}()

	err := d.Promise().Wait(ctx)
	require.ErrorIs(t, err, rejectVal)

	require.NoError(t, <-resolveChan)
	require.ErrorIs(t, <-resolveChan, deferred.ErrPromiseResolved)
}

func CreateResolvedEmptyDeferredAndResolveAfterWait(t *testing.T, ctx context.Context) {
	t.Helper()
	d := deferred.CreateEmptyDeferred()
	require.NoError(t, d.Resolve())
	resolveChan := make(chan error, 1)
	go func() {
		time.Sleep(100 * time.Millisecond)
		resolveChan <- d.Resolve()
		close(resolveChan)
	}()
	err := d.Promise().Wait(ctx)
	require.NoError(t, err)

	require.ErrorIs(t, <-resolveChan, deferred.ErrPromiseResolved)
}

func CreateRejectedEmptyDeferredAndResolveAfterWait(t *testing.T, rejectVal error, ctx context.Context) {
	t.Helper()
	d := deferred.CreateEmptyDeferred()
	require.NoError(t, d.Reject(rejectVal))
	resolveChan := make(chan error, 1)
	go func() {
		time.Sleep(100 * time.Millisecond)
		resolveChan <- d.Resolve()
		close(resolveChan)
	}()

	err := d.Promise().Wait(ctx)
	require.ErrorIs(t, err, rejectVal)

	require.ErrorIs(t, <-resolveChan, deferred.ErrPromiseResolved)
}

func CreateEmptyDeferredGo(t *testing.T, rejectVal error, ctx context.Context) {
	t.Helper()

	// test success
	promise, cancelFn := deferred.GoEmpty(func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		return nil
	})
	defer cancelFn()

	err := promise.Wait(ctx)
	require.NoError(t, err)

	// test reject
	promise, cancelFn = deferred.GoEmpty(func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		return rejectVal
	})
	defer cancelFn()

	err = promise.Wait(ctx)
	require.ErrorIs(t, err, rejectVal)

	// test ???
	promise, cancelFn = deferred.GoEmpty(func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		return rejectVal
	})
	defer cancelFn()

	err = promise.Wait(ctx)
	require.ErrorIs(t, err, rejectVal)

	// test cancel
	ctxCancelChan := make(chan error, 1)
	defer close(ctxCancelChan)

	promise, cancelFn = deferred.GoEmpty(func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		<-ctx.Done()
		ctxCancelChan <- ctx.Err()
		return ctx.Err()
	})
	cancelFn()

	err = promise.Wait(ctx)
	require.ErrorIs(t, err, deferred.ErrPromiseCancelled)

	require.ErrorIs(t, <-ctxCancelChan, context.Canceled)
}

func TestEmptyDeferred(t *testing.T) {
	t.Parallel()

	testError := fmt.Errorf("test error")

	t.Run("Simple 1", func(t *testing.T) {
		t.Parallel()

		testCtx, testCancelFn := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
		defer testCancelFn()

		CreateEmptyDeferredAndResolveAfterWait(t, testCtx)
		CreateResolvedEmptyDeferred(t, testCtx)
		CreateEmptyDeferredAndRejectAfterWait(t, testError, testCtx)
		CreateRejectedEmptyDeferred(t, testError, testCtx)
		CreateEmptyDeferredAndCancelAfterWait(t, testCtx)
		CreateCancelledEmptyDeferred(t, testCtx)
		CreateEmptyDeferredAndMultipleResolveAfterWait(t, testCtx)
		CreateEmptyDeferredAndMultipleRejectAfterWait(t, testError, testCtx)
		CreateResolvedEmptyDeferredAndResolveAfterWait(t, testCtx)
		CreateRejectedEmptyDeferredAndResolveAfterWait(t, testError, testCtx)
		CreateEmptyDeferredGo(t, testError, testCtx)
	})

	t.Run("Edge cases", func(t *testing.T) {
		t.Parallel()

		testCtx, testCancelFn := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
		defer testCancelFn()

		// case 1
		{
			promise, cancelFn := deferred.GoEmpty(func(ctx context.Context) error {
				return nil
			})
			defer cancelFn()

			err := promise.Wait(testCtx)
			require.NoError(t, err)
		}

		// case 2
		{
			promise, cancelFn := deferred.GoEmpty(func(ctx context.Context) error {
				return testError
			})
			defer cancelFn()

			err := promise.Wait(testCtx)
			require.ErrorIs(t, err, testError)
		}
	})
}
