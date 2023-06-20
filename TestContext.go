package gotest

import (
	"fmt"
	tst "testing"
	"runtime/debug"
)

type TestContext interface {
	Fatalf(format string, args ...any)
	Skipf(format string, arg ...any)
	Cleanup(callback func())
}

type realTestContext struct {
	context *tst.T
}

func(context *realTestContext) Fatalf(format string, args ...any) {
	context.context.Logf(format, args...)
	context.context.Logf("This failure occurred at: %s", string(debug.Stack()))
	context.context.FailNow()
}

func(context *realTestContext) Skipf(format string, args ...any) {
	context.context.Logf(format, args...)
	context.context.Logf("This skip occurred at: %s", string(debug.Stack()))
	context.context.SkipNow()
}

func(context *realTestContext) Cleanup(callback func()) {
	context.context.Cleanup(callback)
}

type fakeTestMessage struct {
	message string
	failure bool
}

type fakeTestContext struct {
	messages []*fakeTestMessage
	cleanups []func()
}

func(context *fakeTestContext) Fatalf(format string, args ...any) {
	context.messages = append(
		context.messages,
		&fakeTestMessage {
			message: fmt.Sprintf(format, args...),
			failure: true,
		},
	)
}

func(context *fakeTestContext) Skipf(format string, args ...any) {
	context.messages = append(
		context.messages,
		&fakeTestMessage {
			message: fmt.Sprintf(format, args...),
			failure: false,
		},
	)
}

func(context *fakeTestContext) Cleanup(callback func()) {
	context.cleanups = append(context.cleanups, callback)
}

func Use(context *tst.T) TestContext {
	if context == nil {
		return nil
	}
	return &realTestContext {
		context: context,
	}
}

func Abort(context TestContext, assumption bool) func(string, ...any) {
	if context == nil {
		panic("No TestContext provided; this indicates an incorrect usage of a matcher")
	}
	if assumption {
		return context.Skipf
	} else {
		return context.Fatalf
	}
}

type CaptureTestContext struct {
	Fatal bool
	MostRecentFatalMessage string
	Skip bool
	MostRecentSkipMessage string
	cleanups []func()
}

func(context *CaptureTestContext) Fatalf(format string, args ...any) {
	if context != nil {
		context.Fatal = true
		context.MostRecentFatalMessage = fmt.Sprintf(format, args...)
	}
}

func(context *CaptureTestContext) Skipf(format string, args ...any) {
	if context != nil {
		context.Skip = true
		context.MostRecentSkipMessage = fmt.Sprintf(format, args...)
	}
}

func(context *CaptureTestContext) Cleanup(callback func()) {
	if callback != nil {
		context.cleanups = append(context.cleanups, callback)
	}
}

func(context *CaptureTestContext) CleanOwnHouse() {
	defer func() {
		context.cleanups = nil
	}()
	for _, callback := range context.cleanups {
		callback()
	}
}

func(context *CaptureTestContext) GetRecordedFailure(wasAssumption bool) (bool, string) {
	if wasAssumption {
		if context.Skip {
			return true, context.MostRecentSkipMessage
		} else {
			return false, ""
		}
	} else {
		if context.Fatal {
			return true, context.MostRecentFatalMessage
		} else {
			return false, ""
		}
	}
}
