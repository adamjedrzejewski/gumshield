package main

import (
	"github.com/adamjedrzejewski/gumshield/gum"
	"log"
	"path/filepath"
	// https://github.com/hellflame/argparse
)

func main() {
	definitionPath := "./sysvinit.elplan"
	absDefinitionPath, err := filepath.Abs(definitionPath)
	if err != nil {
		log.Fatalln(err)
	}

	pkg, err := gum.ReadDefinitionFromFile(absDefinitionPath)
	if err != nil {
		log.Fatalln(err)
	}

	buildDir := gum.DefaultBuildDir
	fakeRootDir := gum.DefaultFakeRootDir
	tempDir := gum.DefaultTempDir
	outFilePath := pkg.Name + ".tar"
	verbose := true

	if err := gum.Build(pkg, outFilePath, buildDir, fakeRootDir, tempDir, verbose); err != nil {
		log.Fatal(err)
	}
}

/*
   COMMANDS:
   	- build <definition file> - build package from definition file

   	install <archive file> - install package from archive file

   	create definition <definition name> - create package definition

   	show config - show gumshield configuration
   	show installed - list installed packages
   	show package <package name> - show package information
   	show files <package name> - list package files
	show triggers <package name> - show package scripts

   	uninstall <package name> - remove package

   	remote install <package name> - install package from remote repository

   	get definition <package name> - get package definition file from remote repository
   	get sources <package name> - get package sources from remote repository
   	get triggers <package name> - get package scripts from remote repository
   	get archive <package name> - get archive file from remote repository
*/
