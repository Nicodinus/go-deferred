package tests

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/nicodinus/go-deferred"
	"github.com/stretchr/testify/require"
)

func CreateDeferredAndResolveAfterWait[T any](t *testing.T, val T, ctx context.Context) {
	t.Helper()

	d := deferred.CreateDeferred[T]()
	resolveChan := make(chan testChanHelperResult1[T], 1)
	defer close(resolveChan)
	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		time.Sleep(100 * time.Millisecond)
		resolveChan <- testChanHelperResult1[T]{id: testChanHelperResult1ResolveGoroutine, err: d.Resolve(val)}
		wg.Done()
	}()
	d.Promise().OnSuccess(func(val T) {
		resolveChan <- testChanHelperResult1[T]{id: testChanHelperResult1PromiseOnSuccess, val: val}
		wg.Done()
	})
	d.Promise().OnFail(func(err error) {
		resolveChan <- testChanHelperResult1[T]{id: testChanHelperResult1PromiseOnFail, err: err}
		wg.Done()
	})
	d.Promise().OnResolve(func(val T, err error) {
		resolveChan <- testChanHelperResult1[T]{id: testChanHelperResult1PromiseOnResolve, val: val, err: err}
		wg.Done()
	})
	go func() {
		time.Sleep(300 * time.Millisecond)
		resolveChan <- testChanHelperResult1[T]{}
	}()

	require.False(t, d.Promise().IsResolved())
	result, err := d.Promise().Wait(ctx)
	require.True(t, d.Promise().IsResolved())
	require.NoError(t, err)
	require.EqualValues(t, val, result)
	for range 4 {
		tv := <-resolveChan
		switch tv.id {
		case testChanHelperResult1ResolveGoroutine:
			require.NoError(t, tv.err)
		case testChanHelperResult1PromiseOnSuccess:
			require.EqualValues(t, val, tv.val)
		case testChanHelperResult1PromiseOnFail:
			t.Fatal("unexpected call")
		case testChanHelperResult1PromiseOnResolve:
			require.EqualValues(t, val, tv.val)
			require.NoError(t, tv.err)
		}
	}
}

func CreateResolvedDeferred[T any](t *testing.T, val T, ctx context.Context) {
	t.Helper()

	d := deferred.CreateDeferred[T]()
	resolveChan := make(chan testChanHelperResult1[T], 1)
	defer close(resolveChan)
	wg := sync.WaitGroup{}
	wg.Add(3)

	d.Promise().OnSuccess(func(val T) {
		resolveChan <- testChanHelperResult1[T]{id: testChanHelperResult1PromiseOnSuccess, val: val}
		wg.Done()
	})
	d.Promise().OnFail(func(err error) {
		resolveChan <- testChanHelperResult1[T]{id: testChanHelperResult1PromiseOnFail, err: err}
		wg.Done()
	})
	d.Promise().OnResolve(func(val T, err error) {
		resolveChan <- testChanHelperResult1[T]{id: testChanHelperResult1PromiseOnResolve, val: val, err: err}
		wg.Done()
	})
	go func() {
		time.Sleep(300 * time.Millisecond)
		resolveChan <- testChanHelperResult1[T]{}
	}()

	resolveChan <- testChanHelperResult1[T]{id: testChanHelperResult1ResolveGoroutine, err: d.Resolve(val)}
	wg.Done()

	require.True(t, d.Promise().IsResolved())
	result, err := d.Promise().Wait(ctx)
	require.True(t, d.Promise().IsResolved())
	require.NoError(t, err)
	require.EqualValues(t, val, result)
	for range 4 {
		tv := <-resolveChan
		switch tv.id {
		case testChanHelperResult1ResolveGoroutine:
			require.NoError(t, tv.err)
		case testChanHelperResult1PromiseOnSuccess:
			require.EqualValues(t, val, tv.val)
		case testChanHelperResult1PromiseOnFail:
			t.Fatal("unexpected call")
		case testChanHelperResult1PromiseOnResolve:
			require.EqualValues(t, val, tv.val)
			require.NoError(t, tv.err)
		}
	}
}

