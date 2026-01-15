package envcheck

import (
	"fmt"
	"os/exec"
	"strings"
)

type Tool struct {
	Name    string
	Command string
	Args    []string
	Parse   func(string) string // parses version output
	Min     string              // minimum recommended version
}

type ToolStatus struct {
	Name    string
	Found   bool
	Version string
	Warn    string
}

var tools = []Tool{
	{
		Name:    "Go",
		Command: "go",
		Args:    []string{"version"},
		Parse: func(out string) string {
			parts := strings.Fields(out)
			if len(parts) >= 3 {
				return strings.TrimPrefix(parts[2], "go")
			}
			return ""
		},
		Min: "1.20",
	},
	{
		Name:    "Node.js",
		Command: "node",
		Args:    []string{"--version"},
		Parse: func(out string) string {
			return strings.TrimPrefix(strings.TrimSpace(out), "v")
		},
		Min: "16.0.0",
	},
	{
		Name:    "npm",
		Command: "npm",
		Args:    []string{"--version"},
		Parse:  strings.TrimSpace,
		Min:   "8.0.0",
	},
	{
		Name:    "Python",
		Command: "python",
		Args:    []string{"--version"},
		Parse: func(out string) string {
			return strings.TrimPrefix(strings.TrimSpace(out), "Python ")
		},
		Min: "3.8",
	},
	{
		Name:    "Java",
		Command: "java",
		Args:    []string{"-version"},
		Parse: func(out string) string {
			lines := strings.Split(out, "\n")
			if len(lines) > 0 {
				return strings.Trim(strings.Split(lines[0], "\"")[1], "\"")
			}
			return ""
		},
		Min: "11",
	},
	{
		Name:    ".NET",
		Command: "dotnet",
		Args:    []string{"--version"},
		Parse:  strings.TrimSpace,
		Min:   "6.0",
	},
	{
		Name:    "Rust",
		Command: "rustc",
		Args:    []string{"--version"},
		Parse: func(out string) string {
			parts := strings.Fields(out)
			if len(parts) >= 2 {
				return parts[1]
			}
			return ""
		},
		Min: "1.60",
	},
	{
		Name:    "Ruby",
		Command: "ruby",
		Args:    []string{"--version"},
		Parse: func(out string) string {
			parts := strings.Fields(out)
			if len(parts) >= 2 {
				return parts[1]
			}
			return ""
		},
		Min: "2.7",
	},
	{
		Name:    "Docker",
		Command: "docker",
		Args:    []string{"--version"},
		Parse: func(out string) string {
			parts := strings.Fields(out)
			if len(parts) >= 3 {
				return parts[2]
			}
			return ""
		},
		Min: "20.10",
	},
}

func CheckAll() []ToolStatus {
	var results []ToolStatus
	for _, t := range tools {
		cmd := exec.Command(t.Command, t.Args...)
		out, err := cmd.CombinedOutput()
		status := ToolStatus{Name: t.Name}
		if err == nil {
			status.Found = true
			status.Version = t.Parse(string(out))
			if t.Min != "" && status.Version != "" && compareVersion(status.Version, t.Min) < 0 {
				status.Warn = fmt.Sprintf("Version %s is below recommended %s", status.Version, t.Min)
			}
		} else {
			status.Found = false
			status.Warn = "Not found"
		}
		results = append(results, status)
	}
	return results
}

// compareVersion returns -1 if a < b, 0 if a == b, 1 if a > b
func compareVersion(a, b string) int {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")
	for i := 0; i < len(aParts) || i < len(bParts); i++ {
		var ai, bi int
		if i < len(aParts) {
			fmt.Sscanf(aParts[i], "%d", &ai)
		}
		if i < len(bParts) {
			fmt.Sscanf(bParts[i], "%d", &bi)
		}
		if ai < bi {
			return -1
		}
		if ai > bi {
			return 1
		}
	}
	return 0
}
