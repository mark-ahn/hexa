package hexa_test

import (
	"context"
	"log"
	"testing"

	"github.com/mark-ahn/hexa"
)

func Test_Teardown(t *testing.T) {
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
