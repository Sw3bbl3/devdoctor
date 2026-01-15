package plugin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type PluginResult struct {
	Name   string
	Output string
	Err    error
}

func RunAllPlugins(projectPath string) []PluginResult {
	pluginDir := filepath.Join(projectPath, "devdoctor.d")
	files, err := os.ReadDir(pluginDir)
	if err != nil {
		return nil // no plugins
	}
	var results []PluginResult
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		full := filepath.Join(pluginDir, name)
		var cmd *exec.Cmd
		if strings.HasSuffix(name, ".sh") && runtime.GOOS != "windows" {
			cmd = exec.Command("bash", full)
		} else if strings.HasSuffix(name, ".ps1") && runtime.GOOS == "windows" {
			cmd = exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", full)
		} else if strings.HasSuffix(name, ".bat") && runtime.GOOS == "windows" {
			cmd = exec.Command(full)
		} else if strings.HasSuffix(name, ".exe") {
			cmd = exec.Command(full)
		} else {
			continue // skip unknown
		}
		out, err := cmd.CombinedOutput()
		results = append(results, PluginResult{
			Name:   name,
			Output: string(out),
			Err:    err,
		})
	}
	return results
}
