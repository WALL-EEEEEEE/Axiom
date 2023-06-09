package core

import (
	"sync"

	. "github.com/WALL-EEEEEEE/Axiom/core/stream"
)

type Stream struct {
	PassThrough
	name string
}

func NewStream(name string) Stream {
	return Stream{name: name}
}

func (stream Stream) GetName() string {
	return stream.name
}

type Collector chan interface{}

type IExecutor interface {
	Start()
	GetName() string
	Add(task IRunnable)
	List() []string
}

type IRunnable interface {
	Serv
	Run(collector *Collector)
	GetName() string
}

type Task struct {
	stream Stream
	name   string
}

func NewTask(name string) Task {
	return Task{name: name, stream: NewStream(name)}
}

func (task Task) GetName() string {
	return task.name
}

func (task Task) GetType() []ServType {
	return []ServType{TASK}
}

func (task Task) Run(collector *Collector) {
	task.stream.In() <- task.name
}

var _ IRunnable = (*Task)(nil)

type DefaultExecutor struct {
	name      string
	tasks     []IRunnable
	collector *Collector
	broker    *Broker[interface{}]
}

func NewDFExecutor(name string) DefaultExecutor {
	broker := NewBroker[interface{}]()
	collector := Collector(broker.publishCh)
	exec := DefaultExecutor{name: name, collector: &collector, broker: broker}
	return exec
}

func (executor DefaultExecutor) startTask(wg *sync.WaitGroup) {
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

func (executor DefaultExecutor) Start() {
	go executor.broker.Start()
	var task_wg sync.WaitGroup
	executor.startTask(&task_wg)
	task_wg.Wait()
	executor.broker.Close()
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
