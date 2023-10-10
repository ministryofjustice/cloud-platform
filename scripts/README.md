# Creating issues for EKS upgrades

This directory contains scripts to automatically create GitHub issues based on EKS upgrades.

| Script                             | Description                                                                                                                                                                                                                      |
| ---------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [eks-updates.sh](./eks-updates.sh) | Fetches current cluster versions and compares them against supported Kubernetes versions in EKS. It creates a GitHub issue (example)[https://github.com/ministryofjustice/cloud-platform/issues/4857] to track upgrade progress. |
| [eks-addon-updates.sh](./eks-addon-updates.sh) | Fetches cluster add-on versions and compares them against the latest supported add-on version for the Kubernetes version a cluster is running. It creates a GitHub issue (example)[https://github.com/ministryofjustice/cloud-platform/issues/4867] to track upgrade progress.
