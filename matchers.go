package gotest

import (
	"fmt"
	"strings"
	"strconv"
	"reflect"
	"golang.org/x/exp/constraints"
)

func makeComparison[T any](expected T, satisfied string, comparator func(T, T) bool) Matcher[T] {
	return FuncMatcher[T] {
		Func: func(context TestContext, assumption bool, actual T, descriptor Descriptor[T]) {
			if !comparator(actual, expected) {
				Abort(context, assumption)(
					"Expected %s to be %s %s, but was not",
					doDescribe(actual, descriptor),
					satisfied,
					doDescribe(expected, descriptor),
				)
			}
		},
	}
}

func Equals[T comparable](expected T) Matcher[T] {
	return makeComparison(
		expected,
		"equal to",
		func(actual T, expected T) bool {
			return actual == expected
		},
	)
}

func EqualTo[T comparable](expected T) Matcher[T] {
	return Equals(expected)
}

func Unequals[T comparable](unexpected T) Matcher[T] {
	return makeComparison(
		unexpected,
		"unequal to",
		func(actual T, expected T) bool {
			return actual != expected
		},
	)
}

func UnequalTo[T comparable](unexpected T) Matcher[T] {
	return Unequals(unexpected)
}

func Less[T constraints.Ordered](bound T) Matcher[T] {
	return makeComparison(
		bound,
		"less than",
		func(actual T, bound T) bool {
			return actual < bound
		},
	)
}

func LessThan[T constraints.Ordered](bound T) Matcher[T] {
	return Less(bound)
}

func Greater[T constraints.Ordered](bound T) Matcher[T] {
	return makeComparison(
		bound,
		"greater than",
		func(actual T, bound T) bool {
			return actual > bound
		},
	)
}

func GreaterThan[T constraints.Ordered](bound T) Matcher[T] {
	return Greater(bound)
}

func LessOrEqual[T constraints.Ordered](bound T) Matcher[T] {
	return makeComparison(
		bound,
		"less than or equal to",
		func(actual T, bound T) bool {
			return actual <= bound
		},
	)
}

func LessThanOrEqualTo[T constraints.Ordered](bound T) Matcher[T] {
	return LessOrEqual(bound)
}

func GreaterOrEqual[T constraints.Ordered](bound T) Matcher[T] {
	return makeComparison(
		bound,
		"greater than or equal to",
		func(actual T, bound T) bool {
			return actual >= bound
		},
	)
}

func GreaterThanOrEqualTo[T constraints.Ordered](bound T) Matcher[T] {
	return GreaterOrEqual(bound)
}

type indexSet map[int]bool

func containsIndex(seen indexSet, index int) bool {
	_, ok := seen[index]
	return ok
}

func describeIndexSet(set indexSet) string {
	var builder strings.Builder
	var err error
	for index, _ := range set {
		if builder.Len() > 0 {
			_, err = builder.WriteString(", ")
			if err != nil {
				break
			}
		}
		_, err = builder.WriteString(strconv.FormatInt(int64(index), 10))
		if err != nil {
			break
		}
	}
	if err != nil {
		panic("Failed to stringify index list: " + err.Error())
	}
	return builder.String()
}

type CollectionMatcher[ElementT any, CollectionT any] interface {
	Matcher[CollectionT]
	ElementDescribed(descriptor Descriptor[ElementT]) CollectionMatcher[ElementT, CollectionT]
}

type CollectionAllMatcher[ElementT any, CollectionT any] interface {
	CollectionMatcher[ElementT, CollectionT]
	Distinct() CollectionAllMatcher[ElementT, CollectionT]
}

type containsMatcher[ElementT comparable] struct {
	expectedElements []ElementT
	all bool
	distinct bool
	elementDescriptor Descriptor[ElementT]
	multipleNeedles bool
}

func(matcher *containsMatcher[ElementT]) describeActual(
	actual []ElementT,
	sliceDescriptor Descriptor[[]ElementT],
) string {
	if sliceDescriptor == nil {
		sliceDescriptor = DescribeSliceByElement(matcher.elementDescriptor)
	}
	return sliceDescriptor(actual)
}

func(matcher *containsMatcher[ElementT]) describeElement(element ElementT) string {
	if matcher.elementDescriptor == nil {
		return fmt.Sprintf("%v", element)
	} else {
		return matcher.elementDescriptor(element)
	}
}

