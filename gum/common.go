package gum

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

const (
	scriptCommand        = "bash"
	rootUidString        = "0"
	currentDirPathString = "."
)

// SetEnvVars sets environment variables.
func SetEnvVars(buildDir string, fakeRootDir string) error {
	if err := os.Setenv(BuildDirEnvVarName, buildDir); err != nil {
		return err
	}
	if err := os.Setenv(FakeRootDirEnvVarName, fakeRootDir); err != nil {
		return err
	}

	return nil
}

// runScriptInDir executes bash script in separate process.
func runScriptInDir(dir, logic string, verbose bool) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	err = os.Chdir(dir)
	if err != nil {
		return err
	}

	cmd := exec.Command(scriptCommand)
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

// cleanUpDirs removes directory trees at specified paths.
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

// prepareDirs cleans up and creates directories.
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

// listFiles returns listing of all files in directory tree.
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
	err = filepath.Walk(currentDirPathString, func(path string, info fs.FileInfo, err error) error { // TODO: error?
		if path == currentDirPathString {
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

// isElevated checks if process is running as superuser.
func isElevated() error {
	u, err := user.Current()
	if err != nil {
		return err
	}

	if u.Uid != rootUidString {
		return fmt.Errorf("process not running with elevated permission")
	}

	return nil
}

// TODO: this function should be able to copy big files as well
// copyFile is for copying small files.
func copyFile(sourcePath string, destinationPath string) error {
	input, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}
	err = os.WriteFile(destinationPath, input, 0644)
	if err != nil {
		return err
	}

	return nil
}
