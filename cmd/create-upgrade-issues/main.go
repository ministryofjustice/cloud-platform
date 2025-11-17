package main

import (
	"flag"
	"fmt"
	"ministryofjustice/cloud-platform/cmd/create-upgrade-issues/utils"
	"os"

	"github.com/google/go-github/v68/github"
)

var (
	upgradeVersion = flag.String("upgrade-version", "", "the target EKS upgrade version")
	ghAccess       utils.GitHubAccess
	template       string
)

func main() {
	flag.StringVar(&ghAccess.RepoOwner, "owner", "", "the repository owner")
	flag.StringVar(&ghAccess.RepoName, "repo", "", "the repository where to create issues")
	flag.StringVar(&ghAccess.Key, "key", os.Getenv("GHPKEY"), "GitHub App private key")
	flag.StringVar(&ghAccess.AppID, "appid", os.Getenv("GHAID"), "GitHub App ID")
	flag.StringVar(&ghAccess.InstallID, "installid", os.Getenv("GHIID"), "GitHub Installation ID")
	flag.StringVar(&template, "tmp", "template/cloud-platform-k8s-upgrade-template.md", "path to the issue template file")
	flag.Parse()

	if *upgradeVersion == "" {
		fmt.Println("Error: you must specify the k8s upgrade version")
		os.Exit(1)
	}

	if ghAccess.RepoOwner == "" {
		fmt.Println("Error: you must specify the repository owner")
		os.Exit(1)
	}

	if ghAccess.RepoName == "" {
		fmt.Println("Error: you must specify the repository name")
		os.Exit(1)
	}

	if ghAccess.Key == "" {
		fmt.Println("Error: GitHub App private key is required (set GHPKEY env var or use -key flag)")
		os.Exit(1)
	}

	if ghAccess.AppID == "" {
		fmt.Println("Error: GitHub App ID is required (set GHAID env var or use -appid flag)")
		os.Exit(1)
	}

	if ghAccess.InstallID == "" {
		fmt.Println("Error: GitHub Installation ID is required (set GHIID env var or use -installid flag)")
		os.Exit(1)
	}

	// Create issue titles with the parsed upgrade version
	issueTitles := []utils.Issues{
		// List of issue titles with their corresponding integers for assiging labels when issues are created
		{Title: "Planning upgrade to EKS " + *upgradeVersion, Int: 1},
		{Title: "Update vpc-cni to latest supported release for current EKS version", Int: 2},
		{Title: "Update kube-proxy to latest supported release for current EKS version", Int: 3},
		{Title: "Update core-dns to latest supported release for current EKS version", Int: 4},
		{Title: "Test EKS " + *upgradeVersion + " on test cluster", Int: 5},
		{Title: "Test EKS " + *upgradeVersion + " on live-like cluster", Int: 6},
		{Title: "EKS: Upgrade Production clusters to Kubernetes " + *upgradeVersion, Int: 7},
		{Title: "Update vpc-cni to latest supported release for upgraded EKS version", Int: 8},
		{Title: "Update kube-proxy to latest supported release for upgraded EKS version", Int: 9},
		{Title: "Update core-dns to latest supported release for upgraded EKS version", Int: 10},
		{Title: "Review cluster components for upgrading", Int: 11},
		{Title: "Review and upgrade kube-state-metrics for the upgraded EKS version " + *upgradeVersion, Int: 12},
		{Title: "Upgrade cluster autoscaler for EKS version " + *upgradeVersion, Int: 13},
		{Title: "Upgrade cluster descheduler for EKS version " + *upgradeVersion, Int: 14},
		{Title: "Post EKS version " + *upgradeVersion + " Cleanup", Int: 15},
		{Title: "Add deprecated apis from " + *upgradeVersion + " to Gatekeeper", Int: 16},
		{Title: "Update Upgrade runbook and Cluster Upgrade Issue Template", Int: 17},
	}

	client, err := utils.AppClient(ghAccess.Key, ghAccess.AppID, ghAccess.InstallID)
	if err != nil {
		fmt.Println("Error creating GitHub client:", err)
		os.Exit(1)
	}

	// Ensure all required labels exist with proper colors
	fmt.Println("Creating/updating labels...")
	err = utils.EnsureLabelsExist(client, ghAccess.RepoOwner, ghAccess.RepoName, *upgradeVersion)
	if err != nil {
		fmt.Println("Error creating/updating labels:", err)
		os.Exit(1)
	}

	// Read the issue contents from local file
	issueContentBytes, err := os.ReadFile(template)
	if err != nil {
		fmt.Println("Error reading issue template file:", err)
		os.Exit(1)
	}
	issueContentStr := string(issueContentBytes)

	ms, err := utils.CreateMilestone(client, ghAccess.RepoOwner, ghAccess.RepoName, &github.Milestone{
		Title:       github.Ptr("Kubernetes " + *upgradeVersion + " upgrade"),
		Description: github.Ptr(*upgradeVersion + " upgrade for kubernetes"),
		State:       github.Ptr("open"),
	})
	if err != nil {
		fmt.Println("Error creating milestone:", err)
		os.Exit(1)
	}

	issues := utils.ParseIssue(issueContentStr, *upgradeVersion)
	fmt.Printf("Parsed %d issues from template\n", len(issues))

	for i, issue := range issues {
		fmt.Printf("Creating issue %d/%d...\n", i+1, len(issues))
		err := utils.CreateIssue(client, ghAccess, issue, *upgradeVersion, ms.GetNumber(), issueTitles)
		if err != nil {
			fmt.Printf("Error creating issue %d: %v\n", i+1, err)
			os.Exit(1)
		}
		fmt.Printf("Successfully created issue %d\n", i+1)
	}

	fmt.Println("All issues created successfully!")
}