func CreateDeferredAndRejectAfterWait[T any](t *testing.T, val T, rejectVal error, ctx context.Context) {
	t.Helper()

	d := deferred.CreateDeferred[T]()
	resolveChan := make(chan testChanHelperResult1[T], 1)
	defer close(resolveChan)
	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		time.Sleep(100 * time.Millisecond)
		resolveChan <- testChanHelperResult1[T]{id: testChanHelperResult1ResolveGoroutine, err: d.Reject(rejectVal)}
		wg.Done()
	}()
	d.Promise().OnSuccess(func(val T) {
		resolveChan <- testChanHelperResult1[T]{id: testChanHelperResult1PromiseOnSuccess, val: val}
		wg.Done()
	})
	d.Promise().OnFail(func(err error) {
		resolveChan <- testChanHelperResult1[T]{id: testChanHelperResult1PromiseOnFail, err: err}
		wg.Done()
	})
	d.Promise().OnResolve(func(val T, err error) {
		resolveChan <- testChanHelperResult1[T]{id: testChanHelperResult1PromiseOnResolve, val: val, err: err}
		wg.Done()
	})
	go func() {
		time.Sleep(300 * time.Millisecond)
		resolveChan <- testChanHelperResult1[T]{}
	}()

	var zeroVal T
	require.False(t, d.Promise().IsResolved())
	result, err := d.Promise().Wait(ctx)
	require.True(t, d.Promise().IsResolved())
	require.ErrorIs(t, err, rejectVal)
	require.EqualValues(t, zeroVal, result)
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
			require.EqualValues(t, zeroVal, tv.val)
			require.ErrorIs(t, tv.err, rejectVal)
		}
	}
}

func CreateRejectedDeferred[T any](t *testing.T, val T, rejectVal error, ctx context.Context) {
	t.Helper()

	d := deferred.CreateDeferred[T]()
	resolveChan := make(chan testChanHelperResult1[T], 1)
	defer close(resolveChan)
	wg := sync.WaitGroup{}
	wg.Add(3)

	d.Promise().OnSuccess(func(val T) {
		resolveChan <- testChanHelperResult1[T]{id: testChanHelperResult1PromiseOnSuccess, val: val}
		wg.Done()
	})
	d.Promise().OnFail(func(err error) {
		resolveChan <- testChanHelperResult1[T]{id: testChanHelperResult1PromiseOnFail, err: err}
		wg.Done()
	})
	d.Promise().OnResolve(func(val T, err error) {
		resolveChan <- testChanHelperResult1[T]{id: testChanHelperResult1PromiseOnResolve, val: val, err: err}
		wg.Done()
	})
	go func() {
		time.Sleep(300 * time.Millisecond)
		resolveChan <- testChanHelperResult1[T]{}
	}()

	resolveChan <- testChanHelperResult1[T]{id: testChanHelperResult1ResolveGoroutine, err: d.Reject(rejectVal)}
	wg.Done()

	var zeroVal T
	require.True(t, d.Promise().IsResolved())
	result, err := d.Promise().Wait(ctx)
	require.True(t, d.Promise().IsResolved())
	require.ErrorIs(t, err, rejectVal)
	require.EqualValues(t, zeroVal, result)
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
			require.EqualValues(t, zeroVal, tv.val)
			require.ErrorIs(t, tv.err, rejectVal)
		}
	}
}

func CreateDeferredAndCancelAfterWait[T any](t *testing.T, val T, ctx context.Context) {
	t.Helper()
	d := deferred.CreateDeferred[T]()
	resolveChan := make(chan error, 1)
	go func() {
		time.Sleep(100 * time.Millisecond)
		resolveChan <- d.Cancel()
		close(resolveChan)
	}()
	var zeroVal T
	result, err := d.Promise().Wait(ctx)
	require.ErrorIs(t, err, deferred.ErrPromiseCancelled)
	require.EqualValues(t, zeroVal, result)
	require.NoError(t, <-resolveChan)
}

