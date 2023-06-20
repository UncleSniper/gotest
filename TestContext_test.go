package gotest

import (
	"fmt"
	tst "testing"
)

func describeFakeTestMessagePtr(message *fakeTestMessage) string {
	if message == nil {
		return "nil"
	} else {
		return fmt.Sprintf("{ message: \"%s\", failure: %v }", message.message, message.failure)
	}
}

func describeFakeTestMessage(message fakeTestMessage) string {
	return describeFakeTestMessagePtr(&message)
}

func assertFakeTestMessage(subject *Subject[*fakeTestMessage], expectedMessage *fakeTestMessage) {
	subject.MaybeDescribed(describeFakeTestMessagePtr)
	if expectedMessage == nil {
		subject.Is(ZeroValue[*fakeTestMessage]())
		return
	}
	MapProperty[*fakeTestMessage, string](
		subject,
		"message",
		func(message *fakeTestMessage) string {
			return message.message
		},
		AsString,
	).Is(EqualTo(expectedMessage.message))
	MapProperty(
		subject,
		"failure",
		func(message *fakeTestMessage) bool {
			return message.failure
		},
		nil,
	).Is(EqualTo(expectedMessage.failure))
}

func assertFakeTestMessages(subject *Subject[*fakeTestContext], expectedMessages ...fakeTestMessage) {
	subject.MaybeNamed("fakeTestContext{...}")
	msgs := MapProperty(
		subject,
		"messages",
		func(context *fakeTestContext) []*fakeTestMessage {
			return context.messages
		},
		DescribeSliceByElement(describeFakeTestMessagePtr),
	)
	count := len(expectedMessages)
	msgs.Has(SliceLength[*fakeTestMessage](count))
	for index := 0; index < count; index++ {
		assertFakeTestMessage(MapSliceElement(msgs, index, describeFakeTestMessagePtr), &expectedMessages[index])
	}
}

func mapFakeTestContextToCleanups(subject *Subject[*fakeTestContext], expectedLength int) *Subject[[]func()] {
	subject.MaybeNamed("fakeTestContext{...}")
	return MapProperty(
		subject,
		"cleanups",
		func(context *fakeTestContext) []func() {
			return context.cleanups
		},
		nil,
	).Has(SliceLength[func()](expectedLength))
}

func TestFakeTestContext(t *tst.T) {
	c := Use(t)
	fake := &fakeTestContext{}
	fsubj := AssertThat(c, fake)
	// pristine
	assertFakeTestMessages(fsubj)
	mapFakeTestContextToCleanups(fsubj, 0)
	// 1 fatal
	msg0 := fakeTestMessage { "got 42 and 1337", true }
	fake.Fatalf("got %d and %d", 42, 1337)
	assertFakeTestMessages(fsubj, msg0)
	mapFakeTestContextToCleanups(fsubj, 0)
	// 1 fatal, 1 skip
	msg1 := fakeTestMessage { "now have foo and bar", false }
	fake.Skipf("now have %s and %s", "foo", "bar")
	assertFakeTestMessages(fsubj, msg0, msg1)
	mapFakeTestContextToCleanups(fsubj, 0)
	// 2 fatal, 1 skip
	msg2 := fakeTestMessage { "foo", true}
	fake.Fatalf("foo")
	assertFakeTestMessages(fsubj, msg0, msg1, msg2)
	mapFakeTestContextToCleanups(fsubj, 0)
	// 2 fatal, 1 skip, 1 cleanup
	var cleanup0CallCount int
	fake.Cleanup(func() {
		cleanup0CallCount++
	})
	assertFakeTestMessages(fsubj, msg0, msg1, msg2)
	mapFakeTestContextToCleanups(fsubj, 1)
	fake.cleanups[0]()
	AssertThat(c, cleanup0CallCount).Is(EqualTo(1))
	// 2 fatal, 2 skips, 1 cleanup
	msg3 := fakeTestMessage { "bar", false }
	fake.Skipf("bar")
	assertFakeTestMessages(fsubj, msg0, msg1, msg2, msg3)
	mapFakeTestContextToCleanups(fsubj, 1)
	fake.cleanups[0]()
	AssertThat(c, cleanup0CallCount).Is(EqualTo(2))
	// 2 fatal, 2 skips, 2 cleanups
	var cleanup1CallCount int
	fake.Cleanup(func() {
		cleanup1CallCount++
	})
	assertFakeTestMessages(fsubj, msg0, msg1, msg2, msg3)
	mapFakeTestContextToCleanups(fsubj, 2)
	fake.cleanups[0]()
	AssertThat(c, cleanup0CallCount).Is(EqualTo(3))
	AssertThat(c, cleanup1CallCount).Is(EqualTo(0))
	fake.cleanups[1]()
	AssertThat(c, cleanup0CallCount).Is(EqualTo(3))
	AssertThat(c, cleanup1CallCount).Is(EqualTo(1))
}

