# Ministry of Justice Cloud Platform

## About this repository

This is the Ministry of Justice (MOJ) Cloud Platform team's repository for public facing documentation, feature work, enhancements, and issues.

The Cloud Platform team utilises [GitHub issues](https://github.com/ministryofjustice/cloud-platform/issues) to manage their work, and a team [ZenHub board](https://app.zenhub.com/workspaces/cloud-platform-team-5ccb0b8a81f66118c983c189/board) to surface GitHub issues into a product management board.

It's best to search our [GitHub issues](https://github.com/ministryofjustice/cloud-platform/issues) before adding new issues in an effort to reduce duplicates and encourage activity through existing conversations.

### Link checker

This repository has a GitHub action that checks all links in `*.md` files and creates a GitHub issue if the link returns a non-200 status code. If you have a link that doesn't resolve through the public internet (e.g. `127.0.0.1`, `localhost`, or other internal links), please update the `.ignore-links` file including the fully-qualified domain name (FQDN).

## Other Cloud Platform repositories

We have a series of repositories for our work. We have adopted the naming convention of naming each repository starting with `cloud-platform-*`. Where some repositories have similar purposes, we try to follow a common prefix (e.g. `cloud-platform-terraform-*` for Terraform modules). We also [name things](https://technical-guidance.service.justice.gov.uk/documentation/standards/naming-things.html#naming-things) so that users can understand what a repository does through its name.

### Core

| Name | Description |
|-|-|
| [Cloud Platform (this repository)](https://github.com/ministryofjustice/cloud-platform) | Public facing documentation, feature work, enhancements, and issues |
| [Cloud Platform environments](https://github.com/ministryofjustice/cloud-platform-environments) | User-created environments that are hosted on the Cloud Platform |
| [Cloud Platform infrastructure](https://github.com/ministryofjustice/cloud-platform-infrastructure) | Core infrastructure for the Cloud Platform |
| [Cloud Platform user guide](https://github.com/ministryofjustice/cloud-platform-user-guide) | User-focussed documentation for how to get started and use the Cloud Platform |

### Terraform modules

#### User-facing

| Name | Description |
|-|-|
| [Amazon MQ broker](https://github.com/ministryofjustice/cloud-platform-terraform-amq-broker) | Creates an Amazon MQ broker |
| [Database Migration Service (DMS)](https://github.com/ministryofjustice/cloud-platform-terraform-dms) | Creates a DMS replication instance to move data from another database to one inside Cloud Platform |
| [DynamoDB cluster](https://github.com/ministryofjustice/cloud-platform-terraform-dynamodb-cluster) | Creates a simple (i.e. non-global) DynamoDB cluster |
| [ECR credentials](https://github.com/ministryofjustice/cloud-platform-terraform-ecr-credentials) | Creates an ECR repository and associated IAM credentials |
| [ElastiCache cluster](https://github.com/ministryofjustice/cloud-platform-terraform-elasticache-cluster) | Creates an ElastiCache cluster |
| [ElasticSearch cluster](https://github.com/ministryofjustice/cloud-platform-terraform-elasticsearch) | Creates an ElasticSearch cluster |
| [GOV.UK Prototype](https://github.com/ministryofjustice/cloud-platform-terraform-github-prototype) | Creates the associated infrastructure to host a [GOV.UK Prototype](https://govuk-prototype-kit.herokuapp.com/docs) on Cloud Platform |
| [Kubernetes: IAM roles for service accounts (IRSA)](https://github.com/ministryofjustice/cloud-platform-terraform-irsa) | Creates an IAM role for a Kubernetes service account |
| [Kubernetes: service account](https://github.com/ministryofjustice/cloud-platform-terraform-serviceaccount) | Creates a Kubernetes service account, role, and rolebinding within a namespace|
| [Prometheus Pushgateway](https://github.com/ministryofjustice/cloud-platform-terraform-pushgateway) | Creates a [Prometheus Pushgateway](https://prometheus.io/docs/instrumenting/pushing/) |
| [RDS Aurora cluster](https://github.com/ministryofjustice/cloud-platform-terraform-rds-aurora) | Creates an RDS Aurora cluster |
| [RDS instance](https://github.com/ministryofjustice/cloud-platform-terraform-rds-instance) | Creates an RDS instance |
| [S3 bucket](https://github.com/ministryofjustice/cloud-platform-terraform-s3-bucket) | Creates an S3 bucket |
| [SNS topic](https://github.com/ministryofjustice/cloud-platform-terraform-sns-topic) | Creates an SNS topic |
| [SQS queue](https://github.com/ministryofjustice/cloud-platform-terraform-sqs) | Creates an SQS queue |

#### Team-facing

| Name | Description |
|-|-|
| [Bastion](https://github.com/ministryofjustice/cloud-platform-terraform-bastion) | Deploys a bastion instance |
| [CertManager](https://github.com/ministryofjustice/cloud-platform-terraform-certmanager) | Deploys [certmanager](https://cert-manager.io/docs/installation/) for automated TLS certificates |
| [Concourse](https://github.com/ministryofjustice/cloud-platform-terraform-concourse) | Deploys [ConcourseCI](https://concourse-ci.org/) within a Kubernetes cluster |
| [Ingress controller](https://github.com/ministryofjustice/cloud-platform-terraform-ingress-controller) | Deploys an [NGINX ingress controller](https://github.com/kubernetes/ingress-nginx) |
| [Logging](https://github.com/ministryofjustice/cloud-platform-terraform-logging) | Deploys standard logging tools such as [fluentbit](https://fluentbit.io/), etc. |
| [Monitoring](https://github.com/ministryofjustice/cloud-platform-terraform-monitoring) | Deploys standard monitoring tools such as [AlertManager](https://prometheus.io/docs/alerting/latest/alertmanager/), exporters, etc. |

### Other

#### Demonstration and reference applications

| Name | Description |
|-|-|
| [Multi-container app](https://github.com/ministryofjustice/cloud-platform-multi-container-demo-app) | Reference application for multi-container services |
| [Go app](https://github.com/ministryofjustice/cloud-platform-reference-app) | Reference application written in Go |
| [Ruby app](https://github.com/ministryofjustice/cloud-platform-helloworld-ruby-app) | Reference application written in Ruby |

#### Miscelleanous

| Name | Description |
|-|-|
| [Custom error pages](https://github.com/ministryofjustice/cloud-platform-custom-error-pages) | Customised error pages for uncaught routes |
| [Environments checker](https://github.com/ministryofjustice/cloud-platform-environments-checker) | Detects orphaned namespaces and AWS resources |
| [Helm charts](https://github.com/ministryofjustice/cloud-platform-helm-charts) | Custom Cloud Platform helm charts |
| [Kuberos](https://github.com/ministryofjustice/cloud-platform-kuberos) | A fork of [original Kuberos](https://github.com/negz/kuberos), managed by Cloud Platform |
| [Tools image](https://github.com/ministryofjustice/cloud-platform-tools-image) | Docker image containing tools used by pipelines |


## Useful links

It may be useful to look at:

- [Architecture Decision Records (ADRs) for the Cloud Platform](architecture-decision-record)
- [Cloud Platform internal runbooks](https://runbooks.cloud-platform.service.justice.gov.uk)
- [Cloud Platform user guide](https://user-guide.cloud-platform.service.justice.gov.uk)
