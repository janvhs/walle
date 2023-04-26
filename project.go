package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type Project struct {
	Name           string
	Configurations []Configuration
}

func (p *Project) Clean(absolutePath string) {
	fmt.Printf("Cleaning %s\n", p.Name)
	for _, config := range p.Configurations {
		if config.Matches(absolutePath) {
			config.Clean(absolutePath)
		}
	}
}

// TODO: Do I need generics?
type Configuration struct {
	Identifier      Identifier
	RelativeTargets []string
}

func (c *Configuration) MatchesOptimistically(absolutePath string) bool {
	return c.Identifier.MatchesOptimistically(absolutePath)
}

func (c *Configuration) GenerateTargetList(absolutePath string) []string {
	targets := make([]string, 0)
	for _, relativeTarget := range c.RelativeTargets {
		absoluteTarget := filepath.Join(absolutePath, relativeTarget)
		if _, err := os.Stat(absoluteTarget); err == nil {
			targets = append(targets, absoluteTarget)
		}
	}

	return targets
}

func (c *Configuration) Matches(absolutePath string) bool {
	return c.Identifier.Matches(absolutePath)
}

func (c *Configuration) Clean(absolutePath string) {
	for _, relativeTarget := range c.RelativeTargets {
		absoluteTarget := filepath.Join(absolutePath, relativeTarget)

		_, err := os.Stat(absoluteTarget)
		if err != nil {
			return
		}
		// os.Rename(absoluteTarget, absoluteTarget+".bak")
		fmt.Printf("Removing %s\n", absoluteTarget)
	}
}
