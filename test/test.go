package test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type TestCase[T any, E any] struct {
	Name     string
	Error    error
	Input    T
	Expected E
	Check    func(TestCase[T, E])
}

func Run(cases []TestCase[any, any], t *testing.T) {
	for _, t_case := range cases {
		Convey(t_case.Name, t, func() {
			t_case.Check(t_case)
		})
	}
}
