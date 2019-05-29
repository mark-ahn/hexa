package hexa_test

import (
	"context"
	"log"
	"testing"

	hexa "."
)

func Test_Teardown(t *testing.T) {
	t.SkipNow()
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recover from %v", r)
		}
	}()
	hexa.NewStoppableSpawner(context.Background(), []func() hexa.StoppableOne{
		func() hexa.StoppableOne {
			ctx := hexa.NewDContextToStoppable(context.Background())
			go func() {
				select {
				case <-ctx.InDoneNotify():
					log.Printf("1th func closed")
				}
			}()
			return ctx
		},
		func() hexa.StoppableOne {
			panic("go go go")
		},
	})
}
