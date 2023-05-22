package main

import (
	"os"
	"path/filepath"
	"strings"
)

type Identifier interface {
	Matches(potentialPath string) bool
	MatchesOptimistically(potentialRoot string) bool
}

type FileExtensionIdentifier struct {
	Extension string
	Directory string
}

func (i *FileExtensionIdentifier) Matches(potentialPath string) bool {
	potentialPath = filepath.Clean(potentialPath)
	stat, err := os.Stat(potentialPath)
	if err != nil {
		return false
	}

	if !stat.IsDir() {
		return false
	}

	direntList, err := os.ReadDir(potentialPath)
	if err != nil {
		return false
	}

	for _, dirent := range direntList {
		if filepath.Ext(dirent.Name()) == i.Extension {
			return true
		}
	}

	return false
}

func (i *FileExtensionIdentifier) MatchesOptimistically(potentialRoot string) bool {
	potentialPath := filepath.Join(potentialRoot, i.Directory)
	return i.Matches(potentialPath)
}

type FileNameIdentifier struct {
	Name      string
	Directory string
}

func (i *FileNameIdentifier) Matches(potentialPath string) bool {
	potentialPath = filepath.Clean(potentialPath)
	dirPath := filepath.Dir(potentialPath)
	fileName := filepath.Base(potentialPath)

	stat, err := os.Stat(potentialPath)
	if err != nil {
		return false
	}

	// Checking a suffix on a cleaned path can match dir1/dirN/file
	dirMatches := strings.HasSuffix(dirPath, filepath.Clean(i.Directory))

	// If no dir is provided, it always matches
	if i.Directory == "" {
		dirMatches = true
	}

	return !stat.IsDir() && dirMatches && i.Name == fileName
}

func (i *FileNameIdentifier) MatchesOptimistically(potentialRoot string) bool {
	potentialPath := filepath.Join(potentialRoot, i.Directory, i.Name)
	return i.Matches(potentialPath)
}
