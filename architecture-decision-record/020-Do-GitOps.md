# Do GitOps

Date: 27/11/2020

## Status

🤔 Proposed

## Introduction

At its basic level, "GitOps" is the practice of describing the system in code, reviewing changes using git requests, and deploying it automatically with CI/CD. Most dev teams do that to a certain extent, but in Cloud Platform we talk specifically about GitOps. This ADR is to crystalize our understanding of what it means to the Cloud Platform, its benefits, where we've applied it in the past, and how far we should take it in emerging challenges.

Out of this scope of this ADR are apps, which are defined and deployed by users.

## What is GitOps?

GitOps relies on existing practices:

* Application code + version control + CI/CD. Storing code in version control and using CI/CD to automatically test and deploy it.
* Infrastructure as Code + version control + CI/CD. In recent years the DevOps movement has extended app dev practices to the infrastructure - everything that the app runs on. We might write terraform code, store it in git and use CI/CD to run tests and deploy it.

These existing practise have lots of benefits:

* Having the whole stack as code in git makes it reproducible, reviewable and auditable.
* Testing and applying it all automatically with CI/CD builds tight coupling between the design in git and the production server. So you can trust that what is in git is what's on the server - git becomes the source of truth. This supports reliable deployments and quick delivery of changes to users.
* Automated tools can suggest improvements, such as linting terraform and checking for updated versions of terraform modules
* Tests are run in an environment that is in line with production, meaning that deployment risk is reduced.

[GitOps](https://rancher.com/blog/2020/gitops-kubernetes-connection/) takes these further:

* applies "everything as code + PRs + CI/CD" to a Kubernetes environment
* includes configuration of monitoring, alerts, logging

All these areas gain similar benefits: reviewable, auditable, testable, auto-suggestions, tight coupling between code and prod environment for low risk deploys.

GitOps also helps us to being able to treat the platform elements as [cattle not pets](/docs/Architecture-Official.md#cattle-not-pets).

## Context

In Cloud Platform, existing implementation of GitOps has meant:

* CP has [Terraform files](https://github.com/ministryofjustice/cloud-platform-infrastructure/tree/main/terraform) and associated [CI/CD pipelines]() [TODO add link] for applying it automatically

* kops and EKS clusters ...

* Helm charts

There remains number of questions and challenges to that model:

* Config needs storing somewhere - Terraform state, kops state, Helm values, secrets, differences between environments

* Is 'ClickOps' ever acceptable? By that we mean managing and provisioning a component manually, often by someone clicking around a web console. It is required by some components that CP has considered. However it's error prone and slow. Different environments will diverge / drift apart.

* 'Chicken and egg' problem with the CI/CD software itself

Cloud Platform also has Kubernetes resources, such as Prometheus, which are defined as Helm charts and config, which we aim to test and deploy with CI/CD.

## Decision

And now we have Kubernetes, we can manage the whole platform in the same way too - it's become known as **GitOps**.

GitOps is our aim. We define the Cloud Platform as code, store that in git and deploy it automatically. We include everything that the CP team manages, including the infrastructure, Kubernetes clusters and Kubernetes resources. And we support service teams to use GitOps for their apps.

Opportunities to aim for:

* Test the terraform and Helm charts, for example with Terratest
* Automated checks of terraform, Helm charts, for example with terraform lint and checking for updated Helm charts
* [TODO any more?]

[TODO address config, ClickOps and 'Chicken and egg' questions.]

GitOps is not suitable for some types of state, such as secrets and application databases.
