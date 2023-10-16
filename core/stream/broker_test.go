package stream

import (
	"math/rand"
	"sync"
	"testing"

	. "github.com/WALL-EEEEEEE/Axiom/test"
	"github.com/stretchr/testify/assert"
)

func genRandSlice(num int, min int, max int) []int {
	var slice []int

	for i := 0; i < num; i++ {
		it := rand.Intn(max-min) + min
		slice = append(slice, it)
	}
	return slice
}

func TestBroker(t *testing.T) {
	In_Out_Input := genRandSlice(10, 0, 100)
	In_2Out_Input := genRandSlice(10, 0, 100)
	cases := []TestCase[any, any]{
		{
			Name:     "In->Out",
			Input:    In_Out_Input,
			Error:    nil,
			Expected: In_Out_Input,
			Check: func(tc TestCase[any, any]) {
				broker := NewBroker[int]()
				istream := broker.GetInputStream()
				ostream := broker.GetOutputStream()
				go func() {
					defer broker.Close()
					for _, item := range tc.Input.([]int) {
						//t.Logf("Input: %d", item)
						istream.Write(item)
					}
				}()
				expected := ostream.AsArray()
				//t.Logf("Output: %+v", expected)
				assert.Equal(t, tc.Expected, expected)
			},
		},
		{
			Name:     "In->2Out",
			Input:    In_2Out_Input,
			Error:    nil,
			Expected: In_2Out_Input,
			Check: func(tc TestCase[any, any]) {
				broker := NewBroker[int]()
				multi := 2
				var ostreams []Stream[int]
				var ostreamExpects [][]int
				istream := broker.GetInputStream()
				for i := 0; i < multi; i++ {
					ostreams = append(ostreams, broker.GetOutputStream())
					ostreamExpects = append(ostreamExpects, []int{})
				}
				go func() {
					defer broker.Close()
					for _, it := range tc.Input.([]int) {
						istream.Write(it)
						//t.Logf("Input: %d", it)
					}
				}()

				var wg sync.WaitGroup
				for i, ostream := range ostreams {
					wg.Add(1)
					go func(i int, ostream Stream[int]) {
						ostreamExpects[i] = ostream.AsArray()
						defer wg.Done()
					}(i, ostream)
				}
				wg.Wait()
				for _, ostreamExpect := range ostreamExpects {
					assert.Equal(t, tc.Input, ostreamExpect)
				}
			},
		},
	}
	Run(cases, t)
}
