// Package application offers simple application lifecycle framework.
package application

import (
	"context"
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

// Run runs an instance of Application by the application lifecycle with the given ctx and timeouts.
// It returns false if the quit timeout occurs.
func Run(ctx context.Context, app Application, terminateTimeout, quitTimeout time.Duration) bool {
	return RunAll(ctx, []Application{app}, terminateTimeout, quitTimeout)
}

// RunAll runs all instances of Application in common context.Context by the application lifecycle with the given ctx and timeouts.
// It returns false if the quit timeout occurs.
func RunAll(ctx context.Context, apps []Application, terminateTimeout, quitTimeout time.Duration) bool {
	appCtx := xcontext.WithCancelable2(ctx)
	defer appCtx.Cancel()

	termCtx := xcontext.WithCancelable2(xcontext.DelayAfterContext2(appCtx, terminateTimeout))
	defer termCtx.Cancel()

	quittedCh := make(chan struct{})
	go func() {
		lifecycle(appCtx, termCtx, apps)
		close(quittedCh)
	}()
	select {
	case <-quittedCh:
		return true
	case <-xcontext.DelayAfterContext2(termCtx, quitTimeout).Done():
		return false
	}
}

func lifecycle(appCtx, termCtx xcontext.CancelableContext, apps []Application) {
	var wg sync.WaitGroup

	for _, app := range apps {
		wg.Add(1)
		go func(app Application) {
			defer wg.Done()
			app.Start(appCtx)
		}(app)
	}
	wg.Wait()

	if appCtx.Err() == nil {
		for _, app := range apps {
			wg.Add(1)
			go func(app Application) {
				defer wg.Done()
				app.Run(appCtx)
			}(app)
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-appCtx.Done()
			for _, app := range apps {
				wg.Add(1)
				go func(app Application) {
					defer wg.Done()
					app.Terminate(termCtx)
				}(app)
			}
		}()
		wg.Wait()
		termCtx.Cancel()
	}

	for _, app := range apps {
		wg.Add(1)
		go func(app Application) {
			defer wg.Done()
			app.Stop()
		}(app)
	}
	wg.Wait()
}