func(matcher *containsMatcher[ElementT]) describeExpected() string {
	var builder strings.Builder
	var err error
	_, err = builder.WriteRune('[')
	if err == nil {
		for _, expected := range matcher.expectedElements {
			if builder.Len() > 0 {
				_, err = builder.WriteString(", ")
				if err != nil {
					break
				}
			}
			_, err = builder.WriteString(matcher.describeElement(expected))
			if err != nil {
				break
			}
		}
		if err == nil {
			_, err = builder.WriteRune(']')
		}
	}
	if err != nil {
		panic("Failed to stringify element list: " + err.Error())
	}
	return builder.String()
}

func(matcher *containsMatcher[ElementT]) Match(
	context TestContext,
	assumption bool,
	actual []ElementT,
	sliceDescriptor Descriptor[[]ElementT],
) {
	found := make(map[ElementT]indexSet)
	expCount := make(map[ElementT]uint)
	for expIndex, expected := range matcher.expectedElements {
		expCount[expected] = expCount[expected] + 1
		setForExp := found[expected]
		if setForExp == nil {
			setForExp = make(indexSet)
			found[expected] = setForExp
		}
		foundIndex := -1
		for actIndex, actElement := range actual {
			if actElement == expected {
				if !matcher.all {
					return
				}
				if containsIndex(setForExp, actIndex) {
					// previously seen => just keep looking
					if !matcher.distinct {
						foundIndex = actIndex
					}
				} else {
					// never seen
					setForExp[actIndex] = true
					foundIndex = actIndex
				}
			}
		}
		if matcher.all && foundIndex < 0 {
			// need all, but didn't find anything (or, if distinct, anything new)
			var needleIndex string
			if matcher.multipleNeedles {
				needleIndex = fmt.Sprintf(" (index %d in needles)", expIndex)
			}
			var mismatchReason string
			if matcher.distinct && len(setForExp) > 0 {
				mismatchReason = fmt.Sprintf(
					"did not contain it again after previous matches at indices %s",
					describeIndexSet(setForExp),
				)
			} else {
				mismatchReason = "did not contain it"
			}
			Abort(context, assumption)(
				"Expected %s to contain %s%s, but %s",
				matcher.describeActual(actual, sliceDescriptor),
				matcher.describeElement(expected),
				needleIndex,
				mismatchReason,
			)
			return
		}
	}
	if !matcher.all {
		if matcher.multipleNeedles {
			Abort(context, assumption)(
				"Expected %s to contain at least one of %s, but did not contain any",
				matcher.describeActual(actual, sliceDescriptor),
				matcher.describeExpected(),
			)
		} else {
			Abort(context, assumption)(
				"Expected %s to contain %s, but did not contain it",
				matcher.describeActual(actual, sliceDescriptor),
				matcher.describeExpected,
			)
		}
	}
}

func(matcher *containsMatcher[ElementT]) ElementDescribed(
	descriptor Descriptor[ElementT],
) CollectionMatcher[ElementT, []ElementT] {
	if matcher != nil {
		matcher.elementDescriptor = descriptor
	}
	return matcher
}

func(matcher *containsMatcher[ElementT]) Distinct() CollectionAllMatcher[ElementT, []ElementT] {
	if matcher != nil {
		matcher.distinct = true
	}
	return matcher
}

func ContainAll[ElementT comparable](expectedElements []ElementT) CollectionAllMatcher[ElementT, []ElementT] {
	return &containsMatcher[ElementT] {
		expectedElements: expectedElements,
		all: true,
		multipleNeedles: true,
	}
}

func Contain[ElementT comparable](expectedElement ElementT) CollectionAllMatcher[ElementT, []ElementT] {
	return &containsMatcher[ElementT] {
		expectedElements: []ElementT {expectedElement},
		all: true,
		multipleNeedles: false,
	}
}

func ContainAny[ElementT comparable](expectedElements []ElementT) CollectionMatcher[ElementT, []ElementT] {
	return &containsMatcher[ElementT] {
		expectedElements: expectedElements,
		all: false,
		multipleNeedles: true,
	}
}

