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

func (stream *Stream[T]) Read() T {
	return <-stream.out
}

func (stream *Stream[T]) Write(item T) {
	stream.out <- item
}
