# Create Kubernetes Upgrade Issue

This script will generate a set of GitHub issues based on the content of the [Cloud Platform Kubernetes Upgrade Template ticket](https://raw.githubusercontent.com/ministryofjustice/cloud-platform/refs/heads/main/.github/ISSUE_TEMPLATE/cloud-platform-k8s-upgrade-template.md)

It works by parsing the template ticket and generating an individual issue for specific element of EKS upgrade work, splitting the issues by "## Issue" headers in the template. 

## Prerequisites

- Go installed on your machine
- GH token with required priveledges for creating issues in the Cloud Platform repo.

## Usage

- Create a [template ticket](https://github.com/ministryofjustice/cloud-platform/issues/new?template=cloud-platform-k8s-upgrade-template.md) for feeding into the script.

- Give the template ticket a title that refers to the target upgrade version. Once created take note of the template issue number.

- Set your GH Token in your environment with:

```sh
export GITHUB_TOKEN="your-token-here"
```

- Run the script, passing the required flags:

```sh
go run main.go -owner="ministryofjustice" -repo="cloud-platform" -issue=template issue number from above -upgrade-version="target-k8s-version"
```
