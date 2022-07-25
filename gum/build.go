package gum

import (
	"archive/tar"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func Build(pkg *PackageDefinition, outputFile, buildDir, fakeRootDir, tempDir string, verbose bool) error {
	absBuildDir, err := toAbsolutePath(buildDir)
	if err != nil {
		return err
	}
	absFakeRootDir, err := toAbsolutePath(fakeRootDir)
	if err != nil {
		return err
	}
	absTempDir, err := toAbsolutePath(tempDir)
	if err != nil {
		return err
	}
	absOutputFile, err := toAbsolutePath(outputFile)
	if err != nil {
		return err
	}

	if err := setEnvVars(absBuildDir, absFakeRootDir); err != nil {
		return err
	}
	if err := prepareDirs(absBuildDir, absFakeRootDir, absTempDir); err != nil {
		return err
	}
	if err := downloadSources(pkg.Sources, absBuildDir); err != nil {
		return err
	}
	if err := runBuild(absBuildDir, pkg.BuildLogic, verbose); err != nil {
		return err
	}
	if err := createArchive(absFakeRootDir, absTempDir, absOutputFile, pkg); err != nil {
		return err
	}
	if err := cleanUpDirs(absBuildDir, absFakeRootDir, absTempDir); err != nil {
		return err
	}

	return nil
}

func runBuild(dir, logic string, verbose bool) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	err = os.Chdir(dir)
	if err != nil {
		return err
	}

	cmd := exec.Command("bash")
	if verbose {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	if _, err := stdin.Write([]byte(logic)); err != nil {
		return err
	}
	if err := stdin.Close(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}

	err = os.Chdir(currentDir)
	if err != nil {
		return err
	}

	return nil
}

func prepareDirs(buildDir, fakeRootDir, tempDir string) error {
	if err := cleanUpDirs(buildDir, fakeRootDir, tempDir); err != nil {
		return err
	}
	if err := os.MkdirAll(buildDir, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(fakeRootDir, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func cleanUpDirs(buildDir, fakeRootDir, tempDir string) error {
	if err := os.RemoveAll(buildDir); err != nil {
		return err
	}
	if err := os.RemoveAll(fakeRootDir); err != nil {
		return err
	}
	if err := os.RemoveAll(tempDir); err != nil {
		return err
	}

	return nil
}

func toAbsolutePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	return filepath.Abs(path)
}

func downloadSources(sources []string, dir string) error {
	for _, source := range sources {
		_, outPath := filepath.Split(source)
		outPath = filepath.Join(dir, outPath)
		if err := downloadFile(source, outPath); err != nil {
			return err
		}
	}

	return nil
}

func downloadFile(url, outPath string) error {
	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func setEnvVars(buildDir string, fakeRootDir string) error {
	if err := os.Setenv(BuildDirEnvVarName, buildDir); err != nil {
		return err
	}
	if err := os.Setenv(FakeRootDirEnvVarName, fakeRootDir); err != nil {
		return err
	}

	return nil
}

func createArchive(fromDir, tempDir, outFile string, pkg *PackageDefinition) error {
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

func listFiles(dir string) ([]string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	err = os.Chdir(dir)
	if err != nil {
		return nil, err
	}

	files := make([]string, 0)
	err = filepath.Walk(".", func(path string, info fs.FileInfo, err error) error { // TODO: error?
		if path == "." {
			return nil
		}
		//if f, err := os.Stat(path); err != nil || f.IsDir() {
		//	return nil
		//}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	err = os.Chdir(currentDir)
	if err != nil {
		return nil, err
	}

	return files, nil
}

//  TODO: fix tar
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
