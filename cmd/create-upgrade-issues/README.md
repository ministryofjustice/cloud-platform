# Create Kubernetes Upgrade Issue

It works by parsing a template file and generating an individual issue for specific element of EKS upgrade work, splitting the issues by "## Issue" headers in the template.

## Prerequisites

- Go installed on your machine
- GH token with required priveledges for creating issues in the Cloud Platform repo.

## Usage
### Local Setup
- Set your GH App credentials in your environment with:

```sh
export GHPKEY="path-to-your-github-app-private-key.pem"
export GHAID="your-github-app-id"
export GHIID="your-github-installation-id"
```

- Run the script, passing the required flags:

```sh
go run main.go -owner="ministryofjustice" -repo="cloud-platform" -upgrade-version="target-k8s-version"
```

### Running Workflow Action
The GitHub Action workflow `create-upgrade-issues.yml` is set up to run this script automatically. It triggers on manual dispatch and requires the following inputs:
- `upgrade-version`: The target EKS upgrade version (e.g., "1.32")