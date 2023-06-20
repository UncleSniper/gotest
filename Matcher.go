package gotest

type Matcher[T any] interface {
	Match(context TestContext, assumption bool, value T, descriptor Descriptor[T])
}

type FuncMatcher[T any] struct {
	Func func(context TestContext, assumption bool, value T, descriptor Descriptor[T])
}

func(matcher FuncMatcher[T]) Match(context TestContext, assumption bool, value T, descriptor Descriptor[T]) {
	if matcher.Func == nil {
		panic("No matcher callback provided")
	} else {
		matcher.Func(context, assumption, value, descriptor)
	}
}

func doDescribe[T any](value T, descriptor Descriptor[T]) string {
	if descriptor == nil {
		panic("No descriptor provided; this indicates an incorrect usage of a matcher")
	} else {
		return descriptor(value)
	}
}
