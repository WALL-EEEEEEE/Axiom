package stream

import (
	"testing"
)

func TestStreamPass(t *testing.T) {
	broker := NewBroker[int]()
	istream := broker.GetInputStream()
	ostream := broker.GetOutputStream()
	cnt := 10
	go func() {
		for i := 0; i < cnt; i++ {
			istream.In() <- i
			t.Logf("Input: %d", i)
		}

	}()
	for i := 0; i < cnt; i++ {
		t.Logf("Output: %d", <-ostream.Out())
	}
}