func CreateCancelledDeferred[T any](t *testing.T, val T, ctx context.Context) {
	t.Helper()
	d := deferred.CreateDeferred[T]()
	require.NoError(t, d.Cancel())
	var zeroVal T
	result, err := d.Promise().Wait(ctx)
	require.ErrorIs(t, err, deferred.ErrPromiseCancelled)
	require.EqualValues(t, zeroVal, result)
}

func CreateDeferredAndMultipleResolveAfterWait[T any](t *testing.T, val T, ctx context.Context) {
	t.Helper()
	d := deferred.CreateDeferred[T]()
	resolveChan := make(chan error, 2)
	go func() {
		time.Sleep(100 * time.Millisecond)
		resolveChan <- d.Resolve(val)
		time.Sleep(100 * time.Millisecond)
		resolveChan <- d.Resolve(val)
		close(resolveChan)
	}()
	result, err := d.Promise().Wait(ctx)
	require.NoError(t, err)
	require.EqualValues(t, val, result)
	require.NoError(t, <-resolveChan)
	require.ErrorIs(t, <-resolveChan, deferred.ErrPromiseResolved)
}

func CreateDeferredAndMultipleRejectAfterWait[T any](t *testing.T, val T, rejectVal error, ctx context.Context) {
	t.Helper()
	d := deferred.CreateDeferred[T]()
	resolveChan := make(chan error, 2)
	go func() {
		time.Sleep(100 * time.Millisecond)
		resolveChan <- d.Reject(rejectVal)
		time.Sleep(100 * time.Millisecond)
		resolveChan <- d.Reject(rejectVal)
		close(resolveChan)
	}()
	var zeroVal T
	result, err := d.Promise().Wait(ctx)
	require.ErrorIs(t, err, rejectVal)
	require.EqualValues(t, zeroVal, result)
	require.NoError(t, <-resolveChan)
	require.ErrorIs(t, <-resolveChan, deferred.ErrPromiseResolved)
}

func CreateResolvedDeferredAndResolveAfterWait[T any](t *testing.T, val T, ctx context.Context) {
	t.Helper()
	d := deferred.CreateDeferred[T]()
	require.NoError(t, d.Resolve(val))
	resolveChan := make(chan error, 1)
	go func() {
		time.Sleep(100 * time.Millisecond)
		resolveChan <- d.Resolve(val)
		close(resolveChan)
	}()
	result, err := d.Promise().Wait(ctx)
	require.NoError(t, err)
	require.EqualValues(t, val, result)
	require.ErrorIs(t, <-resolveChan, deferred.ErrPromiseResolved)
}

func CreateRejectedDeferredAndResolveAfterWait[T any](t *testing.T, val T, rejectVal error, ctx context.Context) {
	t.Helper()
	d := deferred.CreateDeferred[T]()
	require.NoError(t, d.Reject(rejectVal))
	resolveChan := make(chan error, 1)
	go func() {
		time.Sleep(100 * time.Millisecond)
		resolveChan <- d.Resolve(val)
		close(resolveChan)
	}()
	var zeroVal T
	result, err := d.Promise().Wait(ctx)
	require.ErrorIs(t, err, rejectVal)
	require.EqualValues(t, zeroVal, result)
	require.ErrorIs(t, <-resolveChan, deferred.ErrPromiseResolved)
}

