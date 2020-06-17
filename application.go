package application

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Application interface {
	Start()
	Run(ctx Context)
	Terminate(ctx context.Context)
	Stop()
}

type Context interface {
	context.Context
	Terminate()
}

type applicationContext struct {
	context.Context
	context.CancelFunc
}

func (c *applicationContext) Terminate() {
	c.CancelFunc()
}

func Run(app Application, terminateTimeout time.Duration, terminateSignal ...os.Signal) {
	app.Start()

	ctx := new(applicationContext)
	ctx.Context, ctx.CancelFunc = context.WithCancel(context.Background())
	defer ctx.Terminate()
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, terminateSignal...)
		<-ch
		ctx.Terminate()
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		app.Run(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		terminateCtx, terminateCtxCancel := context.WithTimeout(context.Background(), terminateTimeout)
		defer terminateCtxCancel()

		app.Terminate(terminateCtx)
	}()

	wg.Wait()
	app.Stop()
}
