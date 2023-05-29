package core

type IPipe interface {
	Serv
	Run(collector *Collector)
}

type Pipe struct {
	name string
}

func NewPipe(name string) Pipe {
	return Pipe{name: name}
}

func (pipe *Pipe) GetName() string {
	return pipe.name
}

func (pipe *Pipe) GetType() []ServType {
	return []ServType{PIPE}
}