func OfType[KnownT any, ExpectedT any](expectedInstance ExpectedT) Matcher[KnownT] {
	return FuncMatcher[KnownT] {
		Func: func(context TestContext, assumption bool, actual KnownT, descriptor Descriptor[KnownT]) {
			var anyActual any = actual
			_, ok := anyActual.(ExpectedT)
			if !ok {
				Abort(context, assumption)(
					"Expected %s to be of type %s, but was of type %s",
					doDescribe(actual, descriptor),
					reflect.TypeOf(expectedInstance).String(),
					reflect.TypeOf(actual).String(),
				)
			}
		},
	}
}

func ZeroValue[T comparable]() Matcher[T] {
	return FuncMatcher[T] {
		Func: func(context TestContext, assumption bool, actual T, descriptor Descriptor[T]) {
			var zero T
			if actual != zero {
				var zeroValue reflect.Value
				expectedValue := reflect.ValueOf(zero)
				var rend string
				if expectedValue == zeroValue {
					rend = "nil"
				} else {
					rend = expectedValue.String()
				}
				Abort(context, assumption)(
					"Expected %s to be zero value of its type, namely %s, but was not",
					doDescribe(actual, descriptor),
					rend,
				)
			}
		},
	}
}

func NotZeroValue[T comparable]() Matcher[T] {
	return FuncMatcher[T] {
		Func: func(context TestContext, assumption bool, actual T, descriptor Descriptor[T]) {
			var zero T
			if actual == zero {
				Abort(context, assumption)(
					"Expected %s to be distinct from the zero value of its type, but was not",
					doDescribe(actual, descriptor),
				)
			}
		},
	}
}

func Panic() Matcher[func()] {
	return FuncMatcher[func()] {
		Func: func(context TestContext, assumption bool, body func(), descriptor Descriptor[func()]) {
			panicWithPredicate[any](
				nil,
				func(any) bool {
					return true
				},
				func(any) string {
					return "<Boom, boom, shake the room, say whaaaaaaat!?>"
				},
				func(any) string {
					return "with any value"
				},
				false,
			)
		},
	}
}

func PanicWithValue[ExpectedT comparable](
	expectedInstance ExpectedT,
	panicDescriptor Descriptor[ExpectedT],
) Matcher[func()] {
	return panicWithValue(expectedInstance, panicDescriptor, false)
}

func PanicWithValueOrNotAtAll[ExpectedT comparable](
	expectedInstance ExpectedT,
	panicDescriptor Descriptor[ExpectedT],
) Matcher[func()] {
	return panicWithValue(expectedInstance, panicDescriptor, true)
}

func panicWithValue[ExpectedT comparable](
	expectedInstance ExpectedT,
	panicDescriptor Descriptor[ExpectedT],
	noPanicIsOK bool,
) Matcher[func()] {
	return panicWithPredicate[ExpectedT](
		expectedInstance,
		func(gotValue any) bool {
			gotOfCorrectType, ok := gotValue.(ExpectedT)
			return ok && gotOfCorrectType == expectedInstance
		},
		func(unexpected any) string {
			gotOfCorrectType, ok := unexpected.(ExpectedT)
			if panicDescriptor == nil || !ok {
				return fmt.Sprintf("with %v", unexpected)
			} else {
				return panicDescriptor(gotOfCorrectType)
			}
		},
		func(expected ExpectedT) string {
			if panicDescriptor == nil {
				return fmt.Sprintf("with type %s and value %v", reflect.TypeOf(expected).String(), expected)
			} else {
				return panicDescriptor(expected)
			}
		},
		noPanicIsOK,
	)
}

func PanicWithPattern[ExpectedT any](
	expectedInstance ExpectedT,
	panicDescriptor Descriptor[ExpectedT],
	pattern Matcher[ExpectedT],
) Matcher[func()] {
	return panicWithPattern(expectedInstance, pattern, panicDescriptor, false)
}

func PanicWithPatternOrNotAtAll[ExpectedT any](
	expectedInstance ExpectedT,
	panicDescriptor Descriptor[ExpectedT],
	pattern Matcher[ExpectedT],
) Matcher[func()] {
	return panicWithPattern(expectedInstance, pattern, panicDescriptor, true)
}

