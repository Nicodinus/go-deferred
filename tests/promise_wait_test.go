package tests

import (
	"context"
	"testing"
	"time"

	"github.com/nicodinus/go-deferred"

	"github.com/stretchr/testify/require"
)

func TestWaitBeforePromiseResolve(t *testing.T) {
	// Контекст теста с таймаутом
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancelFn()

	//
	testData := "test data 123"
	d := deferred.Create[string]()

	type testChResult struct {
		val *string
		err error
	}
	testCh := make(chan testChResult, 1)
	go func() {
		val, err := d.Promise().Wait(ctx)
		testCh <- testChResult{
			val: val,
			err: err,
		}
		close(testCh)
	}()

	time.Sleep(100 * time.Millisecond)

	err := d.Resolve(&testData)
	require.NoError(t, err)

	result := <-testCh
	require.NoError(t, result.err)
	require.NotNil(t, result.val)
	require.EqualValues(t, testData, *result.val)
}

func TestWaitAfterPromiseResolve(t *testing.T) {
	// Контекст теста с таймаутом
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancelFn()

	//
	testData := "test data 123"
	d := deferred.Create[string]()

	err := d.Resolve(&testData)
	require.NoError(t, err)

	val, err := d.Promise().Wait(ctx)

	require.NoError(t, err)
	require.NotNil(t, val)
	require.EqualValues(t, testData, *val)
}

func TestPromiseDoubleResolve(t *testing.T) {
	// Контекст теста с таймаутом
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancelFn()

	//
	testData := "test data 123"
	d := deferred.Create[string]()

	err := d.Resolve(&testData)
	require.NoError(t, err)

	err = d.Resolve(&testData)
	require.ErrorIs(t, deferred.ErrPromiseResolved, err)

	val, err := d.Promise().Wait(ctx)

	require.NoError(t, err)
	require.NotNil(t, val)
	require.EqualValues(t, testData, *val)
}

func TestPromiseCancelAfterWait(t *testing.T) {
	// Контекст теста с таймаутом
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancelFn()

	//
	d := deferred.Create[string]()

	type testChResult struct {
		val *string
		err error
	}
	testCh := make(chan testChResult, 1)
	go func() {
		val, err := d.Promise().Wait(ctx)
		testCh <- testChResult{
			val: val,
			err: err,
		}
		close(testCh)
	}()

	time.Sleep(100 * time.Millisecond)

	err := d.Cancel()
	require.NoError(t, err)

	result := <-testCh
	require.ErrorIs(t, deferred.ErrPromiseCancelled, result.err)
	require.Nil(t, result.val)
}

func TestPromiseCancelBeforeWait(t *testing.T) {
	// Контекст теста с таймаутом
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancelFn()

	//
	d := deferred.Create[string]()

	err := d.Cancel()
	require.NoError(t, err)

	val, err := d.Promise().Wait(ctx)

	require.ErrorIs(t, deferred.ErrPromiseCancelled, err)
	require.Nil(t, val)
}

func TestWaitContextTimeoutBeforePromiseResolve(t *testing.T) {
	// Контекст теста с таймаутом
	ctx, cancelFn := context.WithTimeout(context.Background(), 0)
	defer cancelFn()

	//
	testData := "test data 123"
	d := deferred.Create[string]()

	type testChResult struct {
		val *string
		err error
	}
	testCh := make(chan testChResult, 1)
	go func() {
		val, err := d.Promise().Wait(ctx)
		testCh <- testChResult{
			val: val,
			err: err,
		}
		close(testCh)
	}()

	time.Sleep(100 * time.Millisecond)

	err := d.Resolve(&testData)
	require.NoError(t, err)

	result := <-testCh
	require.ErrorIs(t, context.DeadlineExceeded, result.err)
	require.Nil(t, result.val)
}

func TestWaitContextTimeoutAfterPromiseResolve(t *testing.T) {
	// Контекст теста с отменой
	ctx, cancelFn := context.WithCancel(context.Background())
	cancelFn()

	//
	testData := "test data 123"
	d := deferred.Create[string]()

	err := d.Resolve(&testData)
	require.NoError(t, err)

	val, err := d.Promise().Wait(ctx)
	require.NoError(t, err)
	require.NotNil(t, val)
	require.EqualValues(t, testData, *val)
}

func TestWaitContextCancelBeforePromiseResolve(t *testing.T) {
	// Контекст теста с отменой
	ctx, cancelFn := context.WithCancel(context.Background())
	cancelFn()

	//
	testData := "test data 123"
	d := deferred.Create[string]()

	type testChResult struct {
		val *string
		err error
	}
	testCh := make(chan testChResult, 1)
	go func() {
		val, err := d.Promise().Wait(ctx)
		testCh <- testChResult{
			val: val,
			err: err,
		}
		close(testCh)
	}()

	time.Sleep(100 * time.Millisecond)

	err := d.Resolve(&testData)
	require.NoError(t, err)

	result := <-testCh
	require.ErrorIs(t, context.Canceled, result.err)
	require.Nil(t, result.val)
}

func TestWaitContextCancelAfterPromiseResolve(t *testing.T) {
	// Контекст теста с отменой
	ctx, cancelFn := context.WithCancel(context.Background())
	cancelFn()

	//
	testData := "test data 123"
	d := deferred.Create[string]()

	err := d.Resolve(&testData)
	require.NoError(t, err)

	val, err := d.Promise().Wait(ctx)
	require.NoError(t, err)
	require.NotNil(t, val)
	require.EqualValues(t, testData, *val)
}

func TestMultipleWaitsBeforePromiseResolve(t *testing.T) {
	// Контекст теста с таймаутом
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancelFn()

	//
	testData := "test data 123"
	const numWaits = 10
	d := deferred.Create[string]()

	type testChResult struct {
		val *string
		err error
	}
	testCh := make(chan testChResult, numWaits)
	defer close(testCh)
	for range numWaits {
		go func() {
			val, err := d.Promise().Wait(ctx)
			testCh <- testChResult{
				val: val,
				err: err,
			}
		}()
	}

	time.Sleep(100 * time.Millisecond)

	err := d.Resolve(&testData)
	require.NoError(t, err)

	for range numWaits {
		result := <-testCh
		require.NoError(t, result.err)
		require.NotNil(t, result.val)
		require.EqualValues(t, testData, *result.val)
	}
}

func TestMultipleWaitsAfterPromiseResolve(t *testing.T) {
	// Контекст теста с таймаутом
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancelFn()

	//
	testData := "test data 123"
	const numWaits = 10
	d := deferred.Create[string]()

	err := d.Resolve(&testData)
	require.NoError(t, err)

	for range numWaits {
		val, err := d.Promise().Wait(ctx)
		require.NoError(t, err)
		require.NotNil(t, val)
		require.EqualValues(t, testData, *val)
	}
}
