package hexa

import (
	"context"
	"log"
	"reflect"
)

type StoppableFactoryInfo struct {
	Name    string
	Factory func() StoppableOne
}

type StoppableSpawner struct {
	factory_list []func() StoppableOne
	cases        []reflect.SelectCase
	stoppables   []StoppableOne
	ctx          context.Context
}

func assign(stoppables []StoppableOne, cases []reflect.SelectCase, factory func() StoppableOne, i int) {
	stoppables[i] = factory()
	cases[i] = reflect.SelectCase{
		Chan: reflect.ValueOf(stoppables[i].DoneNotify()),
		Dir:  reflect.SelectRecv,
	}
}

func NewStoppableSpawner(ctx context.Context, factoryList []func() StoppableOne) *StoppableSpawner {
	cases := make([]reflect.SelectCase, len(factoryList))
	stoppables := make([]StoppableOne, len(factoryList))
	defer func() {
		r := recover()
		if r != nil {
			log.Printf("Recovered for tear-down: %v", r)
			tear_down(stoppables)
			panic(r)
		}
		// log.Printf("no panic %v", r)
	}()
	for i, factory := range factoryList {
		assign(stoppables, cases, factory, i)
	}

	return &StoppableSpawner{
		factory_list: factoryList,
		cases:        cases,
		stoppables:   stoppables,
		ctx:          ctx,
	}
}

func (__ *StoppableSpawner) Serve(respawnHandler func(int)) error {
	for {
		i, _, _ := reflect.Select(append(__.cases,
			reflect.SelectCase{
				Chan: reflect.ValueOf(__.ctx.Done()),
				Dir:  reflect.SelectRecv,
			}))
		recv_ctx := len(__.stoppables) <= i
		if recv_ctx {
			// parent Context
			tear_down(__.stoppables)
			return __.ctx.Err()
		}

		__.stoppables[i].Close()
		switch respawnHandler {
		case nil:
		default:
			respawnHandler(i)
		}
		factory := __.factory_list[i]

		assign(__.stoppables, __.cases, factory, i)
	}
}

func tear_down(stoppables []StoppableOne) {
	for i := len(stoppables) - 1; i >= 0; i -= 1 {
		stoppable := stoppables[i]
		switch stoppable {
		case nil:
		default:
			log.Printf("[HEXA] tear down %vth element", i)
			stoppable.Close()
			<-stoppable.DoneNotify()
		}
	}
}
