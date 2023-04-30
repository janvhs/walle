package main

import (
	"os"
	"path/filepath"
)

type Project struct {
	Name           string
	Configurations []Configuration
}

type Configuration struct {
	Identifier      Identifier
	RelativeTargets []string
}

func (c *Configuration) MatchesOptimistically(absolutePath string) bool {
	return c.Identifier.MatchesOptimistically(absolutePath)
}

func (c *Configuration) Matches(absolutePath string) bool {
	return c.Identifier.Matches(absolutePath)
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
