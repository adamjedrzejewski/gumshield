package main

import (
	"github.com/adamjedrzejewski/gumshield/gum"
	"github.com/hellflame/argparse"
	"log"
	"path/filepath"
)

func registerBuildCommand(parser *argparse.Parser) {
	build := parser.AddCommand("build", "build package from definition file", &argparse.ParserConfig{})
	pkgFile := build.String("", "definition_file", &argparse.Option{Positional: true, Help: "path to package definition file"})
	outFile := build.String("o", "out", &argparse.Option{Help: "path to output package archive"})
	buildDir := build.String("b", "build_dir", &argparse.Option{Help: "path to build directory"})
	fakeRootDir := build.String("f", "fake_root_dir", &argparse.Option{Help: "path to fake root directory"})
	tempDir := build.String("t", "temp_dir", &argparse.Option{Help: "path to temp directory"})
	verbose := build.Flag("v", "verbose", &argparse.Option{Help: "print output from underlying processes"})

	build.InvokeAction = func(bool) {
		absPkgFile, err := filepath.Abs(*pkgFile)
		if err != nil {
			log.Fatal(err)
		}
		pkg, err := gum.ReadDefinitionFromFile(absPkgFile)
		if err != nil {
			log.Fatal(err)
		}

		absOutFile := getOutFile(*outFile, pkg.Name)
		absBuildDir := getBuildDir(*buildDir)
		absFakeRootDir := getFakeRootDir(*fakeRootDir)
		absTempDir := getTempDir(*tempDir)

		err = gum.Build(pkg, absOutFile, absBuildDir, absFakeRootDir, absTempDir, *verbose)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getOutFile(input, pkgName string) string {
	if input == "" {
		input = pkgName + gum.ArchiveFileExtension
	}

	absFile, err := filepath.Abs(input)
	if err != nil {
		log.Fatal(err)
	}
	return absFile
}

func getBuildDir(input string) string {
	if input == "" {
		input = gum.DefaultBuildDir
	}

	absFile, err := filepath.Abs(input)
	if err != nil {
		log.Fatal(err)
	}
	return absFile
}

func getFakeRootDir(input string) string {
	if input == "" {
		input = gum.DefaultFakeRootDir
	}

	absFile, err := filepath.Abs(input)
	if err != nil {
		log.Fatal(err)
	}
	return absFile
}

func getTempDir(input string) string {
	if input == "" {
		input = gum.DefaultTempDir
	}

	absFile, err := filepath.Abs(input)
	if err != nil {
		log.Fatal(err)
	}
	return absFile
}
