# 020. Environments and Pipeline

Date: 2021-02-26

(This ADR is written retrospectively, to explain decisions from 2017/2018)

## Status

🤔 Proposed

## Context

The key proposition of Cloud Platform is to do the "hosting" of services, and we choose [Kubernetes for container management](004-use-kubernetes-for-container-management.md).

In agreeing a good interface for service teams, there several concerns:

* Definitions - teams should be able to specify the workloads and infrastructure they want running.

* Control - teams should be able to use a default hosting configuration, getting things running as simply as with a PaaS. However teams should also have full control over their Kubernetes resources, including pod configuration, lifecycle, network connectivity, etc.

* Multi-tenancy - Service teams' workloads need isolation between their dev and prod environments, and from other service teams' workloads.

## Decisions

1. Teams are offered 'namespaces'. A namespace is the concept of an isolated environment for workloads/resources.

2. A CP namespace is implemented as a Kubernetes namespace and AWS resources (e.g. RDS instance, S3 bucket).

3. Isolation in Kubernetes namespaces is implemented using RBAC and NetworkPolicy:

    * RBAC - teams can only administer k8s resources in their own namespaces
    * NetworkPolicy - containers can only receive traffic from its ingresses and other containers in the same namespace (implemented with a NetworkPolicy, which teams can edit if needed)

4. Isolation between AWS resources is achieved using access control.

    Each ECR repo, or S3 bucket, RDS bucket is made accessible to an IAM User, and the team are provided access key credentials for it.

5. A user defines a namespace in files: YAML (Kubernetes) and Terraform (AWS resources).

    The YAML includes by default: a Namespace and various default limits on resources, pods and networking.

    For deploying a simple workload, teams can include a YAML Deployment etc, so that these get applied automatically by CP's pipeline. Alternatively teams get more control by managing app resources using their namespace credentials - see below.

    The Terraform can specify any AWS resources like S3 buckets, RDS databases, Elasticache. Typically teams specify an ECR repo, so they have somewhere to deploy their images to.

6. The namespace definition is held in GitHub.

    GitHub provides a mechanism for peer-review, automated checks and versioning.

    Other options considered for configuring a namespace do not come with these advantages, for example:
    * a console / web form, implemented as a custom web app (click ops)
    * commands via a CLI or API

    Namespace definitions are stored in the [environments repo](https://github.com/ministryofjustice/cloud-platform-environments)

7. Namespace changes are checked by both a bot and a human from the CP team

    In Kubernetes, cluster-wide privileges are required to apply changes to a Kubernetes Namespace, as well as associated resources: LimitRange, NetworkPolicy and ServiceAccount. These privileges mean that the blast radius is large when applying changes.

    In terms of AWS resources, for common ones like S3 and RDS we provide terraform modules - to abstract away detail and promote best practice (for example, setting default encryption for S3 buckets). However Terraform can specify a huge range of AWS resources, each with multitude options. There are likely ways that one team can disrupt or get access to other teams' AWS services, that we can't anticipate, which is a risk to manage.

    To mitigate these concerns:
    * [automated checks](https://github.com/ministryofjustice/cloud-platform-environments/tree/main/.github/workflows) are used to validate against common problems
    * Human review (by an engineer on the CP team) is also required on PRs, to check against unanticipated problems

8. Pipeline to deploy namespace automatically.

    The "deploy pipeline" is a CI/CD pipeline that applies teams' namespace definitions in the clusters and AWS account. It triggers when the reviewed PR is merged to master.

9. Teams have full control within their Kubernetes namespace

    Users are given access to Kubernetes user credentials (kubecfg) with admin rights to their namespace. This gives them full control over their pods etc. They can deploy with 'kubectl apply' or Helm. They can debug problems with pod starting up, see logs etc.

    Users are also invited to create a ServiceAccount (using their environment YAML), and provide the creds to their CI/CD, for deploying their app.
