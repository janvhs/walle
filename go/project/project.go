package project

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

type Project struct {
	// The name of the project type.
	Name string
	// The name of the file that indicates that the current directory is a project.
	// It identifies the project root.
	Identifiers []string
	// The name of the directories that should be deleted.
	// These have to be subdirectories of the project root.
	TargetDirectories []string
}

func (p *Project) pathIsIdentifier(identifierPath string) bool {
	for _, identifier := range p.Identifiers {
		if identifierPathStat, err := os.Stat(identifierPath); err == nil {
			if !identifierPathStat.IsDir() && filepath.Base(identifierPath) == identifier {
				return true
			}
		}
	}

	return false
}

func (p *Project) pathIsTargetDirectory(targetPath string) bool {
	for _, target := range p.TargetDirectories {
		if targetPathStat, err := os.Stat(targetPath); err == nil {
			if targetPathStat.IsDir() && filepath.Base(targetPath) == target {
				return true
			}
		}
	}

	return false
}

func (p *Project) Clean(identifierPath string, dryRun bool) error {
	if p.pathIsIdentifier(identifierPath) {
		currentDir := filepath.Dir(identifierPath)
		for _, target := range p.TargetDirectories {
			possibleTargetDir := filepath.Join(currentDir, target)

			if p.pathIsTargetDirectory(possibleTargetDir) {
				if dryRun {
					fmt.Printf("Deleting %s in %s \n", target, possibleTargetDir)
				} else {
					os.RemoveAll(possibleTargetDir)
				}
			}
		}
	}

	return nil
}

func (p *Project) TargetDirectoryIsDirectChild(targetPath string) bool {
	if p.pathIsTargetDirectory(targetPath) {
		for _, identifier := range p.Identifiers {
			// If the parent directory contains an identifier, it's a target directory.
			// Skip it for performance.
			currentDir := path.Dir(targetPath)
			possibleIdentifierPath := path.Join(currentDir, identifier)
			if _, err := os.Stat(possibleIdentifierPath); err == nil {
				return true
			}
		}
	}

	return false
}
