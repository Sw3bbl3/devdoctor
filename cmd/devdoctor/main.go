package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/Sw3bbl3/devdoctor/internal/checker"
	"github.com/Sw3bbl3/devdoctor/internal/detector"
	"github.com/Sw3bbl3/devdoctor/internal/reporter"
	"github.com/Sw3bbl3/devdoctor/internal/updater"
)

const version = "0.1.0"

func main() {
	var path string
	var showVersion bool
	var update bool
	var checkUpdate bool
	flag.StringVar(&path, "path", ".", "Path to the project directory to diagnose")
	flag.BoolVar(&showVersion, "version", false, "Print DevDoctor version")
	flag.BoolVar(&update, "update", false, "Update DevDoctor to the latest release")
	flag.BoolVar(&checkUpdate, "check-update", false, "Check if a newer version is available")
	flag.Parse()

	if showVersion {
		fmt.Println("DevDoctor", version)
		return
	}

	if checkUpdate {
		latest, err := updater.LatestVersion()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Update check failed: %v\n", err)
			os.Exit(1)
		}
		if latest == version {
			fmt.Printf("Up to date: %s (%s/%s)\n", version, runtime.GOOS, runtime.GOARCH)
		} else {
			fmt.Printf("New version available: %s (current %s)\n", latest, version)
		}
		return
	}

	if update {
		fmt.Printf("Updating DevDoctor (current %s)...\n", version)
		dest, err := updater.UpdateToLatest(version)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Update failed: %v\n", err)
			os.Exit(1)
		}
		if dest == "" {
			fmt.Println("Already up to date.")
		} else {
			fmt.Printf("Updated successfully: %s\n", dest)
		}
		return
	}

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
