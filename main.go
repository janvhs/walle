package main

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
					// "" is the directory where the file was found in
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

	// TODO: When flagRootPath is not a subpath of the user's homedir or equals the homedir, print a warning

	targetChan := make(chan MatchInfo)
	go scanDirs(flagRootPath, projects, targetChan)

	rl, err := readline.New("")
	if err != nil {
		log.Fatalln(err)
	}

	defer rl.Close()

	for target := range targetChan {
		fmt.Printf("Found a %s project with the following directories:\n", target.ProgrammingLanguage)

		for dir, size := range target.TargetDirs {
			fmt.Printf("  %s, %.3fMiB\n", dir, size)
		}

		rl.SetPrompt(fmt.Sprintf("Do you want to delete these directories? [y/N] "))

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
			for dir := range target.TargetDirs {
				fmt.Printf("Deleting %s ...\n", dir)
				if !flagDry {
					os.RemoveAll(dir)
				}
			}
		}
	}
}

type MatchInfo struct {
	ProgrammingLanguage string
	TargetDirs          map[string]float64
}

func scanDirs(rootPath string, projects []Project, targetChannel chan MatchInfo) {
	defer close(targetChannel)

	knownTargetDirs := make(map[string]interface{}, 0)

	filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		for knownTargetDir := range knownTargetDirs {
			if isPathInDir(path, knownTargetDir) {
				return filepath.SkipDir
			}
		}

		newTargetsWithLanguage := make(map[string]string, 0)
		for _, project := range projects {
			for _, config := range project.Configurations {
				if config.MatchesOptimistically(path) {
					for _, target := range config.GenerateTargetList(path) {
						newTargetsWithLanguage[target] = project.Name
					}
				}
			}
		}

		languageWithTargets := make(map[string][]string, 0)

		for target, lang := range newTargetsWithLanguage {
			languageWithTargets[lang] = append(languageWithTargets[lang], target)
			knownTargetDirs[target] = nil
		}

		for lang, targets := range languageWithTargets {
			matchInfo := MatchInfo{
				ProgrammingLanguage: lang,
				TargetDirs:          make(map[string]float64, 0),
			}

			for _, target := range targets {
				size, _ := calculateDirectorySize(target)
				matchInfo.TargetDirs[target] = size
			}

			targetChannel <- matchInfo
		}

		return nil
	})
}

func calculateDirectorySize(dirRoot string) (float64, error) {
	var size int64
	err := filepath.Walk(dirRoot, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})

	sizeInMiB := float64(size) / 1024 / 1024

	return sizeInMiB, err
}

func isPathInDir(path, dir string) bool {
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
