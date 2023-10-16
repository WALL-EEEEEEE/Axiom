package stream

type Stream[T any] struct {
	tunnel PassThrough[T]
	name   string
	closed bool
}

func NewStream[T any](name string) Stream[T] {
	return Stream[T]{name: name, tunnel: NewPassThrough[T](), closed: false}
}

func (stream Stream[T]) GetName() string {
	return stream.name
}

func (stream *Stream[T]) Read() (item T, ok bool) {
	item, ok = <-stream.tunnel.out
	return
}

func (stream *Stream[T]) Write(item T) {
	stream.tunnel.in <- item
}

func (stream *Stream[T]) From(upstream *Stream[T]) {
	go func() {
		for ot := range upstream.tunnel.out {
			stream.Write(ot)
		}
		stream.Close()
	}()
}

func (stream *Stream[T]) To(sink Sink[T]) {
	stream.tunnel.To(sink)
	sink.wait()
}

func (stream *Stream[T]) IsClosed() bool {
	return stream.closed
}

func (stream *Stream[T]) AsArray() []T {
	var result []T
	for {
		it, ok := stream.Read()
		if !ok {
			break
		}
		result = append(result, it)
	}
	return result
}

func (stream *Stream[T]) Close() {
	if stream.closed {
		return
	}
	close(stream.tunnel.in)
	stream.closed = true
}
