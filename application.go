// Package application offers simple application lifecycle framework.
package application

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"
)

// Application is an interface for handling application lifecycle.
type Application interface {
	Start(ctx Context)
	Run()
	Terminate(ctx context.Context)
	Stop()
}

// Context is a custom implementation of context.Context with Terminate() method to terminate application.
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

// Run runs an Application by application lifecycle with terminateTimeout and terminateSignals.
func Run(app Application, terminateTimeout time.Duration, terminateSignals ...os.Signal) {
	ctx := new(applicationContext)
	ctx.Context, ctx.CancelFunc = context.WithCancel(context.Background())
	defer ctx.Terminate()
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, terminateSignals...)
		<-ch
		ctx.Terminate()
	}()
	app.Start(ctx)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		app.Run()
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
