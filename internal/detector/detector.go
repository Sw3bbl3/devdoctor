package detector

import (
	"os"
	"path/filepath"
)

// ProjectType represents a detected project type
type ProjectType struct {
	Name          string
	ConfigFiles   []string
	RequiredTools []string
}

// DetectorRegistry manages project type detection
type DetectorRegistry struct {
	detectors []func(path string) *ProjectType
}

// NewDetectorRegistry creates a new detector registry
func NewDetectorRegistry() *DetectorRegistry {
	registry := &DetectorRegistry{}
	registry.registerDetectors()
	return registry
}

func (r *DetectorRegistry) registerDetectors() {
	r.detectors = []func(string) *ProjectType{
		detectNodeJS,
		detectPython,
		detectGo,
		detectJava,
		detectRuby,
		detectRust,
		detectDotNet,
		detectDocker,
	}
}

// Detect scans the directory and returns all detected project types
func (r *DetectorRegistry) Detect(path string) []*ProjectType {
	var projects []*ProjectType
	for _, detector := range r.detectors {
		if project := detector(path); project != nil {
			projects = append(projects, project)
		}
	}
	return projects
}

func fileExists(path, filename string) bool {
	_, err := os.Stat(filepath.Join(path, filename))
	return err == nil
}

func detectNodeJS(path string) *ProjectType {
	if fileExists(path, "package.json") {
		return &ProjectType{
			Name:          "Node.js",
			ConfigFiles:   []string{"package.json"},
			RequiredTools: []string{"node", "npm"},
		}
	}
	return nil
}

func detectPython(path string) *ProjectType {
	configFiles := []string{}
	if fileExists(path, "requirements.txt") {
		configFiles = append(configFiles, "requirements.txt")
	}
	if fileExists(path, "setup.py") {
		configFiles = append(configFiles, "setup.py")
	}
	if fileExists(path, "pyproject.toml") {
		configFiles = append(configFiles, "pyproject.toml")
	}
	if fileExists(path, "Pipfile") {
		configFiles = append(configFiles, "Pipfile")
	}

	if len(configFiles) > 0 {
		return &ProjectType{
			Name:          "Python",
			ConfigFiles:   configFiles,
			RequiredTools: []string{"python", "pip"},
		}
	}
	return nil
}

func detectGo(path string) *ProjectType {
	if fileExists(path, "go.mod") {
		return &ProjectType{
			Name:          "Go",
			ConfigFiles:   []string{"go.mod"},
			RequiredTools: []string{"go"},
		}
	}
	return nil
}

func detectJava(path string) *ProjectType {
	configFiles := []string{}
	tools := []string{"java"}

	if fileExists(path, "pom.xml") {
		configFiles = append(configFiles, "pom.xml")
		tools = append(tools, "mvn")
	}
	if fileExists(path, "build.gradle") || fileExists(path, "build.gradle.kts") {
		if fileExists(path, "build.gradle") {
			configFiles = append(configFiles, "build.gradle")
		}
		if fileExists(path, "build.gradle.kts") {
			configFiles = append(configFiles, "build.gradle.kts")
		}
		tools = append(tools, "gradle")
	}

	if len(configFiles) > 0 {
		return &ProjectType{
			Name:          "Java",
			ConfigFiles:   configFiles,
			RequiredTools: tools,
		}
	}
	return nil
}

func detectRuby(path string) *ProjectType {
	if fileExists(path, "Gemfile") {
		return &ProjectType{
			Name:          "Ruby",
			ConfigFiles:   []string{"Gemfile"},
			RequiredTools: []string{"ruby", "bundle"},
		}
	}
	return nil
}

func detectRust(path string) *ProjectType {
	if fileExists(path, "Cargo.toml") {
		return &ProjectType{
			Name:          "Rust",
			ConfigFiles:   []string{"Cargo.toml"},
			RequiredTools: []string{"cargo", "rustc"},
		}
	}
	return nil
}

func detectDotNet(path string) *ProjectType {
	// Check for .csproj, .fsproj, .vbproj, or .sln files
	files, err := os.ReadDir(path)
	if err != nil {
		return nil
	}

	configFiles := []string{}
	for _, file := range files {
		name := file.Name()
		ext := filepath.Ext(name)
		if ext == ".csproj" || ext == ".fsproj" || ext == ".vbproj" || ext == ".sln" {
			configFiles = append(configFiles, name)
		}
	}

	if len(configFiles) > 0 {
		return &ProjectType{
			Name:          ".NET",
			ConfigFiles:   configFiles,
			RequiredTools: []string{"dotnet"},
		}
	}
	return nil
}

func detectDocker(path string) *ProjectType {
	if fileExists(path, "Dockerfile") || fileExists(path, "docker-compose.yml") || fileExists(path, "docker-compose.yaml") {
		configFiles := []string{}
		if fileExists(path, "Dockerfile") {
			configFiles = append(configFiles, "Dockerfile")
		}
		if fileExists(path, "docker-compose.yml") {
			configFiles = append(configFiles, "docker-compose.yml")
		}
		if fileExists(path, "docker-compose.yaml") {
			configFiles = append(configFiles, "docker-compose.yaml")
		}
		return &ProjectType{
			Name:          "Docker",
			ConfigFiles:   configFiles,
			RequiredTools: []string{"docker"},
		}
	}
	return nil
}
