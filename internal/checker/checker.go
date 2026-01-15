package checker

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Sw3bbl3/devdoctor/internal/detector"
)

// Severity levels for issues
type Severity string

const (
	SeverityError   Severity = "ERROR"
	SeverityWarning Severity = "WARNING"
	SeverityInfo    Severity = "INFO"
)

// Issue represents a diagnostic issue
type Issue struct {
	Severity    Severity
	ProjectType string
	Message     string
	Suggestion  string
}

// CheckProject runs all checks for a detected project
func CheckProject(path string, project *detector.ProjectType) []Issue {
	issues := []Issue{}

	// Check if required tools are installed
	for _, tool := range project.RequiredTools {
		if !isCommandAvailable(tool) {
			issues = append(issues, Issue{
				Severity:    SeverityError,
				ProjectType: project.Name,
				Message:     fmt.Sprintf("Required tool '%s' is not installed or not in PATH", tool),
				Suggestion:  getInstallSuggestion(tool),
			})
		}
	}

	// Run project-specific checks
	switch project.Name {
	case "Node.js":
		issues = append(issues, checkNodeJS(path)...)
	case "Python":
		issues = append(issues, checkPython(path)...)
	case "Go":
		issues = append(issues, checkGo(path)...)
	case "Java":
		issues = append(issues, checkJava(path)...)
	case "Ruby":
		issues = append(issues, checkRuby(path)...)
	case "Rust":
		issues = append(issues, checkRust(path)...)
	case ".NET":
		issues = append(issues, checkDotNet(path)...)
	case "Docker":
		issues = append(issues, checkDocker(path)...)
	}

	return issues
}

func isCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func getInstallSuggestion(tool string) string {
	suggestions := map[string]string{
		"node":   "Install Node.js from https://nodejs.org/ or use a version manager like nvm",
		"npm":    "npm is included with Node.js. Install from https://nodejs.org/",
		"python": "Install Python from https://python.org/ or use pyenv for version management",
		"pip":    "pip is included with Python 3.4+. Reinstall Python or install pip separately",
		"go":     "Install Go from https://golang.org/dl/",
		"java":   "Install Java JDK from https://adoptium.net/ or your system package manager",
		"mvn":    "Install Maven from https://maven.apache.org/ or use your system package manager",
		"gradle": "Install Gradle from https://gradle.org/ or use the gradle wrapper (./gradlew)",
		"ruby":   "Install Ruby from https://www.ruby-lang.org/ or use rbenv/rvm",
		"bundle": "Install bundler with: gem install bundler",
		"cargo":  "Install Rust from https://rustup.rs/",
		"rustc":  "Install Rust from https://rustup.rs/",
		"dotnet": "Install .NET SDK from https://dotnet.microsoft.com/download",
		"docker": "Install Docker from https://docs.docker.com/get-docker/",
	}
	if suggestion, ok := suggestions[tool]; ok {
		return suggestion
	}
	return fmt.Sprintf("Please install %s and ensure it's in your PATH", tool)
}

func checkNodeJS(path string) []Issue {
	issues := []Issue{}

	// Check if node_modules exists
	nodeModulesPath := filepath.Join(path, "node_modules")
	if _, err := os.Stat(nodeModulesPath); os.IsNotExist(err) {
		issues = append(issues, Issue{
			Severity:    SeverityWarning,
			ProjectType: "Node.js",
			Message:     "Dependencies not installed (node_modules directory not found)",
			Suggestion:  "Run 'npm install' or 'yarn install' to install dependencies",
		})
	}

	// Check package.json for engines
	packageJSONPath := filepath.Join(path, "package.json")
	data, err := os.ReadFile(packageJSONPath)
	if err == nil {
		var packageJSON map[string]interface{}
		if json.Unmarshal(data, &packageJSON) == nil {
			if engines, ok := packageJSON["engines"].(map[string]interface{}); ok {
				if nodeVersion, ok := engines["node"].(string); ok {
					issues = append(issues, Issue{
						Severity:    SeverityInfo,
						ProjectType: "Node.js",
						Message:     fmt.Sprintf("Project requires Node.js version: %s", nodeVersion),
						Suggestion:  "Verify your Node.js version matches the requirement using 'node --version'",
					})
				}
			}
		}
	}

	return issues
}

