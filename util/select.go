package core

import (
	"context"
	"reflect"

	. "github.com/bobg/go-generics/slices"
)

type SelectMode int

const (
	Send SelectMode = iota
	Recv
)

type SelectParam2 struct {
	MaxBatchSize int
	Mode         SelectMode
}

var DefaultSelectParam2 = SelectParam2{
	MaxBatchSize: 64,
	Mode:         Recv,
}

type SelectParam struct {
	Mode SelectMode
}

var DefaultSelectParam = SelectParam{
	Mode: Recv,
}

/*
	func select4Recv[T any](ctx context.Context, chan1 []chan T, chan2 chan T, i int) {
		select {
		case v, _ := <-chan1[0]:
			chan2 <- v
		case v, _ := <-chan1[1]:
			i = i + 1
			chan2 <- v
		case v, _ := <-chan1[2]:
			i = i + 2
			chan2 <- v
		case v, _ := <-chan1[3]:
			i = i + 3
			chan2 <- v
		case <-ctx.Done():
			break
		}
	}

	func select4Send[T any](ctx context.Context, chanz []chan T, i int) {
		select {
		case r.v, r.ok = <-chanz[0]:
			r.i = i + 0
			res <- r
		case r.v, r.ok = <-chanz[1]:
			r.i = i + 1
			res <- r
		case r.v, r.ok = <-chanz[2]:
			r.i = i + 2
			res <- r
		case r.v, r.ok = <-chanz[3]:
			r.i = i + 3
			res <- r
		case <-ctx.Done():
			break
		}
	}

	func select2Send[T any](ctx context.Context, chanz []chan T, i int) {
		select {
		case r.v, r.ok = <-chanz[0]:
			r.i = i + 0
			res <- r
		case r.v, r.ok = <-chanz[1]:
			r.i = i + 1
			res <- r
		case <-ctx.Done():
			break
		}
	}

	func select2Recv[T any](ctx context.Context, chanz []chan T, i int) {
		select {
		case r.v, r.ok = <-chanz[0]:
			r.i = i + 0
			res <- r
		case r.v, r.ok = <-chanz[1]:
			r.i = i + 1
			res <- r
		case <-ctx.Done():
			break
		}
	}
*/
func select1Send[T any](ctx context.Context, chanz chan T, chan2 chan T, i int) {
	select {
	case v, _ := <-chanz:
		chan2 <- v
	case <-ctx.Done():
		break
	}
}
func select1Recv[T any](ctx context.Context, chanz chan T, chan2 chan T, i int) {
	select {
	case v, _ := <-chanz:
		chan2 <- v
	case <-ctx.Done():
		break
	}
}

/*
	func Select[T any](channels []chan T, chan2 chan T, param ...SelectParam) {
		var _param SelectParam
		if len(param) < 1 {
			_param = DefaultSelectParam
		} else {
			_param = param[0]
		}
		//maxBatchSize := _param.MaxBatchSize
		mode := _param.Mode
		ctx := context.Background()
		var l, i int = 0, 0

		for i < len(channels) {
			l = len(channels) - i
			switch {
				case l > 31 && maxBatchSize >= 32:
					if mode == Recv {
						go select32Recv(ctx, channels[i:i+32], chan2, i)
					} else {
						go select32Send(ctx, channels[i:i+32], chan2, i)
					}
					i += 32
				case l > 15 && maxBatchSize >= 16:
					if mode == Recv {
						go select16Recv(ctx, channels[i:i+16], chan2, i)
					} else {
						go select32Send(ctx, channels[i:i+16], chan2, i)
					}
					i += 16
				case l > 7 && maxBatchSize >= 8:
					if mode == Recv {
						go select8Recv(ctx, channels[i:i+8], chan2, i)
					} else {
						go select8Send(ctx, channels[i:i+8], chan2, i)
					}
					i += 8
				case l > 3 && maxBatchSize >= 4:
					if mode == Recv {
						go select4Recv(ctx, channels[i:i+4], chan2, i)
					} else {
						go select4Send(ctx, channels[i:i+4], chan2, i)
					}
					i += 4
				case l > 1 && maxBatchSize >= 2:
					if mode == Recv {
						go select2Recv(ctx, channels[i:i+2], chan2, i)
					} else {
						go select2Send(ctx, channels[i:i+2], chan2, i)
					}
					i += 2
			case l > 0:
				if mode == Recv {
					go select1Recv(ctx, channels[i], chan2, i)
				} else {
					go select1Send(ctx, channels[i], chan2, i)
				}
				i += 1
			}
		}
	}
*/
func GatherRecv[T any](chans []chan T, callback func(int, T, bool)) {
	_chans := Dup[chan T](chans)
	cases := make([]reflect.SelectCase, len(_chans))
	for i, ch := range _chans {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	for len(_chans) > 0 {
		chosen, value, ok := reflect.Select(cases)
		// ok will be true if the channel has not been closed.
		if !ok {
			_chans = RemoveN[chan T](_chans, chosen, 1)
			cases = RemoveN[reflect.SelectCase](cases, chosen, 1)
			continue
		}
		msg := value.Interface().(T)
		callback(chosen, msg, ok)
	}
}

func GatherSend[T any](chans []chan T, v T) {
	_chans := Dup[chan T](chans)
	cases := make([]reflect.SelectCase, len(_chans))
	for i, ch := range _chans {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectSend, Chan: reflect.ValueOf(ch), Send: reflect.ValueOf(v)}
	}
	for len(_chans) > 0 {
		chosen, _, ok := reflect.Select(cases)
		// ok will be true if the channel has not been closed.
		if !ok {
			_chans = RemoveN[chan T](_chans, chosen, 1)
			cases = RemoveN[reflect.SelectCase](cases, chosen, 1)
			continue
		}
	}
}
func Select[T any](chans []chan T, param SelectParam) (T, bool, chan T) {
	cases := make([]reflect.SelectCase, len(chans))
	for i, ch := range chans {
		if param.Mode == Recv {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
		} else {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectSend, Chan: reflect.ValueOf(ch)}
		}
	}
	i, value, ok := reflect.Select(cases)
	msg := value.Interface().(T)
	ch := chans[i]
	return msg, ok, ch
}
