# Ministry of Justice Cloud Platform Master Repo

## About this Repository
This is the MoJ Cloud Platform team's repository for public facing
documentation, feature work, enhancements, and issues. ZenHub users can add issues
through GitHub Issues or directly on our [ZenHub
board](https://app.zenhub.com/workspaces/cloud-platform-team-5ccb0b8a81f66118c983c189/board)

It's best to search our board before adding new Issues in an effort to
reduce duplicates and encourage activity on existing conversations.

## Links

 * [Architecture Decision Records](architecture-decision-record)
 * [Runbooks]
 * [User Guide]

## Cloud platform repo list

We have a series of repositories for our work that we have listed below. We have adopted the naming convention of starting each repo with `cloud-platform-`. Where some repos have a similar purpose we try to name them similarly, (e.g. `cloud-platform-terraform-*` for terraform modules). We also try to [name things](https://technical-guidance.service.justice.gov.uk/documentation/standards/naming-things.html#naming-things) so that users can infer a basic understanding of what a given thing does from its name.

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
| Bastion | https://github.com/ministryofjustice/cloud-platform-terraform-bastion | A Terraform module to create a bastion inside an existing VPC |
| CertManager | https://github.com/ministryofjustice/cloud-platform-terraform-certmanager | Terraform module that deploys cloud-platform certmanager among another resources (like Ingress) |
| DynamoDB cluster | https://github.com/ministryofjustice/cloud-platform-terraform-dynamodb-cluster | A Terraform module to create a DynamoDB instance and IAM credentials, with optional autoscaling. |
| ECR credentials  | https://github.com/ministryofjustice/cloud-platform-terraform-ecr-credentials | Terraform module which creates ECR credentials and repository on AWS. |
| Elasticsearch cluster | https://github.com/ministryofjustice/cloud-platform-terraform-elasticsearch | A Terraform module to create a VPC based Elasticsearch cluster and Kibana dashboard along with a HTTP Proxy to access it. |
| Elaticache cluster | https://github.com/ministryofjustice/cloud-platform-terraform-elasticache-cluster | A Terraform module that users can create an ElastiCache cluster |
| Ingress controller | https://github.com/ministryofjustice/cloud-platform-terraform-ingress-controller | A Terraform module  that deploys cloud-platform ingress controllers among another resources (like certificates) |
| Logging | https://github.com/ministryofjustice/cloud-platform-terraform-logging | Terraform module that deploys cloud-platform logging solution. It includes components like: fluentd, eventrouter, circle-ci-stats |
| Monitoring | https://github.com/ministryofjustice/cloud-platform-terraform-monitoring | Terraform module that deploys cloud-platform monitoring solution. It has support for components like: proxy, thanos, cloudwatch datasource for grafana, side-car, ecr-exporter |
| RDS Aurora Cluster | https://github.com/ministryofjustice/cloud-platform-terraform-rds-aurora | Terraform module that will create an RDS Aurora Cluster, Primary DB Instance, KMS key, Database Subnet Group and credentials in AWS. |
| RDS instance | https://github.com/ministryofjustice/cloud-platform-terraform-rds-instance | Terraform module that will create an RDS instance, KMS key, Database Subnet Group and credentials in AWS. |
| S3 bucket |  https://github.com/ministryofjustice/cloud-platform-terraform-s3-bucket | Terraform module that will create an S3 bucket and credentials in AWS. |
| SNS topic |  https://github.com/ministryofjustice/cloud-platform-terraform-sns-topic | A Terraform module to create an SNS Topic in AWS, along with an IAM User to access it. |
| SQS resource |  https://github.com/ministryofjustice/cloud-platform-terraform-sqs | A Terraform module to provision SQS resources to the Cloud Platform. |
| ServiceAccount | https://github.com/ministryofjustice/cloud-platform-terraform-serviceaccount | Create a serviceaccount in a Cloud Platform namespace. |

## Other repos

| Name            | Link          | Description         |
|-----------------|---------------|---------------------|
| Cloud Platform Multi Container App | https://github.com/ministryofjustice/cloud-platform-multi-container-demo-app | A demo application featuring multiple components, for use in a deployment tutorial
| Cloud Platform Tools Image | https://github.com/ministryofjustice/cloud-platform-tools-image | Docker image containing all the tooling used by pipelines. |
| Reference App   | https://github.com/ministryofjustice/cloud-platform-reference-app  | A reference application to follow along with the cloud platform [user guide](https://ministryofjustice.github.io/cloud-platform-user-docs/#cloud-platform-user-guide). |
| Helloworld Ruby App | https://github.com/ministryofjustice/cloud-platform-helloworld-ruby-app | Minimal containerised hello world ruby application, to use as an example in the user docs. |
| Port-fowarding Image | https://github.com/ministryofjustice/cloud-platform-port-forward-docker-image | Small Docker image to forward network traffic as a non-root user, as described in the [RDS module instructions][rds-port-forward] |
| CircleCI Stats | https://github.com/ministryofjustice/cloud-platform-circleci-stats | Log CircleCI build and queue times to our Elasticsearch cluster |
| Environments Checker | https://github.com/ministryofjustice/cloud-platform-environments-checker | Code to find and delete 'orphaned' cluster namespaces, with no source code in the [environments repository], and their AWS resources. |
| Custom error pages | https://github.com/ministryofjustice/cloud-platform-custom-error-pages | Docker image which serves custom error pages (e.g. for 404 errors) |
| Namespace usage report | https://github.com/ministryofjustice/cloud-platform-namespace-usage-report | Web application to show cluster namespace usage charts |
| Helm Chart Repository | https://github.com/ministryofjustice/cloud-platform-helm-charts | Helm Chart repository to store internal helm charts used by the platform  |
| Kuberos | https://github.com/ministryofjustice/cloud-platform-kuberos | Kuberos fork that Cloud Platform team maintain |

## Link checker

This repository has a GitHub Action that checks all links in `*.md` files and creates an issue if the link returns 404 or 5xx or similar. If you have a link that doesn't resolve, please update the `.ignore-links` file containing the full FQDN.

[rds-port-forward]: https://github.com/ministryofjustice/cloud-platform-terraform-rds-instance#access-outside-the-cluster
[environments repository]: https://github.com/ministryofjustice/cloud-platform-environments
[Runbooks]: https://runbooks.cloud-platform.service.justice.gov.uk
[User Guide]: https://user-guide.cloud-platform.service.justice.gov.uk