func checkPython(path string) []Issue {
	issues := []Issue{}

	// Check for virtual environment
	venvDirs := []string{"venv", ".venv", "env", ".env"}
	venvExists := false
	for _, dir := range venvDirs {
		if _, err := os.Stat(filepath.Join(path, dir)); err == nil {
			venvExists = true
			break
		}
	}

	if !venvExists {
		issues = append(issues, Issue{
			Severity:    SeverityWarning,
			ProjectType: "Python",
			Message:     "No virtual environment detected",
			Suggestion:  "Create a virtual environment with 'python -m venv venv' and activate it",
		})
	}

	// Check if requirements are installed
	if fileExists := func(name string) bool {
		_, err := os.Stat(filepath.Join(path, name))
		return err == nil
	}; fileExists("requirements.txt") {
		issues = append(issues, Issue{
			Severity:    SeverityInfo,
			ProjectType: "Python",
			Message:     "Found requirements.txt",
			Suggestion:  "Install dependencies with 'pip install -r requirements.txt'",
		})
	}

	return issues
}

func checkGo(path string) []Issue {
	issues := []Issue{}

	// Check if go.sum exists
	if _, err := os.Stat(filepath.Join(path, "go.sum")); os.IsNotExist(err) {
		issues = append(issues, Issue{
			Severity:    SeverityWarning,
			ProjectType: "Go",
			Message:     "go.sum not found - dependencies may not be downloaded",
			Suggestion:  "Run 'go mod download' or 'go mod tidy' to download dependencies",
		})
	}

	// Check vendor directory
	if _, err := os.Stat(filepath.Join(path, "vendor")); err == nil {
		issues = append(issues, Issue{
			Severity:    SeverityInfo,
			ProjectType: "Go",
			Message:     "Using vendored dependencies",
			Suggestion:  "Dependencies are vendored. Run 'go mod vendor' to update if needed",
		})
	}

	return issues
}

func checkJava(path string) []Issue {
	issues := []Issue{}

	// Check for Maven
	if _, err := os.Stat(filepath.Join(path, "pom.xml")); err == nil {
		// Check if .m2 or target exists
		if _, err := os.Stat(filepath.Join(path, "target")); os.IsNotExist(err) {
			issues = append(issues, Issue{
				Severity:    SeverityWarning,
				ProjectType: "Java",
				Message:     "Maven project not built (target directory not found)",
				Suggestion:  "Run 'mvn install' or 'mvn package' to build the project",
			})
		}
	}

	// Check for Gradle
	if _, err := os.Stat(filepath.Join(path, "build.gradle")); err == nil {
		if _, err := os.Stat(filepath.Join(path, "build")); os.IsNotExist(err) {
			issues = append(issues, Issue{
				Severity:    SeverityWarning,
				ProjectType: "Java",
				Message:     "Gradle project not built (build directory not found)",
				Suggestion:  "Run 'gradle build' or './gradlew build' to build the project",
			})
		}
	}

	return issues
}

func checkRuby(path string) []Issue {
	issues := []Issue{}

	// Check if Gemfile.lock exists
	if _, err := os.Stat(filepath.Join(path, "Gemfile.lock")); os.IsNotExist(err) {
		issues = append(issues, Issue{
			Severity:    SeverityWarning,
			ProjectType: "Ruby",
			Message:     "Gemfile.lock not found - dependencies may not be installed",
			Suggestion:  "Run 'bundle install' to install dependencies",
		})
	}

	return issues
}

func checkRust(path string) []Issue {
	issues := []Issue{}

	// Check if Cargo.lock exists
	if _, err := os.Stat(filepath.Join(path, "Cargo.lock")); os.IsNotExist(err) {
		issues = append(issues, Issue{
			Severity:    SeverityInfo,
			ProjectType: "Rust",
			Message:     "Cargo.lock not found",
			Suggestion:  "Run 'cargo build' to build and generate Cargo.lock",
		})
	}

	// Check if target directory exists
	if _, err := os.Stat(filepath.Join(path, "target")); os.IsNotExist(err) {
		issues = append(issues, Issue{
			Severity:    SeverityWarning,
			ProjectType: "Rust",
			Message:     "Project not built (target directory not found)",
			Suggestion:  "Run 'cargo build' to build the project",
		})
	}

	return issues
}

