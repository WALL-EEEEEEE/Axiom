package stream

import (
	"testing"
)

// TestTestTask calls executor to run TestTask, checking
// for a valid return value.
func TestBroker(t *testing.T) {
	test_broker := NewBroker[int]()
	test_broker.Publish(1)
	out_stream := test_broker.GetOutputStream()
	println(<-out_stream.Out())
}
