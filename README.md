# Ministry of Justice Cloud Platform Master Repo
[![Waffle.io - Columns and their card
count](https://badge.waffle.io/ministryofjustice/cloud-platform.svg?columns=all)](https://waffle.io/ministryofjustice/cloud-platform)

## About this Repository
This is the MoJ Cloud Platform team's repository for public facing
documentation, feature work, enhancements, and issues. Waffle users can add issues
through GitHub Issues or directly on our [Waffle
board](https://waffle.io/ministryofjustice/cloud-platform).

It's best to search our board before adding new Issues in an effort to
reduce duplicates and encourage activity on existing conversations.

## Cloud platform repo list

We have a series of repositories for our work that we have listed below. We have adopted the naming convention of starting each repo with `cloud-platform-`. Where some repos have a similar purpose we try to name them similarly, (e.g. `cloud-platform-terraform-*` for terraform modules). We also try to [name things](https://ministryofjustice.github.io/technical-guidance/standards/naming-things/#naming-things) so that users can infer a basic understanding of what a given thing does from its name.

## Our core repos

| Name            | Link          | Description         |
|-----------------|---------------|---------------------|
|Cloud platform   | https://github.com/ministryofjustice/cloud-platform  | MoJ cloud platform team's repository for public facing documentation, feature work, enhancements, and issues.  |
| Cloud platform user docs  | https://github.com/ministryofjustice/cloud-platform-user-docs  | The [documentation](https://ministryofjustice.github.io/cloud-platform-user-docs/#cloud-platform-user-guide) for users of the Ministry of Justice cloud platform. It explains how to deploy and run applications on the cloud platform.  |
| Cloud platform environments  | https://github.com/ministryofjustice/cloud-platform-environments  | Create environments on our Kubernetes clusters by adding your environment definition to this repo.  |
| Reference app   | https://github.com/ministryofjustice/cloud-platform-reference-app  | A reference application to follow along with the cloud platform [user guide](https://ministryofjustice.github.io/cloud-platform-user-docs/#cloud-platform-user-guide). |
| Kubernetes investigations   | https://github.com/ministryofjustice/kubernetes-investigations  | A place to collect our investigations into Kubernetes. |
| Cloud platforms K8s kops  | https://github.com/ministryofjustice/cloud-platforms-k8s-kops  | How we build Kubernetes clusters with Kops.  |
| Cloud platform concourse  | https://github.com/ministryofjustice/cloud-platform-concourse  | Concourse CI for cloud platform Kubernetes clusters.  |
| Cloud platforms AWS account  | https://github.com/ministryofjustice/cloudplatforms-aws-account  | A way of setting up a new AWS account for the cloud platform team.  |


## Terraform modules

| Name            | Link          | Description         |
|-----------------|---------------|---------------------|
| ECR credentials  | https://github.com/ministryofjustice/cloud-platform-terraform-ecr-credentials | Terraform module which creates ECR credentials and repository on AWS. |
| S3 bucket |  https://github.com/ministryofjustice/cloud-platform-terraform-s3-bucket | Terraform module that will create an S3 bucket and credentials in AWS. |
| RDS instance | https://github.com/ministryofjustice/cloud-platform-terraform-rds-instance | Terraform module that will create an RDS instance, KMS key, Database Subnet Group and credentials in AWS. |


## Other repos

| Name            | Link          | Description         |
|-----------------|---------------|---------------------|
| Cloud platform tectonic |  https://github.com/ministryofjustice/cloud-platforms-tectonic  |  A spike to create terraform resources for Tectonic Kubernetes cluster |
