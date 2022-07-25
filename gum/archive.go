package gum

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

func createPackageArchive(fromDir, tempDir, outFile string, pkg *PackageDefinition) error {
	files, err := listFiles(fromDir)
	if err != nil {
		return err
	}
	filesArchivePath := filepath.Join(tempDir, FilesArchiveFileName)

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	err = os.Chdir(fromDir)
	if err != nil {
		return err
	}
	if err := createTarball(filesArchivePath, files); err != nil {
		return err
	}
	err = os.Chdir(currentDir)
	if err != nil {
		return err
	}

	pkg.Files = files
	definitionPath := filepath.Join(tempDir, DefinitionFileName)
	if err := writeDefinition(definitionPath, pkg); err != nil {
		return err
	}

	outFileFiles := []string{
		FilesArchiveFileName,
		DefinitionFileName,
	}

	_ = outFileFiles

	currentDir, err = os.Getwd()
	if err != nil {
		return err
	}
	err = os.Chdir(tempDir)
	if err != nil {
		return err
	}
	if err := createTarball(outFile, outFileFiles); err != nil {
		return err
	}
	err = os.Chdir(currentDir)
	if err != nil {
		return err
	}

	return nil
}

func createTarball(outFile string, files []string) error {
	file, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer file.Close()

	tarWriter := tar.NewWriter(file)
	defer tarWriter.Close()

	for _, filePath := range files {
		err := addFileToTarWriter(filePath, tarWriter)
		if err != nil {
			return err
		}
	}

	return nil
}

func addFileToTarWriter(path string, writer *tar.Writer) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(stat, path)
	if err != nil {
		return err
	}

	header.Name = path
	err = writer.WriteHeader(header)
	if err != nil {
		return err
	}
	if stat.IsDir() {
		return nil
	}

	_, err = io.Copy(writer, file)
	if err != nil {
		return err
	}

	return nil
}

func extractPackageArchive(archivePath, outputDir string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}

	if err := extractTar(outputDir, file); err != nil {
		return err
	}

	return nil
}

func extractTar(dst string, reader io.Reader) error {
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

func writeDefinition(path string, pkg *PackageDefinition) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	content, err := SerializePackageDefinition(pkg)
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(content))
	if err != nil {
		return err
	}

	return nil
}
