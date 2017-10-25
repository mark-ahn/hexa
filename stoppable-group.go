package hexa

import (
	"fmt"
	"reflect"
	"sync"
)

type StoppableGroup struct {
	group     map[string]StoppableOne
	index_map []string
	cases     []reflect.SelectCase
	lock      *sync.Mutex
}

func NewStoppableGroup(group map[string]StoppableOne) *StoppableGroup {
	index_map, cases := make_selectors_of_map(group)
	return &StoppableGroup{
		lock:      &sync.Mutex{},
		group:     group,
		index_map: index_map,
		cases:     cases,
	}
}

func (__ *StoppableGroup) Push(name string, one StoppableOne) error {
	__.lock.Lock()
	defer __.lock.Unlock()
	_, dup := __.group[name]
	if dup {
		return fmt.Errorf("already assigned name '%v'", name)
	}

	__.group[name] = one

	return nil
}

func make_selector(v StoppableOne) reflect.SelectCase {
	return reflect.SelectCase{
		Chan: reflect.ValueOf(v.DoneNotify()),
		Dir:  reflect.SelectRecv,
	}
}

func make_selectors_of_map(stoppableMap map[string]StoppableOne) ([]string, []reflect.SelectCase) {
	index_map := make([]string, 0, len(stoppableMap))
	res := make([]reflect.SelectCase, 0, len(stoppableMap))
	for k, v := range stoppableMap {
		res = append(res, make_selector(v))
		index_map = append(index_map, k)
	}
	return index_map, res
}
