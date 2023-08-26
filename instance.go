package application

import (
	"context"
)

// Instance is a method wrapper to implement Application interface.
type Instance struct {
	StartFunc     func(ctx context.Context, cancel context.CancelFunc)
	RunFunc       func(ctx context.Context, cancel context.CancelFunc)
	TerminateFunc func(ctx context.Context)
	StopFunc      func()
}

func (a *Instance) Start(ctx context.Context, cancel context.CancelFunc) {
	if a.StartFunc != nil {
		a.StartFunc(ctx, cancel)
	}
}

func (a *Instance) Run(ctx context.Context, cancel context.CancelFunc) {
	if a.RunFunc != nil {
		a.RunFunc(ctx, cancel)
	}
}

func (a *Instance) Terminate(ctx context.Context) {
	if a.TerminateFunc != nil {
		a.TerminateFunc(ctx)
	}
}

func (a *Instance) Stop() {
	if a.StopFunc != nil {
		a.StopFunc()
	}
}
