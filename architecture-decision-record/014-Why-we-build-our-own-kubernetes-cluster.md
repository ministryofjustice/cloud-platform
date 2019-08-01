# Why we build our own kubernetes cluster

Date: 01/08/2019

## Status

✅ Accepted

## Context

MoJ Cloud Platform team has decided to use [kubernetes for container management platform](https://github.com/ministryofjustice/cloud-platform/blob/master/architecture-decision-record/004-use-kubernetes-for-container-management.md) following the outcome of MOJ Digital's approach to infrastructure management. The team require the below features needed for the infrastructure management:

- Needed an universal authentiction mechanism to manage users without depending on the cloud provider
- Able to customize the control plane of kubernetes for MOJ requirement on Pod security
- Able to integrate external tools easily and be able to interfere and open the control plane for any custom changes

There are several leading cloud providers who provide managed production-ready kubernetes cluster:

- Amazon Elastic Kubernetes Service (Amazon EKS)
- Azure Kubernetes Service (AKS)
- Google Kubernetes Engine (GKE)

## Decision

* Our service team has good development experience working with AWS services and hence we decided to build our infrastructure in AWS. This made it easier for teams to migrate to the kubernetes platform. 

* When the time MOJ needed to build the kubernetes, Amazon EKS was still in the Alpha stage and was not production ready. Also Amazon EKS require to use IAM for user authentication which will be an overhead for managing users of service teams. 

* Kubernetes(k8s) allows to autheticate using OIDC and therefore it was easy to manage the authentication externally using Auth0.

## Consequences
- Extra management overhead

Despite the fact we are already running a production cluster we will continue to evaluate the managed solutions in the future.
