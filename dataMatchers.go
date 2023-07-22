package gotest

func EqualToBytes(expected ...byte) Matcher[[]byte] {
	return FuncMatcher[[]byte] {
		Func: func(context TestContext, assumption bool, actual []byte, descriptor Descriptor[[]byte]) {
			length := len(actual)
			if length != len(expected) {
				Abort(context, assumption)(
					"Expected %s to have length %d, but had length %d",
					doDescribe(actual, descriptor),
					len(expected),
					len(actual),
				)
				return
			}
			for index := 0; index < length; index++ {
				if actual[index] != expected[index] {
					Abort(context, assumption)(
						"Expected %s to be equivalent to %s, but differed at index %d",
						doDescribe(actual, descriptor),
						DescribeByteSlice(expected),
						index,
					)
					return
				}
			}
		},
	}
}
