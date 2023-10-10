package stream

import (
	"sync"
	"time"
)

// TumblingWindow assigns each element to a window of a specified window size.
// Tumbling windows have a fixed size and do not overlap.
type TumblingWindow[T any] struct {
	sync.Mutex
	windowSize time.Duration
	in         chan T
	out        chan T
	done       chan struct{}
	buffer     []T
}

// Verify TumblingWindow satisfies the Flow interface.
var _ Flow[any] = (*TumblingWindow[any])(nil)

// NewTumblingWindow returns a new TumblingWindow instance.
//
// size is the Duration of generated windows.
func NewTumblingWindow[T any](size time.Duration) *TumblingWindow[T] {
	window := &TumblingWindow[T]{
		windowSize: size,
		in:         make(chan T),
		out:        make(chan T),
		done:       make(chan struct{}),
	}
	go window.receive()
	go window.emit()

	return window
}

// Via streams data through the given flow
func (tw *TumblingWindow[T]) From(flow Flow[T]) Flow[T] {
	go tw.transmit(flow)
	return flow
}

// To streams data to the given sink
func (tw *TumblingWindow[T]) To(sink Sink[T]) {
	tw.transmit(sink)
}

// Out returns an output channel for sending data
func (tw *TumblingWindow[T]) Out() <-chan T {
	return tw.out
}

// In returns an input channel for receiving data
func (tw *TumblingWindow[T]) In() chan<- T {
	return tw.in
}

// submit emitted windows to the next Inlet
func (tw *TumblingWindow[T]) transmit(inlet Inlet[T]) {
	for elem := range tw.Out() {
		inlet.In() <- elem
	}
	close(inlet.In())
}

func (tw *TumblingWindow[T]) receive() {
	for elem := range tw.in {
		tw.Lock()
		tw.buffer = append(tw.buffer, elem)
		tw.Unlock()
	}
	close(tw.done)
	close(tw.out)
}

// emit generates and emits a new window.
func (tw *TumblingWindow[T]) emit() {
	ticker := time.NewTicker(tw.windowSize)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tw.Lock()
			windowSlice := tw.buffer
			tw.buffer = nil
			tw.Unlock()

			// send the window slice to the out chan
			if len(windowSlice) > 0 {
				for _, item := range windowSlice {
					tw.out <- item
				}
			}

		case <-tw.done:
			return
		}
	}
}
