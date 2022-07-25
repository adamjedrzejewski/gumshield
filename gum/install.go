package gum

import (
	"archive/tar"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Install(archivePath string) error {
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

	if err := prepareDirs(absBuildDir, absFakeRootDir, absTempDir); err != nil {
		return err
	}
	if err := extractArchive(absArchivePath, absTempDir); err != nil {
		return err
	}
	pkg, err := ReadDefinitionFromFile(filepath.Join(absTempDir, DefinitionFileName))
	if err := validateInstallDefinition(pkg); err != nil {
		return err
	}
	if err := copyDefinition(pkg.Name, absTempDir, absIndexDir); err != nil {
		return err
	}
	// TODO: run before install script
	if err := extractFiles(absTempDir); err != nil {
		return err
	}
	// TODO: run after install script
	if err := cleanUpDirs(absBuildDir, absFakeRootDir, absTempDir); err != nil {
		return err
	}

	return nil
}

func validateInstallDefinition(pkg *PackageDefinition) error {
	if pkg.BeforeInstallLogic == "" {
		return errors.New("missing before install script")
	}
	if pkg.AfterInstallLogic == "" {
		return errors.New("missing after install script")
	}
	if pkg.UninstallLogic == "" {
		return errors.New("missing uninstall script")
	}
	if pkg.Files == nil || len(pkg.Files) == 0 {
		return errors.New("missing file list")
	}

	return nil
}

func extractFiles(fromDir string) error {
	filesArchive := filepath.Join(fromDir, FilesArchiveFileName)
	archive, err := os.Open(filesArchive)
	if err != nil {
		return err
	}
	return unTar(RootDir, archive)
}

func copyDefinition(name, sourceDir, destinationDir string) error {
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
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return err
	}

	destinationFile := filepath.Join(destinationDir, name+DefinitionFileExtension)
	err = ioutil.WriteFile(destinationFile, input, 0644)
	if err != nil {
		return err
	}

	return nil
}

func extractArchive(archivePath, outputDir string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}

	if err := unTar(outputDir, file); err != nil {
		return err
	}

	return nil
}

func unTar(dst string, reader io.Reader) error {
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		}

		if header == nil {
			continue
		}

		targetPath := filepath.Join(dst, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(targetPath); err != nil {
				if err := os.MkdirAll(targetPath, 0755); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			f, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tarReader); err != nil {
				return err
			}

			if err := f.Close(); err != nil {
				return err
			}
		}
	}
}
