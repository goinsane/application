package application

import (
	"context"

	"github.com/goinsane/xcontext"
)

// Instance is a method wrapper to implement Application interface.
type Instance struct {
	StartFunc     func(ctx xcontext.CancelableContext)
	RunFunc       func(ctx xcontext.CancelableContext)
	TerminateFunc func(ctx context.Context)
	StopFunc      func()
}

func (a *Instance) Start(ctx xcontext.CancelableContext) {
	if a.StartFunc != nil {
		a.StartFunc(ctx)
	}
}

func (a *Instance) Run(ctx xcontext.CancelableContext) {
	if a.RunFunc != nil {
		a.RunFunc(ctx)
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
