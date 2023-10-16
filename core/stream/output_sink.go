package stream

type OutputSink[T any] struct {
	Name    string
	in      chan T
	outfunc func(T)
	close   chan bool
}

func (sink OutputSink[T]) In() chan<- T {
	return sink.in
}

func (sink *OutputSink[T]) doStream() {
	for elem := range sink.in {
		sink.outfunc(elem)
	}
	close(sink.close)
}
func (sink *OutputSink[T]) wait() {
	<-sink.close
}

func NewOutputSink[T any](name string, callback func(T)) OutputSink[T] {
	outputSink := OutputSink[T]{
		Name:    name,
		in:      make(chan T),
		outfunc: callback,
		close:   make(chan bool),
	}
	go outputSink.doStream()
	return outputSink
}

var _ Sink[interface{}] = (*OutputSink[interface{}])(nil)
