package main

import (
	"fmt"
	"github.com/adamjedrzejewski/gumshield/gum"
	"log"
	"os"
	"path/filepath"
	// https://github.com/hellflame/argparse
)

func main() {
	definitionPath := "./test.elplan"
	archivePath := "./test.tar"
	absDefinitionPath, err := filepath.Abs(definitionPath)
	if err != nil {
		log.Fatalln(err)
	}

	pkg, err := gum.ReadDefinitionFromFile(absDefinitionPath)
	if err != nil {
		log.Fatalln(err)
	}

	if len(os.Args) < 2 {
		fmt.Println("action?")
		return
	}

	if os.Args[1] == "build" {
		buildDir := gum.DefaultBuildDir
		fakeRootDir := gum.DefaultFakeRootDir
		tempDir := gum.DefaultTempDir
		outFilePath := pkg.Name + ".tar"
		verbose := true
		if err := gum.Build(pkg, outFilePath, buildDir, fakeRootDir, tempDir, verbose); err != nil {
			log.Fatal(err)
		}
		return
	}

	if os.Args[1] == "install" {
		if err := gum.Install(archivePath, true); err != nil {
			log.Fatalln(err)
		}
		return
	}

	if len(os.Args) < 3 {
		fmt.Println("action?")
		return
	}
	if os.Args[1] == "show" && os.Args[2] == "installed" {
		if err := gum.ShowInstalled(); err != nil {
			log.Fatalln(err)
		}
		return
	}
	if os.Args[1] == "show" && os.Args[2] == "files" {
		if err := gum.ShowFiles("test"); err != nil {
			log.Fatalln(err)
		}
		return
	}
	if os.Args[1] == "show" && os.Args[2] == "package" {
		if err := gum.ShowPackage("test"); err != nil {
			log.Fatalln(err)
		}
		return
	}
	if os.Args[1] == "show" && os.Args[2] == "triggers" {
		if err := gum.ShowTriggers("test"); err != nil {
			log.Fatalln(err)
		}
		return
	}

}

/*
   COMMANDS:
   	- build <definition file> - build package from definition file

   	- install <archive file> - install package from archive file

   	create definition <definition name> - create package definition

   	show config - show gumshield configuration
   	- show installed - list installed packages
   	- show package <package name> - show package information
   	- show files <package name> - list package files
	- show triggers <package name> - show package scripts

   	uninstall <package name> - remove package

   	remote install <package name> - install package from remote repository

   	get definition <package name> - get package definition file from remote repository
   	get sources <package name> - get package sources from remote repository
   	get triggers <package name> - get package scripts from remote repository
   	get archive <package name> - get archive file from remote repository
*/
