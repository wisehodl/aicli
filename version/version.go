package version

// These variables are set by the build process using ldflags
var (
	// Version is the semantic version of the release
	Version = "dev"
	// CommitHash is the git commit hash of the build
	CommitHash = "unknown"
	// BuildDate is the date when the binary was built
	BuildDate = "unknown"
)

// GetVersion returns a formatted version string
func GetVersion() string {
	return Version + " (" + CommitHash + ") built at " + BuildDate
}
