package gotest

import (
	"fmt"
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