func TestUse(t *tst.T) {
	c := Use(t)
	AssertThat(c, Use(nil)).Is(ZeroValue[TestContext]())
	used := Use(t)
	AssertThat(c, used).Is(
		NotZeroValue[TestContext](),
	).AndIs(
		OfType[TestContext, *realTestContext](&realTestContext{}),
	)
	realCtx := used.(*realTestContext)
	MapProperty(
		AssertThat(c, realCtx),
		"context",
		func(ctx *realTestContext) *tst.T {
			return ctx.context
		},
		nil,
	).Is(EqualTo(t))
}

func TestAbort(t *tst.T) {
	c := Use(t)
	fake := &fakeTestContext{}
	fsubj := AssertThat(c, fake)
	ForEach(c, false, true).Do(func(ctx TestContext, assumption bool) {
		AssertThat(ctx, func() {
			Abort(nil, assumption)
		}).Will(
			PanicWithValue("No TestContext provided; this indicates an incorrect usage of a matcher", WithTheString),
		)
	})
	// 1 fatal
	msg0 := fakeTestMessage { "got 42 and 1337", true }
	Abort(fake, false)("got %d and %d", 42, 1337)
	assertFakeTestMessages(fsubj, msg0)
	mapFakeTestContextToCleanups(fsubj, 0)
	// 1 fatal, 1 skip
	msg1 := fakeTestMessage { "now have foo and bar", false }
	Abort(fake, true)("now have %s and %s", "foo", "bar")
	assertFakeTestMessages(fsubj, msg0, msg1)
	mapFakeTestContextToCleanups(fsubj, 0)
}

func describeCaptureTestContext(context *CaptureTestContext) string {
	return fmt.Sprintf(
		"{ Fatal: %v, MostRecentFatalMessage: \"%s\", Skip: %v, MostRecentSkipMessage: \"%s\", " +
			"cleanups: <%d elements> }",
		context.Fatal,
		context.MostRecentFatalMessage,
		context.Skip,
		context.MostRecentSkipMessage,
		len(context.cleanups),
	)
}

func assertCaptureTestContext(
	subject *Subject[*CaptureTestContext],
	fatal bool,
	fatalMessage string,
	skip bool,
	skipMessage string,
	cleanupCount int,
	assertCleanups func(*Subject[[]func()]),
) {
	subject.MaybeDescribed(describeCaptureTestContext)
	MapProperty(
		subject,
		"Fatal",
		func(context *CaptureTestContext) bool {
			return context.Fatal
		},
		nil,
	).Is(EqualTo(fatal))
	MapProperty(
		subject,
		"MostRecentFatalMessage",
		func(context *CaptureTestContext) string {
			return context.MostRecentFatalMessage
		},
		AsString,
	).Is(EqualTo(fatalMessage))
	MapProperty(
		subject,
		"Skip",
		func(context *CaptureTestContext) bool {
			return context.Skip
		},
		nil,
	).Is(EqualTo(skip))
	MapProperty(
		subject,
		"MostRecentSkipMessage",
		func(context *CaptureTestContext) string {
			return context.MostRecentSkipMessage
		},
		AsString,
	).Is(EqualTo(skipMessage))
	cleanups := MapProperty(
		subject,
		"cleanups",
		func(context *CaptureTestContext) []func() {
			return context.cleanups
		},
		nil,
	).Has(SliceLength[func()](cleanupCount))
	if assertCleanups != nil {
		assertCleanups(cleanups)
	}
}

