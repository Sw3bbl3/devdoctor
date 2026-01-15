package reporter

import (
	"fmt"
	"strings"

	"github.com/Sw3bbl3/devdoctor/internal/checker"
	"github.com/Sw3bbl3/devdoctor/internal/detector"
)

// Report outputs the diagnostic results
func Report(path string, projects []*detector.ProjectType, issues []checker.Issue) {
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                         DEVDOCTOR                             ║")
	fmt.Println("║              Project Diagnostic Report                        ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("Scanning: %s\n", path)
	fmt.Println()

	// Show detected project types
	fmt.Println("📋 Detected Project Types:")
	for _, project := range projects {
		fmt.Printf("  ✓ %s\n", project.Name)
		for _, configFile := range project.ConfigFiles {
			fmt.Printf("    - %s\n", configFile)
		}
	}
	fmt.Println()

	// Group issues by severity
	errors := []checker.Issue{}
	warnings := []checker.Issue{}
	infos := []checker.Issue{}

	for _, issue := range issues {
		switch issue.Severity {
		case checker.SeverityError:
			errors = append(errors, issue)
		case checker.SeverityWarning:
			warnings = append(warnings, issue)
		case checker.SeverityInfo:
			infos = append(infos, issue)
		}
	}

	// Report errors
	if len(errors) > 0 {
		fmt.Println("❌ ERRORS (Critical issues that prevent the project from running):")
		fmt.Println(strings.Repeat("─", 65))
		for _, issue := range errors {
			fmt.Printf("\n[%s] %s\n", issue.ProjectType, issue.Message)
			fmt.Printf("   💡 %s\n", issue.Suggestion)
		}
		fmt.Println()
	}

	// Report warnings
	if len(warnings) > 0 {
		fmt.Println("⚠️  WARNINGS (Issues that may cause problems):")
		fmt.Println(strings.Repeat("─", 65))
		for _, issue := range warnings {
			fmt.Printf("\n[%s] %s\n", issue.ProjectType, issue.Message)
			fmt.Printf("   💡 %s\n", issue.Suggestion)
		}
		fmt.Println()
	}

	// Report info
	if len(infos) > 0 {
		fmt.Println("ℹ️  INFORMATION (Helpful tips):")
		fmt.Println(strings.Repeat("─", 65))
		for _, issue := range infos {
			fmt.Printf("\n[%s] %s\n", issue.ProjectType, issue.Message)
			fmt.Printf("   💡 %s\n", issue.Suggestion)
		}
		fmt.Println()
	}

	// Summary
	fmt.Println(strings.Repeat("═", 65))
	if len(issues) == 0 {
		fmt.Println("✅ No issues found! Your project should be ready to run.")
	} else {
		fmt.Printf("Summary: %d error(s), %d warning(s), %d info\n",
			len(errors), len(warnings), len(infos))
		if len(errors) > 0 {
			fmt.Println("\n⚠️  Please resolve the errors above before running the project.")
		} else if len(warnings) > 0 {
			fmt.Println("\n⚠️  Consider addressing the warnings to ensure smooth operation.")
		}
	}
	fmt.Println(strings.Repeat("═", 65))
}
