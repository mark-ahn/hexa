package hexa_test

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/mark-ahn/hexa"
)

func ExampleNewContextStop() {
	some_service := func() hexa.StoppableOne {
		__ := hexa.NewContextStop(context.Background())

		// run some work in parallel
		go func() {
			defer func() {
				// InClose() causes the DoneNotify() channel to be closed.
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
					err := func() error { return nil }()
					if err != nil {
						// if want to terminate the service due to such as internal error,
						// call InBreak() then exit the outer loop.
						__.InBreak(err)
						break loop
					}
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

func Example_newContextStopWithMultipleGoRoutine() {
	some_service := func() hexa.StoppableOne {
		__ := hexa.NewContextStop(context.Background())

		// uses wait group to confirm all of the go routines are finished.
		threads := sync.WaitGroup{}

		threads.Add(1)
		go func() {
			defer threads.Done()
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

		threads.Add(1)
		go func() {
			defer threads.Done()
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

		go func() {
			// wait until all of the threads will be finished.
			threads.Wait()
			// InClose() causes the DoneNotify() channel to be closed.
			// So client code can detect done when go routine has entirely done.
			__.InClose()
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
