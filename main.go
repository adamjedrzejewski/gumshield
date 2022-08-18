package main

import "github.com/hellflame/argparse"

func main() {
	parser := argparse.NewParser("gumshield", "gumshield package manager", nil)

	registerBuildCommand(parser)
	registerInstallCommand(parser)
	registerShowCommand(parser)
	registerUninstallCommand(parser)

	_ = parser.Parse(nil)
}
