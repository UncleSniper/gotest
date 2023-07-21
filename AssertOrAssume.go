package gotest

type AssertOrAssume struct {
	Context TestContext
	Assumption bool
}

func Asserting(context TestContext) AssertOrAssume {
	if context == nil {
		panic("No test context set")
	}
	return AssertOrAssume {
		Context: context,
		Assumption: false,
	}
}

func Assuming(context TestContext) AssertOrAssume {
	if context == nil {
		panic("No test context set")
	}
	return AssertOrAssume {
		Context: context,
		Assumption: true,
	}
}

func(contextPlus AssertOrAssume) Abort(format string, args ...any) {
	if contextPlus.Context == nil {
		panic("No test context set")
	}
	if contextPlus.Assumption {
		contextPlus.Context.Skipf(format, args...)
	} else {
		contextPlus.Context.Fatalf(format, args...)
	}
}

func That[T any](contextPlus AssertOrAssume, value T) *Subject[T] {
	if contextPlus.Context == nil {
		panic("No test context set")
	}
	return &Subject[T] {
		context: contextPlus.Context,
		value: value,
		assumption: contextPlus.Assumption,
	}
}