func panicWithPattern[ExpectedT any](
	expectedInstance ExpectedT,
	pattern Matcher[ExpectedT],
	panicDescriptor Descriptor[ExpectedT],
	noPanicIsOK bool,
) Matcher[func()] {
	if pattern == nil {
		panic("No panic value pattern matcher provided")
	}
	return FuncMatcher[func()] {
		Func: func(context TestContext, assumption bool, body func(), descriptor Descriptor[func()]) {
			capture := &CaptureTestContext {}
			defer capture.CleanOwnHouse()
			var matchFailMessage string
			panicWithPredicate[ExpectedT](
				expectedInstance,
				func(gotValue any) bool {
					gotOfCorrectType, ok := gotValue.(ExpectedT)
					if !ok {
						return false
					}
					pattern.Match(
						capture,
						assumption,
						gotOfCorrectType,
						func(gotInner ExpectedT) string {
							if panicDescriptor == nil {
								return fmt.Sprintf("%v", gotInner)
							} else {
								return panicDescriptor(gotInner)
							}
						},
					)
					var didFail bool
					didFail, matchFailMessage = capture.GetRecordedFailure(assumption)
					return !didFail
				},
				func(unexpected any) string {
					if len(matchFailMessage) == 0 {
						return "with non-matching value"
					} else {
						return "with non-matching value: " + matchFailMessage
					}
				},
				func(expected ExpectedT) string {
					if panicDescriptor == nil {
						return fmt.Sprintf("with type %s", reflect.TypeOf(expected).String())
					} else {
						return panicDescriptor(expected)
					}
				},
				noPanicIsOK,
			)
		},
	}
}

func panicWithPredicate[ExpectedT any](
	expectedInstance ExpectedT,
	isExpected func(any) bool,
	describeUnexpected Descriptor[any],
	expectedPanicDescr Descriptor[ExpectedT],
	noPanicIsOK bool,
) Matcher[func()] {
	return FuncMatcher[func()] {
		Func: func(
			context TestContext,
			assumption bool,
			body func(),
			descriptor Descriptor[func()],
		) {
			if body == nil {
				panic("Cannot determine whether nil body panics")
			}
			if context == nil {
				panic("No test context provided")
			}
			var didNotPanic bool
			var panicValue any
			defer func() {
				panicValue = recover()
				var reason string
				if expectedPanicDescr != nil {
					reason = expectedPanicDescr(expectedInstance)
					if len(reason) > 0 {
						reason = " " + reason
					}
				}
				if didNotPanic {
					if !noPanicIsOK {
						Abort(context, assumption)(
							"Expected %s to panic%s, but did not panic",
							doDescribe(body, descriptor),
							reason,
						)
					}
				} else if !isExpected(panicValue) {
					var unexpectedDescr string
					if describeUnexpected != nil {
						unexpectedDescr = describeUnexpected(panicValue)
					} else {
						unexpectedDescr = fmt.Sprintf("with %v", panicValue)
					}
					Abort(context, assumption)(
						"Expected %s to panic%s, but panicked %s",
						doDescribe(body, descriptor),
						reason,
						unexpectedDescr,
					)
				}
			}()
			body()
			didNotPanic = true
		},
	}
}

func StringLength(expectedLength int) Matcher[string] {
	return FuncMatcher[string] {
		Func: func(context TestContext, assumption bool, actual string, descriptor Descriptor[string]) {
			if len(actual) != expectedLength {
				Abort(context, assumption)(
					"Expected %s to have length %d, but had length %d",
					doDescribe(actual, descriptor),
					expectedLength,
					len(actual),
				)
			}
		},
	}
}

func StringLengthWhichIs(lengthPattern Matcher[int]) Matcher[string] {
	if lengthPattern == nil {
		panic("No string length pattern provided")
	}
	return FuncMatcher[string] {
		Func: func(context TestContext, assumption bool, actual string, descriptor Descriptor[string]) {
			capture := &CaptureTestContext {}
			defer capture.CleanOwnHouse()
			lengthPattern.Match(
				capture,
				assumption,
				len(actual),
				func(gotLength int) string {
					return fmt.Sprintf("length %d", gotLength)
				},
			)
			didFail, matchFailMessage := capture.GetRecordedFailure(assumption)
			if didFail {
				Abort(context, assumption)(
					"Expected %s to have matching length, but had length %d: %s",
					doDescribe(actual, descriptor),
					len(actual),
					matchFailMessage,
				)
			}
		},
	}
}

func StringLengthWhichDoes(lengthPattern Matcher[int]) Matcher[string] {
	return StringLengthWhichIs(lengthPattern)
}

func StringLengthWhichWill(lengthPattern Matcher[int]) Matcher[string] {
	return StringLengthWhichIs(lengthPattern)
}

