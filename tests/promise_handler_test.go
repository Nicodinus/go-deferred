package tests

import (
	"context"
	"testing"
	"time"

	"github.com/nicodinus/go-deferred"
	"github.com/stretchr/testify/require"
)

func TestPromiseHandlerBeforeResolve(t *testing.T) {
	// Контекст теста с таймаутом
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancelFn()

	//
	testData := "test data 123"
	deferred := deferred.Create[string]()

	doneCh := make(chan struct{}, 1)

	deferred.Promise().OnResolve(func(val *string, err error) {
		defer func() {
			doneCh <- struct{}{}
		}()

		require.NoError(t, err)
		require.NotNil(t, val)
		require.EqualValues(t, testData, *val)
	})

	deferred.Resolve(&testData)

	select {
	case <-ctx.Done():
		require.NoError(t, ctx.Err())
		return
	case <-doneCh:
		return
	}

}

func TestPromiseHandlerAfterResolve(t *testing.T) {
	// Контекст теста с таймаутом
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancelFn()

	//
	testData := "test data 123"
	deferred := deferred.Create[string]()

	deferred.Resolve(&testData)

	doneCh := make(chan struct{}, 1)

	deferred.Promise().OnResolve(func(val *string, err error) {
		defer func() {
			doneCh <- struct{}{}
		}()

		require.NoError(t, err)
		require.NotNil(t, val)
		require.EqualValues(t, testData, *val)
	})

	select {
	case <-ctx.Done():
		require.NoError(t, ctx.Err())
		return
	case <-doneCh:
		return
	}

}
