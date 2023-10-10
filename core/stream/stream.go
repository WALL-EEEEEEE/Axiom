package stream

type Stream[T any] struct {
	tunnel PassThrough[T]
	name   string
}

func NewStream[T any](name string) Stream[T] {
	return Stream[T]{name: name, tunnel: NewPassThrough[T]()}
}

func (stream Stream[T]) GetName() string {
	return stream.name
}

func (stream *Stream[T]) Read() T {
	return <-stream.tunnel.out
}

func (stream *Stream[T]) Write(item T) {
	stream.tunnel.in <- item
}

func (stream *Stream[T]) From(upstream *Stream[T]) {
	stream.tunnel.From(&(upstream.tunnel))
}

func (stream *Stream[T]) Close() {
	close(stream.tunnel.in)
}
