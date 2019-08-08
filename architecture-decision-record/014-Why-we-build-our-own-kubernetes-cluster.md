# Why we build our own kubernetes cluster

Date: 01/08/2019

## Status

✅ Accepted

## Context

MoJ Cloud Platform team has decided to use [kubernetes for container management platform](https://github.com/ministryofjustice/cloud-platform/blob/master/architecture-decision-record/004-use-kubernetes-for-container-management.md) following the outcome of MOJ Digital's approach to infrastructure management. The team needed the below features for the infrastructure management:

- An universal authentication mechanism to manage users without depending on the cloud provider
- Able to customize the control plane of kubernetes for MOJ requirement on Pod security
- Able to integrate external tools easily
- Able to manage and configure the control plane for any custom changes

## Decision

There are several leading cloud providers who provide managed production-ready kubernetes cluster:

- Amazon Elastic Kubernetes Service (Amazon EKS)
- Azure Kubernetes Service (AKS)
- Google Kubernetes Engine (GKE)

We decided to host our cluster on AWS because our service team has good development experience working with AWS services. This made it easier for teams to migrate to the kubernetes platform

We decided to manage the kubernetes cluster ourselves rather than using EKS mainly for the below reasons:

- When the time MOJ needed to build the kubernetes, Amazon EKS was still in the Alpha stage and was not production ready. Also Amazon EKS require to use IAM for user authentication which will be an overhead for managing users of service teams

- Kubernetes(k8s) allows to authenticate using OIDC and therefore it was easy to manage the authentication externally using Auth0

## Consequences
- Extra management overhead

Despite the fact we are already running a production cluster we will continue to evaluate the managed solutions in the future.
