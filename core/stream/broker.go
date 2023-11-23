package stream

import (
	"context"
	. "context"
	"fmt"

	. "github.com/WALL-EEEEEEE/Axiom/util"
	"github.com/bobg/go-generics/maps"
	"github.com/bobg/go-generics/slices"
	log "github.com/sirupsen/logrus"
)

type Broker[T any] struct {
	name      string
	stopCh    chan struct{}
	publishCh Stream[T]
	subCh     chan struct {
		Stream[T]
		CancelFunc
	}
	unsubCh chan Stream[T]
	subs    map[Stream[T]]struct{}
}

func NewBroker[T any](name string) Broker[T] {
	broker := Broker[T]{
		name:      name,
		stopCh:    make(chan struct{}),
		publishCh: NewStream[T]("publishCh"),
		subCh: make(chan struct {
			Stream[T]
			CancelFunc
		}),
		unsubCh: make(chan Stream[T]),
		subs:    map[Stream[T]]struct{}{},
	}
	go broker.doStream()
	return broker
}

func (b *Broker[T]) Via(stream Stream[T]) {
	b.publishCh.From(&stream)
}

func (b *Broker[T]) doStream() {
	broadcast_cnt := 1
loop:
	for {
		select {
		case <-b.stopCh:
			b.publishCh.Close()
		case subch := <-b.subCh:
			_chan, cancelFunc := subch.Stream, subch.CancelFunc
			b.subs[_chan] = struct{}{}
			cancelFunc()
		case msgCh := <-b.unsubCh:
			delete(b.subs, msgCh)
		case msg, ok := <-b.publishCh.tunnel.out:
			broadcast_cnt++
			if !ok {
				break loop
			}
			log.Debugf("Broker [%s] Recv: %+v", b.name, msg)
			sub_channels, _ := slices.Map(maps.Keys[Stream[T]](b.subs), func(i int, v Stream[T]) (chan T, error) {
				log.Debugf("Broker [%s] Send to Stream [%s]: %+v", b.name, v.GetName(), msg)
				return v.tunnel.in, nil
			})
			GatherSend(sub_channels, msg)
		}
	}
	for subch := range b.subs {
		subch.Close()
	}
}

func (b *Broker[T]) Close() {
	close(b.stopCh)
}

func (b *Broker[T]) Wait() {
	<-b.stopCh
}

func (b *Broker[T]) Subscribe() Stream[T] {
	name := fmt.Sprintf("%s-out-%s", b.name, "0")
	msgCh := NewStream[T](name)
	context, cancel := WithCancel(context.TODO())
	b.subCh <- struct {
		Stream[T]
		CancelFunc
	}{msgCh, cancel}
	<-context.Done()
	return msgCh
}

func (b *Broker[T]) Unsubscribe(msgCh Stream[T]) {
	b.unsubCh <- msgCh
}

func (b *Broker[T]) GetInputStream() Stream[T] {
	return b.publishCh
}

func (b *Broker[T]) GetOutputStream() Stream[T] {
	return b.Subscribe()
}

func (b *Broker[T]) Send(msg T) {
	b.publishCh.tunnel.in <- msg
}