func TestCaptureTestContextWrite(t *tst.T) {
	c := Use(t)
	capture := &CaptureTestContext{}
	csubj := AssertThat(c, capture)
	// 1 fatal
	capture.Fatalf("got %d and %d", 42, 1337)
	assertCaptureTestContext(csubj, true, "got 42 and 1337", false, "", 0, nil)
	// 1 skip
	capture = &CaptureTestContext{}
	csubj = AssertThat(c, capture)
	capture.Skipf("now have %s and %s", "foo", "bar")
	assertCaptureTestContext(csubj, false, "", true, "now have foo and bar", 0, nil)
	// 1 fatal, 1 skip
	capture = &CaptureTestContext{}
	csubj = AssertThat(c, capture)
	capture.Fatalf("myFatal")
	capture.Skipf("mySkip")
	assertCaptureTestContext(csubj, true, "myFatal", true, "mySkip", 0, nil)
	// 2 fatal, 2 skip
	capture = &CaptureTestContext{}
	csubj = AssertThat(c, capture)
	capture.Skipf("skip0")
	capture.Fatalf("fatal0")
	capture.Fatalf("fatal1")
	capture.Skipf("skip1")
	assertCaptureTestContext(csubj, true, "fatal1", true, "skip1", 0, nil)
	// cleanups
	capture = &CaptureTestContext{}
	csubj = AssertThat(c, capture)
	var cleanup0count int
	var cleanup1count int
	cleanup0 := func() {
		cleanup0count++
	}
	cleanup1 := func() {
		cleanup1count++
	}
	capture.Cleanup(cleanup0)
	capture.Cleanup(cleanup1)
	assertCaptureTestContext(
		csubj,
		false,
		"",
		false,
		"",
		2,
		func(cusubj *Subject[[]func()]) {
			cleanups := cusubj.Value()
			cleanups[0]()
			AssertThat(c, cleanup0count).Is(EqualTo(1))
			AssertThat(c, cleanup1count).Is(EqualTo(0))
			cleanups[1]()
			AssertThat(c, cleanup0count).Is(EqualTo(1))
			AssertThat(c, cleanup1count).Is(EqualTo(1))
		},
	)
}

func TestCaptureTestContextCleanOwnHouse(t *tst.T) {
	c := Use(t)
	capture := &CaptureTestContext{}
	csubj := AssertThat(c, capture)
	// call all, and in order
	var calls []int
	capture.cleanups = []func() {
		func() {
			calls = append(calls, 0)
		},
		func() {
			calls = append(calls, 1)
		},
		func() {
			calls = append(calls, 2)
		},
	}
	assertCaptureTestContext(csubj, false, "", false, "", 3, nil)
	capture.CleanOwnHouse()
	assertCaptureTestContext(csubj, false, "", false, "", 0, nil)
	cls := AssertThat(c, calls)
	cls.Has(SliceLength[int](3))
	MapSliceElement(cls, 0, nil).Is(EqualTo(0))
	MapSliceElement(cls, 1, nil).Is(EqualTo(1))
	MapSliceElement(cls, 2, nil).Is(EqualTo(2))
	// clear queue even on panic, and stop calling cleanups on panic
	calls = nil
	capture.cleanups = []func() {
		func() {
			calls = append(calls, 0)
		},
		func() {
			panic("Boom!")
		},
		func() {
			calls = append(calls, 2)
		},
	}
	AssertThat(c, capture.CleanOwnHouse).Named("CaptureTestContext.CleanOwnHouse()").Will(
		PanicWithValue("Boom!", nil),
	)
	assertCaptureTestContext(csubj, false, "", false, "", 0, nil)
	cls = AssertThat(c, calls)
	cls.Has(SliceLength[int](1))
	MapSliceElement(cls, 0, nil).Is(EqualTo(0))
}

