package stream

import (
	"testing"

	. "github.com/WALL-EEEEEEE/Axiom/test"
	. "github.com/smartystreets/goconvey/convey"
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
				So(tc.Expected, ShouldEqual, stream.GetName())
			},
		},
		{
			Name:     "Read/Write",
			Input:    []int{1, 2, 3, 4},
			Error:    nil,
			Expected: []int{1, 2, 3, 4},
			Check: func(tc TestCase[any, any]) {
				stream := NewStream[int](tc.Name)
				go func() {
					for _, item := range tc.Input.([]int) {
						//t.Logf("stream %s <- %+v", stream.GetName(), item)
						stream.Write(item)
					}
					stream.Close()
				}()
				So(tc.Expected, ShouldEqual, stream.AsArray())
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
				go func() {
					for _, item := range tc.Input.([]int) {
						//t.Logf("stream %s <- %+v", upstream.GetName(), item)
						upstream.Write(item)
					}
					upstream.Close()
				}()
				So(tc.Expected, ShouldEqual, downstream.AsArray())
			},
		},
		{
			Name:     "To",
			Input:    []int{1, 2, 3, 4},
			Error:    nil,
			Expected: []int{1, 2, 3, 4},
			Check: func(tc TestCase[any, any]) {
				stream := NewStream[int](tc.Name)
				var expected []int
				go func() {
					for _, item := range tc.Input.([]int) {
						//t.Logf("stream %s <- %+v", stream.GetName(), item)
						stream.Write(item)
					}
					stream.Close()
				}()
				sink := NewOutputSink[int](tc.Name+"_stdout_sink", func(item int) {
					expected = append(expected, item)
				})
				stream.To(&sink)
				So(tc.Expected, ShouldEqual, expected)
			},
		},
		{
			Name:     "AsArray",
			Input:    []int{1, 2, 3, 4},
			Error:    nil,
			Expected: []int{1, 2, 3, 4},
			Check: func(tc TestCase[any, any]) {
				stream := NewStream[int](tc.Name)
				go func() {
					for _, item := range tc.Input.([]int) {
						//t.Logf("stream %s <- %+v", stream.GetName(), item)
						stream.Write(item)
					}
					stream.Close()
				}()
				So(tc.Expected, ShouldEqual, stream.AsArray())
			},
		},
	}
	Run(cases, t)
}
