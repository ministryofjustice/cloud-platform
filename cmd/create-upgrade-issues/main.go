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

// CreateIssue creates a GitHub issue
func CreateIssue(client *github.Client, ghAccess GitHubAccess, issue string, templateIssueNumber int, upgradeVersion string) error {
	issueParts := strings.SplitN(issue, "\n", 3)
	title := strings.TrimSpace(strings.TrimPrefix(issueParts[1], "###"))
	body := strings.TrimSpace(issueParts[2])

	// Append the body with the template issue number
	body = fmt.Sprintf("%s\n\n Related to: #%d", body, templateIssueNumber)
	label := "eks-" + upgradeVersion + "-upgrade"

	issueRequest := &github.IssueRequest{
		Title:  &title,
		Body:   &body,
		Labels: &[]string{label},
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
	upgradeVersion := flag.String("upgrade-version", "", "the target EKS upgrade version")

	flag.Parse()

	if *upgradeVersion == "" {
		fmt.Println("Error: you must specify the k8s upgrade version")
		os.Exit(1)
	}

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

	issues := ParseIssue(*issueContent.Body, *upgradeVersion)
	for _, issue := range issues {
		//fmt.Println("Creating issue:", issue)
		err := CreateIssue(client, ghAccess, issue, ghAccess.IssueNumber, *upgradeVersion)
		if err != nil {
			fmt.Println("Error creating issue:", err)
			os.Exit(1)
		}
	}
}
