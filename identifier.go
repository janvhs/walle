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

func (i *FileExtensionInDirectoryIdentifier) MatchesOptimistically(potentialRoot string) bool {
	potentialPath := filepath.Join(potentialRoot, i.Directory)
	potentialPath = filepath.Clean(potentialPath)
	return i.Matches(potentialPath)
}

type FileNameInDirectoryIdentifier struct {
	Directory string
	Name      string
}

func (i *FileNameInDirectoryIdentifier) Matches(potentialPath string) bool {
	dirName := filepath.Dir(potentialPath)
	fileName := filepath.Base(potentialPath)

	stat, err := os.Stat(potentialPath)
	if err != nil {
		return false
	}

	return !stat.IsDir() && dirName == i.Directory && fileName == i.Name
}

func (i *FileNameInDirectoryIdentifier) MatchesOptimistically(potentialRoot string) bool {
	potentialPath := filepath.Join(potentialRoot, i.Directory, i.Name)
	potentialPath = filepath.Clean(potentialPath)

	return i.Matches(potentialPath)
}

type FileNameIdentifier struct {
	Name string
}

func (i *FileNameIdentifier) Matches(potentialPath string) bool {
	fileName := filepath.Base(potentialPath)
	stat, err := os.Stat(potentialPath)
	if err != nil {
		return false
	}

	return !stat.IsDir() && i.Name == fileName
}

func (i *FileNameIdentifier) MatchesOptimistically(potentialRoot string) bool {
	potentialPath := filepath.Join(potentialRoot, i.Name)
	return i.Matches(potentialPath)
}