func CreateDeferredGo[T any](t *testing.T, val T, rejectVal error, ctx context.Context) {
	t.Helper()

	var zeroVal T

	// test success
	promise, cancelFn := deferred.Go(func(ctx context.Context) (T, error) {
		time.Sleep(100 * time.Millisecond)
		return val, nil
	})
	defer cancelFn()

	result, err := promise.Wait(ctx)
	require.NoError(t, err)
	require.EqualValues(t, val, result)

	// test reject
	promise, cancelFn = deferred.Go(func(ctx context.Context) (T, error) {
		time.Sleep(100 * time.Millisecond)
		return zeroVal, rejectVal
	})
	defer cancelFn()

	result, err = promise.Wait(ctx)
	require.ErrorIs(t, err, rejectVal)
	require.EqualValues(t, zeroVal, result)

	// test ???
	promise, cancelFn = deferred.Go(func(ctx context.Context) (T, error) {
		time.Sleep(100 * time.Millisecond)
		return val, rejectVal
	})
	defer cancelFn()

	result, err = promise.Wait(ctx)
	require.ErrorIs(t, err, rejectVal)
	require.EqualValues(t, zeroVal, result)

	// test cancel
	ctxCancelChan := make(chan error, 1)
	defer close(ctxCancelChan)

	promise, cancelFn = deferred.Go(func(ctx context.Context) (T, error) {
		time.Sleep(100 * time.Millisecond)
		<-ctx.Done()
		ctxCancelChan <- ctx.Err()
		return val, ctx.Err()
	})
	cancelFn()

	result, err = promise.Wait(ctx)
	require.ErrorIs(t, err, deferred.ErrPromiseCancelled)
	require.EqualValues(t, zeroVal, result)
	require.ErrorIs(t, <-ctxCancelChan, context.Canceled)
}

