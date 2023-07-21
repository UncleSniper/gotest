package gotest

import (
	"fmt"
	"strings"
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
