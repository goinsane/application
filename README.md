# application

[![Go Reference](https://pkg.go.dev/badge/github.com/goinsane/application.svg)](https://pkg.go.dev/github.com/goinsane/application)

Package application offers simple application lifecycle framework.

## Application's lifecycle

This package provides an Application interface to run applications in a lifecycle. The interface has four methods:
Start, Run, Terminate, Stop.

- **Start** method always runs the beginning of the lifecycle even the termination process was triggered.
- **Run** method runs immediately after Start method unless the termination process was triggered.
The termination process may be triggered before entering into Run method.
While in Run method, Terminate method may be running at the same time.
- **Terminate** method runs immediately after the termination process is triggered
unless the termination process was triggered before Run method.
Terminate method always runs after Start method.
While in Terminate method, Run method may be running at the same time.
- **Stop** method runs immediately after Run and Terminate methods were returned
unless the termination process was triggered before Run method. Otherwise Stop method runs after Start method.
The lifecycle ends with Stop method was returned.

## Concepts

### Application Context

Start and Run methods have **ctx** argument that is `Context` from `CancelableContext` of `xcontext`.
These contextes are the exactly same and **ctx** can be named as ***Application Context***.
All application contextes are the exactly same at all of Applications in a single lifecycle.
The calling `Cancel` method of the application context cancels the context, and starts the termination process.

### Terminate Context

***Terminate Context*** is a `Context` of `context` and starts with the given terminate timeout
after the termination process was triggered. Terminate context will be cancelled after the given terminate timeout is up.
If the termination process was triggered before Run method, terminate context will be cancelled
after Start method is ended. Otherwise terminate context will be cancelled after Run and Terminate methods are ended.

### Termination Process

Once the termination process has been started, there is no way back to initial state.
It waits methods of Application interface or timeouts.

Termination process can be triggered with the canceling of application contexts or cancelling initial context.
When termination process was started, Terminate method is called with a terminate context after Start method
unless the termination process was triggered before Run method. Otherwise it doesn't call Terminate method.

Terminate timeout describes timeout of terminate context. Terminate context may be cancelled before the timeout.
The waiting of the quit timeout is started after terminate context were cancelled.

## Building with version

You can save your application's name, version and build into this package below this.

    go build -ldflags "\
    -X=github.com/goinsane/application.name=$NAME \
    -X=github.com/goinsane/application.version=$VERSION \
    -X=github.com/goinsane/application.build=$BUILD \
    " ./...
