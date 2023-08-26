package gotest

import (
	"fmt"
	"strings"
)

type Descriptor[T any] func(T) string

type Subject[T any] struct {
	context TestContext
	value T
	descriptor Descriptor[T]
	assumption bool
}

func(subject *Subject[T]) Value() T {
	if subject == nil {
		panic("Cannot get value of nil Subject")
	}
	return subject.value
}

func AssertThat[T any](context TestContext, value T) *Subject[T] {
	if context == nil {
		panic("No test context provided")
	}
	return &Subject[T]{
		context: context,
		value: value,
		assumption: false,
	}
}

func AssumeThat[T any](context TestContext, value T) *Subject[T] {
	if context == nil {
		panic("No test context set")
	}
	return &Subject[T]{
		context: context,
		value: value,
		assumption: true,
	}
}

func(subject *Subject[T]) Named(name string) *Subject[T] {
	if subject != nil {
		subject.descriptor = func(value T) string {
			return name
		}
	}
	return subject
}

func(subject *Subject[T]) NamedAndValued(name string) *Subject[T] {
	if subject != nil {
		oldDescriptor := subject.descriptor
		subject.descriptor = func(value T) string {
			if oldDescriptor == nil {
				return name
			} else {
				return fmt.Sprintf("%s = %s", name, oldDescriptor(value))
			}
		}
	}
	return subject
}

func(subject *Subject[T]) MaybeNamed(name string) *Subject[T] {
	if subject != nil && subject.descriptor == nil {
		subject.descriptor = func(value T) string {
			return name
		}
	}
	return subject
}

func(subject *Subject[T]) Described(descriptor Descriptor[T]) *Subject[T] {
	if subject != nil {
		subject.descriptor = descriptor
	}
	return subject
}

func(subject *Subject[T]) MaybeDescribed(descriptor Descriptor[T]) *Subject[T] {
	if subject != nil && subject.descriptor == nil {
		subject.descriptor = descriptor
	}
	return subject
}

func(subject *Subject[T]) Describe() string {
	if subject == nil {
		return "<missing subject>"
	} else if subject.descriptor == nil {
		return fmt.Sprintf("%v", subject.value)
	} else {
		return subject.descriptor(subject.value)
	}
}

func DescribeSliceByElement[ElementT any](elementDescriptor Descriptor[ElementT]) Descriptor[[]ElementT] {
	return func(slice []ElementT) string {
		if slice == nil {
			return "nil"
		}
		var builder strings.Builder
		builder.WriteRune('[')
		var elementDescr string
		for index, found := range slice {
			if index > 0 {
				 builder.WriteString(", ")
			}
			if elementDescriptor == nil {
				elementDescr = fmt.Sprintf("%v", found)
			} else {
				elementDescr = elementDescriptor(found)
			}
			builder.WriteString(elementDescr)
		}
		builder.WriteRune(']')
		return builder.String()
	}
}

func DescribeMapByPair[KeyT comparable, ValueT any](
	keyDescriptor Descriptor[KeyT],
	valueDescriptor Descriptor[ValueT],
) Descriptor[map[KeyT]ValueT] {
	return func(theMap map[KeyT]ValueT) string {
		if theMap == nil {
			return "nil"
		}
		var builder strings.Builder
		builder.WriteRune('{')
		var keyDescr, valueDescr string
		hadPair := false
		for key, value := range theMap {
			if hadPair {
				builder.WriteString(", ")
			} else {
				hadPair = true
			}
			if keyDescriptor == nil {
				keyDescr = fmt.Sprintf("%v", key)
			} else {
				keyDescr = keyDescriptor(key)
			}
			builder.WriteString(keyDescr)
			builder.WriteString(": ")
			if valueDescriptor == nil {
				valueDescr = fmt.Sprintf("%v", value)
			} else {
				valueDescr = valueDescriptor(value)
			}
			builder.WriteString(valueDescr)
		}
		builder.WriteRune('}')
		return builder.String()
	}
}

func(subject *Subject[T]) Is(matcher Matcher[T]) *Subject[T] {
	if subject == nil {
		panic("Will not assert missing subject")
	} else if matcher != nil {
		matcher.Match(
			subject.context,
			subject.assumption,
			subject.value,
			func(value T) string {
				if subject.descriptor == nil {
					return fmt.Sprintf("%v", value)
				} else {
					return subject.descriptor(value)
				}
			},
		)
	}
	return subject
}

