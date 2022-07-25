package main

import (
	"github.com/adamjedrzejewski/gumshield/gum"
	"log"
	// https://github.com/hellflame/argparse
)

var samplePackage = gum.PackageDefinition{
	Name:    "sysvinit",
	Version: "2.98",
	Sources: []string{
		"https://download.savannah.gnu.org/releases/sysvinit/sysvinit-2.98.tar.xz",
		"https://www.linuxfromscratch.org/patches/lfs/10.1/sysvinit-2.98-consolidated-1.patch",
	},
	BuildLogic: `set -euo pipefail

unpack_src() {
    tar xf sysvinit-2.98.tar.xz && \
    cd sysvinit-2.98
    return
}

configure() {
    patch -Np1 -i ../sysvinit-2.98-consolidated-1.patch

    return
}

make_install() {
    make
    make ROOT=$GUMSHIELD_FAKE_ROOT_DIR -j1 install

    return
}

unpack_src && configure && make_install`,
}

func main() {
	if err := gum.Build(samplePackage, "sysvinit.tar", gum.DefaultBuildDir, gum.DefaultFakeRootDir, gum.DefaultTempDir, true); err != nil {
		log.Fatal(err)
	}
}

/*
   COMMANDS:
   	build <definition file> - build package from definition file

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
