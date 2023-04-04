package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"bode.fun/walle/project"
)

func main() {
	var dryRun bool = true

	projects := []project.Project{
		{
			Name:              "Node.js",
			Identifiers:       []string{"package.json"},
			TargetDirectories: []string{"node_modules"},
		},
		{
			Name:              "Swift",
			Identifiers:       []string{"Package.swift"},
			TargetDirectories: []string{".build"},
		},
		{
			Name:              "Rust",
			Identifiers:       []string{"Cargo.toml"},
			TargetDirectories: []string{"target"},
		},
	}

	rootPath, homeErr := os.UserHomeDir()
	if homeErr != nil {
		log.Fatal("Home directory not found.")
	}

	if len(os.Args) > 1 {
		rootPath = os.Args[1]
	}

	filepath.WalkDir(rootPath, func(fPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden and git directories
		if d.IsDir() && d.Name()[0] == '.' || d.Name() == ".git" {
			return filepath.SkipDir
		}

		for _, project := range projects {

			// If a normal directory is named e.g. "target", it will not be skipped.
			// Actual build directories will get skipped.
			if project.TargetDirectoryIsDirectChild(fPath) {
				return filepath.SkipDir
			}

			project.Clean(fPath, dryRun)
		}

		return nil
	})
}
