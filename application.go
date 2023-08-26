// Package application offers simple application lifecycle framework.
package application

import (
	"context"
	"sync"
	"time"
)

// Application is an interface for handling application lifecycle.
type Application interface {
	Start(ctx context.Context, cancel context.CancelFunc)
	Run(ctx context.Context, cancel context.CancelFunc)
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
	stopped := make(chan struct{})
	go lifecycle(ctx, apps, terminateTimeout, stopped)
	<-ctx.Done()
	quitCtx, quitCancel := context.WithTimeout(context.Background(), quitTimeout)
	defer quitCancel()
	select {
	case <-stopped:
		return true
	case <-quitCtx.Done():
		return false
	}
}

func lifecycle(ctx context.Context, apps []Application, terminateTimeout time.Duration, stopped chan struct{}) {
	defer close(stopped)

	var wg sync.WaitGroup

	appCtx, appCancel := context.WithCancel(ctx)
	defer appCancel()

	for _, app := range apps {
		wg.Add(1)
		go func(app Application) {
			defer wg.Done()
			app.Start(appCtx, appCancel)
		}(app)
	}
	wg.Wait()

	if appCtx.Err() == nil {
		for _, app := range apps {
			wg.Add(1)
			go func(app Application) {
				defer wg.Done()
				app.Run(appCtx, appCancel)
			}(app)
		}
		<-appCtx.Done()
		func() {
			termCtx, termCancel := context.WithTimeout(context.Background(), terminateTimeout)
			defer termCancel()
			for _, app := range apps {
				wg.Add(1)
				go func(app Application) {
					defer wg.Done()
					app.Terminate(termCtx)
				}(app)
			}
			wg.Wait()
		}()
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
