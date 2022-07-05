# application

[![Go Reference](https://pkg.go.dev/badge/github.com/goinsane/application.svg)](https://pkg.go.dev/github.com/goinsane/application)

Package application offers simple application lifecycle framework.

## Application's lifecycle

This package provides an Application interface to run applications in a lifecycle. The interface has four methods:
Start, Run, Terminate, Stop.

- **Start** method always runs the beginning of the lifecycle even the termination process was triggered.
- **Run** method runs immediately after Start method.
The termination process may be triggered before entering into Run method.
While in Run method, Terminate method may be running at the same time.
- **Terminate** method runs immediately after the termination process is triggered.
Terminate method always runs after Start method.
While in Terminate method, Run method may be running at the same time.
- **Stop** method runs immediately after Run and Terminate methods were returned.
The lifecycle ends with Stop method was returned.

## Application Context

Start and Run methods have **ctx** argument that is *Context* from *CancelableContext* of *xcontext*.
These contextes are the exactly same and **ctx** can be named as ***application context***.
All application contextes are the exactly same at all of Applications in a single lifecycle.
The calling *Cancel* method of the application context cancels the context, and starts the termination process.

## Termination process

Termination process can be triggered with the canceling of application contexts or termination signals.
When termination process started, Terminate method always will be called with a *Context* has the given terminate timeout.
If termination process started before Run method, it waits the ending of Start method to call Terminate method.
Otherwise it calls Terminate method immediately.

## Timeouts

Once the termination process has been started, there is no way back to initial state.
It waits methods of Application interface or timeouts.
Terminate method's context will be cancelled after the given terminate timeout or, when Run and Terminate methods were ended.
The waiting of the quit timeout is started after Terminate method's context were cancelled.

## Building with version

You can save your application's name, version and build into this package below this.

    go build -ldflags "\
    -X=github.com/goinsane/application.name=$NAME \
    -X=github.com/goinsane/application.version=$VERSION \
    -X=github.com/goinsane/application.build=$BUILD \
    " ./...
