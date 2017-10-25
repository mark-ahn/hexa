package hexa

import (
	"context"
	"log"
	"os/exec"
)

func NewCommand(parent context.Context, name string, args ...string) StoppableOne {
	own_ctx := NewDContextToStoppable(parent)
	cmd := exec.CommandContext(own_ctx, name, args...)
	go func() {
		defer own_ctx.InClose()
		err := cmd.Run()
		if err != nil {
			log.Printf("%+v", err)
		}
	}()
	return own_ctx
}
