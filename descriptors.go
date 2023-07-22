package gotest

import (
	"fmt"
	"strings"
	"reflect"
)

func TheString(value string) string {
	return fmt.Sprintf("the string \"%s\"", value)
}

func WithTheString(value string) string {
	return fmt.Sprintf("with the string \"%s\"", value)
}

func AsString(value string) string {
	return fmt.Sprintf("\"%s\"", value)
}

func DescribeByteSlice(bytes []byte) string {
	if len(bytes) == 0 {
		return "zero bytes"
	}
	var builder strings.Builder
	builder.WriteString("bytes")
	for _, b := range bytes {
		builder.WriteString(fmt.Sprintf(" %02X", b))
	}
	return builder.String()
}

func DescribeError(err error) string {
	if err == nil {
		return "nil"
	}
	return fmt.Sprintf("error of type %s with message '%s'", reflect.TypeOf(err).String(), err.Error())
}

// shorthands for Subject with descriptor

func AssertThatError(context TestContext, actual error) *Subject[error] {
	return AssertThat(context, actual).Described(DescribeError)
}

func AssumeThatError(context TestContext, actual error) *Subject[error] {
	return AssumeThat(context, actual).Described(DescribeError)
}

func ThatError(contextPlus AssertOrAssume, actual error) *Subject[error] {
	return That(contextPlus, actual).Described(DescribeError)
}