type captureTestContextGetRecordedFailureParam struct {
	contextFatal bool
	contextFatalMessage string
	contextSkip bool
	contextSkipMessage string
	wasAssumption bool
	expectedFailed bool
	expectedMessage string
}

func describeCaptureTestContextGetRecordedFailureParam(param *captureTestContextGetRecordedFailureParam) string {
	return fmt.Sprintf(
		"{ %s, %s, %s, %s, %s, %s, %s }",
		fmt.Sprintf("contextFatal: %v", param.contextFatal),
		fmt.Sprintf("contextFatalMessage: \"%s\"", param.contextFatalMessage),
		fmt.Sprintf("contextSkip: %v", param.contextSkip),
		fmt.Sprintf("contextSkipMessage: \"%s\"", param.contextSkipMessage),
		fmt.Sprintf("wasAssumption: %v", param.wasAssumption),
		fmt.Sprintf("expectedFailed: %v", param.expectedFailed),
		fmt.Sprintf("expectedMessage: \"%s\"", param.expectedMessage),
	)
}

func TestCaptureTestContextGetRecordedFailure(t *tst.T) {
	c := Use(t)
	ForEach(
		c,
		&captureTestContextGetRecordedFailureParam {
			contextFatal: false,
			contextFatalMessage: "notFatal",
			contextSkip: false,
			contextSkipMessage: "notSkip",
			wasAssumption: false,
			expectedFailed: false,
			expectedMessage: "",
		},
		&captureTestContextGetRecordedFailureParam {
			contextFatal: false,
			contextFatalMessage: "notFatal",
			contextSkip: false,
			contextSkipMessage: "notSkip",
			wasAssumption: true,
			expectedFailed: false,
			expectedMessage: "",
		},
		&captureTestContextGetRecordedFailureParam {
			contextFatal: true,
			contextFatalMessage: "fatalMsg",
			contextSkip: false,
			contextSkipMessage: "notSkip",
			wasAssumption: false,
			expectedFailed: true,
			expectedMessage: "fatalMsg",
		},
		&captureTestContextGetRecordedFailureParam {
			contextFatal: true,
			contextFatalMessage: "fatalMsg",
			contextSkip: false,
			contextSkipMessage: "notSkip",
			wasAssumption: true,
			expectedFailed: false,
			expectedMessage: "",
		},
		&captureTestContextGetRecordedFailureParam {
			contextFatal: false,
			contextFatalMessage: "notFatal",
			contextSkip: true,
			contextSkipMessage: "skipMsg",
			wasAssumption: false,
			expectedFailed: false,
			expectedMessage: "",
		},
		&captureTestContextGetRecordedFailureParam {
			contextFatal: false,
			contextFatalMessage: "notFatal",
			contextSkip: true,
			contextSkipMessage: "skipMsg",
			wasAssumption: true,
			expectedFailed: true,
			expectedMessage: "skipMsg",
		},
		&captureTestContextGetRecordedFailureParam {
			contextFatal: true,
			contextFatalMessage: "fatalMsg",
			contextSkip: true,
			contextSkipMessage: "skipMsg",
			wasAssumption: false,
			expectedFailed: true,
			expectedMessage: "fatalMsg",
		},
		&captureTestContextGetRecordedFailureParam {
			contextFatal: true,
			contextFatalMessage: "fatalMsg",
			contextSkip: true,
			contextSkipMessage: "skipMsg",
			wasAssumption: true,
			expectedFailed: true,
			expectedMessage: "skipMsg",
		},
	).Described(describeCaptureTestContextGetRecordedFailureParam).Do(
		func(ctx TestContext, param *captureTestContextGetRecordedFailureParam) {
			capture := &CaptureTestContext {
				Fatal: param.contextFatal,
				MostRecentFatalMessage: param.contextFatalMessage,
				Skip: param.contextSkip,
				MostRecentSkipMessage: param.contextSkipMessage,
			}
			didFail, failMessage := capture.GetRecordedFailure(param.wasAssumption)
			AssertThat(ctx, didFail).Is(EqualTo(param.expectedFailed))
			AssertThat(ctx, failMessage).Is(EqualTo(param.expectedMessage))
		},
	)
}
