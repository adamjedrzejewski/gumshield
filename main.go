package main

import (
	"github.com/hellflame/argparse"
)

func main() {
	parser := argparse.NewParser("gumshield", "gumshield package manager", nil)

	registerBuildCommand(parser)
	registerInstallCommand(parser)
	registerShowCommand(parser)

	_ = parser.Parse(nil)

	//definitionPath := "./test.elplan"
	//archivePath := "./test.tar"
	//absDefinitionPath, err := filepath.Abs(definitionPath)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//pkg, err := gum.ReadDefinitionFromFile(absDefinitionPath)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//if len(os.Args) < 2 {
	//	fmt.Println("action?")
	//	return
	//}
	//
	//if os.Args[1] == "build" {
	//	buildDir := gum.DefaultBuildDir
	//	fakeRootDir := gum.DefaultFakeRootDir
	//	tempDir := gum.DefaultTempDir
	//	outFilePath := pkg.Name + ".tar"
	//	verbose := true
	//	if err := gum.Build(pkg, outFilePath, buildDir, fakeRootDir, tempDir, verbose); err != nil {
	//		log.Fatal(err)
	//	}
	//	return
	//}
	//
	//if os.Args[1] == "install" {
	//	if err := gum.Install(archivePath, true); err != nil {
	//		log.Fatalln(err)
	//	}
	//	return
	//}
	//
	//if len(os.Args) < 3 {
	//	fmt.Println("action?")
	//	return
	//}
	//if os.Args[1] == "show" && os.Args[2] == "installed" {
	//	if err := gum.ShowInstalled(); err != nil {
	//		log.Fatalln(err)
	//	}
	//	return
	//}
	//if os.Args[1] == "show" && os.Args[2] == "files" {
	//	if err := gum.ShowFiles("test"); err != nil {
	//		log.Fatalln(err)
	//	}
	//	return
	//}
	//if os.Args[1] == "show" && os.Args[2] == "package" {
	//	if err := gum.ShowPackage("test"); err != nil {
	//		log.Fatalln(err)
	//	}
	//	return
	//}
	//if os.Args[1] == "show" && os.Args[2] == "triggers" {
	//	if err := gum.ShowTriggers("test"); err != nil {
	//		log.Fatalln(err)
	//	}
	//	return
	//}

}
