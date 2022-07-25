package gum

type PackageDefinition struct {
	Name           string
	Version        string
	Description    string
	BuildLogic     string
	InstallLogic   string
	UninstallLogic string
	Sources        []string
	Files          []string
}

type PackageMetadata struct {
	Name    string
	Version string
	Sources []string
	/*
		Sources []struct {
			Url      string
			Checksum string
		}
	*/
}
