package stream

import (
	"sync"
	"time"
)

// SessionWindow generates groups of elements by sessions of activity.
// Session windows do not overlap and do not have a fixed start and end time.
type SessionWindow[T any] struct {
	sync.Mutex
	inactivityGap time.Duration
	timer         *time.Timer
	in            chan T
	out           chan T
	done          chan struct{}
	buffer        []T
}

// Verify SessionWindow satisfies the Flow interface.
var _ Flow[any] = (*SessionWindow[any])(nil)

// NewSessionWindow returns a new SessionWindow instance.
//
// inactivityGap is the gap of inactivity that closes a session window when occurred.
func NewSessionWindow[T any](inactivityGap time.Duration) *SessionWindow[T] {
	window := &SessionWindow[T]{
		inactivityGap: inactivityGap,
		timer:         time.NewTimer(inactivityGap),
		in:            make(chan T),
		out:           make(chan T),
		done:          make(chan struct{}),
	}
	go window.emit()
	go window.receive()

	return window
}

// Via streams data through the given flow
func (sw *SessionWindow[T]) From(flow Flow[T]) Flow[T] {
	go sw.transmit(flow)
	return flow
}

// To streams data to the given sink
func (sw *SessionWindow[T]) To(sink Sink[T]) {
	sw.transmit(sink)
}

// Out returns an output channel for sending data
func (sw *SessionWindow[T]) Out() <-chan T {
	return sw.out
}

// In returns an input channel for receiving data
func (sw *SessionWindow[T]) In() chan<- T {
	return sw.in
}

// submit emitted windows to the next Inlet
func (sw *SessionWindow[T]) transmit(inlet Inlet[T]) {
	for elem := range sw.Out() {
		inlet.In() <- elem
	}
	close(inlet.In())
}

func (sw *SessionWindow[T]) receive() {
	for elem := range sw.in {
		sw.Lock()
		sw.buffer = append(sw.buffer, elem)
		sw.timer.Reset(sw.inactivityGap)
		sw.Unlock()
	}
	close(sw.done)
	close(sw.out)
}

// emit generates and emits a new window.
func (sw *SessionWindow[T]) emit() {
	defer sw.timer.Stop()

	for {
		select {
		case <-sw.timer.C:
			sw.Lock()
			windowSlice := sw.buffer
			sw.buffer = nil
			sw.Unlock()

			// send the window slice to the out chan
			if len(windowSlice) > 0 {
				for _, item := range windowSlice {
					sw.out <- item
				}
			}

		case <-sw.done:
			return
		}
	}
}
