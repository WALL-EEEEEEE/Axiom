package core

import (
	"github.com/bobg/go-generics/maps"
	log "github.com/sirupsen/logrus"
)

type Broker[T any] struct {
	stopCh    chan struct{}
	publishCh chan T
	subCh     chan chan T
	unsubCh   chan chan T
	subs      map[chan T]struct{}
}

func NewBroker[T any]() *Broker[T] {
	return &Broker[T]{
		stopCh:    make(chan struct{}),
		publishCh: make(chan T, 1),
		subCh:     make(chan chan T, 1),
		unsubCh:   make(chan chan T, 1),
		subs:      map[chan T]struct{}{},
	}
}

func (b *Broker[T]) Start() {
	broadcast_cnt := 1
	for {
		select {
		case <-b.stopCh:
			log.Debugf("Broker broadcast: %d", broadcast_cnt)
			close(b.publishCh)
			for subch := range b.subs {
				close(subch)
			}
			return
		case msgCh := <-b.subCh:
			b.subs[msgCh] = struct{}{}
		case msgCh := <-b.unsubCh:
			delete(b.subs, msgCh)
		case msg := <-b.publishCh:
			broadcast_cnt++
			resi_subs := maps.Dup(b.subs)
			for len(resi_subs) > 0 {
				for msgCh := range resi_subs {
					// msgCh is buffered, use non-blocking send to protect the broker:
					if len(resi_subs) == 1 {
						msgCh <- msg
						delete(resi_subs, msgCh)
					} else {
						select {
						case msgCh <- msg:
							delete(resi_subs, msgCh)
						default:
						}
					}

				}
			}
		}
	}
}

func (b *Broker[T]) Close() {
	close(b.stopCh)
}

func (b *Broker[T]) Wait() {
	<-b.stopCh
}

func (b *Broker[T]) Subscribe() chan T {
	msgCh := make(chan T, 1)
	b.subCh <- msgCh
	return msgCh
}

func (b *Broker[T]) Unsubscribe(msgCh chan T) {
	b.unsubCh <- msgCh
}

func (b *Broker[T]) Publish(msg T) {
	b.publishCh <- msg
}
