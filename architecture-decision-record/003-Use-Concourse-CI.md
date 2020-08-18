# Use Concourse CI

Date: 28/03/2018

## Status

✅ Accepted

## Context

The cloud platforms team currently use a self-host Jenkins server for CI/CD pipeline. This solution is not cloud native. There is a large amount of custom configuration which has equated to a vast amount of user request tickets and an opaque service which would be very difficult to recover [CPT-364](https://dsdmoj.atlassian.net/browse/CPT-364)  There is also a reoccurring issue which has caused a number of outages [CPT-384](https://dsdmoj.atlassian.net/browse/CPT-384).

Reasons behind this move were:

* Average of almost one week per month spent on debugging, fixing and reviving jenkins
* Unnecessary downtimes
* With the move to Kubernetes a cloud native CI/CD solution is needed, all jobs are written as code
* Non working jenkins was a blocker for product teams

## Decision

Replace self hosted Jenkins with self hosted Concourse CI pipeline

## Consequences

1. Spike into deploying Concourse CI into AWS or Kubernetes [CPT-486](https://dsdmoj.atlassian.net/browse/CPT-486)
2. Decision made to deploy to Kubernetes [CPT-488](https://dsdmoj.atlassian.net/browse/CPT-488)
3. Automate deployment of concourse using Terraform and Helm
4. WIP


~~Product teams have confirmed their success builds and deploys. zero downtime and minimal maintenance from Cloud Platforms side.~~
