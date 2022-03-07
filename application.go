// Package application offers simple application lifecycle framework.
package application

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/goinsane/xcontext"
)

// Application is an interface for handling application lifecycle.
type Application interface {
	Start(ctx Context)
	Run(ctx Context)
	Terminate(ctx context.Context)
	Stop()
}

// Context is a custom implementation of context.Context with Terminate() method to terminate application.
type Context = xcontext.TerminateContext

// Run runs an Application by application lifecycle with terminateTimeout and terminateSignals.
// It returns terminateCtx as done or not done that used in Application.Terminate.
func Run(app Application, terminateTimeout time.Duration, terminateSignals ...os.Signal) (terminateCtx context.Context) {
	return RunAll([]Application{app}, terminateTimeout, terminateSignals...)
}

// RunAll runs all Application's in common Context by application lifecycle with terminateTimeout and terminateSignals.
// It returns terminateCtx anyway if it is done or not done.
func RunAll(apps []Application, terminateTimeout time.Duration, terminateSignals ...os.Signal) (terminateCtx context.Context) {
	ctx := xcontext.WithTerminate2(context.Background())
	defer ctx.Terminate()

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, terminateSignals...)
		<-ch
		ctx.Terminate()
	}()

	var wg sync.WaitGroup

	for _, app := range apps {
		wg.Add(1)
		go func(app Application) {
			defer wg.Done()
			app.Start(ctx)
		}(app)
	}
	wg.Wait()

	for _, app := range apps {
		wg.Add(1)
		go func(app Application) {
			defer wg.Done()
			app.Run(ctx)
		}(app)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		terminateCtx = xcontext.WithTimeout2(context.Background(), terminateTimeout)

		var wg sync.WaitGroup
		defer wg.Wait()
		for _, app := range apps {
			wg.Add(1)
			go func(app Application) {
				defer wg.Done()
				app.Terminate(terminateCtx)
			}(app)
		}
	}()
	wg.Wait()

	for _, app := range apps {
		wg.Add(1)
		go func(app Application) {
			defer wg.Done()
			app.Stop()
		}(app)
	}
	wg.Wait()

	return terminateCtx
}
