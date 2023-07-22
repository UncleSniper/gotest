package gotest

func CopySlice[T any](src []T) []T {
	dest := make([]T, len(src))
	copy(dest, src)
	return dest
}

func CopySliceWith[T any](src []T, elementCopier func(T) T) []T {
	dest := make([]T, len(src))
	if elementCopier == nil {
		elementCopier = ShallowCopy[T]
	}
	for index, element := range src {
		dest[index] = elementCopier(element)
	}
	return dest
}

func ShallowCopy[T any](src T) T {
	return src
}
