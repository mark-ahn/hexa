package hexa

import (
	"context"
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