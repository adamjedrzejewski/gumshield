package main

import (
	"fmt"
	"github.com/adamjedrzejewski/gumshield/gum"
	"github.com/hellflame/argparse"
	"log"
	"os"
	"path/filepath"
)

func registerBuildCommand(parser *argparse.Parser) {
	build := parser.AddCommand("build", "build package from definition file", &argparse.ParserConfig{})
	pkgFile := build.String("", "definition_file", &argparse.Option{Positional: true, Help: "path to package definition file"})
	outFile := build.String("o", "out", &argparse.Option{Help: "path to output package archive"})
	buildDir := build.String("b", "build_dir", &argparse.Option{Help: "path to build directory", Default: gum.DefaultBuildDir})
	fakeRootDir := build.String("f", "fake_root_dir", &argparse.Option{Help: "path to fake root directory", Default: gum.DefaultFakeRootDir})
	tempDir := build.String("t", "temp_dir", &argparse.Option{Help: "path to temp directory", Default: gum.DefaultTempDir})
	sourcesDir := build.String("", "sources_dir", &argparse.Option{Help: "look for sources in this directory, if sound no sources will be downloaded", Default: gum.DefaultTempDir})
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
		absBuildDir, err := filepath.Abs(*buildDir)
		if err != nil {
			log.Fatal(err)
		}
		absFakeRootDir, err := filepath.Abs(*fakeRootDir)
		if err != nil {
			log.Fatal(err)
		}
		absTempDir, err := filepath.Abs(*tempDir)
		if err != nil {
			log.Fatal(err)
		}

		err = gum.Build(pkg, absOutFile, absBuildDir, absFakeRootDir, absTempDir, *verbose, sourcesDir)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func registerInstallCommand(parser *argparse.Parser) {
	install := parser.AddCommand("install", "install package from archive file", &argparse.ParserConfig{})
	pkgFile := install.String("", "archive_file", &argparse.Option{
		Positional: true,
		Help:       "path to package archive file",
		Validate:   validateFile})
	targetDir := install.String("", "target_dir", &argparse.Option{HideEntry: true, Default: gum.RootDir})
	disableIndex := install.Flag("", "disable_index", &argparse.Option{HideEntry: true})
	verbose := install.Flag("v", "verbose", &argparse.Option{Help: "print output from underlying processes"})

	install.InvokeAction = func(bool) {
		absPkgFile, err := filepath.Abs(*pkgFile)
		if err != nil {
			log.Fatal(err)
		}

		err = gum.Install(absPkgFile, *targetDir, *verbose, *disableIndex)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func registerUninstallCommand(parser *argparse.Parser) {
	uninstall := parser.AddCommand("uninstall", "uninstall package", &argparse.ParserConfig{})
	pkg := uninstall.String("", "package_name", &argparse.Option{Positional: true, Help: "package name"})
	verbose := uninstall.Flag("v", "verbose", &argparse.Option{Help: "print output from underlying processes"})

	uninstall.InvokeAction = func(bool) {
		err := gum.Uninstall(*pkg, *verbose)
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

func validateFile(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}
	if stat.IsDir() {
		return fmt.Errorf("%s is a direcotry", path)
	}
	return nil
}
