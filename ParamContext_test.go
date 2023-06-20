package gotest

import (
	"fmt"
	tst "testing"
)

func describeSubject[T any](subject *Subject[T], descriptor Descriptor[T]) string {
	var contextDescr, descriptorDescr, valueDescr string
	if subject.context == nil {
		contextDescr = "nil"
	} else {
		contextDescr = "<...>"
	}
	if subject.descriptor == nil {
		descriptorDescr = "nil"
	} else {
		descriptorDescr = "<...>"
	}
	if descriptor == nil {
		valueDescr = fmt.Sprintf("%v", subject.value)
	} else {
		valueDescr = descriptor(subject.value)
	}
	return fmt.Sprintf(
		"{ %s, %s, %s, %s }",
		fmt.Sprintf("context: %s", contextDescr),
		fmt.Sprintf("value: %s", valueDescr),
		fmt.Sprintf("descriptor: %s", descriptorDescr),
		fmt.Sprintf("assumption: %v", subject.assumption),
	)
}

func TestSubjectValue(t *tst.T) {
	c := Use(t)
	subject := &Subject[int] {
		value: 42,
	}
	AssertThat(c, subject.Value()).Is(EqualTo(42))
}
