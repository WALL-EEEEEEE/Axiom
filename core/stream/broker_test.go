package stream

// TestTestTask calls executor to run TestTask, checking
// for a valid return value.
/*
func TestBrokerPass(t *testing.T) {
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

func TestBrokerMultiOutPass(t *testing.T) {
	broker := NewBroker[int]()
	multi := 2
	var ostreams []Stream[int]
	istream := broker.GetInputStream()
	for i := 0; i < multi; i++ {
		ostreams = append(ostreams, broker.GetOutputStream())
	}
	cnt := 10
	go func() {
		for i := 0; i < cnt; i++ {
			istream.In() <- i
			t.Logf("Input: %d", i)
		}

	}()
	for i := 0; i < cnt; i++ {
		for i, ostream := range ostreams {
			t.Logf("Output-%d: %d", i, <-ostream.Out())
		}
	}
}
*/
