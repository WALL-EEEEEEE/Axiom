package stream

type Inlet[T any] interface {
	In() chan<- T
}

type Outlet[T any] interface {
	Out() <-chan T
}

type Sink[T any] interface {
	Inlet[T]
}
type Source[T any] interface {
	Outlet[T]
	Via(Flow[T]) Flow[T]
}
type Flow[T any] interface {
	Inlet[T]
	Outlet[T]
	From(Flow[T]) Flow[T]
	To(Sink[T])
}
