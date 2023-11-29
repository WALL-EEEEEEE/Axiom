package util

import (
	"testing"

	. "github.com/WALL-EEEEEEE/Axiom/test"
	"github.com/stretchr/testify/assert"
)

func TestGather(t *testing.T) {
	cases := []TestCase[any, any]{
		{
			Name:     "GatherRecv",
			Input:    nil,
			Error:    nil,
			Expected: nil,
			Check: func(tc TestCase[any, any]) {
				var select_chans []chan int
				for i := 0; i < 10; i++ {
					schan := make(chan int)
					select_chans = append(select_chans, schan)
				}
				var expected, actual []int
				go func() {
					for i, c := range select_chans {
						//t.Logf("Send %d to %+v", i, c)
						c <- i
						expected = append(expected, i)
						close(c)
					}
				}()
				GatherRecv[int](select_chans, func(chan_no int, v int, ok bool) {
					actual = append(actual, v)
				})
				assert.Equal(t, expected, actual)
			},
		},
		{
			Name:     "GatherSend",
			Input:    []int{1, 2, 3, 4},
			Error:    nil,
			Expected: []int{1, 2, 3, 4},
			Check: func(tc TestCase[any, any]) {
				var select_chans []chan int
				for i := 0; i < 10; i++ {
					schan := make(chan int)
					select_chans = append(select_chans, schan)
				}
				go func() {
					for _, v := range tc.Input.([]int) {
						GatherSend[int](select_chans, v)
					}
					for _, c := range select_chans {
						close(c)
					}
				}()
				var actual [][]int = make([][]int, len(select_chans))
				for i, _ := range actual {
					actual[i] = make([]int, 0)
				}
				GatherRecv[int](select_chans, func(chan_no int, v int, ok bool) {
					actual[chan_no] = append(actual[chan_no], v)
				})
				for _, s := range actual {
					assert.EqualValues(t, tc.Expected, s)
				}
				//t.Logf("%+v", actual)

			},
		},
	}

	Run(cases, t)
}
