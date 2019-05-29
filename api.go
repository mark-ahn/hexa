package hexa

import (
	"sync"
)

// Experimental api...
type ThreadGroup struct {
	cnt sync.WaitGroup
}

// func NewThreadGroup() *ThreadGroup {
// 	return &ThreadGroup{
// 		cnt: sync.WaitGroup{},
// 	}
// }
func (__ *ThreadGroup) SpawnService(stopCtx *ContextStop, ch <-chan interface{}, f func(interface{}) error) {
	__.cnt.Add(1)
	go func() {
		defer __.cnt.Done()
	loop:
		for {
			select {
			// break loop when client calls Close() method.
			case <-stopCtx.InDoneNotify():
				break loop

			// do some parallel work in service.
			case d := <-ch:
				err := f(d)
				if err != nil {
					stopCtx.InBreak(err)
					continue
				}
			}
		}
	}()
}

func (__ *ThreadGroup) Done(stopCtx *ContextStop, f func()) {
	__.cnt.Wait()
	f()
	stopCtx.InClose()
}
