package hexa

import (
	"context"
	"time"
)

type DContextToStoppable struct {
	inCtx, exCtx           context.Context
	inCanceler, exCanceler context.CancelFunc
}

func NewDContextToStoppable(parent context.Context) *DContextToStoppable {
	inCtx, exCancel := context.WithCancel(parent)
	exCtx, inCancel := context.WithCancel(parent)
	return &DContextToStoppable{
		inCtx, exCtx, inCancel, exCancel,
	}
}

func (__ *DContextToStoppable) Close() {
	__.exCanceler()
}

func (__ *DContextToStoppable) DoneNotify() <-chan struct{} {
	return __.exCtx.Done()
}

func (__ *DContextToStoppable) InClose() {
	__.inCanceler()
}

func (__ *DContextToStoppable) InDoneNotify() <-chan struct{} {
	return __.inCtx.Done()
}

func (__ *DContextToStoppable) Deadline() (deadline time.Time, ok bool) {
	return __.exCtx.Deadline()
}
func (__ *DContextToStoppable) Done() <-chan struct{} {
	return __.DoneNotify()
}
func (__ *DContextToStoppable) Err() error {
	return __.exCtx.Err()
}
func (__ *DContextToStoppable) Value(key interface{}) interface{} {
	return nil
}
