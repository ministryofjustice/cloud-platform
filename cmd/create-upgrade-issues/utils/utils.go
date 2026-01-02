package utils

import (
	"fmt"
	"strings"
)

// ParseIssue parses the given issue and returns a slice of issues
func ParseIssue(issue string, upgradeVersion string) []string {
	var issues []string
	sections := strings.Split(issue, "## Issue")

	for _, section := range sections {
		if section != "" {
			lines := strings.Split(strings.TrimSpace(section), "\n")
			title := strings.TrimSpace(lines[0])
			body := strings.TrimSpace(strings.Join(lines[1:], "\n"))
			body = strings.ReplaceAll(body, "<upgrade-version>", upgradeVersion)
			issues = append(issues, fmt.Sprintf("## %s\n%s", title, body))
		}
	}
	return issues
}
