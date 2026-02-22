package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nicodinus/go-deferred"
	"github.com/stretchr/testify/require"
)

func BenchmarkDeferredGoByteSuccess(b *testing.B) {
	b.ReportAllocs()

	testCtx, testCancelFn := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer testCancelFn()

	var testVal byte
	var zeroVal byte
	var testErr error

	//
	for b.Loop() {
		testVal = byte(RandRange(0, 255))
		if RandRange(0, 1) == 0 {
			testErr = fmt.Errorf("test error %d", testVal)
		} else {
			testErr = nil
		}
		promise, _ := deferred.Go(func(ctx context.Context) (byte, error) {
			return testVal, testErr
		})
		result, err := promise.Wait(testCtx)
		if testErr != nil {
			require.ErrorIs(b, err, testErr)
			require.EqualValues(b, zeroVal, result)
		} else {
			require.NoError(b, err)
			require.EqualValues(b, testVal, result)
		}
	}
}

func BenchmarkDeferredGoPtrSuccess(b *testing.B) {
	b.ReportAllocs()

	testCtx, testCancelFn := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer testCancelFn()

	var testErr error
	var zeroVal *byte

	//
	for b.Loop() {
		testVal := byte(RandRange(0, 255))
		if RandRange(0, 1) == 0 {
			testErr = fmt.Errorf("test error %d", testVal)
		} else {
			testErr = nil
		}
		promise, _ := deferred.Go(func(ctx context.Context) (*byte, error) {
			return &testVal, testErr
		})
		result, err := promise.Wait(testCtx)
		if testErr != nil {
			require.ErrorIs(b, err, testErr)
			require.EqualValues(b, zeroVal, result)
		} else {
			require.NoError(b, err)
			require.EqualValues(b, testVal, *result)
		}
	}
}
