package stream

import (
	"github.com/bobg/go-generics/maps"
	log "github.com/sirupsen/logrus"
)

type Broker[T any] struct {
	stopCh    chan struct{}
	publishCh Stream[T]
	subCh     chan Stream[T]
	unsubCh   chan Stream[T]
	subs      map[Stream[T]]struct{}
}

func NewBroker[T any]() Broker[T] {
	return Broker[T]{
		stopCh:    make(chan struct{}),
		publishCh: NewStream[T]("publishCh"),
		subCh:     make(chan Stream[T], 1),
		unsubCh:   make(chan Stream[T], 1),
		subs:      map[Stream[T]]struct{}{},
	}
}

func (b *Broker[T]) Via(stream Stream[T]) {
	b.publishCh = stream
}

func (b *Broker[T]) Start() {
	broadcast_cnt := 1
	for {
		select {
		case <-b.stopCh:
			log.Debugf("Broker broadcast: %d", broadcast_cnt)
			b.publishCh.Close()
			for subch := range b.subs {
				subch.Close()
			}
			return
		case msgCh := <-b.subCh:
			b.subs[msgCh] = struct{}{}
		case msgCh := <-b.unsubCh:
			delete(b.subs, msgCh)
		case msg := <-b.publishCh.Out():
			broadcast_cnt++
			resi_subs := maps.Dup(b.subs)
			for len(resi_subs) > 0 {
				for msgCh := range resi_subs {
					// msgCh is buffered, use non-blocking send to protect the broker:
					if len(resi_subs) == 1 {
						msgCh.In() <- msg
						delete(resi_subs, msgCh)
					} else {
						select {
						case msgCh.In() <- msg:
							delete(resi_subs, msgCh)
						default:
						}
					}

				}
			}
		}
	}
}

func (b Broker[T]) Close() {
	close(b.stopCh)
}

func (b Broker[T]) Wait() {
	<-b.stopCh
}

func (b Broker[T]) Subscribe() Stream[T] {
	msgCh := NewStream[T]("0")
	b.subCh <- msgCh
	return msgCh
}

func (b Broker[T]) Unsubscribe(msgCh Stream[T]) {
	b.unsubCh <- msgCh
}

func (b *Broker[T]) Publish(msg T) {
	b.publishCh.In() <- msg
}
