package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nicodinus/go-deferred"
	"github.com/stretchr/testify/require"
)

func TestDeferredGo(t *testing.T) {
	// Контекст теста с таймаутом
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancelFn()

	//
	testData := "test data 123"
	promise := deferred.Go(func() (*string, error) {
		return &testData, nil
	})

	result, err := promise.Wait(ctx)
	require.NoError(t, err)
	require.EqualValues(t, testData, *result)

	testErr := errors.New("test error")
	promise2 := deferred.Go(func() (*any, error) {
		return nil, testErr
	})

	result2, err := promise2.Wait(ctx)
	require.ErrorIs(t, err, testErr)
	require.Nil(t, result2)
}
