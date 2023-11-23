package core

import (
	"sync"

	. "github.com/WALL-EEEEEEE/Axiom/core/stream"
)

type IExecutor interface {
	Start()
	GetName() string
	Add(task ITask)
	List() []string
}

type IRunnable interface {
	Serv
	Run()
	GetName() string
}

type ITask interface {
	IRunnable
	run()
	From(ITask)
	Chain(...ITask)
	Publish(interface{})
	GetOutputStream() Stream[interface{}]
	GetInputStream() Stream[interface{}]
}

type Task struct {
	ITask
	broker Broker[interface{}]
	name   string
	nxts   []ITask
}

func NewTask(name string) Task {
	broker := NewBroker[interface{}](name)
	return Task{name: name, broker: broker}
}

func (task *Task) GetName() string {
	return task.name
}

func (task *Task) GetType() []ServType {
	return []ServType{TASK}
}

func (task *Task) GetOutputStream() Stream[interface{}] {
	stream := task.broker.GetOutputStream()
	return stream
}
func (task *Task) GetInputStream() Stream[interface{}] {
	return task.broker.GetInputStream()
}

func (task *Task) From(otask ITask) {
	upstream := otask.GetOutputStream()
	task.broker.Via(upstream)
}

func (task *Task) Chain(nextTasks ...ITask) {
	for _, ntask := range nextTasks {
		ntask.From(task)
		task.nxts = append(task.nxts, ntask)
	}
}
func (task *Task) Send(msg interface{}) {
	task.broker.Send(msg)
}

func (task *Task) run() {
	var wg sync.WaitGroup
	_start := func(itask ITask) {
		defer func() {
			wg.Done()
		}()
		itask.run()
	}
	for _, task := range task.nxts {
		wg.Add(1)
		go _start(task)
	}
	task.Run()
	task.close()
	wg.Wait()
}

func (t *Task) close() {
	t.broker.Close()
}

var _ ITask = (*Task)(nil)

type DefaultExecutor struct {
	name  string
	tasks []ITask
}

func NewDFExecutor(name string) DefaultExecutor {
	exec := DefaultExecutor{name: name}
	return exec
}

func (executor *DefaultExecutor) startTask(wg *sync.WaitGroup) {
	_start := func(task ITask) {
		defer func() {
			wg.Done()
		}()
		task.run()
	}
	for _, task := range executor.tasks {
		wg.Add(1)
		go _start(task)
	}
}

func (executor *DefaultExecutor) Start() {
	var task_wg sync.WaitGroup
	executor.startTask(&task_wg)
	task_wg.Wait()
}

func (executor *DefaultExecutor) Add(task ITask) {
	executor.tasks = append(executor.tasks, task)
}

func (executor *DefaultExecutor) List() []string {
	var names []string
	for _, task := range executor.tasks {
		names = append(names, task.GetName())
	}
	return names
}

func (executor *DefaultExecutor) GetName() string {
	return executor.name
}

var _ IExecutor = (*DefaultExecutor)(nil)

var Exec = NewDFExecutor("default")
