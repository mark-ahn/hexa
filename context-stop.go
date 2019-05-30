package hexa

import (
	"context"
	"sync/atomic"
	"time"
	"unsafe"
)

// ContextStop is a imlementation of StoppableOne interface.
// It also compatible to context.Context interface.
// ContextStop is intended to ease offering StoppableOne interface.
//
// ContextStop has two context interally. One is to receive external close request,
// (by Close() method), another is to inform that the parallel rotine has done.
//
// Close() => causes close internal context => the parallel logic checks
// that interal context is closed => the paralle logic close the external context
// => the client code can detect it by DoneNotify()
// Refer to NewStopContext example.
type ContextStop struct {
	inCtx, exCtx           context.Context
	inCanceler, exCanceler context.CancelFunc
	err                    unsafe.Pointer
}

func NewContextStop(parent context.Context) *ContextStop {
	inCtx, exCanceler := context.WithCancel(parent)
	exCtx, inCanceler := context.WithCancel(parent)
	return &ContextStop{
		inCtx:      inCtx,
		exCtx:      exCtx,
		inCanceler: inCanceler,
		exCanceler: exCanceler,
	}
}

func (__ *ContextStop) Close() {
	__.exCanceler()
}

func (__ *ContextStop) DoneNotify() <-chan struct{} {
	return __.exCtx.Done()
}

func (__ *ContextStop) Err() error {
	e := atomic.LoadPointer(&__.err)
	switch e {
	case nil:
		return __.exCtx.Err()
	default:
		return *(*error)(e)
	}
}

func (__ *ContextStop) InBreak(err error) {
	// only the first error is stored.
	atomic.CompareAndSwapPointer(&__.err, nil, unsafe.Pointer(&err))
	__.exCanceler()
}
func (__ *ContextStop) InClose() {
	// make sure exCandeler() is called.
	__.InBreak(nil)

	__.inCanceler()
}

func (__ *ContextStop) InDoneNotify() <-chan struct{} {
	return __.inCtx.Done()
}

func (__ *ContextStop) Deadline() (deadline time.Time, ok bool) {
	return __.exCtx.Deadline()
}
func (__ *ContextStop) Done() <-chan struct{} {
	return __.DoneNotify()
}

// func (__ *ContextStop) Err() error {
// 	return __.exCtx.Err()
// }

func (__ *ContextStop) Value(key interface{}) interface{} {
	return nil
}
