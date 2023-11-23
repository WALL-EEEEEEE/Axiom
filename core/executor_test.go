package core

import (
	"testing"

	. "github.com/WALL-EEEEEEE/Axiom/test"
	log "github.com/sirupsen/logrus"

	//. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

type MockTask struct {
	mock.Mock
	Task
}

func TestExecutor(t *testing.T) {
	cases := []TestCase[any, any]{
		{
			Name:     "SingleTask",
			Input:    nil,
			Error:    nil,
			Expected: nil,
			Check: func(tc TestCase[any, any]) {
				mock_task := &MockTask{Task: NewTask(tc.Name)}
				mock_task.On("Run").Return().Run(func(args mock.Arguments) {
					log.Info("???")
				})
				Exec.Add(mock_task)
				Exec.Start()
				//So(expected, ShouldEqual, actual)
			},
		},
	}

	Run(cases, t)
}
