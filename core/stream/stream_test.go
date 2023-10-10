package stream

import (
	"testing"

	. "github.com/WALL-EEEEEEE/Axiom/test"
	"github.com/stretchr/testify/assert"
)

func TestStream(t *testing.T) {
	cases := []TestCase[any, any]{
		{
			Name:     "GetName",
			Input:    "Test",
			Error:    nil,
			Expected: "Test",
			Check: func(tc TestCase[any, any]) {
				stream := NewStream[int](tc.Input.(string))
				assert.Equal(t, tc.Expected, stream.GetName())
			},
		},
		{
			Name:     "From",
			Input:    []int{1, 2, 3, 4},
			Error:    nil,
			Expected: []int{1, 2, 3, 4},
			Check: func(tc TestCase[any, any]) {
				upstream := NewStream[int](tc.Name + "_up")
				downstream := NewStream[int](tc.Name + "_down")
				downstream.From(&upstream)
				var expected []int
				go func() {
					for _, item := range tc.Input.([]int) {
						t.Logf("stream %s <- %+v", upstream.GetName(), item)
						upstream.Write(item)
					}
					upstream.Close()
				}()
				for item := downstream.Read(); item != 0; item = downstream.Read() {
					t.Logf("stream %s -> %+v", downstream.GetName(), item)
					expected = append(expected, item)
				}
				assert.Equal(t, tc.Expected, expected)
			},
		},
		{
			Name:     "Read/Write",
			Input:    []int{1, 2, 3, 4},
			Error:    nil,
			Expected: []int{1, 2, 3, 4},
			Check: func(tc TestCase[any, any]) {
				stream := NewStream[int](tc.Name)
				var expected []int
				go func() {
					for _, item := range tc.Input.([]int) {
						t.Logf("stream %s <- %+v", stream.GetName(), item)
						stream.Write(item)
					}
					stream.Close()
				}()
				for item := stream.Read(); item != 0; item = stream.Read() {
					t.Logf("stream %s -> %+v", stream.GetName(), item)
					expected = append(expected, item)
				}
				assert.Equal(t, tc.Expected, expected)
			},
		},
	}
	Run(cases, t)
}
