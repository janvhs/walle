package main

// TODO: Move to old GitHub repo because this is gooood!

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
)

var flagRootPath string
var flagDry bool

func init() {
	flag.StringVar(&flagRootPath, "root", "", "The root where to search from")
	flag.BoolVar(&flagDry, "dry", false, "Run without making changes")
}

func main() {
	flag.Parse()

	projects := []Project{
		{
			Name: "JavaScript",
			Configurations: []Configuration{
				{
					Identifier: &FileNameIdentifier{
						Name: "package.json",
					},
					RelativeTargets: []string{
						"node_modules",
					},
				},
			},
		},
		{
			Name: "php",
			Configurations: []Configuration{
				{
					Identifier: &FileNameIdentifier{
						Name: "composer.json",
					},
					RelativeTargets: []string{
						"vendor",
					},
				},
			},
		},
		{
			Name: "Swift",
			Configurations: []Configuration{
				{
					Identifier: &FileNameIdentifier{
						Name: "Package.swift",
					},
					RelativeTargets: []string{
						".build",
					},
				},
			},
		},
		{
			Name: "Rust",
			Configurations: []Configuration{
				{
					Identifier: &FileNameIdentifier{
						Name: "Cargo.toml",
					},
					RelativeTargets: []string{
						"target",
					},
				},
			},
		},
	}

	if flagRootPath == "" {
		log.Fatalln("Root path can not be empty")
	}

	flagRootPath = filepath.Clean(flagRootPath)
	var err error
	flagRootPath, err = filepath.Abs(flagRootPath)

	if err != nil {
		log.Fatalln(err)
	}

	// Make a map with target and project
	targetDirs := make([]string, 0)

	filepath.WalkDir(flagRootPath, func(path string, d fs.DirEntry, err error) error {
		for _, knownTargetDir := range targetDirs {
			if isInDir(path, knownTargetDir) {
				return filepath.SkipDir
			}
		}
		for _, project := range projects {
			for _, config := range project.Configurations {
				if config.MatchesOptimistically(path) {
					targetDirs = append(targetDirs, config.GenerateTargetList(path)...)
				}
			}
		}

		return nil
	})

	for _, target := range targetDirs {
		rl, err := readline.New(fmt.Sprintf("Delete %s ? [y/N] ", target))
		if err != nil {
			continue
		}
		line, err := rl.Readline()
		if err != nil {
			continue
		}

		result := strings.ToUpper(line)
		if result == "Y" {
			if flagDry {
				fmt.Printf("Deleting %s ...\n", target)
			} else {
				os.RemoveAll(target)
			}
		}
	}
}

func isInDir(path, dir string) bool {
	// TODO: Move to filepath.SplitList
	pathList := strings.Split(path, "/")
	pathListLen := len(pathList)
	dirList := strings.Split(dir, "/")

	if pathListLen < len(dirList) {
		return false
	}

	for i, component := range dirList {
		if i >= pathListLen {
			break
		}

		pathComponent := pathList[i]
		if component != pathComponent {
			return false
		}
	}

	return true
}
