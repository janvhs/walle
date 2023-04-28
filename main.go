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
				{
					Identifier: &FileNameIdentifier{
						Name: "pnpmrc",
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
		{
			Name: "Python",
			Configurations: []Configuration{
				{
					Identifier: &FileExtensionInDirectoryIdentifier{
						Directory: "__pycache__",
						Extension: ".pyc",
					},
					RelativeTargets: []string{"__pycache__"},
				},
				{
					Identifier: &FileNameIdentifier{
						Name: "pyvenv.cfg",
					},
					RelativeTargets: []string{""},
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

	targetChan := make(chan string)
	go scanDirs(flagRootPath, projects, targetChan)

	rl, err := readline.New("")
	if err != nil {
		log.Fatalln(err)
	}

	defer rl.Close()

	for target := range targetChan {
		rl.SetPrompt(fmt.Sprintf("Delete %s ? [y/N] ", target))

		line, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				break
			} else {
				continue
			}
		}

		result := strings.ToUpper(line)
		result = strings.TrimSpace(result)
		if result == "Y" {
			if flagDry {
				fmt.Printf("Deleting %s ...\n", target)
			} else {
				os.RemoveAll(target)
			}
		}
	}
}

func scanDirs(rootPath string, projects []Project, targets chan string) {
	defer close(targets)

	targetDirs := make([]string, 0)

	filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		for _, knownTargetDir := range targetDirs {
			if isInDir(path, knownTargetDir) {
				return filepath.SkipDir
			}
		}

		for _, project := range projects {
			for _, config := range project.Configurations {
				if config.MatchesOptimistically(path) {
					newTargets := config.GenerateTargetList(path)

					for _, target := range newTargets {
						targets <- target
					}

					targetDirs = append(targetDirs, newTargets...)
				}
			}
		}

		return nil
	})
}

func isInDir(path, dir string) bool {
	pathList := strings.Split(path, string(os.PathSeparator))
	pathListLen := len(pathList)
	dirList := strings.Split(dir, string(os.PathSeparator))

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
