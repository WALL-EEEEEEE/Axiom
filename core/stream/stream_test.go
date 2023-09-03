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
			Func: func(tc TestCase[any, any]) {
				stream := NewStream[int](tc.Input.(string))
				print(stream.GetName())
				assert.Equal(t, tc.Expected, stream.GetName())
			},
		},
		{
			Name:     "InOut",
			Input:    []int{1, 2, 3, 4},
			Error:    nil,
			Expected: []int{1, 2, 3, 4},
			Func: func(tc TestCase[any, any]) {
				stream := NewStream[int](tc.Name)
				var expected []int
				go func() {
					for _, item := range tc.Input.([]int) {
						//t.Logf("stream %s <- %+v", stream.GetName(), item)
						stream.In() <- item
					}
					close(stream.In())
				}()
				for item := range stream.Out() {
					//t.Logf("stream %s -> %+v", stream.GetName(), item)
					expected = append(expected, item)
				}
				assert.Equal(t, tc.Expected, expected)
			},
		},
		{
			Name:     "Via",
			Input:    []int{1, 2, 3, 4},
			Error:    nil,
			Expected: []int{1, 2, 3, 4},
			Func: func(tc TestCase[any, any]) {
				stream := NewStream[int](tc.Name)
				var expected []int
				go func() {
					for _, item := range tc.Input.([]int) {
						//t.Logf("stream %s <- %+v", stream.GetName(), item)
						stream.In() <- item
					}
					close(stream.In())
				}()
				for item := range stream.Out() {
					//t.Logf("stream %s -> %+v", stream.GetName(), item)
					expected = append(expected, item)
				}
				assert.Equal(t, tc.Expected, expected)
			},
		},
	}
	Run(cases, t)
}
