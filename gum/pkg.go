package gum

import (
	"bufio"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

const (
	noSection = ""

	descriptionSectionTag = "%%% DESCRIPTION"
	metaSectionTag        = "%%% META"
	installSectionTag     = "%%% INSTALL"
	uninstallSectionTag   = "%%% UNINSTALL"
	buildSectionTag       = "%%% BUILD"
	filesSectionTag       = "%%% FILES"
	tagLikeTerminator     = "%%%"
)

func NewPackageDefinition(name, version string, sources []string, description, buildLogic, installLogic, uninstallLogic string, files []string) *PackageDefinition {
	return &PackageDefinition{
		Name:           name,
		Version:        version,
		Sources:        sources,
		Files:          files,
		BuildLogic:     buildLogic,
		InstallLogic:   installLogic,
		UninstallLogic: uninstallLogic,
		Description:    description,
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
		descriptionSectionTag: {},
		metaSectionTag:        {},
		installSectionTag:     {},
		uninstallSectionTag:   {},
		buildSectionTag:       {},
		filesSectionTag:       {},
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
	installLogic := strings.Join(sections[installSectionTag], "\n")
	uninstallLogic := strings.Join(sections[uninstallSectionTag], "\n")
	files := sections[filesSectionTag]
	metadata, err := getMetadata(strings.Join(sections[metaSectionTag], "\n"))
	if err != nil {
		return nil, err
	}

	return NewPackageDefinition(metadata.Name, metadata.Version, metadata.Sources, description, buildLogic, installLogic, uninstallLogic, files), nil
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

	sb.Write([]byte(installSectionTag))
	sb.Write([]byte("\n"))
	sb.Write([]byte(pkg.InstallLogic))
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

func setSectionTag(currentSection *string, line string) bool {
	if strings.HasPrefix(line, descriptionSectionTag) {
		*currentSection = descriptionSectionTag
		return true
	}
	if strings.HasPrefix(line, metaSectionTag) {
		*currentSection = metaSectionTag
		return true
	}
	if strings.HasPrefix(line, installSectionTag) {
		*currentSection = installSectionTag
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
