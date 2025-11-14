package utils

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/go-github/v68/github"
	"github.com/jferrl/go-githubauth"
	"golang.org/x/oauth2"
)

type GitHubAccess struct {
	RepoOwner string
	RepoName  string
	Key       string
	AppID     string
	InstallID string
}

type Issues struct {
	Title string
	Int   int
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

// CreateIssue creates a GitHub issue
func CreateIssue(client *github.Client, ghAccess GitHubAccess, issue string, upgradeVersion string, milestone int, issueTitles []Issues) error {
	issueParts := strings.SplitN(issue, "\n", 3)
	if len(issueParts) < 3 {
		return fmt.Errorf("invalid issue format: expected at least 3 parts, got %d", len(issueParts))
	}
	title := strings.TrimSpace(strings.TrimPrefix(issueParts[1], "###"))
	body := strings.TrimSpace(issueParts[2])

	// Initialize labels as an empty slice instead of nil
	labels := make([]string, 0)
	found := false
	for _, issue := range issueTitles {
		if issue.Title == title {
			found = true
			switch issue.Int {
			case 1, 2, 3, 4, 5, 6:
				labels = append(labels, "eks-"+upgradeVersion+"-upgrade", "eks-pre-upgrade")
			case 7:
				labels = append(labels, "eks-"+upgradeVersion+"-upgrade")
			case 8, 9, 11, 12, 13, 14, 15, 16:
				labels = append(labels, "eks-"+upgradeVersion+"-upgrade", "eks-post-upgrade")
			default:
				fmt.Println("This issue isnt listed please add to Issues struct via variable issuesTitles, default label added:", title)
				labels = append(labels, "eks-"+upgradeVersion+"-upgrade")
			}
			break
		}
	}

	// If no matching title was found, add a default label
	if !found {
		fmt.Printf("Warning: Issue title not found in predefined list: '%s'\n", title)
		labels = append(labels, "eks-"+upgradeVersion+"-upgrade")
	}

	fmt.Printf("Creating issue '%s' with labels: %v\n", title, labels)

	issueRequest := &github.IssueRequest{
		Title:     &title,
		Body:      &body,
		Labels:    &labels,
		Milestone: &milestone,
	}
	_, _, err := client.Issues.Create(context.Background(), ghAccess.RepoOwner, ghAccess.RepoName, issueRequest)
	return err
}

// CreateOrUpdateLabel creates or updates a label with the specified color
func CreateOrUpdateLabel(client *github.Client, owner, repo, labelName, color, description string) error {
	ctx := context.Background()

	// Try to get existing label first
	existingLabel, _, err := client.Issues.GetLabel(ctx, owner, repo, labelName)
	if err != nil {
		// Label doesn't exist, create it
		label := &github.Label{
			Name:        github.Ptr(labelName),
			Color:       github.Ptr(color),
			Description: github.Ptr(description),
		}
		_, _, err := client.Issues.CreateLabel(ctx, owner, repo, label)
		if err != nil {
			return fmt.Errorf("error creating label '%s': %w", labelName, err)
		}
		fmt.Printf("Created label '%s' with color #%s\n", labelName, color)
	} else {
		// Label exists, update it if color is different
		if existingLabel.GetColor() != color {
			label := &github.Label{
				Name:        github.Ptr(labelName),
				Color:       github.Ptr(color),
				Description: github.Ptr(description),
			}
			_, _, err := client.Issues.EditLabel(ctx, owner, repo, labelName, label)
			if err != nil {
				return fmt.Errorf("error updating label '%s': %w", labelName, err)
			}
			fmt.Printf("Updated label '%s' with color #%s\n", labelName, color)
		}
	}
	return nil
}

// EnsureLabelsExist creates or updates all the required labels with their colors
func EnsureLabelsExist(client *github.Client, owner, repo, upgradeVersion string) error {
	labels := map[string]struct {
		color       string
		description string
	}{
		"eks-" + upgradeVersion + "-upgrade": {
			color:       "0052cc", // Blue
			description: "EKS " + upgradeVersion + " upgrade related issue",
		},
		"eks-pre-upgrade": {
			color:       "fbca04", // Yellow
			description: "Tasks to complete before EKS upgrade",
		},
		"eks-post-upgrade": {
			color:       "0e8a16", // Green
			description: "Tasks to complete after EKS upgrade",
		},
	}

	for labelName, labelInfo := range labels {
		err := CreateOrUpdateLabel(client, owner, repo, labelName, labelInfo.color, labelInfo.description)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateMilestone(client *github.Client, owner, repo string, milestone *github.Milestone) (*github.Milestone, error) {
	ctx := context.Background()
	milestone, _, err := client.Issues.CreateMilestone(ctx, owner, repo, milestone)
	if err != nil {
		return nil, fmt.Errorf("error creating milestone: %w", err)
	}
	return milestone, nil
}
