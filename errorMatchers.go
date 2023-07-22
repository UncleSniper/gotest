package gotest

var NotAnError Matcher[error] = FuncMatcher[error] {
	Func: func(context TestContext, assumption bool, actual error, descriptor Descriptor[error]) {
		if actual != nil {
			if descriptor == nil {
				descriptor = DescribeError
			}
			Abort(context, assumption)(
				"Expected %s to be nil (i.e. no error), but was not",
				doDescribe(actual, descriptor),
			)
		}
	},
}

func anError(context TestContext, assumption bool, actual error, descriptor Descriptor[error]) bool {
	if actual != nil {
		return true
	}
	Abort(context, assumption)("Expected %s to be a (non-nil) error, but was nil", doDescribe(actual, descriptor))
	return false
}

var AnError Matcher[error] = FuncMatcher[error] {
	Func: func(context TestContext, assumption bool, actual error, descriptor Descriptor[error]) {
		anError(context, assumption, actual, descriptor)
	},
}

func ErrorWithMessage(expectedMessage string) Matcher[error] {
	return FuncMatcher[error] {
		Func: func(context TestContext, assumption bool, actual error, descriptor Descriptor[error]) {
			if anError(context, assumption, actual, descriptor) {
				actualMessage := actual.Error()
				if actualMessage != expectedMessage {
					Abort(context, assumption)(
						"Expected %s to have error message '%s', but had '%s'",
						doDescribe(actual, descriptor),
						expectedMessage,
						actualMessage,
					)
				}
			}
		},
	}
}

func ErrorWithPattern(messagePattern Matcher[string]) Matcher[error] {
	return FuncMatcher[error] {
		Func: func(context TestContext, assumption bool, actual error, descriptor Descriptor[error]) {
			if !anError(context, assumption, actual, descriptor) {
				return
			}
			if messagePattern == nil {
				return
			}
			actualMessage := actual.Error()
			capture := &CaptureTestContext {}
			defer capture.CleanOwnHouse()
			messagePattern.Match(
				capture,
				assumption,
				actualMessage,
				AsString,
			)
			didFail, matchFailMessage := capture.GetRecordedFailure(assumption)
			if didFail {
				Abort(context, assumption)(
					"Expected %s to have matching error message, but had '%s': %s",
					doDescribe(actual, descriptor),
					actualMessage,
					matchFailMessage,
				)
			}
		},
	}
}
