# 48. Per Business Unit Clusters

Date: 2026-06-05

## Status

✅ Accepted

## Context

Decision 036 reduces the blast radius with non live and live clusters, but this still leaves the noisy neighbour problem and potential security threats traversing business units.

## Decision

We will give a non-live and a live cluster for each business unit.

## Consequences

- Clusters will be isolated by business unit
- Clusters will be built from the same code so that they are identical to ease maintenance and troubleshooting.
- We will upgrade/add features to all non-live clusters before live clusters, reducing risk.
- Configuration between clusters will be largely identical to ease maintenance and troubleshooting.
- Blast radius will be reduced significately.
- We will need to NAT cluster IPs to ensure we have enough address space.
- Access to clusters and cluster accounts will be achieved with AWS Identity Center.
- Infrastructure and Application pipelines will need to be designed and built to push changes through in the desired order.
