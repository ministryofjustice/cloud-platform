package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
)

// GitHubAccess contains GitHub access information
type GitHubAccess struct {
	AccessToken string
	RepoOwner   string
	RepoName    string
	IssueNumber int
}

// ParseIssue parses the given issue and returns a slice of issues
func ParseIssue(issue string) []string {
	var issues []string
	sections := strings.Split(issue, "## Issue")

	for _, section := range sections {
		if section != "" {
			lines := strings.Split(strings.TrimSpace(section), "\n")
			title := strings.TrimSpace(lines[0])
			body := strings.TrimSpace(strings.Join(lines[1:], "\n"))
			issues = append(issues, fmt.Sprintf("## %s\n%s", title, body))
		}
	}
	return issues
}

// CreateIssue creates a GitHub issue
func CreateIssue(client *github.Client, ghAccess GitHubAccess, issue string, epicIssueNumber int) error {
	issueParts := strings.SplitN(issue, "\n", 3)
	title := strings.TrimSpace(strings.TrimPrefix(issueParts[1], "###"))
	body := strings.TrimSpace(issueParts[2])

	// Append the body with the epic issue number
	body = fmt.Sprintf("%s\n\n Related to: #%d", body, epicIssueNumber)

	issueRequest := &github.IssueRequest{
		Title: &title,
		Body:  &body,
	}
	_, _, err := client.Issues.Create(context.Background(), ghAccess.RepoOwner, ghAccess.RepoName, issueRequest)
	return err
}

func main() {

	var ghAccess GitHubAccess
	ghAccess.AccessToken = os.Getenv("GITHUB_TOKEN")
	flag.StringVar(&ghAccess.RepoOwner, "owner", "ministryofjustice", "the repository to create issues")
	flag.StringVar(&ghAccess.RepoName, "repo", "cloud-platform", "the repository to create issues")
	flag.IntVar(&ghAccess.IssueNumber, "issue", 0, "the issue number to create issues")

	flag.Parse()

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghAccess.AccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Fetch the issue contents from GitHub
	issueContent, _, err := client.Issues.Get(ctx, ghAccess.RepoOwner, ghAccess.RepoName, ghAccess.IssueNumber)
	if err != nil {
		fmt.Println("Error fetching issue template:", err)
		os.Exit(1)
	}

	issues := ParseIssue(*issueContent.Body)
	for _, issue := range issues {
		//fmt.Println("Creating issue:", issue)
		err := CreateIssue(client, ghAccess, issue, ghAccess.IssueNumber)
		if err != nil {
			fmt.Println("Error creating issue:", err)
			os.Exit(1)
		}
	}
}
