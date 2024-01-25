package gotest

import (
	"fmt"
)

type NestedParamContext[ParamsT any, BodyT any] struct {
	testContext TestContext
	Params ParamsT
	bodyExecutor func(TestContext, ParamsT, []BodyT)
}

type paramPack2[Param0T any, Param1T any] struct {
	params0 []Param0T
	params1 []Param1T
}

type paramPack3[Param0T any, Param1T any, Param2T any] struct {
	params0 []Param0T
	params1 []Param1T
	params2 []Param2T
}

type paramPack4[Param0T any, Param1T any, Param2T any, Param3T any] struct {
	params0 []Param0T
	params1 []Param1T
	params2 []Param2T
	params3 []Param3T
}

func ForEach2[Param0T any, Param1T any](
	context TestContext,
	params0 []Param0T,
	params1 []Param1T,
) *NestedParamContext[
	paramPack2[Param0T, Param1T],
	func(TestContext, Param0T, Param1T),
] {
	return &NestedParamContext[
		paramPack2[Param0T, Param1T],
		func(TestContext, Param0T, Param1T),
	] {
		testContext: context,
		Params: paramPack2[Param0T, Param1T] {
			params0: params0,
			params1: params1,
		},
		bodyExecutor: nestedBodyExecutor2[Param0T, Param1T],
	}
}

func nestedBodyExecutor2[Param0T any, Param1T any](
	parentContext TestContext,
	pack paramPack2[Param0T, Param1T],
	bodies []func(TestContext, Param0T, Param1T),
) {
	var param0Index int
	var param0Value Param0T
	innerContext0 := &paramTestContext {
		parentContext: parentContext,
		describeCurrentParam: func() string {
			return fmt.Sprintf(
				"[for parameter %d/%d: %s] ",
				param0Index + 1,
				len(pack.params0),
				fmt.Sprintf("%v", param0Value),
			)
		},
	}
	for param0Index, param0Value = range pack.params0 {
		var param1Index int
		var param1Value Param1T
		innerContext1 := &paramTestContext {
			parentContext: innerContext0,
			describeCurrentParam: func() string {
				return fmt.Sprintf(
					"[for parameter %d/%d: %s] ",
					param1Index + 1,
					len(pack.params1),
					fmt.Sprintf("%v", param1Value),
				)
			},
		}
		for param1Index, param1Value = range pack.params1 {
			for _, body := range bodies {
				if body != nil {
					body(innerContext1, param0Value, param1Value)
				}
			}
			innerContext1.cleanOwnHouse()
		}
		innerContext0.cleanOwnHouse()
	}
}

func ForEach3[Param0T any, Param1T any, Param2T any](
	context TestContext,
	params0 []Param0T,
	params1 []Param1T,
	params2 []Param2T,
) *NestedParamContext[
	paramPack3[Param0T, Param1T, Param2T],
	func(TestContext, Param0T, Param1T, Param2T),
] {
	return &NestedParamContext[
		paramPack3[Param0T, Param1T, Param2T],
		func(TestContext, Param0T, Param1T, Param2T),
	] {
		testContext: context,
		Params: paramPack3[Param0T, Param1T, Param2T] {
			params0: params0,
			params1: params1,
			params2: params2,
		},
		bodyExecutor: nestedBodyExecutor3[Param0T, Param1T, Param2T],
	}
}

func nestedBodyExecutor3[Param0T any, Param1T any, Param2T any](
	parentContext TestContext,
	pack paramPack3[Param0T, Param1T, Param2T],
	bodies []func(TestContext, Param0T, Param1T, Param2T),
) {
	var param0Index int
	var param0Value Param0T
	innerContext0 := &paramTestContext {
		parentContext: parentContext,
		describeCurrentParam: func() string {
			return fmt.Sprintf(
				"[for parameter %d/%d: %s] ",
				param0Index + 1,
				len(pack.params0),
				fmt.Sprintf("%v", param0Value),
			)
		},
	}
	for param0Index, param0Value = range pack.params0 {
		var param1Index int
		var param1Value Param1T
		innerContext1 := &paramTestContext {
			parentContext: innerContext0,
			describeCurrentParam: func() string {
				return fmt.Sprintf(
					"[for parameter %d/%d: %s] ",
					param1Index + 1,
					len(pack.params1),
					fmt.Sprintf("%v", param1Value),
				)
			},
		}
		for param1Index, param1Value = range pack.params1 {
			var param2Index int
			var param2Value Param2T
			innerContext2 := &paramTestContext {
				parentContext: innerContext1,
				describeCurrentParam: func() string {
					return fmt.Sprintf(
						"[for parameter %d/%d: %s] ",
						param2Index + 1,
						len(pack.params2),
						fmt.Sprintf("%v", param2Value),
					)
				},
			}
			for param2Index, param2Value = range pack.params2 {
				for _, body := range bodies {
					if body != nil {
						body(innerContext2, param0Value, param1Value, param2Value)
					}
				}
				innerContext2.cleanOwnHouse()
			}
			innerContext1.cleanOwnHouse()
		}
		innerContext0.cleanOwnHouse()
	}
}

