# application

[![Go Reference](https://pkg.go.dev/badge/github.com/goinsane/application.svg)](https://pkg.go.dev/github.com/goinsane/application)

Package application offers simple application lifecycle framework.

## Running the application lifecycle

You can use `Run` or `RunAll` functions to enter the application lifecycle. Please read go reference.

## Building with version

You can save your application's name, version and build into this package below this.

    go build -ldflags "\
    -X=github.com/goinsane/application.name=$NAME \
    -X=github.com/goinsane/application.version=$VERSION \
    -X=github.com/goinsane/application.build=$BUILD \
    " ./...
