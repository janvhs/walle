package main

import (
	"os"
	"path/filepath"
)

type Identifier interface {
	Matches(potentialPath string) bool
	MatchesOptimistically(basePath string) bool
}

type FileExtensionInDirectoryIdentifier struct {
	Directory string
	Extension string
}

func (i *FileExtensionInDirectoryIdentifier) Matches(potentialPath string) bool {
	dir := filepath.Dir(potentialPath)
	ext := filepath.Ext(potentialPath)

	return i.Directory == dir && i.Extension == ext
}

type FileNameIdentifier struct {
	Name string
}

func (i *FileNameIdentifier) Matches(potentialPath string) bool {
	fileName := filepath.Base(potentialPath)

	return i.Name == fileName
}

func (i *FileNameIdentifier) MatchesOptimistically(potentialRoot string) bool {
	potentialPath := filepath.Join(potentialRoot, i.Name)
	if _, err := os.Stat(potentialPath); err != nil {
		return false
	}

	return true
}
