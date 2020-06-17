package application

var (
	version string
	build   string
)

// Version returns application.version defined with ldflags.
func Version() string {
	return version
}

// Build returns application.build defined with ldflags.
func Build() string {
	return build
}
