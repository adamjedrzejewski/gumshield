package gum

import "path/filepath"

func Install(archivePath string) error {
	absArchivePath, err := filepath.Abs(archivePath)
	if err != nil {
		return err
	}

	if err := prepareDirs(DefaultBuildDir, DefaultFakeRootDir, DefaultTempDir); err != nil {
		return err
	}
	if err := extractArchive(absArchivePath, DefaultTempDir); err != nil {
		return err
	}
	pkg, err := ReadDefinitionFromFile(filepath.Join(DefaultTempDir, ManifestFileName))
	if err := copyDefinition(pkg.Name, DefaultTempDir, DefaultIndexDir); err != nil {
		return err
	}
	if err := copyFiles(DefaultTempDir); err != nil {
		return err
	}
	if err := cleanUpDirs(DefaultBuildDir, DefaultFakeRootDir, DefaultTempDir); err != nil {
		return err
	}

	return nil
}

func copyFiles(fromDir string) error {
	return nil
}

func copyDefinition(name, sourceDir, destinationDir string) error {
	return nil
}

func extractArchive(archivePath, outputDir string) error {
	return nil
}
