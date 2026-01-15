package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Sw3bbl3/devdoctor/internal/checker"
	"github.com/Sw3bbl3/devdoctor/internal/detector"
	"github.com/Sw3bbl3/devdoctor/internal/reporter"
)

func main() {
	var path string
	flag.StringVar(&path, "path", ".", "Path to the project directory to diagnose")
	flag.Parse()

	// Resolve to absolute path
	absPath, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	if path != "." {
		absPath = path
	}

	// Detect project types
	detectors := detector.NewDetectorRegistry()
	detectedProjects := detectors.Detect(absPath)

	if len(detectedProjects) == 0 {
		fmt.Println("No supported project types detected in", absPath)
		fmt.Println("\nDevDoctor currently supports:")
		fmt.Println("  - Node.js (package.json)")
		fmt.Println("  - Python (requirements.txt, setup.py, pyproject.toml)")
		fmt.Println("  - Go (go.mod)")
		fmt.Println("  - Java (pom.xml, build.gradle)")
		fmt.Println("  - Ruby (Gemfile)")
		fmt.Println("  - Rust (Cargo.toml)")
		fmt.Println("  - .NET (*.csproj, *.sln)")
		os.Exit(0)
	}

	// Run checks for each detected project type
	allIssues := []checker.Issue{}
	for _, project := range detectedProjects {
		issues := checker.CheckProject(absPath, project)
		allIssues = append(allIssues, issues...)
	}

	// Report results
	reporter.Report(absPath, detectedProjects, allIssues)

	// Exit with code 1 if there are issues
	if len(allIssues) > 0 {
		os.Exit(1)
	}
}
