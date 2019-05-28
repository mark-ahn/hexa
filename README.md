# hexa

hexa package is intended to support organizing services (parallel rotines) with ease.
Currently only ContextStop() is recommened to use and others are experimental.

ContextStop is a imlementation of StoppableOne interface.
It also compatible to context.Context interface.
ContextStop is intended to ease offering StoppableOne interface.

ContextStop has two context interally. One is to receive external close request,
(by Close() method), another is to inform that the parallel rotine has done.

Example 1)

```go
func service_factory() hexa.StoppableOne {
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
                    // call InBreak()
                    __.InBreak(err)
                    continue loop
                }
            }
        }
    }()

    return __
}
func ExampleNewContextStop() {
	srv := service_factory()
	defer func() {
        // tear down the service
		srv.Close()
		<-srv.DoneNotify()
		fmt.Printf("service is terminated\n")
	}()

	fmt.Printf("do some work\n")
	// Output:
	// do some work
	// service is terminated
}

```

Example 2)

```go
func service_factory() hexa.StoppableOne {
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
func Example_newContextStopWithMultipleGoRoutine() {
	srv := service_factory()
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
```
