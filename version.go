package application

var (
	version string
	build   string
)

func Version() string {
	return version
}

func Build() string {
	return build
}

