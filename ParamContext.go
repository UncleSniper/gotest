package gotest

import (
	"fmt"
)

type ParamContext[ParamT any] struct {
	testContext TestContext
	Params []ParamT
	paramDescriptor Descriptor[ParamT]
}

func ForEach[ParamT any](context TestContext, params ...ParamT) *ParamContext[ParamT] {
	if context == nil {
		panic("No test context provided")
	}
	return &ParamContext[ParamT] {
		testContext: context,
		Params: params,
	}
}

func(context *ParamContext[ParamT]) AndAlso(moreParams ...ParamT) *ParamContext[ParamT] {
	if context != nil {
		context.Params = append(context.Params, moreParams...)
	}
	return context
}

func(context *ParamContext[ParamT]) Described(descriptor Descriptor[ParamT]) *ParamContext[ParamT] {
	if context != nil {
		context.paramDescriptor = descriptor
	}
	return context
}

type ParameterizedTestCase[ParamT any] func(TestContext, ParamT)

type paramTestContext struct {
	parentContext TestContext
	describeCurrentParam func() string
	cleanups []func()
}

func(context *paramTestContext) Fatalf(format string, args ...any) {
	innerMessage := fmt.Sprintf(format, args...)
	context.parentContext.Fatalf("%s%s", context.describeCurrentParam(), innerMessage)
}

func(context *paramTestContext) Skipf(format string, args ...any) {
	innerMessage := fmt.Sprintf(format, args...)
	context.parentContext.Skipf("%s%s", context.describeCurrentParam(), innerMessage)
}

func(context *paramTestContext) Cleanup(callback func()) {
	if callback != nil {
		context.cleanups = append(context.cleanups, callback)
	}
}

func(context *paramTestContext) cleanOwnHouse() {
	defer func() {
		context.cleanups = nil
	}()
	for _, callback := range context.cleanups {
		callback()
	}
}

func(context *ParamContext[ParamT]) Do(bodies ...ParameterizedTestCase[ParamT]) *ParamContext[ParamT] {
	if context == nil || len(context.Params) == 0 {
		return context
	}
	var paramIndex int
	var paramValue ParamT
	innerContext := &paramTestContext {
		parentContext: context.testContext,
		describeCurrentParam: func() string {
			var valueDescr string
			if context.paramDescriptor == nil {
				valueDescr = fmt.Sprintf("%v", paramValue)
			} else {
				valueDescr = context.paramDescriptor(paramValue)
			}
			return fmt.Sprintf("[for parameter %d/%d: %s] ", paramIndex + 1, len(context.Params), valueDescr)
		},
	}
	for paramIndex, paramValue = range context.Params {
		for _, body := range bodies {
			if body != nil {
				body(innerContext, paramValue)
			}
		}
		innerContext.cleanOwnHouse()
	}
	return context
}