func(subject *Subject[T]) AndIs(matcher Matcher[T]) *Subject[T] {
	return subject.Is(matcher)
}

func(subject *Subject[T]) Does(matcher Matcher[T]) *Subject[T] {
	return subject.Is(matcher)
}

func(subject *Subject[T]) AndDoes(matcher Matcher[T]) *Subject[T] {
	return subject.Is(matcher)
}

func(subject *Subject[T]) Will(matcher Matcher[T]) *Subject[T] {
	return subject.Is(matcher)
}

func(subject *Subject[T]) AndWill(matcher Matcher[T]) *Subject[T] {
	return subject.Is(matcher)
}

func(subject *Subject[T]) Has(matcher Matcher[T]) *Subject[T] {
	return subject.Is(matcher)
}

func(subject *Subject[T]) AndHas(matcher Matcher[T]) *Subject[T] {
	return subject.Is(matcher)
}

func(subject *Subject[T]) YaKnow(matcher Matcher[T]) *Subject[T] {
	return subject.Is(matcher)
}

func(subject *Subject[T]) AndYaKnow(matcher Matcher[T]) *Subject[T] {
	return subject.Is(matcher)
}

func MapChain[OldT any, NewT any](
	subject *Subject[OldT],
	valueTransform func(OldT) NewT,
	descriptionTransform func(string, NewT) string,
) *Subject[NewT] {
	oldDescription := subject.Describe()
	var newDescriptor Descriptor[NewT]
	if descriptionTransform == nil {
		newDescriptor = func(newValue NewT) string {
			return oldDescription
		}
	} else {
		newDescriptor = func(newValue NewT) string {
			return descriptionTransform(oldDescription, newValue)
		}
	}
	return &Subject[NewT] {
		context: subject.context,
		value: valueTransform(subject.value),
		descriptor: newDescriptor,
		assumption: subject.assumption,
	}
}

func MapCompose[OldT any, NewT any](
	subject *Subject[OldT],
	valueTransform func(OldT) NewT,
	descriptionTransform func(OldT, Descriptor[OldT], NewT) string,
) *Subject[NewT] {
	var newDescriptor Descriptor[NewT]
	if descriptionTransform == nil {
		oldDescription := subject.Describe()
		newDescriptor = func(newValue NewT) string {
			return oldDescription
		}
	} else {
		newDescriptor = func(newValue NewT) string {
			return descriptionTransform(
				subject.value,
				func(oldValue OldT) string {
					if subject.descriptor == nil {
						return fmt.Sprintf("%v", oldValue)
					} else {
						return subject.descriptor(oldValue)
					}
				},
				newValue,
			)
		}
	}
	return &Subject[NewT] {
		context: subject.context,
		value: valueTransform(subject.value),
		descriptor: newDescriptor,
		assumption: subject.assumption,
	}
}

func MapProperty[OwnerT any, PropertyT any](
	subject *Subject[OwnerT],
	propertyName string,
	valueTransform func(OwnerT) PropertyT,
	propertyDescriptor Descriptor[PropertyT],
) *Subject[PropertyT] {
	return MapChain[OwnerT, PropertyT](
		subject,
		valueTransform,
		func(ownerDescr string, propertyValue PropertyT) string {
			var propDescr string
			if propertyDescriptor == nil {
				propDescr = fmt.Sprintf("%v", propertyValue)
			} else {
				propDescr = propertyDescriptor(propertyValue)
			}
			return fmt.Sprintf("(%s).%s == %s", ownerDescr, propertyName, propDescr)
		},
	)
}

func MapSliceElement[ElementT any](
	subject *Subject[[]ElementT],
	index int,
	elementDescriptor Descriptor[ElementT],
) *Subject[ElementT] {
	if index < 0 {
		panic(fmt.Sprintf("Negative slice index: %d", index))
	}
	if index >= len(subject.Value()) {
		panic(fmt.Sprintf("Slice index out of bounds: %d >= %d", index, len(subject.value)))
	}
	return MapChain[[]ElementT, ElementT](
		subject,
		func(slice []ElementT) ElementT {
			return slice[index]
		},
		func(ownerDescr string, element ElementT) string {
			var elemDescr string
			if elementDescriptor == nil {
				elemDescr = fmt.Sprintf("%v", element)
			} else {
				elemDescr = elementDescriptor(element)
			}
			return fmt.Sprintf("(%s)[%d] == %s", ownerDescr, index, elemDescr)
		},
	)
}
