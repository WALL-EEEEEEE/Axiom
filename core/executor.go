package core

import (
	"sync"

	. "github.com/WALL-EEEEEEE/Axiom/core/stream"
)

type IExecutor interface {
	Start()
	GetName() string
	Add(task IRunnable)
	List() []string
}

type IRunnable interface {
	Serv
	Run()
	GetName() string
}

type Task struct {
	name   string
	nxts   []Task
}

func NewTask(name string) Task {
	broker := NewBroker[interface{}]()
}

func (task Task) GetName() string {
	return task.name
}

func (task Task) GetType() []ServType {
	return []ServType{TASK}
}

}

func (task Task) Run() {
	var wg sync.WaitGroup
	_start := func(task Task) {
		defer func() {
			wg.Done()
		}()
		task.Run()
	}
	for _, task := range task.nxts {
		wg.Add(1)
		go _start(task)
	}
	wg.Wait()
}

var _ IRunnable = (*Task)(nil)

type DefaultExecutor struct {
}

func NewDFExecutor(name string) DefaultExecutor {
	return exec
}

func (executor DefaultExecutor) startTask(wg *sync.WaitGroup) {
	_start := func(task IRunnable) {
		defer func() {
			wg.Done()
		}()
		task.Run()
	}
	for _, task := range executor.tasks {
		wg.Add(1)
		go _start(task)
	}
}

func (executor DefaultExecutor) Start() {
	var task_wg sync.WaitGroup
	executor.startTask(&task_wg)
	task_wg.Wait()
}

func (executor DefaultExecutor) Add(task IRunnable) {
	executor.tasks = append(executor.tasks, task)
}

func (executor DefaultExecutor) List() []string {
	var names []string
	for _, task := range executor.tasks {
		names = append(names, task.GetName())
	}
	return names
}

func (executor DefaultExecutor) GetName() string {
	return executor.name
}

var _ IExecutor = (*DefaultExecutor)(nil)

var Exec = NewDFExecutor("default")
