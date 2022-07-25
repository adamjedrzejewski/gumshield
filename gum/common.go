package gum

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

func SetEnvVars(buildDir string, fakeRootDir string) error {
	if err := os.Setenv(BuildDirEnvVarName, buildDir); err != nil {
		return err
	}
	if err := os.Setenv(FakeRootDirEnvVarName, fakeRootDir); err != nil {
		return err
	}

	return nil
}

func runScriptInDir(dir, logic string, verbose bool) error {
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