func StringLengthWhichHas(lengthPattern Matcher[int]) Matcher[string] {
	return StringLengthWhichIs(lengthPattern)
}

func StringLengthWhichYaKnow(lengthPattern Matcher[int]) Matcher[string] {
	return StringLengthWhichIs(lengthPattern)
}

type sliceLengthMatcher[ElementT any] struct {
	expectedLength int
	lengthMatcher Matcher[int]
	elementDescriptor Descriptor[ElementT]
}

func(matcher *sliceLengthMatcher[ElementT]) describeActual(
	actual []ElementT,
	sliceDescriptor Descriptor[[]ElementT],
) string {
	if sliceDescriptor == nil {
		sliceDescriptor = DescribeSliceByElement(matcher.elementDescriptor)
	}
	return sliceDescriptor(actual)
}

func(matcher *sliceLengthMatcher[ElementT]) ElementDescribed(
	descriptor Descriptor[ElementT],
) CollectionMatcher[ElementT, []ElementT] {
	if matcher != nil {
		matcher.elementDescriptor = descriptor
	}
	return matcher
}

func(matcher *sliceLengthMatcher[ElementT]) Match(
	context TestContext,
	assumption bool,
	actual []ElementT,
	sliceDescriptor Descriptor[[]ElementT],
) {
	if matcher.lengthMatcher == nil {
		if len(actual) != matcher.expectedLength {
			Abort(context, assumption)(
				"Expected %s to have length %d, but had length %d",
				matcher.describeActual(actual, sliceDescriptor),
				matcher.expectedLength,
				len(actual),
			)
		}
	} else {
		capture := &CaptureTestContext {}
		defer capture.CleanOwnHouse()
		matcher.lengthMatcher.Match(
			capture,
			assumption,
			len(actual),
			func(gotLength int) string {
				return fmt.Sprintf("length %d", gotLength)
			},
		)
		didFail, matchFailMessage := capture.GetRecordedFailure(assumption)
		if didFail {
			Abort(context, assumption)(
				"Expected %s to have matching length, but had length %d: %s",
				matcher.describeActual(actual, sliceDescriptor),
				len(actual),
				matchFailMessage,
			)
		}
	}
}

func SliceLength[ElementT any](expectedLength int) CollectionMatcher[ElementT, []ElementT] {
	return &sliceLengthMatcher[ElementT] {
		expectedLength: expectedLength,
	}
}

func SliceLengthWhichIs[ElementT any](lengthPattern Matcher[int]) CollectionMatcher[ElementT, []ElementT] {
	if lengthPattern == nil {
		panic("No slice length pattern provided")
	}
	return &sliceLengthMatcher[ElementT] {
		lengthMatcher: lengthPattern,
	}
}

func SliceLengthWhichDoes[ElementT any](lengthPattern Matcher[int]) CollectionMatcher[ElementT, []ElementT] {
	return SliceLengthWhichIs[ElementT](lengthPattern)
}

func SliceLengthWhichWill[ElementT any](lengthPattern Matcher[int]) CollectionMatcher[ElementT, []ElementT] {
	return SliceLengthWhichIs[ElementT](lengthPattern)
}

func SliceLengthWhichHas[ElementT any](lengthPattern Matcher[int]) CollectionMatcher[ElementT, []ElementT] {
	return SliceLengthWhichIs[ElementT](lengthPattern)
}

func SliceLengthWhichYaKnow[ElementT any](lengthPattern Matcher[int]) CollectionMatcher[ElementT, []ElementT] {
	return SliceLengthWhichIs[ElementT](lengthPattern)
}

func Nil[T any](nilChecker func(T) bool) Matcher[T] {
	return FuncMatcher[T] {
		Func: func(context TestContext, assumption bool, actual T, descriptor Descriptor[T]) {
			if !nilChecker(actual) {
				Abort(context, assumption)(
					"Expected %s to be nil, but was not",
					doDescribe(actual, descriptor),
				)
			}
		},
	}
}

func NotNil[T any](nilChecker func(T) bool) Matcher[T] {
	return FuncMatcher[T] {
		Func: func(context TestContext, assumption bool, actual T, descriptor Descriptor[T]) {
			if nilChecker(actual) {
				Abort(context, assumption)(
					"Expected %s to be non-nil, but was not",
					doDescribe(actual, descriptor),
				)
			}
		},
	}
}

func IsPointerNil[T any](pointer *T) bool {
	return pointer == nil
}
