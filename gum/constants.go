package gum

const (
	DefaultBuildDir    = "/tmp/gumshield/build"
	DefaultFakeRootDir = "/tmp/gumshield/fake_root"
	DefaultTempDir     = "/tmp/gumshield/temp"
	DefaultIndexDir    = "/var/lib/gumshield"
	RootDir            = "/"

	// DefaultConfigFile = "/etc/gumshield" // TODO: config

	DefinitionFileName   = "manifest"
	FilesArchiveFileName = "files.tar"

	DefinitionFileExtension = ".elplan"

	BuildDirEnvVarName    = "GUMSHIELD_BUILD_DIR"
	FakeRootDirEnvVarName = "GUMSHIELD_FAKE_ROOT_DIR"
)
