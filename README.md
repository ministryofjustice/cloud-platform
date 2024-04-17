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

| Name                                                                                                | Description                                                                   |
| --------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------- |
| [Cloud Platform (this repository)](https://github.com/ministryofjustice/cloud-platform)             | Public facing documentation, feature work, enhancements, and issues           |
| [Cloud Platform environments](https://github.com/ministryofjustice/cloud-platform-environments)     | User-created environments that are hosted on the Cloud Platform               |
| [Cloud Platform infrastructure](https://github.com/ministryofjustice/cloud-platform-infrastructure) | Core infrastructure for the Cloud Platform                                    |
| [Cloud Platform user guide](https://github.com/ministryofjustice/cloud-platform-user-guide)         | User-focussed documentation for how to get started and use the Cloud Platform |

### Terraform modules

#### User-facing

Find an up to date list in our [user guide](https://user-guide.cloud-platform.service.justice.gov.uk/documentation/deploying-an-app/add-aws-resources.html#available-modules)

#### Team-facing

| Name                                                                                                                                              | Description                                                                                                                         |
| ------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------- |
| [Auth0](https://github.com/ministryofjustice/cloud-platform-terraform-auth0)                                                                      | Creates the auth0 clients for the Kubernetes server and its components                                                              |
| [AWS Read-Only - SSO](https://github.com/ministryofjustice/cloud-platform-terraform-aws-sso)                                                      | Allow web console logins using Github credentials via SAML                                                                          |
| [AWS Account Baselines](https://github.com/ministryofjustice/cloud-platform-terraform-awsaccounts-baselines)                                      | Holds security and operational baselines implemented in Cloud Platform AWS accounts                                                 |
| [Bastion](https://github.com/ministryofjustice/cloud-platform-terraform-bastion)                                                                  | Deploys a bastion instance                                                                                                          |
| [CertManager](https://github.com/ministryofjustice/cloud-platform-terraform-certmanager)                                                          | Deploys [certmanager](https://cert-manager.io/docs/installation/) for automated TLS certificates                                    |
| [Cluster Autoscaler](https://github.com/ministryofjustice/cloud-platform-terraform-cluster-autoscaler)                                            | Deploys [Cluster Autoscaler](https://github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler)                               |
| [Concourse](https://github.com/ministryofjustice/cloud-platform-terraform-concourse)                                                              | Deploys [ConcourseCI](https://concourse-ci.org/) within a Kubernetes cluster                                                        |
| [Descheduler](https://github.com/ministryofjustice/cloud-platform-terraform-descheduler)                                                          | Deploys [Descheduler](https://github.com/kubernetes-sigs/descheduler#descheduler-for-kubernetes)                                    |
| [EKS Addons](https://github.com/ministryofjustice/cloud-platform-terraform-eks-add-ons)                                                           | Deploys Cloud Platform EKS Add ons                                                                                                  |
| [EFS CSI - ** NOT IN USE **](https://github.com/ministryofjustice/cloud-platform-terraform-efs-csi)                                               | Enables AWS EFS (NFS compatible) storage backend for Kubernetes                                                                     |
| [EKS CSI Storage](https://github.com/ministryofjustice/cloud-platform-terraform-eks-csi)                                                          | Enables EKS CSI storage backend for Kubernetes (EBS volumes)                                                                        |
| [External DNS](https://github.com/ministryofjustice/cloud-platform-terraform-external-dns)                                                        | Deploys external-dns to control DNS records dynamically                                                                             |
| [Global Auth0](https://github.com/ministryofjustice/cloud-platform-terraform-global-resources-auth0)                                              | Deploys Auth0 actions globally for auth0 tenant                                                                                     |
| [IAM Configuration](https://github.com/ministryofjustice/cloud-platform-terraform-awsaccounts-iam)                                                | Holds Cloud Platform team IAM configuration for AWS Accounts                                                                        |
| [Ingress controller](https://github.com/ministryofjustice/cloud-platform-terraform-ingress-controller)                                            | Deploys an [NGINX ingress controller](https://github.com/kubernetes/ingress-nginx)                                                  |
| [Kuberhealthy](https://github.com/ministryofjustice/cloud-platform-terraform-kuberhealthy)                                                        | Deploys [Kuberhealthy Operator](https://github.com/kuberhealthy/kuberhealthy) and custom checks                                     |
| [Kuberos](https://github.com/ministryofjustice/cloud-platform-terraform-kuberos)                                                                  | Deploys kuberos which enables users to authenticate to the cluster                                                                  |
| [Logging](https://github.com/ministryofjustice/cloud-platform-terraform-logging)                                                                  | Deploys standard logging tools such as [fluentbit](https://fluentbit.io/), etc.                                                     |
| [Monitoring](https://github.com/ministryofjustice/cloud-platform-terraform-monitoring)                                                            | Deploys standard monitoring tools such as [AlertManager](https://prometheus.io/docs/alerting/latest/alertmanager/), exporters, etc. |
| [Gatekeeper](https://github.com/ministryofjustice/cloud-platform-terraform-gatekeeper) Deploys Gatekeeper policy controller and required policies |
| [Starter Pack](https://github.com/ministryofjustice/cloud-platform-terraform-starter-pack)                                                        | Deploys Helloworld and multicontainer app                                                                                           |
| [Trivy Operator](https://github.com/ministryofjustice/cloud-platform-terraform-trivy-operator)                                                    | Deploys [Trivy Operator](https://aquasecurity.github.io/trivy-operator/v0.1.5/operator/installation/helm/)                          |
| [Velero](https://github.com/ministryofjustice/cloud-platform-terraform-velero) Deploys velero to manage backup and restore                        |
| [VPC Flow logs](https://github.com/ministryofjustice/cloud-platform-terraform-flow-logs)                                                          | Enables AWS Flow logs to capture information about the IP traffic going to and from network interfaces in VPC.                      |

### Other

#### Demonstration and reference applications

| Name                                                                                                | Description                                        |
| --------------------------------------------------------------------------------------------------- | -------------------------------------------------- |
| [Multi-container app](https://github.com/ministryofjustice/cloud-platform-multi-container-demo-app) | Reference application for multi-container services |
| [Go app](https://github.com/ministryofjustice/cloud-platform-reference-app)                         | Reference application written in Go                |
| [Ruby app](https://github.com/ministryofjustice/cloud-platform-helloworld-ruby-app)                 | Reference application written in Ruby              |

#### Miscelleanous

| Name                                                                                             | Description                                                                              |
| ------------------------------------------------------------------------------------------------ | ---------------------------------------------------------------------------------------- |
| [Canary](https://github.com/ministryofjustice/cloud-platform-terraform-canary)                   | Deploys AWS Synthetics Canary resource                                                   |
| [Custom error pages](https://github.com/ministryofjustice/cloud-platform-custom-error-pages)     | Customised error pages for uncaught routes                                               |
| [Environments checker](https://github.com/ministryofjustice/cloud-platform-environments-checker) | Detects orphaned namespaces and AWS resources                                            |
| [Helm charts](https://github.com/ministryofjustice/cloud-platform-helm-charts)                   | Custom Cloud Platform helm charts                                                        |
| [Kuberos](https://github.com/ministryofjustice/cloud-platform-kuberos)                           | A fork of [original Kuberos](https://github.com/negz/kuberos), managed by Cloud Platform |
| [Tools image](https://github.com/ministryofjustice/cloud-platform-tools-image)                   | Docker image containing tools used by pipelines                                          |

## Useful links

It may be useful to look at:

- [Architecture Decision Records (ADRs) for the Cloud Platform](architecture-decision-record)
- [Cloud Platform internal runbooks](https://runbooks.cloud-platform.service.justice.gov.uk)
- [Cloud Platform user guide](https://user-guide.cloud-platform.service.justice.gov.uk)
