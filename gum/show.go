package gum

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
)

func ShowInstalled() error {
	packages, err := readPackages()
	if err != nil {
		return err
	}

	for _, v := range packages {
		fmt.Println(v.Name)
	}

	return nil
}

func ShowConfig() {

}

func ShowFiles(packageName string) error {
	pkg, err := getPackage(packageName)
	if err != nil {
		return err
	}

	for _, file := range pkg.Files {
		fmt.Println(file)
	}

	return nil
}

func ShowPackage(packageName string) error {
	pkg, err := getPackage(packageName)
	if err != nil {
		return err
	}

	fmt.Println("name:", pkg.Name)
	fmt.Println("version:", pkg.Version)
	fmt.Println("description:", pkg.Description)
	fmt.Println("files:")
	for _, file := range pkg.Files {
		fmt.Println(file)
	}

	return nil
}

func ShowTriggers(packageName string) error {
	pkg, err := getPackage(packageName)
	if err != nil {
		return err
	}

	fmt.Println("build:")
	fmt.Println(pkg.BuildLogic)
	fmt.Println("before install:")
	fmt.Println(pkg.BeforeInstallLogic)
	fmt.Println("after install:")
	fmt.Println(pkg.AfterInstallLogic)
	fmt.Println("uninstall:")
	fmt.Println(pkg.UninstallLogic)

	return nil
}

func readPackages() ([]*PackageDefinition, error) {
	files, err := ioutil.ReadDir(DefaultIndexDir)
	if err != nil {
		return nil, err
	}

	packages := make([]*PackageDefinition, 0)
	for _, file := range files {
		filePath := path.Join(DefaultIndexDir, file.Name())
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		pkg, err := ParsePackageDefinition(string(content))
		if err != nil {
			return nil, err
		}

		packages = append(packages, pkg)
	}

	return packages, nil
}

func getPackage(packageName string) (*PackageDefinition, error) {
	packages, err := readPackages()
	if err != nil {
		return nil, err
	}
	for _, pkg := range packages {
		if pkg.Name == packageName {
			return pkg, nil
		}
	}

	return nil, errors.New("no such package")
}
