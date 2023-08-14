package stream

type Stream[T any] struct {
	PassThrough[T]
	name string
}

func NewStream[T any](name string) Stream[T] {
	return Stream[T]{name: name, PassThrough: NewPassThrough[T]()}
}

func (stream Stream[T]) GetName() string {
	return stream.name
}
