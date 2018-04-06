# Use ECR As Container Registry

Date: 06/04/2018

## Status

Pending

## Context

The cloud platforms team currently use Docker Registry for storage of docker images. This solution is self-hosted, needs patching, and occasionally downtime. 

example of issue [CPT-~~~](https://dsdmoj.atlassian.net/browse/CPT-~~~).


Reasons behind this move were:

* Reduce risk area by migrating to a managed service.
* Unnecessary downtimes
* Good value for money, preference for open source
* Can or is able to run addons i.e. security scanners

## Decision

Replace self hosted Docker Registry to managed ECR

## Consequences

1. Spike into what container registry service to use [CPT-479](https://dsdmoj.atlassian.net/browse/CPT-479). Nexus (self hosted), Artifactory (SaaS), ECR (managed)
2. Establish how concourse will push to the ECR container registry [CPT-CPT-511](https://dsdmoj.atlassian.net/browse/CPT-CPT-511)
3. Establish how concourse will push to the ~Nexus~ or Artifactory container registry [CPT-CPT-512](https://dsdmoj.atlassian.net/browse/CPT-CPT-512)
4. Decision made to use ECR [CPT-479](https://dsdmoj.atlassian.net/browse/CPT-479)
5. WIP 


~~Product teams have confirmed their success builds and deploys. zero downtime and minimal maintenance from Cloud Platforms side.~~