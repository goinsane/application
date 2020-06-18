package application

var (
	name    string
	version string
	build   string
)

// Name returns application.name defined with ldflags.
func Name() string {
	return name
}

// Version returns application.version defined with ldflags.
func Version() string {
	return version
}

// Build returns application.build defined with ldflags.
func Build() string {
	return build
}