func ForEach4[Param0T any, Param1T any, Param2T any, Param3T any](
	context TestContext,
	params0 []Param0T,
	params1 []Param1T,
	params2 []Param2T,
	params3 []Param3T,
) *NestedParamContext[
	paramPack4[Param0T, Param1T, Param2T, Param3T],
	func(TestContext, Param0T, Param1T, Param2T, Param3T),
] {
	return &NestedParamContext[
		paramPack4[Param0T, Param1T, Param2T, Param3T],
		func(TestContext, Param0T, Param1T, Param2T, Param3T),
	] {
		testContext: context,
		Params: paramPack4[Param0T, Param1T, Param2T, Param3T] {
			params0: params0,
			params1: params1,
			params2: params2,
			params3: params3,
		},
		bodyExecutor: nestedBodyExecutor4[Param0T, Param1T, Param2T, Param3T],
	}
}

func nestedBodyExecutor4[Param0T any, Param1T any, Param2T any, Param3T any](
	parentContext TestContext,
	pack paramPack4[Param0T, Param1T, Param2T, Param3T],
	bodies []func(TestContext, Param0T, Param1T, Param2T, Param3T),
) {
	var param0Index int
	var param0Value Param0T
	innerContext0 := &paramTestContext {
		parentContext: parentContext,
		describeCurrentParam: func() string {
			return fmt.Sprintf(
				"[for parameter %d/%d: %s] ",
				param0Index + 1,
				len(pack.params0),
				fmt.Sprintf("%v", param0Value),
			)
		},
	}
	for param0Index, param0Value = range pack.params0 {
		var param1Index int
		var param1Value Param1T
		innerContext1 := &paramTestContext {
			parentContext: innerContext0,
			describeCurrentParam: func() string {
				return fmt.Sprintf(
					"[for parameter %d/%d: %s] ",
					param1Index + 1,
					len(pack.params1),
					fmt.Sprintf("%v", param1Value),
				)
			},
		}
		for param1Index, param1Value = range pack.params1 {
			var param2Index int
			var param2Value Param2T
			innerContext2 := &paramTestContext {
				parentContext: innerContext1,
				describeCurrentParam: func() string {
					return fmt.Sprintf(
						"[for parameter %d/%d: %s] ",
						param2Index + 1,
						len(pack.params2),
						fmt.Sprintf("%v", param2Value),
					)
				},
			}
			for param2Index, param2Value = range pack.params2 {
				var param3Index int
				var param3Value Param3T
				innerContext3 := &paramTestContext {
					parentContext: innerContext2,
					describeCurrentParam: func() string {
						return fmt.Sprintf(
							"[for parameter %d/%d: %s] ",
							param3Index + 1,
							len(pack.params3),
							fmt.Sprintf("%v", param3Value),
						)
					},
				}
				for param3Index, param3Value = range pack.params3 {
					for _, body := range bodies {
						if body != nil {
							body(innerContext3, param0Value, param1Value, param2Value, param3Value)
						}
					}
					innerContext3.cleanOwnHouse()
				}
				innerContext2.cleanOwnHouse()
			}
			innerContext1.cleanOwnHouse()
		}
		innerContext0.cleanOwnHouse()
	}
}

func(context *NestedParamContext[ParamsT, BodyT]) Do(bodies ...BodyT) *NestedParamContext[ParamsT, BodyT] {
	if context == nil {
		return nil
	}
	context.bodyExecutor(context.testContext, context.Params, bodies)
	return context
}
