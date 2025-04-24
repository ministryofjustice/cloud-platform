package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/v68/github"
	"github.com/jferrl/go-githubauth"
	"golang.org/x/oauth2"
)

// GitHubAccess contains GitHub access information
type GitHubAccess struct {
	RepoOwner   string
	RepoName    string
	IssueNumber int
}

var (
	key       = flag.String("key", os.Getenv("GHPKEY"), "GitHub App private key")
	appid     = flag.String("appid", os.Getenv("GHAID"), "GitHub App ID")
	installid = flag.String("installid", os.Getenv("GHIID"), "GitHub Installation ID")
	ghAccess  GitHubAccess
)

func main() {
	flag.StringVar(&ghAccess.RepoOwner, "owner", "ministryofjustice", "the repository to create issues")
	flag.StringVar(&ghAccess.RepoName, "repo", "cloud-platform", "the repository to create issues")
	flag.IntVar(&ghAccess.IssueNumber, "issue", 0, "the issue number to create issues")
	upgradeVersion := flag.String("upgrade-version", "", "the target EKS upgrade version")

	flag.Parse()

	if *upgradeVersion == "" {
		fmt.Println("Error: you must specify the k8s upgrade version")
		os.Exit(1)
	}

	client, err := AppClient(*key, *appid, *installid)
	if err != nil {
		fmt.Println("Error creating GitHub client:", err)
		os.Exit(1)
	}
	ctx := context.Background()

	// Fetch the issue contents from GitHub
	issueContent, _, err := client.Issues.Get(ctx, ghAccess.RepoOwner, ghAccess.RepoName, ghAccess.IssueNumber)
	if err != nil {
		fmt.Println("Error fetching issue template:", err)
		os.Exit(1)
	}

	ms, err := CreateMilestone(client, ghAccess.RepoOwner, ghAccess.RepoName, &github.Milestone{
		Title:       github.Ptr("Kubernetes " + *upgradeVersion + " upgrade"),
		Description: github.Ptr(*upgradeVersion + " upgrade for kubernetes"),
		State:       github.Ptr("open"),
	})
	if err != nil {
		fmt.Println("Error creating milestone:", err)
		os.Exit(1)
	}

	issues := ParseIssue(*issueContent.Body, *upgradeVersion)
	for _, issue := range issues {
		err := CreateIssue(client, ghAccess, issue, ghAccess.IssueNumber, *upgradeVersion, ms.GetNumber())
		if err != nil {
			fmt.Println("Error creating issue:", err)
			os.Exit(1)
		}
	}
}

func AppClient(key, appid, installid string) (*github.Client, error) {
	privateKey := []byte(key)

	appIDInt, err := strconv.ParseInt(appid, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error converting app ID to int64: %w", err)
	}

	installIDInt, err := strconv.ParseInt(installid, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error converting installation ID to int64: %w", err)
	}

	appTokenSource, err := githubauth.NewApplicationTokenSource(appIDInt, privateKey)
	if err != nil {
		return nil, fmt.Errorf("error creating application token source: %w", err)
	}

	installationTokenSource := githubauth.NewInstallationTokenSource(installIDInt, appTokenSource)

	oauthHttpClient := oauth2.NewClient(context.Background(), installationTokenSource)

	client := github.NewClient(oauthHttpClient)
	return client, nil
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
func CreateIssue(client *github.Client, ghAccess GitHubAccess, issue string, templateIssueNumber int, upgradeVersion string, milestone int) error {
	issueParts := strings.SplitN(issue, "\n", 3)
	title := strings.TrimSpace(strings.TrimPrefix(issueParts[1], "###"))
	body := strings.TrimSpace(issueParts[2])

	// Append the body with the template issue number
	body = fmt.Sprintf("%s\n\n Related to: #%d", body, templateIssueNumber)
	label := "eks-" + upgradeVersion + "-upgrade"

	issueRequest := &github.IssueRequest{
		Title:     &title,
		Body:      &body,
		Labels:    &[]string{label},
		Milestone: &milestone,
	}
	_, _, err := client.Issues.Create(context.Background(), ghAccess.RepoOwner, ghAccess.RepoName, issueRequest)
	return err
}

func CreateMilestone(client *github.Client, owner, repo string, milestone *github.Milestone) (*github.Milestone, error) {
	ctx := context.Background()
	milestone, _, err := client.Issues.CreateMilestone(ctx, owner, repo, milestone)
	if err != nil {
		return nil, fmt.Errorf("error creating milestone: %w", err)
	}
	return milestone, nil
}
