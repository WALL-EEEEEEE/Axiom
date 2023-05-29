package core

import (
	"sync"
)

type Collector chan interface{}

type IExecutor interface {
	Start()
	GetName() string
	Add(task IRunnable)
	AddPipe(pipe Pipe)
	List() []string
}

type IRunnable interface {
	Serv
	Run(collector *Collector)
	GetName() string
}

type Task struct {
	name string
}

func NewTask(name string) Task {
	return Task{name: name}
}

func (task *Task) GetName() string {
	return task.name
}

func (task *Task) GetType() []ServType {
	return []ServType{TASK}
}

type DefaultExecutor struct {
	name      string
	tasks     []IRunnable
	pipes     []IPipe
	collector *Collector
	broker    *Broker[interface{}]
}

func NewDFExecutor(name string) DefaultExecutor {
	broker := NewBroker[interface{}]()
	collector := Collector(broker.publishCh)
	exec := DefaultExecutor{name: name, collector: &collector, broker: broker}
	return exec
}

func (executor *DefaultExecutor) startTask(wg *sync.WaitGroup) {
	_start := func(task IRunnable) {
		defer func() {
			wg.Done()
		}()
		task.Run(executor.collector)
	}
	for _, task := range executor.tasks {
		wg.Add(1)
		go _start(task)
	}
}

func (executor *DefaultExecutor) startPipe(wg *sync.WaitGroup) {
	_start := func(pipe IPipe) {
		defer func() {
			wg.Done()
		}()
		collector := Collector(executor.broker.Subscribe())
		pipe.Run(&collector)
	}
	for _, pipe := range executor.pipes {
		wg.Add(1)
		go _start(pipe)
	}
}

func (executor *DefaultExecutor) Start() {
	go executor.broker.Start()
	var task_wg sync.WaitGroup
	var pipe_wg sync.WaitGroup
	executor.startTask(&task_wg)
	executor.startPipe(&pipe_wg)
	task_wg.Wait()
	executor.broker.Close()
	pipe_wg.Wait()
}

func (executor *DefaultExecutor) Add(task IRunnable) {
	executor.tasks = append(executor.tasks, task)
}

func (executor *DefaultExecutor) AddPipe(pipe IPipe) {
	executor.pipes = append(executor.pipes, pipe)
}
func (executor *DefaultExecutor) List() []string {
	var names []string
	for _, task := range executor.tasks {
		names = append(names, task.GetName())
	}
	return names
}

var Exec = NewDFExecutor("default")
