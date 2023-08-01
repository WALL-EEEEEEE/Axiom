package stream

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
