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
	Start(ctx xcontext.CancelableContext)
	Run(ctx xcontext.CancelableContext)
	Terminate(ctx context.Context)
	Stop()
}

// Run runs an Application by the application lifecycle with timeouts and terminate signals.
// It returns false if the quit timeout occurs.
func Run(app Application, terminateTimeout, quitTimeout time.Duration, terminateSignals ...os.Signal) bool {
	return RunAll([]Application{app}, terminateTimeout, quitTimeout, terminateSignals...)
}

// RunAll runs all Application's in common Context by the application lifecycle with timeouts and terminate signals.
// It returns false if the quit timeout occurs.
func RunAll(apps []Application, terminateTimeout, quitTimeout time.Duration, terminateSignals ...os.Signal) bool {
	ctx := xcontext.WithCancelable2(context.Background())
	defer ctx.Cancel()
	terminateCtx, terminateCtxCancel := xcontext.DelayAfterContext(ctx, terminateTimeout)
	defer terminateCtxCancel()
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, terminateSignals...)
		<-ch
		ctx.Cancel()
	}()

	quittedCh := make(chan struct{})
	go func() {
		lifecycle(ctx, apps, terminateCtx)
		close(quittedCh)
	}()
	select {
	case <-quittedCh:
		return true
	case <-xcontext.DelayAfterContext2(terminateCtx, quitTimeout).Done():
		return false
	}
}

func lifecycle(ctx xcontext.CancelableContext, apps []Application, terminateCtx context.Context) {
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
}
