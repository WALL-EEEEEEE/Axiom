package stream

// PassThrough retransmits incoming elements as is.
//
// in  -- 1 -- 2 ---- 3 -- 4 ------ 5 --
//
// out -- 1 -- 2 ---- 3 -- 4 ------ 5 --
type PassThrough[T any] struct {
	in  chan T
	out chan T
}

// Verify PassThrough satisfies the Flow interface.
var _ Flow[interface{}] = (*PassThrough[interface{}])(nil)

// NewPassThrough returns a new PassThrough instance.
func NewPassThrough[T any]() PassThrough[T] {
	passThrough := PassThrough[T]{
		in:  make(chan T),
		out: make(chan T),
	}
	go passThrough.doStream()
	return passThrough
}

// From streams data through the given flow
func (pt *PassThrough[T]) From(flow Flow[T]) Flow[T] {
	go func() {
		for elem := range flow.Out() {
			pt.in <- elem
		}
		close(pt.in)
	}()
	return pt
}

// To streams data to the given sink
func (pt *PassThrough[T]) To(sink Sink[T]) {
	pt.transmit(sink)
}

// Out returns an output channel for sending data
func (pt *PassThrough[T]) Out() <-chan T {
	return pt.out
}

// In returns an input channel for receiving data
func (pt *PassThrough[T]) In() chan<- T {
	return pt.in
}

func (pt *PassThrough[T]) transmit(inlet Inlet[T]) {
	for elem := range pt.out {
		inlet.In() <- elem
	}
	close(inlet.In())
}

func (pt *PassThrough[T]) doStream() {
	for elem := range pt.in {
		pt.out <- elem
	}
	close(pt.out)
}
