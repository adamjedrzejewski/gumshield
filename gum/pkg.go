package gum

import (
	"bufio"
	"errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const (
	noSection = ""

	descriptionSectionTag   = "%%% DESCRIPTION"
	metaSectionTag          = "%%% META"
	beforeInstallSectionTag = "%%% BEFORE INSTALL"
	afterInstallSectionTag  = "%%% AFTER INSTALL"
	uninstallSectionTag     = "%%% UNINSTALL"
	buildSectionTag         = "%%% BUILD"
	filesSectionTag         = "%%% FILES"
	tagLikeTerminator       = "%%%"
)

func NewPackageDefinition(name, version string, sources []string, description, buildLogic, beforeInstallLogic, afterInstallLogic, uninstallLogic string, files []string) *PackageDefinition {
	return &PackageDefinition{
		Name:               name,
		Version:            version,
		Sources:            sources,
		Files:              files,
		BuildLogic:         buildLogic,
		BeforeInstallLogic: beforeInstallLogic,
		AfterInstallLogic:  afterInstallLogic,
		UninstallLogic:     uninstallLogic,
		Description:        description,
	}
}

func ReadDefinitionFromFile(path string) (*PackageDefinition, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ParsePackageDefinition(string(content))
}

func ParsePackageDefinition(content string) (*PackageDefinition, error) {
	sections := map[string][]string{
		descriptionSectionTag:   {},
		metaSectionTag:          {},
		beforeInstallSectionTag: {},
		afterInstallSectionTag:  {},
		uninstallSectionTag:     {},
		buildSectionTag:         {},
		filesSectionTag:         {},
	}

	currentSection := noSection
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if setSectionTag(&currentSection, line) {
			continue
		}
		if currentSection == noSection {
			continue
		}

		addToSection(&sections, currentSection, line)
	}

	description := strings.Join(sections[descriptionSectionTag], "\n")
	buildLogic := strings.Join(sections[buildSectionTag], "\n")
	beforeInstallLogic := strings.Join(sections[beforeInstallSectionTag], "\n")
	afterInstallLogic := strings.Join(sections[afterInstallSectionTag], "\n")
	uninstallLogic := strings.Join(sections[uninstallSectionTag], "\n")
	files := sections[filesSectionTag]
	metadata, err := getMetadata(strings.Join(sections[metaSectionTag], "\n"))
	if err != nil {
		return nil, err
	}

	return NewPackageDefinition(
		metadata.Name,
		metadata.Version,
		metadata.Sources,
		description,
		buildLogic,
		beforeInstallLogic,
		afterInstallLogic,
		uninstallLogic,
		files,
	), nil
}

func SerializePackageDefinition(pkg *PackageDefinition) (string, error) {
	sb := strings.Builder{}
	meta := PackageMetadata{
		Name:    pkg.Name,
		Version: pkg.Version,
		Sources: pkg.Sources,
	}

	if pkg.Description != "" {
		sb.Write([]byte(descriptionSectionTag))
		sb.Write([]byte("\n"))
		sb.Write([]byte(pkg.Description))
		sb.Write([]byte("\n"))
	}

	sb.Write([]byte(metaSectionTag))
	sb.Write([]byte("\n"))
	metaYaml, err := yaml.Marshal(meta)
	if err != nil {
		return "", err
	}
	sb.Write(metaYaml)
	sb.Write([]byte("\n"))

	sb.Write([]byte(buildSectionTag))
	sb.Write([]byte("\n"))
	sb.Write([]byte(pkg.BuildLogic))
	sb.Write([]byte("\n"))

	sb.Write([]byte(beforeInstallSectionTag))
	sb.Write([]byte("\n"))
	sb.Write([]byte(pkg.BeforeInstallLogic))
	sb.Write([]byte("\n"))

	sb.Write([]byte(afterInstallSectionTag))
	sb.Write([]byte("\n"))
	sb.Write([]byte(pkg.AfterInstallLogic))
	sb.Write([]byte("\n"))

	sb.Write([]byte(uninstallSectionTag))
	sb.Write([]byte("\n"))
	sb.Write([]byte(pkg.UninstallLogic))
	sb.Write([]byte("\n"))

	if pkg.Files != nil && len(pkg.Files) > 0 {
		sb.Write([]byte(filesSectionTag))
		sb.Write([]byte("\n"))
		sb.Write([]byte(strings.Join(pkg.Files, "\n")))
		sb.Write([]byte("\n"))
	}

	return sb.String(), nil
}

func ValidateInstalledDefinition(pkg *PackageDefinition) error {
	if pkg.BeforeInstallLogic == "" {
		return errors.New("missing before install script")
	}
	if pkg.AfterInstallLogic == "" {
		return errors.New("missing after install script")
	}
	if pkg.UninstallLogic == "" {
		return errors.New("missing uninstall script")
	}
	if pkg.Files == nil || len(pkg.Files) == 0 {
		return errors.New("missing file list")
	}

	return nil
}

func readPackagesFromIndex() ([]*PackageDefinition, error) {
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

func getPackageFromIndex(packageName string) (*PackageDefinition, error) {
	packages, err := readPackagesFromIndex()
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

func isInstalled(packageName string) error {
	packages, err := readPackagesFromIndex()
	if err != nil {
		return err
	}

	for _, pkg := range packages {
		if pkg.Name == packageName {
			return errors.New("package already installed")
		}
	}

	return nil
}

// TODO: move adding to package index to separete function

func removePackageFromIndex(pkgName string) error {
	filePath := path.Join(DefaultIndexDir, pkgName+DefinitionFileExtension)
	_, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	err = os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}

func setSectionTag(currentSection *string, line string) bool {
	if strings.HasPrefix(line, descriptionSectionTag) {
		*currentSection = descriptionSectionTag
		return true
	}
	if strings.HasPrefix(line, metaSectionTag) {
		*currentSection = metaSectionTag
		return true
	}
	if strings.HasPrefix(line, beforeInstallSectionTag) {
		*currentSection = beforeInstallSectionTag
		return true
	}
	if strings.HasPrefix(line, afterInstallSectionTag) {
		*currentSection = afterInstallSectionTag
		return true
	}
	if strings.HasPrefix(line, uninstallSectionTag) {
		*currentSection = uninstallSectionTag
		return true
	}
	if strings.HasPrefix(line, buildSectionTag) {
		*currentSection = buildSectionTag
		return true
	}
	if strings.HasPrefix(line, filesSectionTag) {
		*currentSection = filesSectionTag
		return true
	}
	if strings.HasPrefix(line, tagLikeTerminator) {
		*currentSection = noSection
		return true
	}

	return false
}

func addToSection(sections *map[string][]string, currentSection, line string) {
	section, ok := (*sections)[currentSection]
	if !ok {
		return
	}

	(*sections)[currentSection] = append(section, line)
}

func getMetadata(yamlContent string) (PackageMetadata, error) {
	metadata := PackageMetadata{}
	err := yaml.Unmarshal([]byte(yamlContent), &metadata)

	return metadata, err
}
