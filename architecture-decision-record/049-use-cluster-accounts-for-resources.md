# 49. Use Cluster Accounts for Resources

Date: 2026-06-05

## Status

✅ Accepted

## Context

We will need somewhere to put user AWS resources.

## Options

1. Use the per business unit cluster accounts - Less AWS accounts and associated overhead, apps sit next to resources, blast radius is limited to business unit. Recommend approach by AWS.
2. Use a new AWS non-live and live resources accounts - More overhead, blast radius is still large.
3. User new per buiness unit resources accounts - Much more over head, increased distance between app and resouces, good isolation.

## Decision

Option 1 - We will use the cluster accounts for AWS resources.  This is a good balance between limiting blast radius and keeping the apps close to the resources.

## Consequences

- We need to ensure we have enough IPs in the cluster accounts for existing services and growth.
