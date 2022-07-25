package gum

const (
	DefaultBuildDir    = "/tmp/gumshield/build"
	DefaultFakeRootDir = "/tmp/gumshield/fake_root"
	DefaultTempDir     = "/tmp/gumshield/temp"
	// DefaultIndexDir    = "/var/lib/gumshield" // TODO: index

	// DefaultConfigFile = "/etc/gumshield" // TODO: config

	ManifestFileName     = "manifest"
	FilesArchiveFileName = "files.tar"

	BuildDirEnvVarName    = "GUMSHIELD_BUILD_DIR"
	FakeRootDirEnvVarName = "GUMSHIELD_FAKE_ROOT_DIR"
)
