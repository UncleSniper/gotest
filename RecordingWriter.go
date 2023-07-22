package gotest

import (
	"io"
	"fmt"
	"errors"
)

type RecordingWriter struct {
	contextPlus AssertOrAssume
	expected []ExpectedWrite
	Actual [][]byte
}

type ExpectedWrite struct {
	Count int
	Error error
	Matcher Matcher[[]byte]
}

func NewRecordingWriter(contextPlus AssertOrAssume, expected ...ExpectedWrite) *RecordingWriter {
	if contextPlus.Context == nil {
		panic("No test context set")
	}
	return &RecordingWriter {
		contextPlus: contextPlus,
		expected: expected,
	}
}

func(writer *RecordingWriter) Write(bytes []byte) (int, error) {
	index := len(writer.Actual)
	expectCount := len(writer.expected)
	if index >= expectCount {
		writer.contextPlus.Abort("Expected only %d writes, but got at least one more", expectCount)
		return 0, errors.New(fmt.Sprintf("Expected only %d writes, but got at least one more", expectCount))
	}
	expected := writer.expected[index]
	writer.Actual = append(writer.Actual, CopySlice(bytes))
	if expected.Matcher != nil {
		expected.Matcher.Match(
			writer.contextPlus.Context,
			writer.contextPlus.Assumption,
			bytes,
			func(actual []byte) string {
				return fmt.Sprintf("write #%d of %s", index, DescribeByteSlice(actual))
			},
		)
	}
	return expected.Count, expected.Error
}

func(writer *RecordingWriter) Verify() {
	expected := len(writer.expected)
	actual := len(writer.Actual)
	if actual < expected {
		writer.contextPlus.Abort("Expected %d writes, but got only %d", expected, actual)
	} else if actual > expected {
		writer.contextPlus.Abort("Expected only %d writes, but got %d", expected, actual)
	}
}

var _ io.Writer = &RecordingWriter{}
