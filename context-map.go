package hexa

import (
	"context"
)

type DContextToStoppable struct {
	inCtx, exCtx           context.Context
	inCanceler, exCanceler context.CancelFunc
}

func NewDContextToStoppable(parent context.Context) *DContextToStoppable {
	inCtx, exCancel := context.WithCancel(parent)
	exCtx, inCancel := context.WithCancel(inCtx)
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
