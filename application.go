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

// Run runs an instance of Application by the application lifecycle with timeouts and terminate signals.
// It returns false if the quit timeout occurs.
func Run(app Application, terminateTimeout, quitTimeout time.Duration, terminateSignals ...os.Signal) bool {
	return RunAll([]Application{app}, terminateTimeout, quitTimeout, terminateSignals...)
}

// RunAll runs all instances of Application in common Context by the application lifecycle with timeouts and terminate signals.
// It returns false if the quit timeout occurs.
func RunAll(apps []Application, terminateTimeout, quitTimeout time.Duration, terminateSignals ...os.Signal) bool {
	appCtx := xcontext.WithCancelable2(context.Background())
	defer appCtx.Cancel()

	termCtx := xcontext.WithCancelable2(xcontext.DelayAfterContext2(appCtx, terminateTimeout))
	defer termCtx.Cancel()

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, terminateSignals...)
		<-ch
		appCtx.Cancel()
	}()

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

	for _, app := range apps {
		wg.Add(1)
		go func(app Application) {
			defer wg.Done()
			app.Stop()
		}(app)
	}
	wg.Wait()
}
