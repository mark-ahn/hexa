package hexa_test

import (
	"context"
	"fmt"
	"hexa"
	"time"
)

func ExampleNewContextStop() {
	some_service := func() hexa.StoppableOne {
		__ := hexa.NewContextStop(context.Background())

		// run some work in parallel
		go func() {
			defer func() {
				// InClose() causes DoneNotify() channel close.
				// So client code can detect done when go routine has entirely done.
				__.InClose()
			}()
		loop:
			for {
				select {
				// break loop when client calls Close() method.
				case <-__.InDoneNotify():
					break loop

				// do some parallel work in service.
				case <-time.After(time.Microsecond):
				}
			}
		}()
		return __
	}

	srv := some_service()
	defer func() {
		srv.Close()
		<-srv.DoneNotify()
		fmt.Printf("service is terminated\n")
	}()

	fmt.Printf("do some work\n")
	// Output:
	// do some work
	// service is terminated
}
