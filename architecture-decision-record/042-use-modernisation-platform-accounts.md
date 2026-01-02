# 42. Use Modernisation Platform accounts

Date: 2026-01-02

## Status

✅ Accepted

## Context

The Cloud Platform currently runs in a single AWS account. This was fine whilst the platform was smaller, but we are now hitting multiple AWS limits and the risk associated with blast radius of a single account grows as more services are hosted in there.

AWS accounts for the MoJ are now provisioned via the Modernisation Platform, independant AWS accounts such as the current Cloud Platform account are not part of the strategic vision of all AWS accounts being part of strategic platforms.

## Decision

When moving the Cloud Platform to a multi cluster architecture we will create multiple new accounts in the Modernisation Platform. The old Cloud Platform AWS account will then be decommissioned or migrated into the Modernisation Platform (decision to be made depending on if there are resources which cannot be migrated).

We will use an AWS account per live and non-live clusters (as per ADR 36 Multi Cluster decision).

## Consequences

- We will use the Modernisation Platform security baselines
- We will use the Modernisation Platform SSO for engineer access
- We will use the Modernisation deployment pipelines where possible
- We will use other Modernisation provided services where possible
