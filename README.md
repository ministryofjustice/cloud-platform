# Ministry of Justice Cloud Platform Master Repo

## About this Repository
This is the MoJ Cloud Platform team's repository for public facing
documentation, feature work, enhancements, and issues. Waffle users can add issues
through GitHub Issues or directly on our [Waffle
board](https://waffle.io/ministryofjustice/cloud-platform).

It's best to search our board before adding new Issues in an effort to
reduce duplicates and encourage activity on existing conversations.

## Links

 * [Architecture Decision Records](architecture-decision-record)

## Cloud platform repo list

We have a series of repositories for our work that we have listed below. We have adopted the naming convention of starting each repo with `cloud-platform-`. Where some repos have a similar purpose we try to name them similarly, (e.g. `cloud-platform-terraform-*` for terraform modules). We also try to [name things](https://ministryofjustice.github.io/technical-guidance/standards/naming-things/#naming-things) so that users can infer a basic understanding of what a given thing does from its name.

## Our core repos

| Name            | Link          | Description         |
|-----------------|---------------|---------------------|
| Cloud Platform   | https://github.com/ministryofjustice/cloud-platform  | MoJ Cloud Platform team's repository for public facing documentation, feature work, enhancements, and issues.  |
| Cloud Platform Concourse  | https://github.com/ministryofjustice/cloud-platform-concourse  | Concourse CI for cloud platform Kubernetes clusters.  |
| Cloud Platform Environments  | https://github.com/ministryofjustice/cloud-platform-environments  | Create environments in the Cloud Platform by adding your environment definition to this repo.  |
| Cloud Platform Infrastructure   | https://github.com/ministryofjustice/cloud-platform-infrastructure  | Contains all infrastructure as code to create a Cloud Platform cluster. |
| Cloud Platform User Guide | https://github.com/ministryofjustice/cloud-platform-user-guide  | The [documentation](https://user-guide.cloud-platform.service.justice.gov.uk) for users of the Ministry of Justice Cloud Platform. It explains how to deploy and run applications on the Cloud Platform.  |

## Terraform modules

| Name            | Link          | Description         |
|-----------------|---------------|---------------------|
| Amazon MQ broker | https://github.com/ministryofjustice/cloud-platform-terraform-amq-broker | AWS MQ broker instance and credentials for the Cloud Platform |
| ECR credentials  | https://github.com/ministryofjustice/cloud-platform-terraform-ecr-credentials | Terraform module which creates ECR credentials and repository on AWS. |
| Elaticache cluster | https://github.com/ministryofjustice/cloud-platform-terraform-elasticache-cluster | A Terraform module that users can create an ElastiCache cluster |
| RDS instance | https://github.com/ministryofjustice/cloud-platform-terraform-rds-instance | Terraform module that will create an RDS instance, KMS key, Database Subnet Group and credentials in AWS. |
| S3 bucket |  https://github.com/ministryofjustice/cloud-platform-terraform-s3-bucket | Terraform module that will create an S3 bucket and credentials in AWS. |
| SQS resource |  https://github.com/ministryofjustice/cloud-platform-terraform-sqs | A Terraform module to provision SQS resources to the Cloud Platform. |

## Other repos

| Name            | Link          | Description         |
|-----------------|---------------|---------------------|
| Cloud Platform Multi Container App | https://github.com/ministryofjustice/cloud-platform-multi-container-demo-app | A demo application featuring multiple components, for use in a deployment tutorial
| Cloud Platform Smoke Tests | https://github.com/ministryofjustice/cloud-platform-smoke-tests | Cluster specific smoke tests to ensure a working environment. |
| Cloud Platform Tools Image | https://github.com/ministryofjustice/cloud-platform-tools-image | Docker image containing all the tooling used by pipelines. |
| Django Reference App   | https://github.com/ministryofjustice/cloud-platform-reference-app  | A reference application to follow along with the cloud platform [user guide](https://ministryofjustice.github.io/cloud-platform-user-docs/#cloud-platform-user-guide). |
| Helloworld Ruby App | https://github.com/ministryofjustice/cloud-platform-helloworld-ruby-app | Minimal containerised hello world ruby application, to use as an example in the user docs. |

