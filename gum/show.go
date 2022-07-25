package gum

import (
	"fmt"
)

func ShowInstalled() error {
	packages, err := readPackagesFromIndex()
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
	pkg, err := getPackageFromIndex(packageName)
	if err != nil {
		return err
	}

	for _, file := range pkg.Files {
		fmt.Println(file)
	}

	return nil
}

func ShowPackage(packageName string) error {
	pkg, err := getPackageFromIndex(packageName)
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
	pkg, err := getPackageFromIndex(packageName)
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
