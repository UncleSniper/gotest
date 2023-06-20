package gotest

import (
	"fmt"
	tst "testing"
)

func TestFuncMatcherMatch(t *tst.T) {
	c := Use(t)
	var gaveContext TestContext = &fakeTestContext{}
	var gotContext TestContext
	var gotAssumption bool
	var gotValue int
	var descriptorArg int
	var matcherCallCount int
	var descriptorCallCount int
	funcMatcher := FuncMatcher[int] {
		Func: func(context TestContext, assumption bool, value int, descriptor Descriptor[int]) {
			matcherCallCount++
			gotContext = context
			gotAssumption = assumption
			gotValue = value
			descriptor(1337)
		},
	}
	funcMatcher.Match(
		gaveContext,
		true,
		42,
		func(describee int) string {
			descriptorCallCount++
			descriptorArg = describee
			return fmt.Sprintf("describing %d", describee)
		},
	)
	AssertThat(c, gotContext).Is(EqualTo(gaveContext))
	AssertThat(c, gotAssumption).Is(EqualTo(true))
	AssertThat(c, gotValue).Is(EqualTo(42))
	AssertThat(c, descriptorArg).Is(EqualTo(1337))
	AssertThat(c, matcherCallCount).Is(EqualTo(1))
	AssertThat(c, descriptorCallCount).Is(EqualTo(1))
}

func TestFuncMatcherMatchWithNilCallback(t *tst.T) {
	c := Use(t)
	var gaveContext TestContext = &fakeTestContext{}
	funcMatcher := FuncMatcher[int] {}
	descriptor := func(int) string {
		return ""
	}
	AssertThat(c, func() {
		funcMatcher.Match(gaveContext, false, 42, descriptor)
	}).Will(PanicWithValue("No matcher callback provided", WithTheString))
}

func TestDoDescribe(t *tst.T) {
	c := Use(t)
	AssertThat(c, func() {
		doDescribe(0, nil)
	}).Will(PanicWithValue("No descriptor provided; this indicates an incorrect usage of a matcher", WithTheString))
	var descriptorCallCount int
	descriptor := func(value int) string {
		descriptorCallCount++
		return fmt.Sprintf("the int %d", value)
	}
	AssertThat(c, doDescribe(42, descriptor)).Is(EqualTo("the int 42"))
	AssertThat(c, descriptorCallCount).Is(EqualTo(1))
}
