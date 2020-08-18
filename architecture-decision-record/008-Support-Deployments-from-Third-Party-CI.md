# Support Deployments from Third Party CI

Date: 03/05/2018

## Status

✅ Accepted

## Context

The cloud platforms team [is transitioning to Concourse CI](003-Use-Concourse-CI.md) internally. Product teams should also be able to use it, however, given that:

- Switching to Concourse CI would require familiarity with how it works and incur additional overhead
- We have not yet developed a streamlined approach to deployments through Concourse CI in order to confidently and properly offer support
- Product teams already use third party CI systems

We think it would be good as a starting point to make it easy for teams to deploy directly from the third party CI systems that teams are already using, rather than requiring the deployments to be implemented in Concourse CI.

## Decision

We will support deploying applications to the Cloud Platform from third party CI systems and will offer documentation on how to do so, at least for the most commonly used CI systems.

## Consequences

1. Cloud Platform adoption will be easier and quicker for product teams, since the required changes in their pipelines will be few.
2. Product teams can still use Concourse CI, but do not depend on it to get their applications deployed.
3. The Cloud Platform team has time to design an alternative based on Concourse CI.
