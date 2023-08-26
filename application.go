// Package application offers simple application lifecycle framework.
package application

import (
	"context"
	"sync"
	"time"
)

// Application is an interface for handling the application lifecycle.
// Run and RunAll functions trigger the application lifecycle, and enter methods of the Application interface with the given order.
type Application interface {
	// Start is always called when the lifecycle starts.
	// ctx is the application context and can be canceled by the cancel function.
	Start(ctx context.Context, cancel context.CancelFunc)
	// Run is called after the Start method if the application context has not been canceled.
	// ctx is the application context and can be canceled by the cancel function.
	Run(ctx context.Context, cancel context.CancelFunc)
	// Terminate is called after the application context was cancelled, with the terminate context which has the given terminate timeout.
	// If the application context is canceled before the Run method, Terminate is not called.
	// ctx is the terminate context.
	Terminate(ctx context.Context)
	// Stop is always called when the lifecycle ends even the Run method was not called.
	Stop()
}

// Run runs an instance of Application by the application lifecycle with the given ctx and timeouts.
// It returns false if the quit timeout occurs.
// Quit timeout has to be greater than terminate timeout. And it starts after the application context was canceled.
func Run(ctx context.Context, app Application, terminateTimeout, quitTimeout time.Duration) bool {
	return RunAll(ctx, []Application{app}, terminateTimeout, quitTimeout)
}

// RunAll runs all instances of Application in common context.Context by the application lifecycle with the given ctx and timeouts.
// It returns false if the quit timeout occurs.
// Quit timeout has to be greater than terminate timeout. And it starts after the application context was canceled.
func RunAll(ctx context.Context, apps []Application, terminateTimeout, quitTimeout time.Duration) bool {
	appCtx, appCancel := context.WithCancel(ctx)
	defer appCancel()

	stopped := make(chan struct{})

	go lifecycle(appCtx, appCancel, apps, terminateTimeout, stopped)

	<-appCtx.Done()

	quitCtx, quitCancel := context.WithTimeout(context.Background(), quitTimeout)
	defer quitCancel()
	select {
	case <-stopped:
		return true
	case <-quitCtx.Done():
		return false
	}
}

func lifecycle(appCtx context.Context, appCancel context.CancelFunc, apps []Application, terminateTimeout time.Duration, stopped chan struct{}) {
	var wg sync.WaitGroup

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
	close(stopped)
}
