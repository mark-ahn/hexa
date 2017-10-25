package hexa

import (
	"context"
	"io"
	"log"
	"os/exec"
)

type StdPlumber interface {
	StderrPipe() (io.ReadCloser, error)
	StdinPipe() (io.WriteCloser, error)
	StdoutPipe() (io.ReadCloser, error)
}

func NewCommand(parent context.Context, pipeHandler func(plumber StdPlumber), name string, args ...string) StoppableOne {
	own_ctx := NewDContextToStoppable(parent)
	cmd := exec.CommandContext(own_ctx, name, args...)

	switch pipeHandler {
	case nil:
	default:
		pipeHandler(cmd)
	}

	go func() {
		defer own_ctx.InClose()
		err := cmd.Run()
		if err != nil {
			log.Printf("%+v", err)
		}
	}()
	return own_ctx
}