func checkDotNet(path string) []Issue {
	issues := []Issue{}

	// Check for bin/obj directories
	hasBin := false
	hasObj := false

	if _, err := os.Stat(filepath.Join(path, "bin")); err == nil {
		hasBin = true
	}
	if _, err := os.Stat(filepath.Join(path, "obj")); err == nil {
		hasObj = true
	}

	if !hasBin && !hasObj {
		issues = append(issues, Issue{
			Severity:    SeverityWarning,
			ProjectType: ".NET",
			Message:     "Project not built (bin/obj directories not found)",
			Suggestion:  "Run 'dotnet restore' and 'dotnet build' to build the project",
		})
	}

	return issues
}

func checkDocker(path string) []Issue {
	issues := []Issue{}

	// Check if Docker daemon is running
	if isCommandAvailable("docker") {
		cmd := exec.Command("docker", "info")
		if err := cmd.Run(); err != nil {
			issues = append(issues, Issue{
				Severity:    SeverityError,
				ProjectType: "Docker",
				Message:     "Docker daemon is not running",
				Suggestion:  "Start Docker Desktop or the Docker daemon",
			})
		}
	}

	// Check for .env file if docker-compose is present
	hasCompose := false
	if _, err := os.Stat(filepath.Join(path, "docker-compose.yml")); err == nil {
		hasCompose = true
	}
	if _, err := os.Stat(filepath.Join(path, "docker-compose.yaml")); err == nil {
		hasCompose = true
	}

	if hasCompose {
		// Look for .env.example or .env.sample
		hasEnvExample := false
		if _, err := os.Stat(filepath.Join(path, ".env.example")); err == nil {
			hasEnvExample = true
		}
		if _, err := os.Stat(filepath.Join(path, ".env.sample")); err == nil {
			hasEnvExample = true
		}

		if hasEnvExample {
			if _, err := os.Stat(filepath.Join(path, ".env")); os.IsNotExist(err) {
				issues = append(issues, Issue{
					Severity:    SeverityWarning,
					ProjectType: "Docker",
					Message:     "Environment file (.env) not found but example exists",
					Suggestion:  "Copy .env.example to .env and configure your environment variables",
				})
			}
		}
	}

	return issues
}

// Check for common environment files
func checkEnvironmentFiles(path string) []Issue {
	issues := []Issue{}

	// Look for .env.example or .env.sample
	hasEnvExample := false
	exampleFile := ""
	for _, filename := range []string{".env.example", ".env.sample", "env.example"} {
		if _, err := os.Stat(filepath.Join(path, filename)); err == nil {
			hasEnvExample = true
			exampleFile = filename
			break
		}
	}

	if hasEnvExample {
		if _, err := os.Stat(filepath.Join(path, ".env")); os.IsNotExist(err) {
			issues = append(issues, Issue{
				Severity:    SeverityWarning,
				ProjectType: "General",
				Message:     "Environment file (.env) not found",
				Suggestion:  fmt.Sprintf("Copy %s to .env and configure your environment variables", exampleFile),
			})
		}
	}

	return issues
}

// Helper function to parse version requirements
func parseVersionRequirement(requirement string) string {
	// Remove common prefixes and return cleaned version
	cleaned := strings.TrimSpace(requirement)
	cleaned = strings.TrimPrefix(cleaned, "^")
	cleaned = strings.TrimPrefix(cleaned, "~")
	cleaned = strings.TrimPrefix(cleaned, ">=")
	cleaned = strings.TrimPrefix(cleaned, "<=")
	cleaned = strings.TrimPrefix(cleaned, ">")
	cleaned = strings.TrimPrefix(cleaned, "<")
	return cleaned
}
