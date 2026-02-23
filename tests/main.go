package tests

type testChanHelperResult1Id int

const (
	testChanHelperResult1Undefined testChanHelperResult1Id = iota
	testChanHelperResult1ResolveGoroutine
	testChanHelperResult1PromiseOnSuccess
	testChanHelperResult1PromiseOnFail
	testChanHelperResult1PromiseOnResolve
)

type testChanHelperResult1[T any] struct {
	id  testChanHelperResult1Id
	val T
	err error
}

type testStringer struct {
	val string
}

func (s *testStringer) String() string {
	return s.val
}
