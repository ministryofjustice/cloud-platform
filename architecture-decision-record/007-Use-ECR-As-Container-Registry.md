# Use ECR As Container Registry

Date: 06/04/2018

## Status

✅ Accepted

## Context

The cloud platforms team currently use Docker Registry for storage of docker images. This solution is self-hosted, needs regular patching, and occasionally has downtime.

Example of an issue [CPT-274](https://dsdmoj.atlassian.net/browse/CPT-274).


We want to update the container registry to avoid some of the problems we have been seeing with it. The container registry will also be a key part of our new [Kubernetes based infrastructure](https://github.com/ministryofjustice/cloud-platform/blob/master/architecture-decision-record/004-use-kubernetes-for-container-management.md).


The criteria for selecting a new solution included:

* Finding a solution that would work with well GitHub based identity which is one of our [major architectural choices](https://github.com/ministryofjustice/cloud-platform/blob/master/architecture-decision-record/006-Use-github-as-user-directory.md)
* Decreasing the amount of day to day ops work, perhaps by using a managed service
* Good value for money
* Preference for open source tooling
* Can or is able to run addons i.e. security scanners
* We can make images that we host public by default

From this process we tentatively chose ECR. Unfortunately none of the SaaS registries in common use support Github authentication, but as ECR uses IAM for authentication, and IAM supports federated identity with Github via Auth0, it will in theory be possible to define ECR access policies that reference Github teams and roles.

## Decision

We will replace our self hosted Docker Registry to managed Elastic Container Registry (ECR).

## Consequences

We will go forward with implementing ECR as the container registry that we use. Our first aim will be to set this up as part of the [new build pipeline]().

@Kerin has done some proof of concept work on defining ECR access policies that reference Github teams and roles but we need to establish how we will map GitHub IDs to IAM roles we need to push and pull images before this decision is fully validated.

Following that we will replace the existing container registry with ECR.

## Notes

Some of the work that went into this decision:

1. Spike into what container registry service to use [CPT-479](https://dsdmoj.atlassian.net/browse/CPT-479). Nexus (self hosted), Artifactory (SaaS), ECR (managed)
2. Establish how concourse will push to the ECR container registry [CPT-511](https://dsdmoj.atlassian.net/browse/CPT-511)
3. Establish how concourse will push to the ~Nexus~ or Artifactory container registry [CPT-512](https://dsdmoj.atlassian.net/browse/CPT-512)
4. Decision made to use ECR [CPT-479](https://dsdmoj.atlassian.net/browse/CPT-479)
