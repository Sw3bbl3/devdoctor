package checker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Sw3bbl3/devdoctor/internal/detector"
)

func TestCheckNodeJS(t *testing.T) {
	tmpDir := t.TempDir()

	// Test without node_modules
	issues := checkNodeJS(tmpDir)
	if len(issues) == 0 {
		t.Error("Expected issues when node_modules is missing")
	}

	// Create node_modules
	nodeModulesPath := filepath.Join(tmpDir, "node_modules")
	if err := os.Mkdir(nodeModulesPath, 0755); err != nil {
		t.Fatal(err)
	}

	issues = checkNodeJS(tmpDir)
	hasNodeModulesWarning := false
	for _, issue := range issues {
		if issue.Message == "Dependencies not installed (node_modules directory not found)" {
			hasNodeModulesWarning = true
		}
	}
	if hasNodeModulesWarning {
		t.Error("Should not have node_modules warning when directory exists")
	}
}

func TestCheckPython(t *testing.T) {
	tmpDir := t.TempDir()

	issues := checkPython(tmpDir)
	hasVenvWarning := false
	for _, issue := range issues {
		if issue.Message == "No virtual environment detected" {
			hasVenvWarning = true
		}
	}
	if !hasVenvWarning {
		t.Error("Expected virtual environment warning")
	}

	// Create venv
	venvPath := filepath.Join(tmpDir, "venv")
	if err := os.Mkdir(venvPath, 0755); err != nil {
		t.Fatal(err)
	}

	issues = checkPython(tmpDir)
	hasVenvWarning = false
	for _, issue := range issues {
		if issue.Message == "No virtual environment detected" {
			hasVenvWarning = true
		}
	}
	if hasVenvWarning {
		t.Error("Should not have venv warning when venv exists")
	}
}

func TestCheckGo(t *testing.T) {
	tmpDir := t.TempDir()

	// Test without go.sum
	issues := checkGo(tmpDir)
	hasGoSumWarning := false
	for _, issue := range issues {
		if issue.Message == "go.sum not found - dependencies may not be downloaded" {
			hasGoSumWarning = true
		}
	}
	if !hasGoSumWarning {
		t.Error("Expected go.sum warning")
	}

	// Create go.sum
	goSumPath := filepath.Join(tmpDir, "go.sum")
	if err := os.WriteFile(goSumPath, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	issues = checkGo(tmpDir)
	hasGoSumWarning = false
	for _, issue := range issues {
		if issue.Message == "go.sum not found - dependencies may not be downloaded" {
			hasGoSumWarning = true
		}
	}
	if hasGoSumWarning {
		t.Error("Should not have go.sum warning when go.sum exists")
	}
}

func TestCheckJava(t *testing.T) {
	tmpDir := t.TempDir()

	// Create pom.xml
	pomPath := filepath.Join(tmpDir, "pom.xml")
	if err := os.WriteFile(pomPath, []byte("<project></project>"), 0644); err != nil {
		t.Fatal(err)
	}

	issues := checkJava(tmpDir)
	hasTargetWarning := false
	for _, issue := range issues {
		if issue.Message == "Maven project not built (target directory not found)" {
			hasTargetWarning = true
		}
	}
	if !hasTargetWarning {
		t.Error("Expected target directory warning")
	}

	// Create target directory
	targetPath := filepath.Join(tmpDir, "target")
	if err := os.Mkdir(targetPath, 0755); err != nil {
		t.Fatal(err)
	}

	issues = checkJava(tmpDir)
	hasTargetWarning = false
	for _, issue := range issues {
		if issue.Message == "Maven project not built (target directory not found)" {
			hasTargetWarning = true
		}
	}
	if hasTargetWarning {
		t.Error("Should not have target warning when target exists")
	}
}

func TestCheckRuby(t *testing.T) {
	tmpDir := t.TempDir()

	issues := checkRuby(tmpDir)
	hasGemfileLockWarning := false
	for _, issue := range issues {
		if issue.Message == "Gemfile.lock not found - dependencies may not be installed" {
			hasGemfileLockWarning = true
		}
	}
	if !hasGemfileLockWarning {
		t.Error("Expected Gemfile.lock warning")
	}
}

func TestCheckRust(t *testing.T) {
	tmpDir := t.TempDir()

	issues := checkRust(tmpDir)
	if len(issues) == 0 {
		t.Error("Expected issues for unbuilt Rust project")
	}

	// Create target directory
	targetPath := filepath.Join(tmpDir, "target")
	if err := os.Mkdir(targetPath, 0755); err != nil {
		t.Fatal(err)
	}

	issues = checkRust(tmpDir)
	hasTargetWarning := false
	for _, issue := range issues {
		if issue.Message == "Project not built (target directory not found)" {
			hasTargetWarning = true
		}
	}
	if hasTargetWarning {
		t.Error("Should not have target warning when target exists")
	}
}

func TestCheckDotNet(t *testing.T) {
	tmpDir := t.TempDir()

	issues := checkDotNet(tmpDir)
	hasBuildWarning := false
	for _, issue := range issues {
		if issue.Message == "Project not built (bin/obj directories not found)" {
			hasBuildWarning = true
		}
	}
	if !hasBuildWarning {
		t.Error("Expected build warning")
	}

	// Create bin directory
	binPath := filepath.Join(tmpDir, "bin")
	if err := os.Mkdir(binPath, 0755); err != nil {
		t.Fatal(err)
	}

	issues = checkDotNet(tmpDir)
	hasBuildWarning = false
	for _, issue := range issues {
		if issue.Message == "Project not built (bin/obj directories not found)" {
			hasBuildWarning = true
		}
	}
	if hasBuildWarning {
		t.Error("Should not have build warning when bin exists")
	}
}

func TestCheckProject(t *testing.T) {
	tmpDir := t.TempDir()

	project := &detector.ProjectType{
		Name:          "Node.js",
		ConfigFiles:   []string{"package.json"},
		RequiredTools: []string{"node", "npm"},
	}

	issues := CheckProject(tmpDir, project)
	if len(issues) == 0 {
		t.Error("Expected some issues for a fresh Node.js project")
	}
}

func TestGetInstallSuggestion(t *testing.T) {
	tests := []struct {
		tool string
		want string
	}{
		{"node", "Install Node.js from https://nodejs.org/ or use a version manager like nvm"},
		{"python", "Install Python from https://python.org/ or use pyenv for version management"},
		{"go", "Install Go from https://golang.org/dl/"},
		{"unknown", "Please install unknown and ensure it's in your PATH"},
	}

	for _, tt := range tests {
		t.Run(tt.tool, func(t *testing.T) {
			got := getInstallSuggestion(tt.tool)
			if got != tt.want {
				t.Errorf("getInstallSuggestion(%s) = %s, want %s", tt.tool, got, tt.want)
			}
		})
	}
}
