package project

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Project struct {
	// The name of the project type.
	Name string

	Configurations []ProjectConfiguration
}

type ProjectConfiguration struct {
	Identifiers       []Identifier
	TargetDirectories []string
}

type IdentifierType string

const (
	IdentifierTypeFileName               IdentifierType = "file"
	IdentifierTypeFileExtension          IdentifierType = "extension"
	IdentifierTypeFileNameInDir          IdentifierType = "fileInDir"
	IdentifierTypeExtensionInSpecificDir IdentifierType = "extensionInDir"
)

// Identifies the project root
type Identifier struct {
	identifierType      IdentifierType
	name                string
	extension           string
	parentDirectory     string
	relationToTargetDir string
}

func (i *Identifier) IsIdentifier(identifierPath string) bool {
	identifierStat, identifierStatErr := os.Stat(identifierPath)

	if identifierStatErr != nil {
		return false
	}

	fileBase := filepath.Base(identifierPath)
	fileExt := filepath.Ext(identifierPath)
	fileName := strings.TrimSuffix(fileBase, fileExt)
	parentDir := filepath.Base(filepath.Dir(identifierPath))

	nameIsEqual := fileName == i.name
	extensionIsEqual := fileExt == i.extension
	parentDirIsEqual := parentDir == i.parentDirectory

	isFile := !identifierStat.IsDir()

	switch i.identifierType {
	case IdentifierTypeFileName:
		return isFile && nameIsEqual && extensionIsEqual
	case IdentifierTypeFileExtension:
		return isFile && extensionIsEqual
	case IdentifierTypeFileNameInDir:
		return isFile && nameIsEqual && extensionIsEqual && parentDirIsEqual
	case IdentifierTypeExtensionInSpecificDir:
		return isFile && extensionIsEqual && parentDirIsEqual
	}

	return false
}

// TODO: project can have multiple project configurations
// These have relative paths to the target folders
