# application

[![Go Reference](https://pkg.go.dev/badge/github.com/goinsane/application.svg)](https://pkg.go.dev/github.com/goinsane/application)

Package application offers simple application lifecycle framework.

### Building with version

    go build -ldflags "\
    -X=github.com/goinsane/application.name=$NAME \
    -X=github.com/goinsane/application.version=$VERSION \
    -X=github.com/goinsane/application.build=$BUILD \
    " ./...
