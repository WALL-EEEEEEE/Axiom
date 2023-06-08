package core

type Inlet interface {
	In() chan<- interface{}
}

type Outlet interface {
	Out() <-chan interface{}
}

type Sink interface {
	Inlet
}
type Source interface {
	Outlet
	Via(Flow) Flow
}
type Flow interface {
	Inlet
	Outlet
	Via(Flow) Flow
	To(Sink)
}
