// TODO: Add comments where they are useful
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/chzyer/readline"
)

var flagDry bool

// TODO: Extract styles to file like charm does it
var messageStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(""))

var scanningStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("8"))

var langStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("2"))

var dirListStyle = lipgloss.NewStyle().
	PaddingLeft(2).
	Foreground(lipgloss.Color("6"))

var destructiveQuestion = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("5"))

var strongWarningStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("1"))

var errorStyle = strongWarningStyle

func init() {
	flag.BoolVar(&flagDry, "dry", false, "Run without making changes")
}

type MatchInfo struct {
	ProgrammingLanguage string
	TargetDirs          map[string]float64
}

// Custom usage output
// TODO: Add version
func usage() {
	description := "Keep your computer tidy with this little helper 🤖"
	fmt.Fprintf(flag.CommandLine.Output(), "%s\n\nUsage:\n  %s [path]\n\nFlags:\n", description, os.Args[0])

	flag.PrintDefaults()
}

func errorUsage(err error) {
	flag.Usage()
	fmt.Fprintf(flag.CommandLine.Output(), "Error:\n  %s\n", errorStyle.Render(fmt.Sprint(err)))
	os.Exit(2)
}

func main() {
	// TODO: Add a version
	flag.Usage = usage
	flag.Parse()
	flag.CommandLine.Args()

	var err error
	flagRootPath, err := handleRootPath(flag.Arg(0))

	if err != nil {
		errorUsage(err)
	}

	// TODO: Extract this to a serialisable format and embed it with embed.FS
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
					Identifier: &FileExtensionIdentifier{
						// TODO: It is kind of annoying that __pycache__ comes up in every subdirectory. Maybe check if the parent dir already contains a pycache and include it into the same list? Idk I don't use python so I will not implement it.
						Directory: "__pycache__",
						Extension: ".pyc",
					},
					RelativeTargets: []string{"__pycache__"},
				},
				{
					Identifier: &FileNameIdentifier{
						Name: "pyvenv.cfg",
					},
					// if no directory is provided,
					// "" is the directory where the file was found in
					RelativeTargets: []string{""},
				},
			},
		},
	}

	targetChan := make(chan MatchInfo)
	go scanDirs(flagRootPath, projects, targetChan)

	fmt.Println(scanningStyle.Render("Scanning..."))

	// TODO: Change the UI to a bubbletea UI
	for target := range targetChan {
		lipgloss.DefaultRenderer().Output().ClearLines(1)
		err := handleTarget(target)
		if err != nil {
			if err == readline.ErrInterrupt {
				break
			} else {
				continue
			}
		}

		fmt.Println()
		fmt.Println(scanningStyle.Render("Scanning..."))
	}
	lipgloss.DefaultRenderer().Output().ClearLines(1)

	// TODO: After everything is done, print the total size, that is cleaned
}

// TODO: Replace readline with stdio
func handleTarget(target MatchInfo) error {
	// Currently this can never fail
	rl, err := readline.New(destructiveQuestion.Render("Do you want to delete these directories? [y/N] "))
	if err != nil {
		return err
	}

	defer rl.Close()

	fmt.Println(messageStyle.Render(fmt.Sprintf("Found a %s project with the following directories:", langStyle.Render(target.ProgrammingLanguage))))

	for dir, size := range target.TargetDirs {
		fmt.Println(dirListStyle.Render(fmt.Sprintf("- %s, %.3f MiB", dir, size)))
	}

	line, err := rl.Readline()
	if err != nil {
		return err
	}

	result := strings.ToUpper(line)
	result = strings.TrimSpace(result)
	if result == "Y" {
		for dir := range target.TargetDirs {
			fmt.Println(strongWarningStyle.Render(fmt.Sprintf("Deleting %s", dir)))
			if !flagDry {
				os.RemoveAll(dir)
			}
		}
	}
	return nil
}

func handleRootPath(rootPath string) (string, error) {
	// TODO: When rootPath is not a subpath of the user's homedir or equals the homedir, print an error and ask for force flag
	var err error
	var resolvedPath string
	if rootPath == "" {
		resolvedPath, err = os.Getwd()
		if err != nil {
			return resolvedPath, err
		}
	}

	resolvedPath = filepath.Clean(resolvedPath)
	resolvedPath, err = filepath.Abs(resolvedPath)

	if err != nil {
		return resolvedPath, err
	}

	return resolvedPath, nil
}

func scanDirs(rootPath string, projects []Project, targetChannel chan MatchInfo) {
	defer close(targetChannel)

	knownTargetDirs := make(map[string]interface{}, 0)

	// Walk
	filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err := skipIfDirIsKnown(&knownTargetDirs, path); err != nil {
			return err
		}

		var newTargetsWithLanguage map[string]string
		knownTargetDirs, newTargetsWithLanguage = collectNewTargets(path, knownTargetDirs, projects)

		// Make the language the primary key
		languageWithTargets := make(map[string][]string, 0)
		for target, lang := range newTargetsWithLanguage {
			languageWithTargets[lang] = append(languageWithTargets[lang], target)
		}

		// Notify the channel subscribers
		notifyTargetSubscribers(languageWithTargets, targetChannel)

		return nil
	})
}

func skipIfDirIsKnown(knownTargetDirs *map[string]interface{}, path string) error {
	// Skip if it is a target dir
	for knownTargetDir := range *knownTargetDirs {
		if isPathInDir(path, knownTargetDir) {
			return filepath.SkipDir
		}
	}

	return nil
}

func collectNewTargets(path string, knownTargetDirs map[string]interface{}, projects []Project) (updatedKnownTargetDirs map[string]interface{}, newTargetsWithLanguage map[string]string) {
	targetsWithLanguage := make(map[string]string, 0)
	for _, project := range projects {
		for _, config := range project.Configurations {
			if config.MatchesOptimistically(path) {
				// Add to targets to lists
				for _, target := range config.GenerateTargetList(path) {
					targetsWithLanguage[target] = project.Name
					knownTargetDirs[target] = nil
				}
			}
		}
	}

	return knownTargetDirs, targetsWithLanguage
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

func notifyTargetSubscribers(languageWithTargets map[string][]string, targetChannel chan MatchInfo) {
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
