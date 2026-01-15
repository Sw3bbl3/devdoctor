package detector

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectNodeJS(t *testing.T) {
	tmpDir := t.TempDir()
	packageJSON := filepath.Join(tmpDir, "package.json")
	if err := os.WriteFile(packageJSON, []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}

	project := detectNodeJS(tmpDir)
	if project == nil {
		t.Fatal("Expected Node.js project to be detected")
	}
	if project.Name != "Node.js" {
		t.Errorf("Expected project name to be 'Node.js', got %s", project.Name)
	}
	if len(project.ConfigFiles) != 1 || project.ConfigFiles[0] != "package.json" {
		t.Errorf("Expected config file 'package.json', got %v", project.ConfigFiles)
	}
}

func TestDetectPython(t *testing.T) {
	tests := []struct {
		name       string
		files      []string
		wantDetect bool
	}{
		{"requirements.txt", []string{"requirements.txt"}, true},
		{"setup.py", []string{"setup.py"}, true},
		{"pyproject.toml", []string{"pyproject.toml"}, true},
		{"Pipfile", []string{"Pipfile"}, true},
		{"no files", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			for _, file := range tt.files {
				filePath := filepath.Join(tmpDir, file)
				if err := os.WriteFile(filePath, []byte(""), 0644); err != nil {
					t.Fatal(err)
				}
			}

			project := detectPython(tmpDir)
			if tt.wantDetect && project == nil {
				t.Fatal("Expected Python project to be detected")
			}
			if !tt.wantDetect && project != nil {
				t.Fatal("Expected Python project not to be detected")
			}
		})
	}
}

func TestDetectGo(t *testing.T) {
	tmpDir := t.TempDir()
	goMod := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test"), 0644); err != nil {
		t.Fatal(err)
	}

	project := detectGo(tmpDir)
	if project == nil {
		t.Fatal("Expected Go project to be detected")
	}
	if project.Name != "Go" {
		t.Errorf("Expected project name to be 'Go', got %s", project.Name)
	}
}

func TestDetectJava(t *testing.T) {
	tests := []struct {
		name  string
		files []string
	}{
		{"Maven", []string{"pom.xml"}},
		{"Gradle", []string{"build.gradle"}},
		{"Gradle Kotlin", []string{"build.gradle.kts"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			for _, file := range tt.files {
				filePath := filepath.Join(tmpDir, file)
				if err := os.WriteFile(filePath, []byte(""), 0644); err != nil {
					t.Fatal(err)
				}
			}

			project := detectJava(tmpDir)
			if project == nil {
				t.Fatal("Expected Java project to be detected")
			}
			if project.Name != "Java" {
				t.Errorf("Expected project name to be 'Java', got %s", project.Name)
			}
		})
	}
}

func TestDetectRuby(t *testing.T) {
	tmpDir := t.TempDir()
	gemfile := filepath.Join(tmpDir, "Gemfile")
	if err := os.WriteFile(gemfile, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	project := detectRuby(tmpDir)
	if project == nil {
		t.Fatal("Expected Ruby project to be detected")
	}
	if project.Name != "Ruby" {
		t.Errorf("Expected project name to be 'Ruby', got %s", project.Name)
	}
}

func TestDetectRust(t *testing.T) {
	tmpDir := t.TempDir()
	cargoToml := filepath.Join(tmpDir, "Cargo.toml")
	if err := os.WriteFile(cargoToml, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	project := detectRust(tmpDir)
	if project == nil {
		t.Fatal("Expected Rust project to be detected")
	}
	if project.Name != "Rust" {
		t.Errorf("Expected project name to be 'Rust', got %s", project.Name)
	}
}

func TestDetectDotNet(t *testing.T) {
	tests := []struct {
		name string
		file string
	}{
		{"C# project", "test.csproj"},
		{"F# project", "test.fsproj"},
		{"VB project", "test.vbproj"},
		{"Solution", "test.sln"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filePath := filepath.Join(tmpDir, tt.file)
			if err := os.WriteFile(filePath, []byte(""), 0644); err != nil {
				t.Fatal(err)
			}

			project := detectDotNet(tmpDir)
			if project == nil {
				t.Fatal("Expected .NET project to be detected")
			}
			if project.Name != ".NET" {
				t.Errorf("Expected project name to be '.NET', got %s", project.Name)
			}
		})
	}
}

func TestDetectDocker(t *testing.T) {
	tests := []struct {
		name  string
		files []string
	}{
		{"Dockerfile only", []string{"Dockerfile"}},
		{"docker-compose.yml", []string{"docker-compose.yml"}},
		{"docker-compose.yaml", []string{"docker-compose.yaml"}},
		{"Both", []string{"Dockerfile", "docker-compose.yml"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			for _, file := range tt.files {
				filePath := filepath.Join(tmpDir, file)
				if err := os.WriteFile(filePath, []byte(""), 0644); err != nil {
					t.Fatal(err)
				}
			}

			project := detectDocker(tmpDir)
			if project == nil {
				t.Fatal("Expected Docker project to be detected")
			}
			if project.Name != "Docker" {
				t.Errorf("Expected project name to be 'Docker', got %s", project.Name)
			}
		})
	}
}

func TestDetectorRegistry(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create multiple project type markers
	os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte("{}"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test"), 0644)

	registry := NewDetectorRegistry()
	projects := registry.Detect(tmpDir)

	if len(projects) != 2 {
		t.Errorf("Expected 2 projects to be detected, got %d", len(projects))
	}

	// Check that both Node.js and Go are detected
	foundNodeJS := false
	foundGo := false
	for _, p := range projects {
		if p.Name == "Node.js" {
			foundNodeJS = true
		}
		if p.Name == "Go" {
			foundGo = true
		}
	}

	if !foundNodeJS {
		t.Error("Expected Node.js to be detected")
	}
	if !foundGo {
		t.Error("Expected Go to be detected")
	}
}
