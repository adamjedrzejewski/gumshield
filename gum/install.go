package gum

import (
	"errors"
	"os"
	"path/filepath"
)

func Install(archivePath, targetDir string, verbose, disableIndex bool) error {
	err := isElevated()
	if err != nil {
		return err
	}

	absArchivePath, err := filepath.Abs(archivePath)
	if err != nil {
		return err
	}
	absBuildDir, err := filepath.Abs(DefaultBuildDir)
	if err != nil {
		return err
	}
	absFakeRootDir, err := filepath.Abs(DefaultFakeRootDir)
	if err != nil {
		return err
	}
	absTempDir, err := filepath.Abs(DefaultTempDir)
	if err != nil {
		return err
	}
	absIndexDir, err := filepath.Abs(DefaultIndexDir)
	if err != nil {
		return err
	}

	if _, err := os.Stat(absIndexDir); os.IsNotExist(err) {
		os.MkdirAll(absIndexDir, 0755)
	}
	if err := SetEnvVars(absBuildDir, absFakeRootDir); err != nil {
		return err
	}
	if err := prepareDirs(absBuildDir, absFakeRootDir, absTempDir); err != nil {
		return err
	}
	if err := extractPackageArchive(absArchivePath, absTempDir); err != nil {
		return err
	}
	pkg, err := ReadDefinitionFromFile(filepath.Join(absTempDir, DefinitionFileName))
	if err := isInstalled(pkg.Name); err != nil {
		return err
	}
	if err := ValidateInstalledDefinition(pkg); err != nil {
		return err
	}
	if !disableIndex {
		if err := copyDefinitionToIndex(pkg.Name, absTempDir, absIndexDir); err != nil {
			return err
		}
	}
	if pkg.BeforeInstallLogic != "" {
		if err := runScriptInDir(DefaultTempDir, pkg.BeforeInstallLogic, verbose); err != nil {
			return err
		}
	}
	if err := extractFilesToDir(absTempDir, targetDir); err != nil {
		return err
	}
	if pkg.AfterInstallLogic != "" {
		if err := runScriptInDir(DefaultTempDir, pkg.AfterInstallLogic, verbose); err != nil {
			return err
		}
	}
	if err := cleanUpDirs(absBuildDir, absFakeRootDir, absTempDir); err != nil {
		return err
	}

	return nil
}

func extractFilesToDir(fromDir, toDir string) error {
	filesArchive := filepath.Join(fromDir, FilesArchiveFileName)
	archive, err := os.Open(filesArchive)
	if err != nil {
		return err
	}
	return extractTar(toDir, archive)
}

func copyDefinitionToIndex(name, sourceDir, destinationDir string) error {
	fileInfo, err := os.Stat(destinationDir)
	if errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(destinationDir, 0755)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	if !fileInfo.IsDir() {
		return errors.New(destinationDir + " is not a directory")
	}

	sourceFile := filepath.Join(sourceDir, DefinitionFileName)
	input, err := os.ReadFile(sourceFile)
	if err != nil {
		return err
	}

	destinationFile := filepath.Join(destinationDir, name+DefinitionFileExtension)
	err = os.WriteFile(destinationFile, input, 0644)
	if err != nil {
		return err
	}

	return nil
}
