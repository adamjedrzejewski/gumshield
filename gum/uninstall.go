package gum

import (
	"os"
	"path/filepath"
)

func Uninstall(packageName string, verbose bool) error {
	err := isElevated()
	if err != nil {
		return err
	}

	pkg, err := getPackageFromIndex(packageName)
	if err != nil {
		return err
	}

	if err := ValidateInstalledDefinition(pkg); err != nil {
		return err
	}
	if pkg.UninstallLogic != "" {
		if err := runScriptInDir(DefaultTempDir, pkg.UninstallLogic, verbose); err != nil {
			return err
		}
	}
	if err := removeRegularPackageFiles(pkg.Files); err != nil {
		return err
	}
	if err := removePackageDirectoriesIfEmpty(pkg.Files); err != nil {
		return err
	}
	if err := removePackageFromIndex(pkg.Name); err != nil {
		return err
	}
	return nil
}

func removePackageDirectoriesIfEmpty(dirs []string) error {
	for _, dir := range dirs {
		path := filepath.Join(RootDir, dir)
		info, err := os.Stat(path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		if !info.IsDir() {
			continue
		}

		items, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		if items == nil || len(items) == 0 {
			err := os.RemoveAll(path)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func removeRegularPackageFiles(files []string) error {
	for _, file := range files {
		path := filepath.Join(RootDir, file)
		info, err := os.Stat(path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		if info.IsDir() {
			continue
		}

		err = os.Remove(path)
		if err != nil {
			return err
		}
	}
	return nil
}
