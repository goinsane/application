# application

[![Go Reference](https://pkg.go.dev/badge/github.com/goinsane/application.svg)](https://pkg.go.dev/github.com/goinsane/application)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=goinsane_application&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=goinsane_application)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=goinsane_application&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=goinsane_application)

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
