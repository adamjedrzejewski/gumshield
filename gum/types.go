package gum

type PackageDefinition struct {
	Name       string
	Version    string
	Sources    []string
	BuildLogic string
}

type PackageManifest struct {
	Package PackageDefinition
	Files   []string
}
