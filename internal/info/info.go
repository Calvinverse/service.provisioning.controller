package info

var (
	sha1      string
	buildTime string
	version   string
)

// BuildTime returns the time the application was built.
func BuildTime() string {
	return buildTime
}

// Revision returns the GIT SHA for the commit off which the application was built.
func Revision() string {
	return sha1
}

// Version returns the version number for the application.
func Version() string {
	return version
}