func TestDeferred(t *testing.T) {
	t.Parallel()

	type testStruct1 struct {
		field1 string
		field2 int
	}

	testError := fmt.Errorf("test error")

	t.Run("With string", func(t *testing.T) {
		t.Parallel()

		testCtx, testCancelFn := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
		defer testCancelFn()

		testData := "yolo 1337"
		CreateDeferredAndResolveAfterWait(t, testData, testCtx)
		CreateResolvedDeferred(t, testData, testCtx)
		CreateDeferredAndRejectAfterWait(t, testData, testError, testCtx)
		CreateRejectedDeferred(t, testData, testError, testCtx)
		CreateDeferredAndCancelAfterWait(t, testData, testCtx)
		CreateCancelledDeferred(t, testData, testCtx)
		CreateDeferredAndMultipleResolveAfterWait(t, testData, testCtx)
		CreateDeferredAndMultipleRejectAfterWait(t, testData, testError, testCtx)
		CreateResolvedDeferredAndResolveAfterWait(t, testData, testCtx)
		CreateRejectedDeferredAndResolveAfterWait(t, testData, testError, testCtx)
		CreateDeferredGo(t, testData, testError, testCtx)
	})

	t.Run("With int", func(t *testing.T) {
		t.Parallel()

		testCtx, testCancelFn := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
		defer testCancelFn()

		testData := 1337
		CreateDeferredAndResolveAfterWait(t, testData, testCtx)
		CreateResolvedDeferred(t, testData, testCtx)
		CreateDeferredAndRejectAfterWait(t, testData, testError, testCtx)
		CreateRejectedDeferred(t, testData, testError, testCtx)
		CreateDeferredAndCancelAfterWait(t, testData, testCtx)
		CreateCancelledDeferred(t, testData, testCtx)
		CreateDeferredAndMultipleResolveAfterWait(t, testData, testCtx)
		CreateDeferredAndMultipleRejectAfterWait(t, testData, testError, testCtx)
		CreateResolvedDeferredAndResolveAfterWait(t, testData, testCtx)
		CreateRejectedDeferredAndResolveAfterWait(t, testData, testError, testCtx)
		CreateDeferredGo(t, testData, testError, testCtx)
	})

	t.Run("With struct", func(t *testing.T) {
		t.Parallel()

		testCtx, testCancelFn := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
		defer testCancelFn()

		testData := testStruct1{
			field1: "testDataStruct",
			field2: 456,
		}
		CreateDeferredAndResolveAfterWait(t, testData, testCtx)
		CreateResolvedDeferred(t, testData, testCtx)
		CreateDeferredAndRejectAfterWait(t, testData, testError, testCtx)
		CreateRejectedDeferred(t, testData, testError, testCtx)
		CreateDeferredAndCancelAfterWait(t, testData, testCtx)
		CreateCancelledDeferred(t, testData, testCtx)
		CreateDeferredAndMultipleResolveAfterWait(t, testData, testCtx)
		CreateDeferredAndMultipleRejectAfterWait(t, testData, testError, testCtx)
		CreateResolvedDeferredAndResolveAfterWait(t, testData, testCtx)
		CreateRejectedDeferredAndResolveAfterWait(t, testData, testError, testCtx)
		CreateDeferredGo(t, testData, testError, testCtx)
	})

	t.Run("With pointer to struct", func(t *testing.T) {
		t.Parallel()

		testCtx, testCancelFn := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
		defer testCancelFn()

		testData := &testStruct1{
			field1: "testDataStructPtr",
			field2: 1337,
		}
		CreateDeferredAndResolveAfterWait(t, testData, testCtx)
		CreateResolvedDeferred(t, testData, testCtx)
		CreateDeferredAndRejectAfterWait(t, testData, testError, testCtx)
		CreateRejectedDeferred(t, testData, testError, testCtx)
		CreateDeferredAndCancelAfterWait(t, testData, testCtx)
		CreateCancelledDeferred(t, testData, testCtx)
		CreateDeferredAndMultipleResolveAfterWait(t, testData, testCtx)
		CreateDeferredAndMultipleRejectAfterWait(t, testData, testError, testCtx)
		CreateResolvedDeferredAndResolveAfterWait(t, testData, testCtx)
		CreateRejectedDeferredAndResolveAfterWait(t, testData, testError, testCtx)
		CreateDeferredGo(t, testData, testError, testCtx)
	})

	t.Run("With nil pointer to struct", func(t *testing.T) {
		t.Parallel()

		testCtx, testCancelFn := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
		defer testCancelFn()

		var testData *testStruct1
		CreateDeferredAndResolveAfterWait(t, testData, testCtx)
		CreateResolvedDeferred(t, testData, testCtx)
		CreateDeferredAndRejectAfterWait(t, testData, testError, testCtx)
		CreateRejectedDeferred(t, testData, testError, testCtx)
		CreateDeferredAndCancelAfterWait(t, testData, testCtx)
		CreateCancelledDeferred(t, testData, testCtx)
		CreateDeferredAndMultipleResolveAfterWait(t, testData, testCtx)
		CreateDeferredAndMultipleRejectAfterWait(t, testData, testError, testCtx)
		CreateResolvedDeferredAndResolveAfterWait(t, testData, testCtx)
		CreateRejectedDeferredAndResolveAfterWait(t, testData, testError, testCtx)
		CreateDeferredGo(t, testData, testError, testCtx)
	})

	t.Run("Edge cases", func(t *testing.T) {
		t.Parallel()

		testCtx, testCancelFn := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
		defer testCancelFn()

		// case 1
		{
			promise, cancelFn := deferred.Go(func(ctx context.Context) (net.Conn, error) {
				var conn net.Conn
				return conn, nil
			})
			defer cancelFn()

			result, err := promise.Wait(testCtx)
			require.NoError(t, err)
			require.Nil(t, result)
		}

		// case 2
		{
			promise, cancelFn := deferred.Go(func(ctx context.Context) (net.Conn, error) {
				return nil, nil
			})
			defer cancelFn()

			result, err := promise.Wait(testCtx)
			require.NoError(t, err)
			require.Nil(t, result)
		}

		// case 3
		{
			promise, cancelFn := deferred.Go(func(ctx context.Context) (net.Conn, error) {
				var conn net.Conn
				return conn, fmt.Errorf("test error")
			})
			defer cancelFn()

			result, err := promise.Wait(testCtx)
			require.EqualError(t, err, "test error")
			require.Nil(t, result)
		}

		// case 4
		{
			promise, cancelFn := deferred.Go(func(ctx context.Context) (net.Conn, error) {
				return nil, fmt.Errorf("test error")
			})
			defer cancelFn()

			result, err := promise.Wait(testCtx)
			require.EqualError(t, err, "test error")
			require.Nil(t, result)
		}

		// case 5
		{
			promise, cancelFn := deferred.Go(func(ctx context.Context) (fmt.Stringer, error) {
				return &testStringer{val: "test 123"}, fmt.Errorf("test error")
			})
			defer cancelFn()

			result, err := promise.Wait(testCtx)
			require.EqualError(t, err, "test error")
			require.Nil(t, result)
		}
	})
}
