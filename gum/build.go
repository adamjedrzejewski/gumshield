package gum

import (
	"path/filepath"
)

func Build(pkg *PackageDefinition, outputFile, buildDir, fakeRootDir, tempDir string, verbose bool) error {
	absBuildDir, err := filepath.Abs(buildDir)
	if err != nil {
		return err
	}
	absFakeRootDir, err := filepath.Abs(fakeRootDir)
	if err != nil {
		return err
	}
	absTempDir, err := filepath.Abs(tempDir)
	if err != nil {
		return err
	}
	absOutputFile, err := filepath.Abs(outputFile)
	if err != nil {
		return err
	}

	if err := SetEnvVars(absBuildDir, absFakeRootDir); err != nil {
		return err
	}
	if err := prepareDirs(absBuildDir, absFakeRootDir, absTempDir); err != nil {
		return err
	}
	if err := downloadSources(pkg.Sources, absBuildDir); err != nil {
		return err
	}
	if err := runScriptInDir(absBuildDir, pkg.BuildLogic, verbose); err != nil {
		return err
	}
	if err := createPackageArchive(absFakeRootDir, absTempDir, absOutputFile, pkg); err != nil {
		return err
	}
	if err := cleanUpDirs(absBuildDir, absFakeRootDir, absTempDir); err != nil {
		return err
	}

	return nil
}
