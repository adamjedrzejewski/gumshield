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

func registerInstallCommand(parser *argparse.Parser) {
	install := parser.AddCommand("install", "install package from archive file", &argparse.ParserConfig{})
	pkgFile := install.String("", "archive_file", &argparse.Option{Positional: true, Help: "path to package archive file"})
	verbose := install.Flag("v", "verbose", &argparse.Option{Help: "print output from underlying processes"})

	install.InvokeAction = func(bool) {
		absPkgFile, err := filepath.Abs(*pkgFile)
		if err != nil {
			log.Fatal(err)
		}

		err = gum.Install(absPkgFile, *verbose)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func registerShowCommand(parser *argparse.Parser) {
	show := parser.AddCommand("show", "display information", &argparse.ParserConfig{})

	registerShowInstalledCommand(show)
	registerShowFilesCommand(show)
	registerShowPackageCommand(show)
	registerShowTriggersCommand(show)
	registerShowConfigCommand(show)
}

func registerShowConfigCommand(parser *argparse.Parser) {
	installed := parser.AddCommand("config", "show gumshield config", &argparse.ParserConfig{DisableDefaultShowHelp: true})

	installed.InvokeAction = func(bool) {
		err := gum.ShowConfig()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func registerShowTriggersCommand(parser *argparse.Parser) {
	pkg := parser.AddCommand("triggers", "show package triggers", &argparse.ParserConfig{})
	pkgName := pkg.String("", "package_name", &argparse.Option{Positional: true, Help: "package name"})

	pkg.InvokeAction = func(bool) {
		err := gum.ShowTriggers(*pkgName)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func registerShowPackageCommand(parser *argparse.Parser) {
	pkg := parser.AddCommand("package", "show package information", &argparse.ParserConfig{})
	pkgName := pkg.String("", "package_name", &argparse.Option{Positional: true, Help: "package name"})

	pkg.InvokeAction = func(bool) {
		err := gum.ShowPackage(*pkgName)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func registerShowFilesCommand(parser *argparse.Parser) {
	files := parser.AddCommand("files", "show package files", &argparse.ParserConfig{})
	pkgName := files.String("", "package_name", &argparse.Option{Positional: true, Help: "package name"})

	files.InvokeAction = func(bool) {
		err := gum.ShowFiles(*pkgName)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func registerShowInstalledCommand(parser *argparse.Parser) {
	installed := parser.AddCommand("installed", "show installed packages", &argparse.ParserConfig{DisableDefaultShowHelp: true})

	installed.InvokeAction = func(bool) {
		err := gum.ShowInstalled()
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
