package stream

type OutputSink[T any] struct {
	Name    string
	in      chan T
	outfunc func(T)
}

func (sink OutputSink[T]) In() chan<- T {
	return sink.in
}

func (sink *OutputSink[T]) doStream() {
	for elem := range sink.in {
		sink.outfunc(elem)
	}
}

func NewOutputSink[T any](name string, callback func(T)) OutputSink[T] {
	outputSink := OutputSink[T]{
		Name:    name,
		in:      make(chan T),
		outfunc: callback,
	}
	go outputSink.doStream()
	return outputSink
}

var _ Sink[interface{}] = (*OutputSink[interface{}])(nil)
