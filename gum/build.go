package gum

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func Build(pkg *PackageDefinition, outputFile, buildDir, fakeRootDir, tempDir string, verbose bool, sourcesDir *string) error {
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
	if err := getSourcesFromLocalDir(sourcesDir, absBuildDir, pkg.Sources); err != nil {
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

func getSourcesFromLocalDir(dir *string, outDir string, sources []string) error {
	if dir == nil {
		return nil
	}
	absSourcesDir, err := filepath.Abs(*dir)
	if err != nil {
		return err
	}

	for _, source := range sources {
		_, fileName := filepath.Split(source)
		sourceFile := filepath.Join(absSourcesDir, fileName)

		if _, err := os.Stat(sourceFile); err != nil {
			continue
		}

		destinationFile := filepath.Join(outDir, fileName)
		if err := copyFile(sourceFile, destinationFile); err != nil {
			return err
		}
	}

	return nil
}

// TODO: refactor
func copyFile(sourceFile string, destinationFile string) error {
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(destinationFile, input, 0644)
	if err != nil {
		return err
	}

	return nil
}
